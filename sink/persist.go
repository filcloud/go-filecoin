package sink

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

func Persist() error {
	tx := sink.db.Begin()
	defer tx.End()

	// TODO

	return tx.Commit(true)
}
