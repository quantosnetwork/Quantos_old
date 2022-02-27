package crypto

import (
	"errors"
	"go.dedis.ch/kyber/v3"

	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/xof/blake2xb"
)

type HardenedKeys struct {
	Group   kyber.Group
	PubKey  kyber.Point
	PrivKey kyber.Scalar
	Suite   *edwards25519.SuiteEd25519
}

var PrivKey kyber.Scalar
var PubKey kyber.Point

func RestorePrivateKey(b []byte) {
	err := PrivKey.UnmarshalBinary(b)
	if err != nil {
		return
	}
}

func GenerateHardenedKeys() *HardenedKeys {
	rng := blake2xb.New(nil)
	suite := edwards25519.NewBlakeSHA256Ed25519WithRand(rng)
	h := &HardenedKeys{}
	sk := suite.Scalar().Pick(rng)   // private key
	pk := suite.Point().Mul(sk, nil) // public key
	h.PrivKey = sk
	h.PubKey = pk
	h.Suite = suite
	return h
}

func GenerateAndVerifySharedKeys(h1 *HardenedKeys, h2 *HardenedKeys) (secret string, err error) {

	S1 := h1.Suite.Point().Mul(h1.PrivKey, h2.PubKey)
	S2 := h2.Suite.Point().Mul(h2.PrivKey, h1.PubKey)

	if !S1.Equal(S2) {
		err = errors.New("shared secrets exchange didn't work")
		return "", err
	}

	return S1.String(), nil

}
