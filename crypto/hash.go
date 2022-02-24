package crypto

import (
	"crypto"
	"encoding/binary"
	"hash"
	"math/big"
)

type Hash hash.Hash

type Hex struct {
	bigN  *big.Int
	bit  uint
	words []big.Word
	buf   []byte
}

func (h Hex) New() *Hex {

	buf := make([]byte,40)
	for i := range buf {
		buf[i] = 0
	}
	hex := &Hex{}
	hex.buf = buf[:]

	hex.bigN = new(big.Int)
	hex.bit = hex.bigN.Bit(0)
	hex.words = hex.bigN.Bits()
	return hex

}

func (h *Hex) ToString() string {
	return string(h.ToBytes())
}

func (h *Hex) ToBytes() []byte {
	return h.buf
}

func HexFromUint64(i uint64) *Hex {
		h := new(Hex)
		hex := h.New()
		b := make([]byte,40)
		binary.LittleEndian.PutUint64(b, i)
		copy(hex.buf, b)
		hex.bigN.SetBytes(hex.buf)
		return hex
}

func HexFromBytes(b []byte) *Hex {
	h := new(Hex)
	hex := h.New()
	b2 := make([]byte,40)
	copy(b2, b)
	copy(hex.buf, b2)
	hex.bigN.SetBytes(hex.buf)
	return hex
}

var hashFn = crypto.SHA3_256

func (h *Hex) Hash() []byte {

	hasher := hashFn.New()
	hasher.Write(h.buf)
	hasher.Write(h.bigN.Bytes())
	return hasher.Sum(nil)

}

func HashFromHex(h *Hex, content []byte) []byte {
	return HexFromBytes(content).Hash()
}


