package p2p

import (
	"github.com/quantosnetwork/Quantosp2p/go-libp2p-core/crypto"
	p2phost "github.com/quantosnetwork/Quantosp2p/go-libp2p-core/host"
	"github.com/quantosnetwork/Quantosp2p/go-libp2p-core/peer"
	pstore "github.com/quantosnetwork/Quantosp2p/go-libp2p-core/peerstore"
	"sync"
)

type P2P struct {
	ListenersLocal *Listeners
	ListenersP2P   *Listeners
	Streams        *StreamRegistry
	identity       peer.ID
	peerHost       p2phost.Host
	peerstore      pstore.Peerstore
}

func (p2p *P2P) CheckIfProtoExists(proto string) bool {
	protos := p2p.peerHost.Mux().Protocols()

	for _, p := range protos {
		if p != proto {
			continue
		}
		return true
	}
	return false
}
func New(peerHost p2phost.Host, peerstore pstore.Peerstore) *P2P {

	p := &P2P{

		peerHost:  peerHost,
		peerstore: peerstore,

		ListenersLocal: newListenersLocal(),
		ListenersP2P:   newListenersP2P(peerHost),

		Streams: &StreamRegistry{
			Mutex:       sync.Mutex{},
			Streams:     map[uint64]*Stream{},
			conns:       map[peer.ID]int{},
			nextID:      0,
			ConnManager: peerHost.ConnManager(),
		},
	}
	p.GenerateNewIDFromNewKyberKeys()
	return p
}

func (p2p *P2P) GenerateNewIDFromNewKyberKeys() {

	_, pk, err := crypto.GenerateKyberKey()
	if err != nil {
		panic(err)
	}
	key, err := peer.IDFromPublicKey(pk)
	if err != nil {
		panic(err)
	}

	p2p.identity = key

}