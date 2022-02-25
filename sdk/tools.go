package sdk

import (
	"github.com/holiman/uint256"
	"math/big"
)

func Uint256StringFromBytes(b []byte) string {
	b1 := new(big.Int)
	b1.SetBytes(b)
	s, _ := uint256.FromBig(b1)
	return s.String()
}

func Uint256BytesFromHex(hex string) []byte {
	s, _ := uint256.FromHex(hex)
	return s.Bytes()
}

func MakeNewUint256(b []byte) *uint256.Int {
	b1 := new(big.Int)
	b1.SetBytes(b)
	s, _ := uint256.FromBig(b1)
	return s
}
