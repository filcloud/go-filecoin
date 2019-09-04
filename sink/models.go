package sink

import (
	"github.com/filecoin-project/go-filecoin/protocol/storage/storagedeal"
	"time"

	"github.com/filecoin-project/go-filecoin/actor"
	"github.com/filecoin-project/go-filecoin/types"
)

type Actor struct {
	Code         string
	Head         string
	Nonce        uint64
	BalanceBytes []byte
	Balance      float64
}

func BuildActor(a *actor.Actor) Actor {
	return Actor{
		Code:         a.Code.String(),
		Head:         a.Head.String(),
		Nonce:        uint64(a.Nonce),
		BalanceBytes: a.Balance.Bytes(),
		Balance:      attoToFloat64(a.Balance),
	}
}

/*
type TipSet struct {
	ID     string `gorm:"primary_key"`
	Height uint64 `gorm:"index"`
}

func BuildTipSet(t types.TipSet) TipSet {
	h, err := t.Height()
	if err != nil {
		panic(err)
	}
	return TipSet{
		ID:     t.String(),
		Height: h,
	}
}
*/

type Block struct {
	Cid             string `gorm:"primary_key"`
	Miner           string `gorm:"index"`
	Tickets         []byte
	Parents         string
	ParentWeight    uint64
	Height          uint64 `gorm:"index"`
	Messages        string
	StateRoot       string
	MessageReceipts string
	ElectionProof   []byte
	Timestamp       time.Time `gorm:"index"`
}

func BuildBlock(b types.Block) Block {
	return Block{
		Cid:             b.Cid().String(),
		Miner:           b.Miner.String(),
		Tickets:         toCbor(b.Tickets),
		Parents:         b.Parents.String(),
		ParentWeight:    uint64(b.ParentWeight),
		Height:          uint64(b.Height),
		Messages:        b.Messages.String(),
		StateRoot:       b.StateRoot.String(),
		MessageReceipts: b.MessageReceipts.String(),
		ElectionProof:   b.ElectionProof,
		Timestamp:       time.Unix(int64(uint64(b.Timestamp)), 0).UTC(),
	}
}

type Message struct {
	Cid        string `gorm:"primary_key"`
	To         string `gorm:"index:idx_to,idx_to_method"`
	From       string `gorm:"index"`
	Nonce      uint64
	Value      float64 `gorm:"index"`
	ValueBytes []byte
	Method     string `gorm:"index:idx_to_method"`
	Params     []byte

	GasPrice float64 `json:"gasPrice"`
	GasLimit uint64  `json:"gasLimit"`

	Signature []byte

	// receipt
	ExitCode uint8
	Return   []byte
	Gas      float64
}

func BuildMessage(m types.SignedMessage, r types.MessageReceipt) Message {
	cid, err := m.Cid()
	if err != nil {
		panic(err)
	}
	return Message{
		Cid:        cid.String(),
		To:         m.To.String(),
		From:       m.From.String(),
		Nonce:      uint64(m.Nonce),
		Value:      attoToFloat64(m.Value),
		ValueBytes: m.Value.Bytes(),
		Method:     m.Method,
		Params:     m.Params,

		GasPrice: attoToFloat64(m.GasPrice),
		GasLimit: uint64(m.GasLimit),

		Signature: m.Signature,


		ExitCode: r.ExitCode,
		Return:   toCbor(r.Return),
		Gas:      attoToFloat64(r.GasAttoFIL),
	}
}

/*
type MessageReceipt struct {
	ExitCode uint8
	Return   []byte
	Gas      float64
}

func BuildMessageReceipt(m types.MessageReceipt) MessageReceipt {
	return MessageReceipt{
		ExitCode: m.ExitCode,
		Return:   toCbor(m.Return),
		Gas:      attoToFloat64(m.GasAttoFIL),
	}
}
*/

type Deal struct {
	ProposalCid Cid
	Payer       Address
	Miner       Address
	State       DealState
	PieceRef    Cid
	Size        uint64
	TotalPrice  float64
	Duration    uint64
}

func BuildDeal(d storagedeal.Deal) Deal {
	if d.Proposal.Size.GreaterThan(types.NewBytesAmount(1<<64 - 1)) {
		panic("too big size")
	}
	return Deal{
		ProposalCid: Cid(d.Response.ProposalCid),
		Payer:       Address(d.Proposal.Payment.Payer),
		Miner:       Address(d.Miner),
		State:       DealState(d.Response.State),
		PieceRef:    Cid(d.Proposal.PieceRef),
		Size:        d.Proposal.Size.Uint64(),
		TotalPrice:  attoToFloat64(d.Proposal.TotalPrice),
		Duration:    d.Proposal.Duration,
	}
}
