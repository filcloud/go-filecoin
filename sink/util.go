package sink

import (
	"database/sql/driver"
	"fmt"
	"math/big"
	"strconv"

	"github.com/filecoin-project/go-filecoin/address"
	"github.com/filecoin-project/go-filecoin/protocol/storage/storagedeal"
	"github.com/filecoin-project/go-filecoin/types"
	"github.com/ipfs/go-cid"
	cbor "github.com/ipfs/go-ipld-cbor"
)

type Address address.Address

func (a *Address) Scan(src interface{}) error {
	addr, err := address.NewFromString(src.(string))
	if err != nil {
		return err
	}
	*a = Address(addr)
	return nil
}

func (a Address) Value() (driver.Value, error) {
	return address.Address(a).String(), nil
}

type Cid cid.Cid

func (c *Cid) Scan(src interface{}) error {
	s, err := cid.Decode(src.(string))
	if err != nil {
		return err
	}
	*c = Cid(s)
	return nil
}

func (c Cid) Value() (driver.Value, error) {
	return cid.Cid(c).String(), nil
}

type DealState storagedeal.State

func (s *DealState) Scan(src interface{}) error {
	var state storagedeal.State
	switch src.(string) {
	case "unset":
		state = storagedeal.Unset
	case "rejected":
		state = storagedeal.Rejected
	case "accepted":
		state = storagedeal.Accepted
	case "started":
		state = storagedeal.Started
	case "failed":
		state = storagedeal.Failed
	case "staged":
		state = storagedeal.Staged
	case "complete":
		state = storagedeal.Complete
	default:
		return fmt.Errorf("unrecognized %s", src)
	}
	*s = DealState(state)
	return nil
}

func (s DealState) Value() (driver.Value, error) {
	return storagedeal.State(s).String(), nil
}

type IntSet types.IntSet

func (i *IntSet) Scan(src interface{}) error {
	var s types.IntSet
	err := fromCbor(src.([]byte), &s)
	if err != nil {
		return err
	}
	*i = IntSet(s)
	return nil
}

func (i IntSet) Value() (driver.Value, error) {
	s := types.IntSet(i)
	return toCbor(s), nil
}

func toCbor(v interface{}) []byte {
	obj, err := cbor.WrapObject(v, types.DefaultHashFunction, -1)
	if err != nil {
		panic(err)
	}
	return obj.RawData()
}

func fromCbor(b []byte, v interface{}) error {
	return cbor.DecodeInto(b, v)
}

func attoToFloat64(v types.AttoFIL) float64 {
	f, err := strconv.ParseFloat(v.String(), 64)
	if err != nil {
		panic(err)
	}
	return f
}

func checkUint64(b *big.Int) {
	if !b.IsUint64() {
		panic("too big size")
	}
}
