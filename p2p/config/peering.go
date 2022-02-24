package config

import (
	"Quantos/p2p/go-libp2p-core/peer"
)

// Peering configures the peering service.
type Peering struct {
	// Peers lists the nodes to attempt to stay connected with.
	Peers []peer.AddrInfo
}