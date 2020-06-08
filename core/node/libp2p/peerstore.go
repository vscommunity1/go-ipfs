package libp2p

import (
	"context"

<<<<<<< HEAD
	"github.com/libp2p/go-libp2p-peerstore"
=======
	"github.com/libp2p/go-libp2p-core/peerstore"
>>>>>>> fix: close peerstore on stop
>>>>>>> 795845ea3e69d475f7eeab37fa155ed9964486ee
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
