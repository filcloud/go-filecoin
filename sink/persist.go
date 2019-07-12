package sink

import (
	"context"
	"fmt"
	"strings"

	"github.com/cochainio/orm/bulk_insert"
	"github.com/filecoin-project/go-filecoin/address"
)

type Cache struct {
	// History
	TipSets      []TipSet
	Blocks       []Block
	Messages     []Message
	SendMessages []SendMessage

	// State
	Actors    []Actor
	Miners    []Miner
	MinerAsks []MinerAsk

	// History, but local, not whole network
	Deals           []Deal
	PaymentChannels []PaymentChannel
}

func (c *Cache) IsEmpty() bool {
	return len(c.TipSets) == 0 && len(c.Blocks) == 0 &&
		len(c.Messages) == 0 && len(c.SendMessages) == 0 &&
		len(c.Actors) == 0 && len(c.Miners) == 0 &&
		len(c.MinerAsks) == 0 && len(c.Deals) == 0 &&
		len(c.PaymentChannels) == 0
}

func Persist(ctx context.Context) error {
	tx := sink.db.Begin()
	defer tx.End()

	c := sink.cache
	p := sink.porcelain

	heightMin := c.TipSets[0]
	heightMax := c.TipSets[len(c.TipSets)-1]
	_, _ = heightMin, heightMax
	var heights []uint64
	for _, t := range c.TipSets {
		heights = append(heights, t.Height)
	}

	err := tx.Where("Height IN (?)", heights).Delete(&TipSet{}).Error
	if err != nil {
		return err
	}
	err = tx.Where("Height IN (?)", heights).Delete(&Block{}).Error
	if err != nil {
		return err
	}
	err = tx.Where("Height IN (?)", heights).Delete(&Message{}).Error
	if err != nil {
		return err
	}
	err = tx.Where("Height IN (?)", heights).Delete(&SendMessage{}).Error
	if err != nil {
		return err
	}

	err = tx.BulkCreate(c.TipSets, bulk_insert.ReplaceOpt(true))
	if err != nil {
		return err
	}
	err = tx.BulkCreate(c.Blocks, bulk_insert.ReplaceOpt(true))
	if err != nil {
		return err
	}
	err = tx.BulkCreate(c.Messages, bulk_insert.ReplaceOpt(true))
	if err != nil {
		return err
	}
	err = tx.BulkCreate(c.SendMessages, bulk_insert.ReplaceOpt(true))
	if err != nil {
		return err
	}

	var existingActors, addActors []Actor
	var removeActors []Address
	err = tx.Where("UpdatedHeight >= ? AND UpdatedHeight <= ?", heightMin.Height, heightMax.Height).Find(&existingActors).Error
	if err != nil {
		return err
	}
	for _, a := range existingActors {
		actor, err := p.ActorGet(ctx, address.Address(a.Address))
		if err == nil {
			addActors = append(addActors, BuildActor(address.Address(a.Address), actor))
		} else if strings.Contains(err.Error(), "no actor at address") { // TODO:
			removeActors = append(removeActors, a.Address)
		} else {
			return err
		}
	}
	for _, a := range c.Actors {
		var existing bool
		for _, aa := range addActors {
			if a.Address == aa.Address {
				existing = true
				break
			}
		}
		if !existing {
			addActors = append(addActors, a)
			for i, r := range removeActors {
				if a.Address == r {
					or := removeActors
					removeActors = append(or[:i])
					if i+1 < len(or) {
						removeActors = append(removeActors, or[i+1:]...)
					}
					break
				}
			}
		}
	}
	err = tx.Where("Address IN (?)", removeActors).Delete(&Actor{}).Error
	if err != nil {
		return err
	}
	err = tx.BulkCreate(addActors, bulk_insert.ReplaceOpt(true))
	if err != nil {
		return err
	}

	var existingMiners, addMiners []Miner
	var removeMiners []Address
	err = tx.Where("UpdatedHeight >= ? AND UpdatedHeight <= ?", heightMin.Height, heightMax.Height).Find(&existingMiners).Error
	if err != nil {
		return err
	}
	for _, m := range existingMiners {
		state, err := p.MinerGetState(ctx, address.Address(m.Miner))
		if err == nil {
			addMiners = append(addMiners, BuildMiner(address.Address(m.Miner), &state))
		} else if strings.Contains(err.Error(), "failed to get To actor") {
			removeMiners = append(removeMiners, m.Miner)
		} else {
			return err
		}
	}
	for _, a := range c.Miners {
		var existing bool
		for _, aa := range addMiners {
			if a.Miner == aa.Miner {
				existing = true
				break
			}
		}
		if !existing {
			addMiners = append(addMiners, a)
			for i, r := range removeMiners {
				if a.Miner == r {
					or := removeMiners
					removeMiners = append(or[:i])
					if i+1 < len(or) {
						removeMiners = append(removeMiners, or[i+1:]...)
					}
					break
				}
			}
		}
	}
	err = tx.Where("Miner IN (?)", removeMiners).Delete(&Miner{}).Error
	if err != nil {
		return err
	}
	err = tx.BulkCreate(addMiners, bulk_insert.ReplaceOpt(true))
	if err != nil {
		return err
	}

	var existingAsks, addAsks []MinerAsk
	var removeAsksConditions []string
	var removeAsksArgs []interface{}
	err = tx.Where("UpdatedHeight >= ? AND UpdatedHeight <= ?", heightMin.Height, heightMax.Height).Find(&existingAsks).Error
	if err != nil {
		return err
	}
	for _, a := range existingAsks {
		ask, err := p.MinerGetAsk(ctx, address.Address(a.Miner), a.ID)
		if err == nil {
			addAsks = append(addAsks, BuildMinerAsk(address.Address(a.Miner), &ask))
		} else if strings.Contains(err.Error(), "no ask was found") { // TODO:
			removeAsksConditions = append(removeAsksConditions, "(?,?)")
			removeAsksArgs = append(removeAsksArgs, a.Miner, a.ID)
		} else {
			return err
		}
	}
	for _, a := range c.MinerAsks {
		var existing bool
		for _, aa := range addAsks {
			if a.Miner == aa.Miner {
				existing = true
				break
			}
		}
		if !existing {
			addAsks = append(addAsks, a)
			for i := 0; i < len(removeAsksArgs); i += 2 {
				if a.Miner == removeAsksArgs[i].(Address) && a.ID == removeAsksArgs[i+1].(uint64) {
					oldRemoveAsksConditions := removeAsksConditions
					oldRemoveAsksArgs := removeAsksArgs
					removeAsksConditions = append(oldRemoveAsksConditions[:i])
					removeAsksArgs = append(oldRemoveAsksArgs[:i])
					if i+2 < len(oldRemoveAsksConditions) {
						removeAsksConditions = append(removeAsksConditions, oldRemoveAsksConditions[i+2:]...)
						removeAsksArgs = append(removeAsksArgs, oldRemoveAsksArgs[i+2:]...)
					}
					break
				}
			}
		}
	}
	if len(removeAsksConditions) > 0 {
		cond := strings.Join(removeAsksConditions, ",")
		err = tx.Exec(fmt.Sprintf("DELETE FROM MinerAsk WHERE (Miner, ID) IN (%s)", cond), removeAsksArgs...).Error
		if err != nil {
			return err
		}
	}
	err = tx.BulkCreate(addAsks, bulk_insert.ReplaceOpt(true))
	if err != nil {
		return err
	}

	return tx.Commit(true)
}
