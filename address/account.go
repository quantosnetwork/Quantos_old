package address

import (
	"github.com/google/uuid"
	"go.uber.org/atomic"
)

type Account struct {
	ID uuid.UUID
	address *KeyManager
	Address string
	loadedMaster string
	Lock atomic.Bool
	Wallet interface{}
}

type account interface {
	New() *Account
	GetWalletAddress() *string
	Unlock() bool
	Lock() bool
	GetAllTransactions() map[string]interface{}
	GetOneTransaction(txid string) interface{}

}
