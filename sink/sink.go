package sink

import "github.com/filecoin-project/go-filecoin/types"

type Sink struct {
}

var sink *Sink

func Init() {
	sink = &Sink{}
}

func HandleBlock(b types.Block) {
}

func HandleMessage(m types.SignedMessage, r types.MessageReceipt) {
}
