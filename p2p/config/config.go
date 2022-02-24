package config

import (
	_default "Quantos/p2p/config/default"
	"Quantos/p2p/go-libp2p-core/crypto"
	"Quantos/p2p/go-libp2p-core/peer"
	"encoding/base64"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/mitchellh/go-homedir"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Identity  Identity
	Datastore DataStore
	Addresses Addresses
	Mounts    Mounts
	Discovery Discovery
	Bootstrap []string
	Routing   Routing
	QNS Qns
	Peering Peering
	DNS DNS
	Plugins Plugins
	API API
	AutoNAT AutoNATConfig
	Gateway Gateway
	Chain Blockchain
	Consensus Consensus
	PubSub  PubsubConfig
	Pinning Pinning
	Swarm SwarmConfig
	SecretsVault interface{}
}


const (
	DefaultPathName = ".quantosnetwork"
	DefaultPathRoot = "~/"+DefaultPathName
	DefaultConfigFile = "config"
	EnvDir = "QUANTOS_PATH"
)

func PathRoot() (string, error) {
	dir := os.Getenv(EnvDir)
	var err error
	if len(dir) == 0 {
		dir, err = homedir.Expand(DefaultPathRoot)
	}
	return dir, err
}

func Path(configroot, extension string) (string, error) {
	if len(configroot) == 0 {
		dir, err := PathRoot()
		if err != nil {
			return "", err
		}
		return filepath.Join(dir, extension), nil

	}
	return filepath.Join(configroot, extension), nil
}
func (c *Config) BootstrapPeers() ([]peer.AddrInfo, error) {
	return _default.ParseBootstrapPeers(c.Bootstrap)
}

func Init(out io.Writer) (*Config, error) {
	identity, err := CreateIdentity(out)
	if err != nil {
		return nil, err
	}

	return InitWithIdentity(identity)
}

func InitWithIdentity(identity Identity) (*Config, error) {
	bootstrapPeers, err := _default.DefaultBootstrapPeers()
	if err != nil {
		return nil, err
	}
	datastore := DefaultDatastoreConfig()

	conf := &Config{
		API: API{
			HTTPHeaders: map[string][]string{},
		},
		Addresses: addressesConfig(),
		Datastore: datastore,
		Bootstrap: _default.BootstrapPeerStrings(bootstrapPeers),
		Identity: identity,
		Discovery: Discovery{
			MDNS: MDNS{
				Enabled: true,
				Interval: 10,
			},
		},
		Routing: Routing{
			Type: "dht",
		},
		Mounts: Mounts{
			QFS: "/qfs",
			QNS: "/qns",
			QCNT: "/qcnt",
			QWAL: "/qwallets",
			QLive: "/quantos",
			QTest: "/quantostest",
			QLocal: "/qlocal",
		},
		QNS: Qns{
			ResolveCacheSize: 128,
		},
		Gateway: Gateway{
			RootRedirect: "",
			Writable:     false,
			NoFetch:      false,
			PathPrefixes: []string{},
			HTTPHeaders: map[string][]string{
				"Access-Control-Allow-Origin":  {"*"},
				"Access-Control-Allow-Methods": {"GET"},
				"Access-Control-Allow-Headers": {"X-Requested-With", "Range", "User-Agent"},
			},
			APICommands: []string{},
		},
		Swarm: SwarmConfig{
			ConnMgr: ConnMgr{
				LowWater:    DefaultConnMgrLowWater,
				HighWater:   DefaultConnMgrHighWater,
				GracePeriod: DefaultConnMgrGracePeriod.String(),
				Type:        "basic",
			},
		},
		Pinning: Pinning{
			RemoteServices: map[string]RemotePinningService{},
		},
		DNS: DNS{
			Resolvers: map[string]string{},
		},
	}
	return conf, nil

}

// DefaultConnMgrHighWater is the default value for the connection managers
// 'high water' mark
const DefaultConnMgrHighWater = 900

// DefaultConnMgrLowWater is the default value for the connection managers 'low
// water' mark
const DefaultConnMgrLowWater = 600

// DefaultConnMgrGracePeriod is the default value for the connection managers
// grace period
const DefaultConnMgrGracePeriod = time.Second * 20


func addressesConfig() Addresses {
	return Addresses{
		Swarm: []string{
			"/ip4/0.0.0.0/tcp/4001",
			"/ip6/::/tcp/4001",
			"/ip4/0.0.0.0/udp/4001/quic",
			"/ip6/::/udp/4001/quic",
		},
		Announce:       []string{},
		AppendAnnounce: []string{},
		NoAnnounce:     []string{},
		API:            Strings{"/ip4/127.0.0.1/tcp/5001"},
		Gateway:        Strings{"/ip4/127.0.0.1/tcp/8080"},
	}
}

// DefaultDatastoreConfig is an internal function exported to aid in testing.
func DefaultDatastoreConfig() DataStore {
	return DataStore{
		StorageMax:         "50GB",
		StorageGCWatermark: 90, // 90%
		GCPeriod:           "1h",
		BloomFilterSize:    0,
		Spec:               flatfsSpec(),
	}
}

func badgerSpec() map[string]interface{} {
return map[string]interface{}{
"type":   "measure",
"prefix": "badger.datastore",
"child": map[string]interface{}{
"type":       "badgerds",
"path":       "badgerds",
"syncWrites": false,
"truncate":   true,
},
}
}

func flatfsSpec() map[string]interface{} {
	return map[string]interface{}{
		"type": "mount",
		"mounts": []interface{}{
			map[string]interface{}{
				"mountpoint": "/blocks",
				"type":       "measure",
				"prefix":     "flatfs.datastore",
				"child": map[string]interface{}{
					"type":      "flatfs",
					"path":      "blocks",
					"sync":      true,
					"shardFunc": "/repo/flatfs/shard/v1/next-to-last/2",
				},
			},
			map[string]interface{}{
				"mountpoint": "/",
				"type":       "measure",
				"prefix":     "leveldb.datastore",
				"child": map[string]interface{}{
					"type":        "levelds",
					"path":        "datastore",
					"compression": "none",
				},
			},
		},
	}
}

func CreateIdentity(out io.Writer) (Identity, error) {
	ident := Identity{}
	var sk crypto.PrivKey
	var pk crypto.PubKey

	fmt.Fprintf(out, "generating post-quantum quantos-kyber-crystal-schnorr keypair...\n")
	priv, pub, err := crypto.GenerateKyberKey()
	if err != nil {
		return ident, err
	}

	sk = priv
	pk = pub

	fmt.Fprintf(out, "done\n")
	skbytes, err := crypto.MarshalPrivateKey(sk)
	if err != nil {
		return ident, err
	}
	spew.Dump(sk)
	ident.PrivKey = base64.StdEncoding.EncodeToString(skbytes)
	id, err := peer.IDFromPublicKey(pk)
	if err != nil {
		return ident, err
	}
	spew.Dump(id)
	ident.PeerID = id.Pretty()
	fmt.Fprintf(out, "peer identity: %s\n", ident.PeerID)
	return ident, nil

}