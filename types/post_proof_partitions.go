package types

import "fmt"

type PoStProofPartitions uint64

const (
	TestPoStPartitions = PoStProofPartitions(iota)
	OnePoStPartition
)

func (p PoStProofPartitions) ProofSize() *BytesAmount {
	switch p {
	case TestPoStPartitions:
		return NewBytesAmount(192)
	case OnePoStPartition:
		return NewBytesAmount(192)
	default:
		panic(fmt.Sprintf("PoStProofPartitions#ProofSize(): unsupported value %v", p))
	}
}
