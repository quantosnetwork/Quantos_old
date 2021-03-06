package crypto

import (
	"encoding/hex"
	"go.dedis.ch/kyber/v3/sign/schnorr"
)

func (h *HardenedKeys) Sign(msg []byte) []byte {
	sign, err := schnorr.Sign(h.Suite, h.PrivKey, msg)
	if err != nil {
		return nil
	}
	return sign
}

type Signature []byte

func (hs Signature) String() string {
	return hex.EncodeToString(hs)
}

func (h *HardenedKeys) VerifySignature(msg, signature []byte) bool {
	err := schnorr.Verify(h.Suite, h.PubKey, msg, signature)
	if err != nil {
		return false
	}
	return true
}
