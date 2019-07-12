package sink

import (
	"time"

	"github.com/filecoin-project/go-filecoin/actor"
	"github.com/filecoin-project/go-filecoin/actor/builtin/miner"
	"github.com/filecoin-project/go-filecoin/actor/builtin/paymentbroker"
	"github.com/filecoin-project/go-filecoin/address"
	"github.com/filecoin-project/go-filecoin/protocol/storage/storagedeal"
	"github.com/filecoin-project/go-filecoin/types"
)

type Actor struct {
	Address      Address `gorm:"primary_key;type:varchar(41);"`
	Code         string
	Head         string
	Nonce        uint64
	BalanceBytes []byte
	Balance      float64

	UpdatedHeight uint64 `gorm:"index"`
	UpdatedAt     time.Time
}

func BuildActor(addr address.Address, a *actor.Actor) Actor {
	return Actor{
		Address:      Address(addr),
		Code:         a.Code.String(),
		Head:         a.Head.String(),
		Nonce:        uint64(a.Nonce),
		BalanceBytes: a.Balance.Bytes(),
		Balance:      attoToFloat64(a.Balance),
	}
}

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
	Cid         string `gorm:"primary_key"` // cid of signed message
	SuccinctCid string `gorm:"index"`       // cid of message
	To          string `gorm:"index:idx_to,idx_to_method;"`
	From        string `gorm:"index;"`
	Nonce       uint64
	Value       float64 `gorm:"index"`
	ValueBytes  []byte
	Method      string `gorm:"index:idx_to_method"`
	Params      []byte

	GasPrice float64 `json:"gasPrice"`
	GasLimit uint64  `json:"gasLimit"`

	Signature []byte

	// receipt
	ExitCode uint8
	Return   []byte
	Gas      float64

	Height uint64 `gorm:"index"`
}

