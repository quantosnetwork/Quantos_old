package address

import (
	"bytes"
	"encoding/base32"
	"encoding/binary"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/iancoleman/strcase"
	"github.com/zeebo/blake3"
	"lukechampine.com/frand"
	"strings"
	"unsafe"
)

/*

	QBIT Addresses

	Purpose: Wallet addresses
	Size: 36 bytes
	Prefix: 0x (0x = byte representation)

	Structure of Address

	[All Qbit address starts with]
    [4]byte  uint32 [ 0x38A ]
	[4]byte crc32 IEEE of kyber public key generated at the same time as the p2p peer ID
	[20]byte of cryptographically safe random bytes data
	[8]byte as a random nonce to avoid addresses collisions

	0x is then appended to the address to mark it as bytes

	@Suggestion we could but qbit: [address here] to differentiate from other 0x address types

	906 / 38A = cointype

*/

func SliceToArray32(bytes []byte) *[32]uint8 { return (*[32]uint8)(unsafe.Pointer(&bytes[0])) }
func SliceToArray64(bytes []byte) *[64]uint8 { return (*[64]uint8)(unsafe.Pointer(&bytes[0])) }

type QBITAddress struct {
	seed            []byte
	words           *[16]uint32
	network         [2]byte
	protocolVersion [2]byte
	prefix          uint32 // uint32(906)
	checksumIEEE    uint32
	context         uint32 // blake3 context of the address
	Signature       []byte
}

var ZEROADDRESS string
var ZBYTES [32]byte

func GenerateNewQbitAddress(networkID [2]byte, version [2]byte, prefix uint32, context uint32) *QBITAddress {

	var addr QBITAddress
	add := NewQBITAddress(networkID, version, prefix, context, &addr)
	return add

}

func (q *QBITAddress) Hash() []byte {
	qs, _ := json.Marshal(q.seed)
	hash := blake3.Sum256(qs)
	return hash[:]

}

func (q *QBITAddress) String() string {
	h := q.Hash()
	str := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(h)
	str = Reverse(str)
	str = strings.ToLower(str)
	str = strcase.ToCamel(str)
	str = Add0xPrefix(str)
	return str
}

func QBITAddressFromAddressString(str string) string {
	strsplit := strings.Split(str, "0x")
	revAddrStr := strsplit[0]
	addrStr := Reverse(revAddrStr)
	return addrStr

}

func Add0xPrefix(addr string) string {
	return "0x" + addr
}

func NewQBITAddress(networkID [2]byte, version [2]byte, prefix uint32, context uint32, addr *QBITAddress) *QBITAddress {

	q := addr
	q.network = networkID
	q.protocolVersion = version
	q.prefix = prefix
	q.context = context
	seed1 := generateRandomBytes()
	seed2 := generateRandomBytes()
	seed := make([][]byte, 2)
	seed[0] = seed1[:]
	seed[1] = seed2[:]
	theSeed := bytes.Join(seed, nil)
	q.seed = theSeed
	q.words = new([16]uint32)

	return q
}

func getZeroAddress(netID [2]byte, version [2]byte, prefix uint32, context uint32) *QBITAddress {
	q := GenerateNewQbitAddress(netID, version, prefix, context)

	buf := make([]byte, 64)
	b := fillBufferWithZeros(buf)
	spew.Dump(b[:])
	q.seed = b[:]
	q.context = context

	return q

}

func ZeroAddress(netID [2]byte, version [2]byte, prefix uint32, context uint32) string {
	addr := getZeroAddress(netID, version, prefix, context)
	return addr.String()
}

func fillBufferWithZeros(buf []byte) []byte {
	for i, _ := range buf {
		buf[i] = 0
	}
	return buf
}

func generateRandomBytes() [32]byte {
	b := frand.Entropy256()
	frand.Read(b[:])
	return b

}

func BytesToWords(bytes *[64]uint8, words *[16]uint32) {
	words[0] = binary.LittleEndian.Uint32(bytes[0*4:])
	words[1] = binary.LittleEndian.Uint32(bytes[1*4:])
	words[2] = binary.LittleEndian.Uint32(bytes[2*4:])
	words[3] = binary.LittleEndian.Uint32(bytes[3*4:])
	words[4] = binary.LittleEndian.Uint32(bytes[4*4:])
	words[5] = binary.LittleEndian.Uint32(bytes[5*4:])
	words[6] = binary.LittleEndian.Uint32(bytes[6*4:])
	words[7] = binary.LittleEndian.Uint32(bytes[7*4:])
	words[8] = binary.LittleEndian.Uint32(bytes[8*4:])
	words[9] = binary.LittleEndian.Uint32(bytes[9*4:])
	words[10] = binary.LittleEndian.Uint32(bytes[10*4:])
	words[11] = binary.LittleEndian.Uint32(bytes[11*4:])
	words[12] = binary.LittleEndian.Uint32(bytes[12*4:])
	words[13] = binary.LittleEndian.Uint32(bytes[13*4:])
	words[14] = binary.LittleEndian.Uint32(bytes[14*4:])
	words[15] = binary.LittleEndian.Uint32(bytes[15*4:])

}

func WordsToBytes(words *[16]uint32, bytes []byte) {
	bytes = bytes[:64]
	binary.LittleEndian.PutUint32(bytes[0*4:1*4], words[0])
	binary.LittleEndian.PutUint32(bytes[1*4:2*4], words[1])
	binary.LittleEndian.PutUint32(bytes[2*4:3*4], words[2])
	binary.LittleEndian.PutUint32(bytes[3*4:4*4], words[3])
	binary.LittleEndian.PutUint32(bytes[4*4:5*4], words[4])
	binary.LittleEndian.PutUint32(bytes[5*4:6*4], words[5])
	binary.LittleEndian.PutUint32(bytes[6*4:7*4], words[6])
	binary.LittleEndian.PutUint32(bytes[7*4:8*4], words[7])
	binary.LittleEndian.PutUint32(bytes[8*4:9*4], words[8])
	binary.LittleEndian.PutUint32(bytes[9*4:10*4], words[9])
	binary.LittleEndian.PutUint32(bytes[10*4:11*4], words[10])
	binary.LittleEndian.PutUint32(bytes[11*4:12*4], words[11])
	binary.LittleEndian.PutUint32(bytes[12*4:13*4], words[12])
	binary.LittleEndian.PutUint32(bytes[13*4:14*4], words[13])
	binary.LittleEndian.PutUint32(bytes[14*4:15*4], words[14])
	binary.LittleEndian.PutUint32(bytes[15*4:16*4], words[15])
}

func KeyFromBytes(key []byte, out *[8]uint32) {
	key = key[:32]
	out[0] = binary.LittleEndian.Uint32(key[0:])
	out[1] = binary.LittleEndian.Uint32(key[4:])
	out[2] = binary.LittleEndian.Uint32(key[8:])
	out[3] = binary.LittleEndian.Uint32(key[12:])
	out[4] = binary.LittleEndian.Uint32(key[16:])
	out[5] = binary.LittleEndian.Uint32(key[20:])
	out[6] = binary.LittleEndian.Uint32(key[24:])
	out[7] = binary.LittleEndian.Uint32(key[28:])
}

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
