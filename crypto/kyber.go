package crypto

import (
	"errors"
	"go.dedis.ch/kyber/v3"

	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/xof/blake2xb"

)

type HardenedKeys struct {
	group kyber.Group
	pubKey kyber.Point
	privKey kyber.Scalar
	suite *edwards25519.SuiteEd25519

}

func GemerateHardenedKeys() {
	rng := blake2xb.New(nil)
	suite := edwards25519.NewBlakeSHA256Ed25519WithRand(rng)
	h := &HardenedKeys{}
	sk := suite.Scalar().Pick(rng) // private key
	pk := suite.Point().Mul(sk, nil) // public key
	h.privKey = sk
	h.pubKey = pk
	h.suite = suite
}

func GenerateAndVerifySharedKeys(h1 *HardenedKeys, h2 *HardenedKeys) (secret string, err error) {

	S1 := h1.suite.Point().Mul(h1.privKey, h2.pubKey)
	S2 := h2.suite.Point().Mul(h2.privKey, h1.pubKey)

	if !S1.Equal(S2) {
		err = errors.New("shared secrets exchange didn't work")
		return "", err
	}

	return S1.String(), nil

}