package types

import "fmt"

type PoRepProofPartitions uint64

const (
	TestPoRepProofPartitions = PoRepProofPartitions(iota)
	TwoPoRepPartitions
)

func (p PoRepProofPartitions) ProofSize() *BytesAmount {
	switch p {
	case TestPoRepProofPartitions:
		return NewBytesAmount(192)
	case TwoPoRepPartitions:
		return NewBytesAmount(384)
	default:
		panic(fmt.Sprintf("PoRepProofPartitions#ProofSize(): unsupported value %v", p))
	}
}
