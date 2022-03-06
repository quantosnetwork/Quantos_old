package sdk

import (
	"github.com/quantosnetwork/Quantosaddress"
	"github.com/quantosnetwork/Quantoscrypto"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
	"net"
)

const SDKVERSION = 1

type QuantosSDK interface {
	Version() int
	KeyManager() KeyManager
	AccountManager() AccountManager
	Protocol() Protocol
	BlockchainManager()
	Wire()
	GetChainConfig() map[string]string
}

type KeyManager interface {
	GenerateKeyPair() *crypto.HardenedKeys
	SignItem(item []byte) []byte
	ValidateSignedItem(item []byte, sig []byte) (bool, error)
	GetPublicKey() kyber.Point
	GetPrivateKey() kyber.Scalar
	GetCurrentSuite() suites.Suite
	ExchangeSecrets() bool
	Encrypt()
	Decrypt()
}
type AccountManager interface {
	CreateNewAccount() *address.Account
	GetAddressFromAccount() address.Address
	GetAccountFromAddress(addr address.Address) *address.Account
	Wallet()
	Keys() KeyManager
	Authenticate(args ...[]byte) bool
}
type Protocol interface {
	P2P()
	KeyExchangeProtocol()
}

type VM interface {
}

type API interface {
	Connection() net.Conn
	Auth()
	Synchronize()
	GetEndpoints() map[string]string
	Disconnect()
	Close()
	Send(dataType string, value interface{})
	SendTo(uri string, dataType string, value interface{})
	Subscribe(topic string) chan interface{}
	Info()
	Stats()
	Metrics()
}
