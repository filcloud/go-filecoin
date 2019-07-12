package types

import (
	"math/big"

	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/filecoin-project/go-filecoin/address"
)

func init() {
	cbor.RegisterCborType(cbor.BigIntAtlasEntry)
	cbor.RegisterCborType(MinerState{})
	cbor.RegisterCborType(Ask{})
}

// Ask is a price advertisement by the miner
type Ask struct {
	Price  AttoFIL
	Expiry *BlockHeight
	ID     *big.Int
}

// State is the miner actors storage.
type MinerState struct {
	// Owner is the address of the account that owns this miner. Income and returned
	// collateral are paid to this address. This address is also allowed todd change the
	// worker address for the miner.
	Owner address.Address

	// Worker is the address of the worker account for this miner.
	// This will be the key that is used to sign blocks created by this miner, and
	// sign messages sent on behalf of this miner to commit sectors, submit PoSts, and
	// other day to day miner activities.
	Worker address.Address

	// PeerID references the libp2p identity that the miner is operating.
	PeerID peer.ID

	// ActiveCollateral is the amount of collateral currently committed to live
	// storage.
	ActiveCollateral AttoFIL

	// Asks is the set of asks this miner has open
	Asks      []*Ask
	NextAskID *big.Int

	// SectorCommitments maps sector id to commitments, for all sectors this
	// miner has committed.  Sector ids are removed from this collection
	// when they are included in the done or fault parameters of submitPoSt.
	// Due to a bug in refmt, the sector id-keys need to be
	// stringified.
	//
	// See also: https://github.com/polydawn/refmt/issues/35
	SectorCommitments SectorSet

	// Faults reported since last PoSt
	CurrentFaultSet IntSet

	// Faults reported since last PoSt, but too late to be included in the current PoSt
	NextFaultSet IntSet

	// NextDoneSet is a set of sector ids reported during the last PoSt
	// submission as being 'done'.  The collateral for them is still being
	// held until the next PoSt submission in case early sector removal
	// penalization is needed.
	NextDoneSet IntSet

	// ProvingSet is the set of sector ids of sectors this miner is
	// currently required to prove.
	ProvingSet IntSet

	LastUsedSectorID uint64

	// ProvingPeriodEnd is the block height at the end of the current proving period.
	// This is the last round in which a proof will be considered to be on-time.
	ProvingPeriodEnd *BlockHeight

	// The amount of space proven to the network by this miner in the
	// latest proving period.
	Power *BytesAmount

	// SectorSize is the amount of space in each sector committed to the network
	// by this miner.
	SectorSize *BytesAmount

	// SlashedSet is a set of sector ids that have been slashed
	SlashedSet IntSet

	// SlashedAt is the time at which this miner was slashed
	SlashedAt *BlockHeight

	// OwedStorageCollateral is the collateral for sectors that have been slashed.
	// This collateral can be collected from arbitrated deals, but not de-pledged.
	OwedStorageCollateral AttoFIL
}
