package core

import (
	"context"
	"sync"

	"github.com/ipfs/go-ipfs/core/bootstrap"
	"github.com/ipfs/go-ipfs/core/node"

	"github.com/ipfs/go-metrics-interface"
	"go.uber.org/fx"
)

type BuildCfg = node.BuildCfg // Alias for compatibility until we properly refactor the constructor interface

// NewNode constructs and returns an IpfsNode using the given cfg.
func NewNode(ctx context.Context, cfg *BuildCfg) (*IpfsNode, error) {
	ctx = metrics.CtxScope(ctx, "ipfs")

	n := &IpfsNode{
		ctx: ctx,
	}

	app := fx.New(
		node.IPFS(ctx, cfg),

		fx.NopLogger,
		fx.Extract(n),
	)

	var once sync.Once
	var stopErr error
	n.stop = func() error {
		once.Do(func() {
			stopErr = app.Stop(context.Background())
		})
		return stopErr
	}
	n.IsOnline = cfg.Online

	go func() {
		// Note that some services use contexts to signal shutting down, which is
		// very suboptimal. This needs to be here until that's addressed somehow
		<-ctx.Done()
		err := n.stop()
		if err != nil {
			log.Error("failure on stop: ", err)
		}
	}()

<<<<<<< HEAD
	if app.Err() != nil {
		return nil, app.Err()
=======
	if cfg.Online {
		do := setupDiscoveryOption(rcfg.Discovery)
		pubsub := cfg.getOpt("pubsub")
		ipnsps := cfg.getOpt("ipnsps")
		mplex := cfg.getOpt("mplex")
		follow := cfg.getOpt("follow")
		if err := n.startOnlineServices(ctx, cfg.Routing, hostOption, do, pubsub, ipnsps, mplex, follow); err != nil {
			return err
		}
	} else {
		n.Exchange = offline.Exchange(n.Blockstore)
		n.Routing = offroute.NewOfflineRouter(n.Repo.Datastore(), n.RecordValidator)
		n.Namesys = namesys.NewNameSystem(n.Routing, n.Repo.Datastore(), 0)
>>>>>>> origin/feat/ipns-follow
	}

	if err := app.Start(ctx); err != nil {
		return nil, err
	}

	// TODO: How soon will bootstrap move to libp2p?
	if !cfg.Online {
		return n, nil
	}

	return n, n.Bootstrap(bootstrap.DefaultBootstrapConfig)
}
