package sink

import "github.com/filecoin-project/go-filecoin/types"

type Sink struct {
	currentTipSet types.TipSet
}

var sink *Sink

func Init() {
	sink = &Sink{}
}

func BeginHandleTipSet() {
}

func EndHandleTipSet() {
}

func HandleBlock(b types.Block) {
}

func HandleMessage(m types.SignedMessage, r types.MessageReceipt) {
}
