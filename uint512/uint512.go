package uint512

import (
	"github.com/holiman/uint256"
	"math/big"
)

type Uint512 struct {
	a, b *uint256.Int
}

func NewUint512(a1, b1 []byte) *Uint512 {

	ab1 := ToBigInt(a1)
	ab2 := ToBigInt(b1)
	u1, _ := uint256.FromBig(ab1)
	u2, _ := uint256.FromBig(ab2)
	return &Uint512{
		u1, u2,
	}
}

func ToBigInt(bb []byte) *big.Int {

	bbb := new(big.Int)
	bbb.SetBytes(bb)
	return bbb

}

func (u512 *Uint512) Mul(a, b *uint256.Int) (r0, r1 *uint256.Int) {
	mmu := uint256.NewInt(0)
	mm := mmu.MulMod(a, b, uint256.NewInt(1))
	r0 = mmu.Mul(a, b)
	r1 = mmu.Sub(mm, r0)
	r0.Lt(mm)
	return
}
