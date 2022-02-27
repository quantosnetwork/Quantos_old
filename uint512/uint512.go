package uint512

import (
	"Quantos/crypto"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/holiman/uint256"
	"github.com/zeebo/blake3"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/sign/schnorr"
	"lukechampine.com/frand"
	"math/big"
	"time"
)

type Uint512 struct {
	a, b *uint256.Int
}

type Int struct {
	b []byte
}

func NewUint512FromBytes(a1, b1 []byte) *Uint512 {

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

func (u512 *Uint512) Merge() *Int {
	buf := make([][]byte, 2)
	buf[0] = u512.a.Bytes()
	buf[1] = u512.b.Bytes()
	nInt := bytes.Join(buf, nil)
	var bufInt Int
	copy(bufInt.b, nInt[:])
	return &bufInt
}

func (uInt *Int) ToUint512Struct() *Uint512 {

	a := uInt.Bytes()[0:255]
	b := uInt.Bytes()[255:]
	bb := new(big.Int).SetBytes(a[:])
	ba := new(big.Int).SetBytes(b[:])
	bba, _ := uint256.FromBig(bb)
	bbb, _ := uint256.FromBig(ba)
	return &Uint512{bba, bbb}

}

func NewIntFromUint64s(a, b uint64) (*Uint512, *Int) {

	aa := uint256.NewInt(a)
	bb := uint256.NewInt(b)
	_uint := &Uint512{aa, bb}
	return _uint, _uint.Merge()

}

func NewBlankUint512() (*Uint512, *Int) {
	aa := uint256.NewInt(0)
	bb := uint256.NewInt(0)
	_uint := &Uint512{aa, bb}
	return _uint, _uint.Merge()

}

func (u512 *Uint512) FromBig(a, b *big.Int) *Int {
	aa, _ := uint256.FromBig(a)
	bb, _ := uint256.FromBig(b)
	_uint := &Uint512{aa, bb}
	return _uint.Merge()
}

func (uInt *Int) ToHex() string {

	return hex.EncodeToString(uInt.Bytes())

}

func (uInt *Int) ToString() string {
	return "0x" + uInt.ToHex()
}

func (uInt *Int) Bytes() []byte {
	return uInt.b
}

func (uInt *Int) Hash() []byte {
	hasher := blake3.Hasher{}
	hasher.Write(uInt.Bytes()[:])
	sum := hasher.Sum(nil)
	return sum
}

func (uInt *Int) KeyedSignedHash(key kyber.Scalar) (hash []byte, sig []byte) {

	// to make a key we take the sha256 of the private key (kyber.Scalar)
	k, _ := key.MarshalBinary()
	h := sha256.Sum256(k)
	hk, _ := blake3.NewKeyed(h[:])
	hk.Write(uInt.Bytes()[:])
	hash = hk.Sum(nil)
	var suite *edwards25519.SuiteEd25519
	sig, _ = schnorr.Sign(suite, key, hash)
	return
}

func (uInt *Int) Sign(key kyber.Scalar, msg []byte) []byte {
	var suite *edwards25519.SuiteEd25519
	sig, _ := schnorr.Sign(suite, key, msg)
	return sig
}

func (uInt *Int) VerifySignedContent(key kyber.Point, msg []byte, sig []byte) bool {
	var group kyber.Group
	err := schnorr.Verify(group, key, msg, sig)
	if err != nil {
		return false
	}
	return true
}

type address struct {
	*Uint512
}

type Address struct {
	Raw             *address
	pk              []byte
	sk              []byte
	ssk             []byte
	group           kyber.Group
	suite           *edwards25519.SuiteEd25519
	Signature       []byte
	TimestampSigned []byte
	Timestamp       int64
}

func (addr *address) Create() *Address {
	seed := make([]byte, 512)
	frand.Read(seed)
	seedA := seed[:255]
	seedB := seed[255:]
	u := new(Uint512)
	sab := new(big.Int).SetBytes(seedA)
	sbb := new(big.Int).SetBytes(seedB)
	u.a, _ = uint256.FromBig(sab)
	u.b, _ = uint256.FromBig(sbb)
	add := new(Address)
	add.Raw = &address{u}

	k := crypto.GenerateHardenedKeys()
	add.pk, _ = k.PubKey.MarshalBinary()
	add.sk, _ = k.PrivKey.MarshalBinary()
	add.group = k.Group
	add.suite = k.Suite
	_ = k.PrivKey.UnmarshalBinary(add.sk)

	// we sign the public key
	now := time.Now().UnixNano()
	pubbytes, _ := k.PubKey.MarshalBinary()
	add.Signature = k.Sign(pubbytes)
	n256 := uint256.NewInt(uint64(now))
	add.TimestampSigned = k.Sign(n256.Bytes())
	add.Timestamp = now
	return add
}

func (addr *Address) Serialize() []byte {

	toSerialize := addr.Raw.a.String() + "-" + hex.EncodeToString(addr.pk) + "-" + hex.EncodeToString(addr.
		Signature) + "-" + hex.
		EncodeToString(addr.TimestampSigned)
	return []byte(toSerialize)
}

func (addr *Address) Master() []byte {
	privateKey := addr.sk
	keyed, _ := blake3.NewKeyed(privateKey)
	keyed.Write(addr.Serialize())
	return keyed.Sum(nil)

}

func (addr *Address) Derive() string {
	buf := make([]byte, 32)
	frand.Read(buf)
	bb := new(big.Int).SetBytes(buf).String()
	walletBytes := make([]byte, 40)
	blake3.DeriveKey("qbit-address-"+bb, addr.Serialize(), walletBytes)
	wbb := new(big.Int).SetBytes(walletBytes)
	u, _ := uint256.FromBig(wbb)
	return u.String()
}
