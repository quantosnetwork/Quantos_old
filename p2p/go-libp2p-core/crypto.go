package go_libp2p_core

import (
	pcrypto "github.com/quantosnetwork/Quantos/protocol/p2p/go-libp2p-core/crypto"
)

var _ pcrypto.PrivKey
var _ pcrypto.PubKey

type PrivKey interface {
	pcrypto.Key
	Sign([]byte) ([]byte, error)
	GetPublic() pcrypto.PubKey
}