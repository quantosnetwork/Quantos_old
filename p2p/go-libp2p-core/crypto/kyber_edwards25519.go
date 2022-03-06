package crypto

import (
	"bytes"
	"github.com/quantosnetwork/Quantos/protocol/p2p/go-libp2p-core/crypto/pb"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/sign/schnorr"
	"go.dedis.ch/kyber/v3/xof/blake2xb"
	"lukechampine.com/frand"
)

type HardenedKeys struct {
	group   kyber.Group
	pubKey  kyber.Point
	privKey kyber.Scalar
	suite   *edwards25519.SuiteEd25519
}

func GenerateHardenedKeys() *HardenedKeys {
	seed := make([]byte, 128)
	frand.Read(seed)
	rng := blake2xb.New(seed)

	suite := edwards25519.NewBlakeSHA256Ed25519WithRand(rng)
	h := &HardenedKeys{}
	sk := suite.Scalar().Pick(rng)   // private key
	pk := suite.Point().Mul(sk, nil) // public key
	h.privKey = sk
	h.pubKey = pk
	h.suite = suite
	return h
}

// KyberPrivateKey is a Scalar
type KyberPrivateKey struct {
	k kyber.Scalar
	s *edwards25519.SuiteEd25519
}

func (k KyberPrivateKey) Equals(key Key) bool {
	b1, _ := k.Raw()
	b2, _ := key.Raw()
	return bytes.Compare(b1, b2) == 0
}

func (k KyberPrivateKey) Raw() ([]byte, error) {
	return k.k.MarshalBinary()
}

func (k KyberPrivateKey) Type() pb.KeyType {
	return pb.KeyType_KYBER
}

func (k KyberPrivateKey) Sign(bytes []byte) ([]byte, error) {
	sign, err := schnorr.Sign(k.s, k.k, bytes)
	if err != nil {
		return nil, nil
	}
	return sign, nil
}

func (k KyberPrivateKey) GetPublic() PubKey {
	pub := k.s.Point().Mul(k.k, nil)
	return &KyberPublicKey{k: pub, s: k.s}
}

// KyberPublicKey is a Point
type KyberPublicKey struct {
	k kyber.Point
	s *edwards25519.SuiteEd25519
}

func (k KyberPublicKey) Equals(key Key) bool {
	b1, _ := k.Raw()
	b2, _ := key.Raw()
	return bytes.Compare(b1, b2) == 0
}

func (k KyberPublicKey) Raw() ([]byte, error) {
	return k.k.MarshalBinary()
}

func (k KyberPublicKey) Type() pb.KeyType {
	return pb.KeyType_KYBER
}

func (k KyberPublicKey) Verify(data []byte, sig []byte) (bool, error) {
	err := schnorr.Verify(k.s, k.k, data, sig)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GenerateKyberKey() (PrivKey, PubKey, error) {
	h := GenerateHardenedKeys()
	return &KyberPrivateKey{k: h.privKey, s: h.suite}, &KyberPublicKey{k: h.pubKey, s: h.suite}, nil
}

func UnmarshalKyberPublicKey(data []byte) (PubKey, error) {
	var KPub KyberPublicKey
	err := KPub.k.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}
	return &KyberPublicKey{
		k: KPub.k,
	}, nil
}

func UnmarshalKyberPrivateKey(data []byte) (PrivKey, error) {
	var KPriv KyberPrivateKey
	err := KPriv.k.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}
	return &KyberPrivateKey{
		k: KPriv.k,
	}, nil
}

// GenerateEd25519Key generates a new ed25519 private and public key pair.
/*func GenerateEd25519Key(src io.Reader) (PrivKey, PubKey, error) {
	pub, priv, err := ed25519.GenerateKey(src)
	if err != nil {
		return nil, nil, err
	}

	return &Ed25519PrivateKey{
			k: priv,
		},
		&Ed25519PublicKey{
			k: pub,
		},
		nil
}

// Type of the private key (Ed25519).
func (k *Ed25519PrivateKey) Type() pb.KeyType {
	return pb.KeyType_Ed25519
}

// Raw private key bytes.
func (k *Ed25519PrivateKey) Raw() ([]byte, error) {
	// The Ed25519 private key contains two 32-bytes curve points, the private
	// key and the public key.
	// It makes it more efficient to get the public key without re-computing an
	// elliptic curve multiplication.
	buf := make([]byte, len(k.k))
	copy(buf, k.k)

	return buf, nil
}

func (k *Ed25519PrivateKey) pubKeyBytes() []byte {
	return k.k[ed25519.PrivateKeySize-ed25519.PublicKeySize:]
}

// Equals compares two ed25519 private keys.
func (k *Ed25519PrivateKey) Equals(o Key) bool {
	edk, ok := o.(*Ed25519PrivateKey)
	if !ok {
		return basicEquals(k, o)
	}

	return subtle.ConstantTimeCompare(k.k, edk.k) == 1
}

// GetPublic returns an ed25519 public key from a private key.
func (k *Ed25519PrivateKey) GetPublic() PubKey {
	return &Ed25519PublicKey{k: k.pubKeyBytes()}
}

// Sign returns a signature from an input message.
func (k *Ed25519PrivateKey) Sign(msg []byte) ([]byte, error) {
	return ed25519.Sign(k.k, msg), nil
}

// Type of the public key (Ed25519).
func (k *Ed25519PublicKey) Type() pb.KeyType {
	return pb.KeyType_Ed25519
}

// Raw public key bytes.
func (k *Ed25519PublicKey) Raw() ([]byte, error) {
	return k.k, nil
}

// Equals compares two ed25519 public keys.
func (k *Ed25519PublicKey) Equals(o Key) bool {
	edk, ok := o.(*Ed25519PublicKey)
	if !ok {
		return basicEquals(k, o)
	}

	return bytes.Equal(k.k, edk.k)
}

// Verify checks a signature agains the input data.
func (k *Ed25519PublicKey) Verify(data []byte, sig []byte) (bool, error) {
	return ed25519.Verify(k.k, data, sig), nil
}

// UnmarshalEd25519PublicKey returns a public key from input bytes.
func UnmarshalEd25519PublicKey(data []byte) (PubKey, error) {
	if len(data) != 32 {
		return nil, errors.New("expect ed25519 public key data size to be 32")
	}

	return &Ed25519PublicKey{
		k: ed25519.PublicKey(data),
	}, nil
}

// UnmarshalEd25519PrivateKey returns a private key from input bytes.
func UnmarshalEd25519PrivateKey(data []byte) (PrivKey, error) {
	switch len(data) {
	case ed25519.PrivateKeySize + ed25519.PublicKeySize:
		// Remove the redundant public key. See issue #36.
		redundantPk := data[ed25519.PrivateKeySize:]
		pk := data[ed25519.PrivateKeySize-ed25519.PublicKeySize : ed25519.PrivateKeySize]
		if subtle.ConstantTimeCompare(pk, redundantPk) == 0 {
			return nil, errors.New("expected redundant ed25519 public key to be redundant")
		}

		// No point in storing the extra data.
		newKey := make([]byte, ed25519.PrivateKeySize)
		copy(newKey, data[:ed25519.PrivateKeySize])
		data = newKey
	case ed25519.PrivateKeySize:
	default:
		return nil, fmt.Errorf(
			"expected ed25519 data size to be %d or %d, got %d",
			ed25519.PrivateKeySize,
			ed25519.PrivateKeySize+ed25519.PublicKeySize,
			len(data),
		)
	}

	return &Ed25519PrivateKey{
		k: ed25519.PrivateKey(data),
	}, nil
}*/
