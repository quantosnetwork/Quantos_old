package address

type Purpose uint32

const (
	PurposeBIP44 Purpose = 0x8000002C // 44' BIP44
	PurposeBIP49 Purpose = 0x80000031 // 49' BIP49
	PurposeBIP84 Purpose = 0x80000054 // 84' BIP84
)

type CoinType = uint32

const (
	// QTO = QuantusOS Coin
	QTO = 0x8000038a //906
)

const (
	Apostrophe uint32 = 0x80000000 // 0'
)

type AddressContext struct {
	prefix          string
	purpose         int
	cointype        int
	account         int
	change          int
	index           int
	derivationPath  string
	derivationBytes []byte
	createdAtBlock  string
	merkleRoot      string
	isMaster        bool
	loadedMaster    []byte
	validatorBytes  []byte
	checksum        [4]byte
	signedTimestamp string
	signature       string
}
