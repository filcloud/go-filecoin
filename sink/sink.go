package sink

import (
	"github.com/cochainio/orm"
	"github.com/filecoin-project/go-filecoin/actor"
	"github.com/filecoin-project/go-filecoin/actor/builtin/miner"
	"github.com/filecoin-project/go-filecoin/actor/builtin/storagemarket"
	"github.com/filecoin-project/go-filecoin/address"
	"github.com/filecoin-project/go-filecoin/consensus"
	"github.com/filecoin-project/go-filecoin/state"
	"github.com/filecoin-project/go-filecoin/types"
	"github.com/filecoin-project/go-filecoin/vm"
)

type Sink struct {
	db *orm.DB

	inHandling      bool
	messagesInBlock bool
	cache           *Cache
}

var sink *Sink

func Init() {
	sink = &Sink{
		db: orm.Singleton,
	}
	consensus.MarkMessagesInBlock = MarkMessagesInBlock
	consensus.HandleMessagesInBlock = HandleMessagesInBlock
	consensus.HandleMessagesInTipSet = HandleMessagesInTipSet
	vm.HandleSendMessage = HandleSendMessage
	storagemarket.HandleCreateMiner = HandleCreateMiner
	miner.HandleAddMinerAsk = HandleAddMinerAsk
	state.HandleSetActor = HandleSetActor
}

func Begin() {
	sink.cache = &Cache{}
}

func End() {
}

func Commit() error {
	if !sink.cache.IsEmpty() {
		err := Persist()
		if err != nil {
			return err
		}
	}
	sink.cache = &Cache{}
	return nil
}

func BeginTipSet(tipSet types.TipSet) {
	sink.inHandling = true
	sink.messagesInBlock = false

	sink.cache.TipSets = append(sink.cache.TipSets, BuildTipSet(tipSet))
}

func EndTipSet() {
	sink.inHandling = false
}

func MarkMessagesInBlock() {
	sink.messagesInBlock = true
}

func HandleMessagesInBlock(b *types.Block, r consensus.ApplyMessagesResponse) {
	if !sink.messagesInBlock {
		return
	}
}

func HandleMessagesInTipSet(b *types.Block, r consensus.ApplyMessagesResponse) {
	if sink.messagesInBlock {
		return
	}
}

func HandleSendMessage(m *types.Message) {
	if !sink.inHandling {
		return
	}
	sink.cache.SendMessages = append(sink.cache.SendMessages, BuildSendMessage(m))
}

func HandleCreateMiner(miner address.Address, m *miner.State) {
	if !sink.inHandling {
		return
	}
	sink.cache.Miners = append(sink.cache.Miners, BuildMiner(miner, m))
}

func HandleAddMinerAsk(miner address.Address, a *miner.Ask) {
	if !sink.inHandling {
		return
	}
	sink.cache.MinerAsks = append(sink.cache.MinerAsks, BuildMinerAsk(miner, a))
}

func HandleSetActor(a address.Address, actor *actor.Actor) {
	if !sink.inHandling {
		return
	}
	sink.cache.Actors = append(sink.cache.Actors, BuildActor(a, actor))
}
