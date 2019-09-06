package libp2p

import (
	"context"

<<<<<<< HEAD
	"github.com/libp2p/go-libp2p-peerstore"
=======
	"github.com/libp2p/go-libp2p-core/peerstore"
>>>>>>> fix: close peerstore on stop
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	"go.uber.org/fx"
)

func Peerstore(lc fx.Lifecycle) peerstore.Peerstore {
	pstore := pstoremem.NewPeerstore()
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return pstore.Close()
		},
	})

	return pstore
}
