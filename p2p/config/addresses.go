package config

type Addresses struct {
	Swarm          []string // addresses for the swarm to listen on
	Announce       []string // swarm addresses to announce to the network, if len > 0 replaces auto detected addresses
	AppendAnnounce []string // similar to Announce but doesn't overwride auto detected addresses, they are just appended
	NoAnnounce     []string // swarm addresses not to announce to the network
	API            Strings  // address for the local API (RPC)
	Gateway        Strings  // address to listen on for Quantos HTTP object Gateway
}
