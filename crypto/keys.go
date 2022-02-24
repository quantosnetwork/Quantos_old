package crypto

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"log"
	"lukechampine.com/frand"
	_ "go.dedis.ch/kyber/v3"
	_ "go.dedis.ch/kyber/v3/suites"
)


const (
	// SizePublicKey is the size in bytes of a nodes/peers public key.
	SizePublicKey = ed25519.PublicKeySize

	// SizePrivateKey is the size in bytes of a nodes/peers private key.
	SizePrivateKey = ed25519.PrivateKeySize

	// SizeSignature is the size in bytes of a cryptographic signature.
	SizeSignature = ed25519.SignatureSize
)

type PublicKey [SizePublicKey]byte
type PrivateKey [SizePrivateKey]byte
type SharedKey interface{}


var (
	// ZeroPublicKey is the zero-value for a node/peer public key.
	ZeroPublicKey PublicKey

	// ZeroPrivateKey is the zero-value for a node/peer private key.
	ZeroPrivateKey PrivateKey

	// ZeroSignature is the zero-value for a cryptographic signature.
	ZeroSignature Signature
)


type Keys struct {
	pk PublicKey
	sk PrivateKey
	sharedKey SharedKey
	kSig Signature
}

func GenerateKeys(rand frand.RNG) *Keys {
	pub, priv, err := ed25519.GenerateKey(&rand)
	if err != nil {
		log.Fatal(err)
	}
	k := &Keys{}
	copy(k.pk[:], pub)
	copy(k.sk[:], priv)

	return k

}

func (k PublicKey) String() string {
	return hex.EncodeToString(k[:])
}

func (k PublicKey) Address() string {
	return "0x"+k.String()
}

func (k PublicKey) ToJSON() ([]byte, error) {
	return json.Marshal(k.String())
}

