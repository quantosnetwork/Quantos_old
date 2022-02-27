package address

import (
	"github.com/google/uuid"
	"go.uber.org/atomic"
)

type Address interface{}

type Account struct {
	ID             uuid.UUID
	address        *Address
	Address        string
	loadedMaster   string
	Lock           atomic.Bool
	Wallet         interface{}
	CreatedAtBlock uint32
}

type account interface {
	New() *Account
	GetWalletAddress() *string
	Unlock() bool
	Lock() bool
	GetAllTransactions() map[string]interface{}
	GetOneTransaction(txid string) interface{}
}

type AccountState struct {
	Nonce               []byte
	Balance             string
	MerkleRoot          []byte
	VMCodeHashToExecute []byte
}
