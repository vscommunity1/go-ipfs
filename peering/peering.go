package peering

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

<<<<<<< HEAD
// Seed the random number generator.
//
// We don't need good randomness, but we do need randomness.
const (
	// maxBackoff is the maximum time between reconnect attempts.
	maxBackoff = 10 * time.Minute
	// The backoff will be cut off when we get within 10% of the actual max.
	// If we go over the max, we'll adjust the delay down to a random value
	// between 90-100% of the max backoff.
	maxBackoffJitter = 10 // %
	connmgrTag       = "ipfs-peering"
=======
// maxBackoff is the maximum time between reconnect attempts.
const (
	maxBackoff = 10 * time.Minute
	connmgrTag = "ipfs-peering"
>>>>>>> feat: implement peering service
	// This needs to be sufficient to prevent two sides from simultaneously
	// dialing.
	initialDelay = 5 * time.Second
)

var logger = log.Logger("peering")

type state int

const (
	stateInit state = iota
	stateRunning
	stateStopped
)

// peerHandler keeps track of all state related to a specific "peering" peer.
type peerHandler struct {
	peer   peer.ID
	host   host.Host
	ctx    context.Context
	cancel context.CancelFunc

<<<<<<< HEAD
	mu             sync.Mutex
	addrs          []multiaddr.Multiaddr
	reconnectTimer *time.Timer
=======
	mu    sync.Mutex
	addrs []multiaddr.Multiaddr
	timer *time.Timer
>>>>>>> feat: implement peering service

	nextDelay time.Duration
}

<<<<<<< HEAD
// setAddrs sets the addresses for this peer.
func (ph *peerHandler) setAddrs(addrs []multiaddr.Multiaddr) {
	// Not strictly necessary, but it helps to not trust the calling code.
	addrCopy := make([]multiaddr.Multiaddr, len(addrs))
	copy(addrCopy, addrs)

	ph.mu.Lock()
	defer ph.mu.Unlock()
	ph.addrs = addrCopy
}

// getAddrs returns a shared slice of addresses for this peer. Do not modify.
func (ph *peerHandler) getAddrs() []multiaddr.Multiaddr {
	ph.mu.Lock()
	defer ph.mu.Unlock()
	return ph.addrs
}

// stop permanently stops the peer handler.
func (ph *peerHandler) stop() {
	ph.cancel()

	ph.mu.Lock()
	defer ph.mu.Unlock()
	if ph.reconnectTimer != nil {
		ph.reconnectTimer.Stop()
		ph.reconnectTimer = nil
=======
func (ph *peerHandler) stop() {
	ph.mu.Lock()
	defer ph.mu.Unlock()

	if ph.timer != nil {
		ph.timer.Stop()
		ph.timer = nil
>>>>>>> feat: implement peering service
	}
}

func (ph *peerHandler) nextBackoff() time.Duration {
<<<<<<< HEAD
	if ph.nextDelay < maxBackoff {
		ph.nextDelay += ph.nextDelay/2 + time.Duration(rand.Int63n(int64(ph.nextDelay)))
	}

	// If we've gone over the max backoff, reduce it under the max.
	if ph.nextDelay > maxBackoff {
		ph.nextDelay = maxBackoff
		// randomize the backoff a bit (10%).
		ph.nextDelay -= time.Duration(rand.Int63n(int64(maxBackoff) * maxBackoffJitter / 100))
	}

=======
	// calculate the timeout
	if ph.nextDelay < maxBackoff {
		ph.nextDelay += ph.nextDelay/2 + time.Duration(rand.Int63n(int64(ph.nextDelay)))
	}
>>>>>>> feat: implement peering service
	return ph.nextDelay
}

func (ph *peerHandler) reconnect() {
	// Try connecting
<<<<<<< HEAD
	addrs := ph.getAddrs()
=======

	ph.mu.Lock()
	addrs := append(([]multiaddr.Multiaddr)(nil), ph.addrs...)
	ph.mu.Unlock()

>>>>>>> feat: implement peering service
	logger.Debugw("reconnecting", "peer", ph.peer, "addrs", addrs)

	err := ph.host.Connect(ph.ctx, peer.AddrInfo{ID: ph.peer, Addrs: addrs})
	if err != nil {
		logger.Debugw("failed to reconnect", "peer", ph.peer, "error", err)
		// Ok, we failed. Extend the timeout.
		ph.mu.Lock()
<<<<<<< HEAD
		if ph.reconnectTimer != nil {
			// Only counts if the reconnectTimer still exists. If not, a
			// connection _was_ somehow established.
			ph.reconnectTimer.Reset(ph.nextBackoff())
=======
		if ph.timer != nil {
			// Only counts if the timer still exists. If not, a
			// connection _was_ somehow established.
			ph.timer.Reset(ph.nextBackoff())
>>>>>>> feat: implement peering service
		}
		// Otherwise, someone else has stopped us so we can assume that
		// we're either connected or someone else will start us.
		ph.mu.Unlock()
	}

	// Always call this. We could have connected since we processed the
	// error.
	ph.stopIfConnected()
}

func (ph *peerHandler) stopIfConnected() {
	ph.mu.Lock()
	defer ph.mu.Unlock()

<<<<<<< HEAD
	if ph.reconnectTimer != nil && ph.host.Network().Connectedness(ph.peer) == network.Connected {
		logger.Debugw("successfully reconnected", "peer", ph.peer)
		ph.reconnectTimer.Stop()
		ph.reconnectTimer = nil
=======
	if ph.timer != nil && ph.host.Network().Connectedness(ph.peer) == network.Connected {
		logger.Debugw("successfully reconnected", "peer", ph.peer)
		ph.timer.Stop()
		ph.timer = nil
>>>>>>> feat: implement peering service
		ph.nextDelay = initialDelay
	}
}

// startIfDisconnected is the inverse of stopIfConnected.
func (ph *peerHandler) startIfDisconnected() {
	ph.mu.Lock()
	defer ph.mu.Unlock()

<<<<<<< HEAD
	if ph.reconnectTimer == nil && ph.host.Network().Connectedness(ph.peer) != network.Connected {
		logger.Debugw("disconnected from peer", "peer", ph.peer)
		// Always start with a short timeout so we can stagger things a bit.
		ph.reconnectTimer = time.AfterFunc(ph.nextBackoff(), ph.reconnect)
=======
	if ph.timer == nil && ph.host.Network().Connectedness(ph.peer) != network.Connected {
		logger.Debugw("disconnected from peer", "peer", ph.peer)
		// Always start with a short timeout so we can stagger things a bit.
		ph.timer = time.AfterFunc(ph.nextBackoff(), ph.reconnect)
>>>>>>> feat: implement peering service
	}
}

// PeeringService maintains connections to specified peers, reconnecting on
// disconnect with a back-off.
type PeeringService struct {
	host host.Host

	mu    sync.RWMutex
	peers map[peer.ID]*peerHandler
<<<<<<< HEAD
	state state
=======

	ctx    context.Context
	cancel context.CancelFunc
	state  state
>>>>>>> feat: implement peering service
}

// NewPeeringService constructs a new peering service. Peers can be added and
// removed immediately, but connections won't be formed until `Start` is called.
func NewPeeringService(host host.Host) *PeeringService {
<<<<<<< HEAD
	return &PeeringService{host: host, peers: make(map[peer.ID]*peerHandler)}
=======
	ps := &PeeringService{host: host, peers: make(map[peer.ID]*peerHandler)}
	ps.ctx, ps.cancel = context.WithCancel(context.Background())
	return ps
>>>>>>> feat: implement peering service
}

// Start starts the peering service, connecting and maintaining connections to
// all registered peers. It returns an error if the service has already been
// stopped.
func (ps *PeeringService) Start() error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	switch ps.state {
	case stateInit:
		logger.Infow("starting")
	case stateRunning:
		return nil
	case stateStopped:
		return errors.New("already stopped")
	}
	ps.host.Network().Notify((*netNotifee)(ps))
	ps.state = stateRunning
	for _, handler := range ps.peers {
		go handler.startIfDisconnected()
	}
	return nil
}

// Stop stops the peering service.
func (ps *PeeringService) Stop() error {
<<<<<<< HEAD
=======
	ps.cancel()
>>>>>>> feat: implement peering service
	ps.host.Network().StopNotify((*netNotifee)(ps))

	ps.mu.Lock()
	defer ps.mu.Unlock()

<<<<<<< HEAD
	switch ps.state {
	case stateInit, stateRunning:
=======
	if ps.state == stateRunning {
>>>>>>> feat: implement peering service
		logger.Infow("stopping")
		for _, handler := range ps.peers {
			handler.stop()
		}
<<<<<<< HEAD
		ps.state = stateStopped
=======
>>>>>>> feat: implement peering service
	}
	return nil
}

// AddPeer adds a peer to the peering service. This function may be safely
// called at any time: before the service is started, while running, or after it
// stops.
//
// Add peer may also be called multiple times for the same peer. The new
// addresses will replace the old.
func (ps *PeeringService) AddPeer(info peer.AddrInfo) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if handler, ok := ps.peers[info.ID]; ok {
		logger.Infow("updating addresses", "peer", info.ID, "addrs", info.Addrs)
<<<<<<< HEAD
		handler.setAddrs(info.Addrs)
=======
		handler.addrs = info.Addrs
>>>>>>> feat: implement peering service
	} else {
		logger.Infow("peer added", "peer", info.ID, "addrs", info.Addrs)
		ps.host.ConnManager().Protect(info.ID, connmgrTag)

		handler = &peerHandler{
			host:      ps.host,
			peer:      info.ID,
			addrs:     info.Addrs,
			nextDelay: initialDelay,
		}
<<<<<<< HEAD
		handler.ctx, handler.cancel = context.WithCancel(context.Background())
		ps.peers[info.ID] = handler
		switch ps.state {
		case stateRunning:
			go handler.startIfDisconnected()
		case stateStopped:
			// We still construct everything in this state because
			// it's easier to reason about. But we should still free
			// resources.
			handler.cancel()
=======
		handler.ctx, handler.cancel = context.WithCancel(ps.ctx)
		ps.peers[info.ID] = handler
		if ps.state == stateRunning {
			go handler.startIfDisconnected()
>>>>>>> feat: implement peering service
		}
	}
}

// RemovePeer removes a peer from the peering service. This function may be
// safely called at any time: before the service is started, while running, or
// after it stops.
func (ps *PeeringService) RemovePeer(id peer.ID) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if handler, ok := ps.peers[id]; ok {
		logger.Infow("peer removed", "peer", id)
		ps.host.ConnManager().Unprotect(id, connmgrTag)

		handler.stop()
<<<<<<< HEAD
=======
		handler.cancel()
>>>>>>> feat: implement peering service
		delete(ps.peers, id)
	}
}

type netNotifee PeeringService

func (nn *netNotifee) Connected(_ network.Network, c network.Conn) {
	ps := (*PeeringService)(nn)

	p := c.RemotePeer()
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if handler, ok := ps.peers[p]; ok {
		// use a goroutine to avoid blocking events.
		go handler.stopIfConnected()
	}
}
func (nn *netNotifee) Disconnected(_ network.Network, c network.Conn) {
	ps := (*PeeringService)(nn)

	p := c.RemotePeer()
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if handler, ok := ps.peers[p]; ok {
		// use a goroutine to avoid blocking events.
		go handler.startIfDisconnected()
	}
}
func (nn *netNotifee) OpenedStream(network.Network, network.Stream)     {}
func (nn *netNotifee) ClosedStream(network.Network, network.Stream)     {}
func (nn *netNotifee) Listen(network.Network, multiaddr.Multiaddr)      {}
func (nn *netNotifee) ListenClose(network.Network, multiaddr.Multiaddr) {}
