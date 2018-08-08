package types

import (
	"crypto/ecdsa"
	"fmt"

	cbor "gx/ipfs/QmSyK1ZiAP98YvnxsTfQpb669V2xeTHRbG4Y6fgKS3vVSd/go-ipld-cbor"

	"github.com/filecoin-project/go-filecoin/crypto"
	cu "github.com/filecoin-project/go-filecoin/crypto/util"
)

func init() {
	cbor.RegisterCborType(KeyInfo{})
}

// KeyInfo is a key and its type used for signing
type KeyInfo struct {
	// Private key as bytes
	PrivateKey []byte `json:"privateKey"`
	// Curve used to generate private key
	Curve string `json:"curve"`
}

// Unmarshal decodes raw cbor bytes into KeyInfo.
func (ki *KeyInfo) Unmarshal(b []byte) error {
	return cbor.DecodeInto(b, ki)
}

// Marshal KeyInfo into bytes.
func (ki *KeyInfo) Marshal() ([]byte, error) {
	return cbor.DumpObject(ki)
}

// Key returns the private key of KeyInfo
func (ki *KeyInfo) Key() []byte {
	return ki.PrivateKey
}

// Type returns the type of curve used to generate the private key
func (ki *KeyInfo) Type() string {
	return ki.Curve
}

// Equals returns true if the KeyInfo is equal to other.
func (ki *KeyInfo) Equals(other *KeyInfo) bool {
	if ki == nil && other == nil {
		return true
	}
	if ki == nil || other == nil {
		return false
	}
	if ki.Curve != other.Curve {
		return false
	}
	if len(ki.PrivateKey) != len(other.PrivateKey) {
		return false
	}
	for i := range ki.PrivateKey {
		if ki.PrivateKey[i] != other.PrivateKey[i] {
			return false
		}
	}
	return true
}

// Address returns the address for this keyinfo
func (ki *KeyInfo) Address() (Address, error) {
	prv, err := crypto.BytesToECDSA(ki.Key())
	if err != nil {
		return Address{}, nil
	}

	pub, ok := prv.Public().(*ecdsa.PublicKey)
	if !ok {
		// means a something is wrong with key generation
		return Address{}, fmt.Errorf("unknown public key type")
	}

	addrHash, err := AddressHash(cu.SerializeUncompressed(pub))
	if err != nil {
		return Address{}, err
	}

	// TODO: Use the address type we are running on from the config.
	return NewMainnetAddress(addrHash), nil
}
