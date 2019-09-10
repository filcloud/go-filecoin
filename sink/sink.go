package sink

import (
	"github.com/cochainio/orm"
	"github.com/filecoin-project/go-filecoin/actor/builtin/miner"
	"github.com/filecoin-project/go-filecoin/actor/builtin/storagemarket"
	"github.com/filecoin-project/go-filecoin/address"
	"github.com/filecoin-project/go-filecoin/consensus"
	"github.com/filecoin-project/go-filecoin/types"
	"github.com/filecoin-project/go-filecoin/vm"
)

type Sink struct {
	currentTipSet   types.TipSet
	widen           bool
	messagesInBlock bool
	db              *orm.DB
	tx              *orm.TX
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
}

func BeginTipSet(tipSet types.TipSet) {
	sink.currentTipSet = tipSet
	sink.widen = false
	sink.messagesInBlock = false
	sink.tx = sink.db.Begin()
}

func EndTipSet() {
	sink.currentTipSet = types.UndefTipSet
	sink.tx.End()
	sink.tx = nil
}

func CommitTipSet() {
	_ = sink.tx.Commit()
}

func Widen() {
	sink.widen = true
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
}

func HandleCreateMiner(miner address.Address, m *miner.State) {

}

func HandleAddMinerAsk(miner address.Address, a *miner.Ask) {
}
