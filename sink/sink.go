package sink

import (
	"context"
	"fmt"
	"os"

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

type PorcelainAPI interface {
	ActorGet(ctx context.Context, addr address.Address) (*actor.Actor, error)
	MinerGetAsk(ctx context.Context, minerAddr address.Address, askID uint64) (miner.Ask, error)
	MinerGetState(ctx context.Context, minerAddr address.Address) (miner.State, error)
}

type Sink struct {
	db        *orm.DB
	porcelain PorcelainAPI

	inHandling      bool
	messagesInBlock bool
	cache           *Cache
}

var sink *Sink

func Init(porcelain PorcelainAPI) {
	dsn := os.Getenv("FIL_SINK_DSN")
	if dsn == "" {
		return // no sink
	}
	orm.Instantiate(dsn)
	sink = &Sink{
		db:        orm.Singleton,
		porcelain: porcelain,
	}
	sink.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&Actor{}, &TipSet{}, &Block{}, &Message{}, &SendMessage{}, &Miner{}, &MinerAsk{}, &Deal{}, &PaymentChannel{})

	consensus.MarkMessagesInBlock = MarkMessagesInBlock
	consensus.HandleMessagesInBlock = HandleMessagesInBlock
	consensus.HandleMessagesInTipSet = HandleMessagesInTipSet
	vm.HandleSendMessage = HandleSendMessage
	storagemarket.HandleCreateMiner = HandleCreateMiner
	miner.HandleAddMinerAsk = HandleAddMinerAsk
	state.HandleSetActor = HandleSetActor
}

func Begin() {
	if sink == nil {
		return
	}
	sink.cache = &Cache{}
}

func End() {
	if sink == nil {
		return
	}
}

func Commit() error {
	if sink == nil {
		return nil
	}
	if !sink.cache.IsEmpty() {
		err := Persist(context.Background())
		if err != nil {
			return err
		}
	}
	sink.cache = &Cache{}
	return nil
}

func BeginTipSet(tipSet types.TipSet) {
	if sink == nil {
		return
	}
	sink.inHandling = true
	sink.messagesInBlock = false

	sink.cache.TipSets = append(sink.cache.TipSets, BuildTipSet(tipSet))
}

func EndTipSet() {
	if sink == nil {
		return
	}
	sink.inHandling = false
}

func MarkMessagesInBlock() {
	if sink == nil {
		return
	}
	sink.messagesInBlock = true
}

func HandleMessagesInBlock(b *types.Block, r consensus.ApplyMessagesResponse) {
	if sink == nil || !sink.messagesInBlock {
		return
	}
	fmt.Println("######HandleMessagesInBlock", b, r)
}

func HandleMessagesInTipSet(b *types.Block, r consensus.ApplyMessagesResponse) {
	if sink == nil || sink.messagesInBlock {
		return
	}
	fmt.Println("#####HandleMessagesInTipSet", b, r)
}

func HandleSendMessage(m *types.Message) {
	if sink == nil || !sink.inHandling {
		return
	}
	sink.cache.SendMessages = append(sink.cache.SendMessages, BuildSendMessage(m))
	fmt.Println("#####HandleSendMessage", m)
}

func HandleCreateMiner(miner address.Address, m *miner.State) {
	if sink == nil || !sink.inHandling {
		return
	}
	sink.cache.Miners = append(sink.cache.Miners, BuildMiner(miner, m))
	fmt.Println("#####HandleCreateMiner", miner, m)
}

func HandleAddMinerAsk(miner address.Address, a *miner.Ask) {
	if sink == nil || !sink.inHandling {
		return
	}
	sink.cache.MinerAsks = append(sink.cache.MinerAsks, BuildMinerAsk(miner, a))
	fmt.Println("#####HandleAddMinerAsk", miner, a)
}

func HandleSetActor(a address.Address, actor *actor.Actor) {
	if sink == nil || !sink.inHandling {
		return
	}
	sink.cache.Actors = append(sink.cache.Actors, BuildActor(a, actor))
	fmt.Println("#####HandleSetActor", a, actor)
}
