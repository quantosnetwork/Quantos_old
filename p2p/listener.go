package p2p

import (
	p2phost "github.com/quantosnetwork/Quantos/protocol/p2p/go-libp2p-core/host"
	net "github.com/quantosnetwork/Quantos/protocol/p2p/go-libp2p-core/network"
	"github.com/quantosnetwork/Quantos/protocol/p2p/go-libp2p-core/protocol"
	"context"
	"errors"

	"fmt"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	"sync"
)

var maPrefix = "/" + ma.ProtocolWithCode(ma.P_P2P).Name + "/"

type Listener interface {
	Protocol() protocol.ID
	ListenAddress() ma.Multiaddr
	TargetAddress() ma.Multiaddr
	key() string
	close()
}

type Listeners struct {
	sync.RWMutex
	Listeners map[string]Listener
}

func newListenersLocal() *Listeners {
	return &Listeners{Listeners: map[string]Listener{}}
}

func newListenersP2P(host p2phost.Host) *Listeners {
	reg := &Listeners{
		Listeners: map[string]Listener{},
	}

	host.SetStreamHandlerMatch("/x/", func(p string) bool {
		reg.RLock()
		defer reg.RUnlock()
		_, ok := reg.Listeners[p]
		return ok
	}, func(stream net.Stream) {
		reg.RLock()
		defer reg.RUnlock()
		l := reg.Listeners[string(stream.Protocol())]
		if l != nil {
			go l.(*remoteListener).handleStream(stream)
		}
	})
	return reg
}

type remoteListener struct {
	p2p *P2P

	// Application proto identifier.
	proto protocol.ID

	// Address to proxy the incoming connections to
	addr ma.Multiaddr

	// reportRemote if set to true makes the handler send '<base58 remote peerid>\n'
	// to target before any data is forwarded
	reportRemote bool
}

// ForwardRemote creates new p2p listener
func (p2p *P2P) ForwardRemote(ctx context.Context, proto protocol.ID, addr ma.Multiaddr, reportRemote bool) (Listener, error) {
	listener := &remoteListener{
		p2p: p2p,

		proto: proto,
		addr:  addr,

		reportRemote: reportRemote,
	}

	if err := p2p.ListenersP2P.Register(listener); err != nil {
		return nil, err
	}

	return listener, nil
}

func (l *remoteListener) handleStream(remote net.Stream) {
	local, err := manet.Dial(l.addr)
	if err != nil {
		_ = remote.Reset()
		return
	}

	peer := remote.Conn().RemotePeer()

	if l.reportRemote {
		if _, err := fmt.Fprintf(local, "%s\n", peer.Pretty()); err != nil {
			_ = remote.Reset()
			return
		}
	}

	peerMa, err := ma.NewMultiaddr(maPrefix + peer.Pretty())
	if err != nil {
		_ = remote.Reset()
		return
	}

	stream := &Stream{
		Protocol: l.proto,

		OriginAddr: peerMa,
		TargetAddr: l.addr,
		peer:       peer,

		Local:  local,
		Remote: remote,

		Registry: l.p2p.Streams,
	}

	l.p2p.Streams.Register(stream)
}

func (l *remoteListener) Protocol() protocol.ID {
	return l.proto
}

func (l *remoteListener) ListenAddress() ma.Multiaddr {
	addr, err := ma.NewMultiaddr(maPrefix + l.p2p.identity.Pretty())
	if err != nil {
		panic(err)
	}
	return addr
}

func (l *remoteListener) TargetAddress() ma.Multiaddr {
	return l.addr
}

func (l *remoteListener) close() {}

func (l *remoteListener) key() string {
	return string(l.proto)
}

func (r *Listeners) Register(l Listener) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.Listeners[l.key()]; ok {
		return errors.New("listener already registered")
	}
	r.Listeners[l.key()] = l
	return nil
}

func (r *Listeners) Close(matchFunc func(listener Listener) bool) int {
	todo := make([]Listener, 0)
	r.Lock()
	for _, l := range r.Listeners {
		if !matchFunc(l) {
			continue
		}

		if _, ok := r.Listeners[l.key()]; ok {
			delete(r.Listeners, l.key())
			todo = append(todo, l)
		}
	}
	r.Unlock()

	for _, l := range todo {
		l.close()
	}

	return len(todo)
}