func BuildMessage(m types.SignedMessage, r types.MessageReceipt) Message {
	cid, err := m.Cid()
	if err != nil {
		panic(err)
	}
	succinctCid, err := m.Message.Cid()
	if err != nil {
		panic(err)
	}
	return Message{
		Cid:         cid.String(),
		SuccinctCid: succinctCid.String(),
		To:          m.To.String(),
		From:        m.From.String(),
		Nonce:       uint64(m.Nonce),
		Value:       attoToFloat64(m.Value),
		ValueBytes:  m.Value.Bytes(),
		Method:      m.Method,
		Params:      m.Params,

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

type SendMessage struct {
	Cid    string  `gorm:"primary_key"` // cid of message
	To     string  `gorm:"index:idx_to,idx_to_method"`
	From   string  `gorm:"index"`
	Value  float64 `gorm:"index"`
	Method string  `gorm:"index:idx_to_method"`

	Height uint64 `gorm:"index"`
}

func BuildSendMessage(m *types.Message) SendMessage {
	cid, err := m.Cid()
	if err != nil {
		panic(err)
	}
	return SendMessage{
		Cid:    cid.String(),
		To:     m.To.String(),
		From:   m.From.String(),
		Value:  attoToFloat64(m.Value),
		Method: m.Method,
	}
}

type Miner struct {
	Miner                 Address `gorm:"primary_key;type:varchar(41);"`
	Owner                 Address `gorm:"index;type:varchar(41);"`
	Worker                Address `gorm:"index;type:varchar(41);"`
	PeerID                string
	ActiveCollateral      float64
	NextAskID             uint64
	ProvingSet            IntSet `gorm:"type:varbinary(255);"`
	LastUsedSectorID      uint64
	ProvingPeriodEnd      uint64
	Power                 uint64
	SectorSize            uint64
	SlashedAt             uint64
	OwedStorageCollateral float64

	UpdatedHeight uint64 `gorm:"index"`
	UpdatedAt     time.Time
}

func BuildMiner(miner address.Address, m *miner.State) Miner {
	var NextAskID, ProvingPeriodEnd, Power, SectorSize, SlashedAt uint64
	if m.NextAskID != nil {
		checkUint64(m.NextAskID)
		NextAskID = m.NextAskID.Uint64()
	}
	if m.ProvingPeriodEnd != nil {
		checkUint64(m.ProvingPeriodEnd.AsBigInt())
		ProvingPeriodEnd = m.ProvingPeriodEnd.AsBigInt().Uint64()
	}
	if m.Power != nil {
		checkUint64(m.Power.BigInt())
		Power = m.Power.Uint64()
	}
	if m.SectorSize != nil {
		checkUint64(m.SectorSize.BigInt())
		SectorSize = m.SectorSize.Uint64()
	}
	if m.SlashedAt != nil {
		checkUint64(m.SlashedAt.AsBigInt())
		SlashedAt = m.SlashedAt.AsBigInt().Uint64()
	}
	return Miner{
		Miner:                 Address(miner),
		Owner:                 Address(m.Owner),
		Worker:                Address(m.Worker),
		PeerID:                m.PeerID.String(),
		ActiveCollateral:      attoToFloat64(m.ActiveCollateral),
		NextAskID:             NextAskID,
		ProvingSet:            IntSet(m.ProvingSet),
		LastUsedSectorID:      m.LastUsedSectorID,
		ProvingPeriodEnd:      ProvingPeriodEnd,
		Power:                 Power,
		SectorSize:            SectorSize,
		SlashedAt:             SlashedAt,
		OwedStorageCollateral: attoToFloat64(m.OwedStorageCollateral),
	}
}

type MinerAsk struct {
	Miner  Address `gorm:"primary_key;type:varchar(41);"`
	ID     uint64  `gorm:"primary_key;type:varchar(41);"`
	Price  float64
	Expiry uint64

	UpdatedHeight uint64 `gorm:"index"`
	UpdatedAt     time.Time
}

func BuildMinerAsk(miner address.Address, a *miner.Ask) MinerAsk {
	checkUint64(a.ID)
	checkUint64(a.Expiry.AsBigInt())
	return MinerAsk{
		Miner:  Address(miner),
		ID:     a.ID.Uint64(),
		Price:  attoToFloat64(a.Price),
		Expiry: a.Expiry.AsBigInt().Uint64(),
	}
}

// Since deals are only stored on client and miner side, `Deal` should only be used in testing environment.
type Deal struct {
	ProposalCid Cid       `gorm:"primary_key;type:varchar(41);"`
	Payer       Address   `gorm:"type:varchar(41);"`
	Miner       Address   `gorm:"type:varchar(41);"`
	State       DealState `gorm:"type:varchar(20);"`
	PieceRef    Cid       `gorm:"type:varchar(46);"`
	Size        uint64
	TotalPrice  float64
	Duration    uint64
}

func BuildDeal(d storagedeal.Deal) Deal {
	checkUint64(d.Proposal.Size.BigInt())
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

type PaymentChannel struct {
	ChannelID      uint64  `gorm:"primary_key"`
	Payer          Address `gorm:"primary_key;type:varchar(41);"`
	Target         Address `gorm:"type:varchar(41);"`
	Amount         float64
	AmountRedeemed float64
	AgreedEol      uint64
	Eol            uint64
	Redeemed       bool
}

func BuildPaymentChannel(payer address.Address, channelID uint64, p paymentbroker.PaymentChannel) PaymentChannel {
	checkUint64(p.AgreedEol.AsBigInt())
	checkUint64(p.Eol.AsBigInt())
	return PaymentChannel{
		Payer:          Address(payer),
		ChannelID:      channelID,
		Target:         Address(p.Target),
		Amount:         attoToFloat64(p.Amount),
		AmountRedeemed: attoToFloat64(p.AmountRedeemed),
		AgreedEol:      p.AgreedEol.AsBigInt().Uint64(),
		Eol:            p.Eol.AsBigInt().Uint64(),
		Redeemed:       p.Redeemed,
	}
}
