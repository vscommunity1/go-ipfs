// +build linux

package sockets

import (
	"net"
	"sync"

	activation "github.com/coreos/go-systemd/activation"
	logging "github.com/ipfs/go-log"
	manet "github.com/multiformats/go-multiaddr-net"
)

var log = logging.Logger("socket-activation")

var socketsMu sync.Mutex
var sockets map[string][]manet.Listener

func initSockets() {
	if sockets != nil {
		return
	}
	nlisteners, err := activation.ListenersWithNames()
	// Do this before checking the error. We need this to be non-nil so we
	// don't try again.
	sockets = make(map[string][]manet.Listener, len(nlisteners))
	if err != nil {
		log.Errorf("error parsing systemd sockets: %s", err)
		return
	}
	for name, nls := range nlisteners {
		mls := make([]manet.Listener, 0, len(nls))
		for _, nl := range nls {
			ml, err := manet.WrapNetListener(nl)
			if err != nil {
				log.Errorf("error converting a systemd-socket to a multiaddr listener: %s", err)
				nl.Close()
				continue
			}
			mls = append(mls, ml)
		}
		sockets[name] = mls
	}
}

func mapListeners(nls []net.Listener) ([]manet.Listener, error) {
	mls := make([]manet.Listener, len(nls))
	return mls, nil
}

// TakeSockets takes the sockets associated with the given name.
func TakeSockets(name string) ([]manet.Listener, error) {
	socketsMu.Lock()
	defer socketsMu.Unlock()
	initSockets()

	s := sockets[name]
	delete(sockets, name)

	return s, nil
}
