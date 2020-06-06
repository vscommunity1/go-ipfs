/*
Package core implements the IpfsNode object and related methods.

Packages underneath core/ provide a (relatively) stable, low-level API
to carry out most IPFS-related tasks.  For more details on the other
interfaces and how core/... fits into the bigger IPFS picture, see:

  $ godoc github.com/ipfs/go-ipfs
*/
package core

import (
	"context"
	"io"

<<<<<<< HEAD
	"github.com/ipfs/go-filestore"
	"github.com/ipfs/go-ipfs-pinner"

	bserv "github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-graphsync"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	exchange "github.com/ipfs/go-ipfs-exchange-interface"
	"github.com/ipfs/go-ipfs-provider"
	ipld "github.com/ipfs/go-ipld-format"
	logging "github.com/ipfs/go-log"
	mfs "github.com/ipfs/go-mfs"
	resolver "github.com/ipfs/go-path/resolver"
	goprocess "github.com/jbenet/goprocess"
	connmgr "github.com/libp2p/go-libp2p-core/connmgr"
	ic "github.com/libp2p/go-libp2p-core/crypto"
	p2phost "github.com/libp2p/go-libp2p-core/host"
	metrics "github.com/libp2p/go-libp2p-core/metrics"
	peer "github.com/libp2p/go-libp2p-core/peer"
	pstore "github.com/libp2p/go-libp2p-core/peerstore"
	routing "github.com/libp2p/go-libp2p-core/routing"
	ddht "github.com/libp2p/go-libp2p-kad-dht/dual"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	psrouter "github.com/libp2p/go-libp2p-pubsub-router"
	record "github.com/libp2p/go-libp2p-record"
	"github.com/libp2p/go-libp2p/p2p/discovery"
	p2pbhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/ipfs/go-ipfs/core/bootstrap"
	"github.com/ipfs/go-ipfs/core/node"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/fuse/mount"
	"github.com/ipfs/go-ipfs/namesys"
=======
	version "github.com/ipfs/go-ipfs"
	rp "github.com/ipfs/go-ipfs/exchange/reprovide"
	filestore "github.com/ipfs/go-ipfs/filestore"
	mount "github.com/ipfs/go-ipfs/fuse/mount"
	namecache "github.com/ipfs/go-ipfs/namecache"
	namesys "github.com/ipfs/go-ipfs/namesys"
>>>>>>> origin/feat/ipns-follow
	ipnsrp "github.com/ipfs/go-ipfs/namesys/republisher"
	"github.com/ipfs/go-ipfs/p2p"
	"github.com/ipfs/go-ipfs/peering"
	"github.com/ipfs/go-ipfs/repo"
)

var log = logging.Logger("core")

// IpfsNode is IPFS Core module. It represents an IPFS instance.
type IpfsNode struct {

	// Self
	Identity peer.ID // the local node's identity

	Repo repo.Repo

	// Local node
	Pinning         pin.Pinner             // the pinning manager
	Mounts          Mounts                 `optional:"true"` // current mount state, if any.
	PrivateKey      ic.PrivKey             `optional:"true"` // the local node's private Key
	PNetFingerprint libp2p.PNetFingerprint `optional:"true"` // fingerprint of private network

	// Services
	Peerstore       pstore.Peerstore          `optional:"true"` // storage for other Peer instances
	Blockstore      bstore.GCBlockstore       // the block store (lower level)
	Filestore       *filestore.Filestore      `optional:"true"` // the filestore blockstore
	BaseBlocks      node.BaseBlocks           // the raw blockstore, no filestore wrapping
	GCLocker        bstore.GCLocker           // the locker used to protect the blockstore during gc
	Blocks          bserv.BlockService        // the block service, get/add blocks.
	DAG             ipld.DAGService           // the merkle dag service, get/add objects.
	Resolver        *resolver.Resolver        // the path resolution system
	Reporter        *metrics.BandwidthCounter `optional:"true"`
	Discovery       discovery.Service         `optional:"true"`
	FilesRoot       *mfs.Root
	RecordValidator record.Validator

	// Online
<<<<<<< HEAD
	PeerHost      p2phost.Host            `optional:"true"` // the network host (server+client)
	Peering       peering.PeeringService  `optional:"true"`
	Filters       *ma.Filters             `optional:"true"`
	Bootstrapper  io.Closer               `optional:"true"` // the periodic bootstrapper
	Routing       routing.Routing         `optional:"true"` // the routing system. recommend ipfs-dht
	Exchange      exchange.Interface      // the block exchange + strategy (bitswap)
	Namesys       namesys.NameSystem      // the name system, resolves paths to hashes
	Provider      provider.System         // the value provider system
	IpnsRepub     *ipnsrp.Republisher     `optional:"true"`
	GraphExchange graphsync.GraphExchange `optional:"true"`
=======
	PeerHost     p2phost.Host        // the network host (server+client)
	Bootstrapper io.Closer           // the periodic bootstrapper
	Routing      routing.IpfsRouting // the routing system. recommend ipfs-dht
	Exchange     exchange.Interface  // the block exchange + strategy (bitswap)
	Namesys      namesys.NameSystem  // the name system, resolves paths to hashes
	Namecache    namecache.NameCache // the name system follower cache
	Reprovider   *rp.Reprovider // the value reprovider system
	IpnsRepub    *ipnsrp.Republisher
>>>>>>> origin/feat/ipns-follow

	PubSub   *pubsub.PubSub             `optional:"true"`
	PSRouter *psrouter.PubsubValueStore `optional:"true"`
	DHT      *ddht.DHT                  `optional:"true"`
	P2P      *p2p.P2P                   `optional:"true"`

	Process goprocess.Process
	ctx     context.Context

	stop func() error

	// Flags
	IsOnline bool `optional:"true"` // Online is set when networking is enabled.
	IsDaemon bool `optional:"true"` // Daemon is set when running on a long-running daemon.
}

// Mounts defines what the node's mount state is. This should
// perhaps be moved to the daemon or mount. It's here because
// it needs to be accessible across daemon requests.
type Mounts struct {
	Ipfs mount.Mount
	Ipns mount.Mount
}

<<<<<<< HEAD
// Close calls Close() on the App object
=======
func (n *IpfsNode) startOnlineServices(ctx context.Context, routingOption RoutingOption, hostOption HostOption, do DiscoveryOption, pubsub, ipnsps, mplex, follow bool) error {
	if n.PeerHost != nil { // already online.
		return errors.New("node already online")
	}

	if n.PrivateKey == nil {
		return fmt.Errorf("private key not available")
	}

	// get undialable addrs from config
	cfg, err := n.Repo.Config()
	if err != nil {
		return err
	}

	var libp2pOpts []libp2p.Option
	for _, s := range cfg.Swarm.AddrFilters {
		f, err := mamask.NewMask(s)
		if err != nil {
			return fmt.Errorf("incorrectly formatted address filter in config: %s", s)
		}
		libp2pOpts = append(libp2pOpts, libp2p.FilterAddresses(f))
	}

	if !cfg.Swarm.DisableBandwidthMetrics {
		// Set reporter
		n.Reporter = metrics.NewBandwidthCounter()
		libp2pOpts = append(libp2pOpts, libp2p.BandwidthReporter(n.Reporter))
	}

	swarmkey, err := n.Repo.SwarmKey()
	if err != nil {
		return err
	}

	if swarmkey != nil {
		protec, err := pnet.NewProtector(bytes.NewReader(swarmkey))
		if err != nil {
			return fmt.Errorf("failed to configure private network: %s", err)
		}
		n.PNetFingerprint = protec.Fingerprint()
		go func() {
			t := time.NewTicker(30 * time.Second)
			<-t.C // swallow one tick
			for {
				select {
				case <-t.C:
					if ph := n.PeerHost; ph != nil {
						if len(ph.Network().Peers()) == 0 {
							log.Warning("We are in private network and have no peers.")
							log.Warning("This might be configuration mistake.")
						}
					}
				case <-n.Process().Closing():
					t.Stop()
					return
				}
			}
		}()

		libp2pOpts = append(libp2pOpts, libp2p.PrivateNetwork(protec))
	}

	addrsFactory, err := makeAddrsFactory(cfg.Addresses)
	if err != nil {
		return err
	}
	if !cfg.Swarm.DisableRelay {
		addrsFactory = composeAddrsFactory(addrsFactory, filterRelayAddrs)
	}
	libp2pOpts = append(libp2pOpts, libp2p.AddrsFactory(addrsFactory))

	connm, err := constructConnMgr(cfg.Swarm.ConnMgr)
	if err != nil {
		return err
	}
	libp2pOpts = append(libp2pOpts, libp2p.ConnectionManager(connm))

	libp2pOpts = append(libp2pOpts, makeSmuxTransportOption(mplex))

	if !cfg.Swarm.DisableNatPortMap {
		libp2pOpts = append(libp2pOpts, libp2p.NATPortMap())
	}

	// disable the default listen addrs
	libp2pOpts = append(libp2pOpts, libp2p.NoListenAddrs)

	if cfg.Swarm.DisableRelay {
		// Enabled by default.
		libp2pOpts = append(libp2pOpts, libp2p.DisableRelay())
	} else {
		relayOpts := []circuit.RelayOpt{circuit.OptDiscovery}
		if cfg.Swarm.EnableRelayHop {
			relayOpts = append(relayOpts, circuit.OptHop)
		}
		libp2pOpts = append(libp2pOpts, libp2p.EnableRelay(relayOpts...))
	}

	// explicitly enable the default transports
	libp2pOpts = append(libp2pOpts, libp2p.DefaultTransports)

	if cfg.Experimental.QUIC {
		libp2pOpts = append(libp2pOpts, libp2p.Transport(quic.NewTransport))
	}

	// enable routing
	libp2pOpts = append(libp2pOpts, libp2p.Routing(func(h p2phost.Host) (routing.PeerRouting, error) {
		r, err := routingOption(ctx, h, n.Repo.Datastore(), n.RecordValidator)
		n.Routing = r
		return r, err
	}))

	// enable autorelay
	if cfg.Swarm.EnableAutoRelay {
		libp2pOpts = append(libp2pOpts, libp2p.EnableAutoRelay())
	}

	peerhost, err := hostOption(ctx, n.Identity, n.Peerstore, libp2pOpts...)

	if err != nil {
		return err
	}

	n.PeerHost = peerhost

	if err := n.startOnlineServicesWithHost(ctx, routingOption, pubsub, ipnsps); err != nil {
		return err
	}

	// Ok, now we're ready to listen.
	if err := startListening(n.PeerHost, cfg); err != nil {
		return err
	}

	n.P2P = p2p.NewP2P(n.Identity, n.PeerHost, n.Peerstore)

	if follow {
		n.Namecache = namecache.NewNameCache(ctx, n.Namesys, n.DAG)
		n.Namecache, err = namecache.NewPersistentCache(n.Namecache, n.Repo.Datastore())
		if err != nil {
			return err
		}
	}

	// setup local discovery
	if do != nil {
		service, err := do(ctx, n.PeerHost)
		if err != nil {
			log.Error("mdns error: ", err)
		} else {
			service.RegisterNotifee(n)
			n.Discovery = service
		}
	}

	return n.Bootstrap(DefaultBootstrapConfig)
}

func constructConnMgr(cfg config.ConnMgr) (ifconnmgr.ConnManager, error) {
	switch cfg.Type {
	case "":
		// 'default' value is the basic connection manager
		return connmgr.NewConnManager(config.DefaultConnMgrLowWater, config.DefaultConnMgrHighWater, config.DefaultConnMgrGracePeriod), nil
	case "none":
		return nil, nil
	case "basic":
		grace, err := time.ParseDuration(cfg.GracePeriod)
		if err != nil {
			return nil, fmt.Errorf("parsing Swarm.ConnMgr.GracePeriod: %s", err)
		}

		return connmgr.NewConnManager(cfg.LowWater, cfg.HighWater, grace), nil
	default:
		return nil, fmt.Errorf("unrecognized ConnMgr.Type: %q", cfg.Type)
	}
}

func (n *IpfsNode) startLateOnlineServices(ctx context.Context) error {
	cfg, err := n.Repo.Config()
	if err != nil {
		return err
	}

	var keyProvider rp.KeyChanFunc

	switch cfg.Reprovider.Strategy {
	case "all":
		fallthrough
	case "":
		keyProvider = rp.NewBlockstoreProvider(n.Blockstore)
	case "roots":
		keyProvider = rp.NewPinnedProvider(n.Pinning, n.DAG, true)
	case "pinned":
		keyProvider = rp.NewPinnedProvider(n.Pinning, n.DAG, false)
	default:
		return fmt.Errorf("unknown reprovider strategy '%s'", cfg.Reprovider.Strategy)
	}
	n.Reprovider = rp.NewReprovider(ctx, n.Routing, keyProvider)

	reproviderInterval := kReprovideFrequency
	if cfg.Reprovider.Interval != "" {
		dur, err := time.ParseDuration(cfg.Reprovider.Interval)
		if err != nil {
			return err
		}

		reproviderInterval = dur
	}

	go n.Reprovider.Run(reproviderInterval)

	return nil
}

func makeAddrsFactory(cfg config.Addresses) (p2pbhost.AddrsFactory, error) {
	var annAddrs []ma.Multiaddr
	for _, addr := range cfg.Announce {
		maddr, err := ma.NewMultiaddr(addr)
		if err != nil {
			return nil, err
		}
		annAddrs = append(annAddrs, maddr)
	}

	filters := mafilter.NewFilters()
	noAnnAddrs := map[string]bool{}
	for _, addr := range cfg.NoAnnounce {
		f, err := mamask.NewMask(addr)
		if err == nil {
			filters.AddDialFilter(f)
			continue
		}
		maddr, err := ma.NewMultiaddr(addr)
		if err != nil {
			return nil, err
		}
		noAnnAddrs[maddr.String()] = true
	}

	return func(allAddrs []ma.Multiaddr) []ma.Multiaddr {
		var addrs []ma.Multiaddr
		if len(annAddrs) > 0 {
			addrs = annAddrs
		} else {
			addrs = allAddrs
		}

		var out []ma.Multiaddr
		for _, maddr := range addrs {
			// check for exact matches
			ok, _ := noAnnAddrs[maddr.String()]
			// check for /ipcidr matches
			if !ok && !filters.AddrBlocked(maddr) {
				out = append(out, maddr)
			}
		}
		return out
	}, nil
}

func makeSmuxTransportOption(mplexExp bool) libp2p.Option {
	const yamuxID = "/yamux/1.0.0"
	const mplexID = "/mplex/6.7.0"

	ymxtpt := &yamux.Transport{
		AcceptBacklog:          512,
		ConnectionWriteTimeout: time.Second * 10,
		KeepAliveInterval:      time.Second * 30,
		EnableKeepAlive:        true,
		MaxStreamWindowSize:    uint32(1024 * 512),
		LogOutput:              ioutil.Discard,
	}

	if os.Getenv("YAMUX_DEBUG") != "" {
		ymxtpt.LogOutput = os.Stderr
	}

	muxers := map[string]smux.Transport{yamuxID: ymxtpt}
	if mplexExp {
		muxers[mplexID] = mplex.DefaultTransport
	}

	// Allow muxer preference order overriding
	order := []string{yamuxID, mplexID}
	if prefs := os.Getenv("LIBP2P_MUX_PREFS"); prefs != "" {
		order = strings.Fields(prefs)
	}

	opts := make([]libp2p.Option, 0, len(order))
	for _, id := range order {
		tpt, ok := muxers[id]
		if !ok {
			log.Warning("unknown or duplicate muxer in LIBP2P_MUX_PREFS: %s", id)
			continue
		}
		delete(muxers, id)
		opts = append(opts, libp2p.Muxer(id, tpt))
	}

	return libp2p.ChainOptions(opts...)
}

func setupDiscoveryOption(d config.Discovery) DiscoveryOption {
	if d.MDNS.Enabled {
		return func(ctx context.Context, h p2phost.Host) (discovery.Service, error) {
			if d.MDNS.Interval == 0 {
				d.MDNS.Interval = 5
			}
			return discovery.NewMdnsService(ctx, h, time.Duration(d.MDNS.Interval)*time.Second, discovery.ServiceTag)
		}
	}
	return nil
}

// HandlePeerFound attempts to connect to peer from `PeerInfo`, if it fails
// logs a warning log.
func (n *IpfsNode) HandlePeerFound(p pstore.PeerInfo) {
	log.Warning("trying peer info: ", p)
	ctx, cancel := context.WithTimeout(n.Context(), discoveryConnTimeout)
	defer cancel()
	if err := n.PeerHost.Connect(ctx, p); err != nil {
		log.Warning("Failed to connect to peer found by discovery: ", err)
	}
}

// startOnlineServicesWithHost  is the set of services which need to be
// initialized with the host and _before_ we start listening.
func (n *IpfsNode) startOnlineServicesWithHost(ctx context.Context, routingOption RoutingOption, enablePubsub bool, enableIpnsps bool) error {
	cfg, err := n.Repo.Config()
	if err != nil {
		return err
	}

	if cfg.Swarm.EnableAutoNATService {
		var opts []libp2p.Option
		if cfg.Experimental.QUIC {
			opts = append(opts, libp2p.DefaultTransports, libp2p.Transport(quic.NewTransport))
		}

		svc, err := autonat.NewAutoNATService(ctx, n.PeerHost, opts...)
		if err != nil {
			return err
		}
		n.AutoNAT = svc
	}

	if enablePubsub || enableIpnsps {
		var service *pubsub.PubSub

		var pubsubOptions []pubsub.Option
		if cfg.Pubsub.DisableSigning {
			pubsubOptions = append(pubsubOptions, pubsub.WithMessageSigning(false))
		}

		if cfg.Pubsub.StrictSignatureVerification {
			pubsubOptions = append(pubsubOptions, pubsub.WithStrictSignatureVerification(true))
		}

		switch cfg.Pubsub.Router {
		case "":
			fallthrough
		case "floodsub":
			service, err = pubsub.NewFloodSub(ctx, n.PeerHost, pubsubOptions...)

		case "gossipsub":
			service, err = pubsub.NewGossipSub(ctx, n.PeerHost, pubsubOptions...)

		default:
			err = fmt.Errorf("Unknown pubsub router %s", cfg.Pubsub.Router)
		}

		if err != nil {
			return err
		}
		n.PubSub = service
	}

	// this code is necessary just for tests: mock network constructions
	// ignore the libp2p constructor options that actually construct the routing!
	if n.Routing == nil {
		r, err := routingOption(ctx, n.PeerHost, n.Repo.Datastore(), n.RecordValidator)
		if err != nil {
			return err
		}
		n.Routing = r
		n.PeerHost = rhost.Wrap(n.PeerHost, n.Routing)
	}

	// TODO: I'm not a fan of type assertions like this but the
	// `RoutingOption` system doesn't currently provide access to the
	// IpfsNode.
	//
	// Ideally, we'd do something like:
	//
	// 1. Add some fancy method to introspect into tiered routers to extract
	//    things like the pubsub router or the DHT (complicated, messy,
	//    probably not worth it).
	// 2. Pass the IpfsNode into the RoutingOption (would also remove the
	//    PSRouter case below.
	// 3. Introduce some kind of service manager? (my personal favorite but
	//    that requires a fair amount of work).
	if dht, ok := n.Routing.(*dht.IpfsDHT); ok {
		n.DHT = dht
	}

	if enableIpnsps {
		n.PSRouter = psrouter.NewPubsubValueStore(
			ctx,
			n.PeerHost,
			n.Routing,
			n.PubSub,
			n.RecordValidator,
		)
		n.Routing = rhelpers.Tiered{
			Routers: []routing.IpfsRouting{
				// Always check pubsub first.
				&rhelpers.Compose{
					ValueStore: &rhelpers.LimitedValueStore{
						ValueStore: n.PSRouter,
						Namespaces: []string{"ipns"},
					},
				},
				n.Routing,
			},
			Validator: n.RecordValidator,
		}
	}

	// setup exchange service
	bitswapNetwork := bsnet.NewFromIpfsHost(n.PeerHost, n.Routing)
	n.Exchange = bitswap.New(ctx, bitswapNetwork, n.Blockstore)

	size, err := n.getCacheSize()
	if err != nil {
		return err
	}

	// setup name system
	n.Namesys = namesys.NewNameSystem(n.Routing, n.Repo.Datastore(), size)

	// setup ipns republishing
	return n.setupIpnsRepublisher()
}

// getCacheSize returns cache life and cache size
func (n *IpfsNode) getCacheSize() (int, error) {
	cfg, err := n.Repo.Config()
	if err != nil {
		return 0, err
	}

	cs := cfg.Ipns.ResolveCacheSize
	if cs == 0 {
		cs = DefaultIpnsCacheSize
	}
	if cs < 0 {
		return 0, fmt.Errorf("cannot specify negative resolve cache size")
	}
	return cs, nil
}

func (n *IpfsNode) setupIpnsRepublisher() error {
	cfg, err := n.Repo.Config()
	if err != nil {
		return err
	}

	n.IpnsRepub = ipnsrp.NewRepublisher(n.Namesys, n.Repo.Datastore(), n.PrivateKey, n.Repo.Keystore())

	if cfg.Ipns.RepublishPeriod != "" {
		d, err := time.ParseDuration(cfg.Ipns.RepublishPeriod)
		if err != nil {
			return fmt.Errorf("failure to parse config setting IPNS.RepublishPeriod: %s", err)
		}

		if !u.Debug && (d < time.Minute || d > (time.Hour*24)) {
			return fmt.Errorf("config setting IPNS.RepublishPeriod is not between 1min and 1day: %s", d)
		}

		n.IpnsRepub.Interval = d
	}

	if cfg.Ipns.RecordLifetime != "" {
		d, err := time.ParseDuration(cfg.Ipns.RecordLifetime)
		if err != nil {
			return fmt.Errorf("failure to parse config setting IPNS.RecordLifetime: %s", err)
		}

		n.IpnsRepub.RecordLifetime = d
	}

	n.Process().Go(n.IpnsRepub.Run)

	return nil
}

// Process returns the Process object
func (n *IpfsNode) Process() goprocess.Process {
	return n.proc
}

// Close calls Close() on the Process object
>>>>>>> origin/feat/ipns-follow
func (n *IpfsNode) Close() error {
	return n.stop()
}

// Context returns the IpfsNode context
func (n *IpfsNode) Context() context.Context {
	if n.ctx == nil {
		n.ctx = context.TODO()
	}
	return n.ctx
}

// Bootstrap will set and call the IpfsNodes bootstrap function.
func (n *IpfsNode) Bootstrap(cfg bootstrap.BootstrapConfig) error {
	// TODO what should return value be when in offlineMode?
	if n.Routing == nil {
		return nil
	}

	if n.Bootstrapper != nil {
		n.Bootstrapper.Close() // stop previous bootstrap process.
	}

	// if the caller did not specify a bootstrap peer function, get the
	// freshest bootstrap peers from config. this responds to live changes.
	if cfg.BootstrapPeers == nil {
		cfg.BootstrapPeers = func() []peer.AddrInfo {
			ps, err := n.loadBootstrapPeers()
			if err != nil {
				log.Warn("failed to parse bootstrap peers from config")
				return nil
			}
			return ps
		}
	}

	var err error
	n.Bootstrapper, err = bootstrap.Bootstrap(n.Identity, n.PeerHost, n.Routing, cfg)
	return err
}

func (n *IpfsNode) loadBootstrapPeers() ([]peer.AddrInfo, error) {
	cfg, err := n.Repo.Config()
	if err != nil {
		return nil, err
	}

	return cfg.BootstrapPeers()
}

type ConstructPeerHostOpts struct {
	AddrsFactory      p2pbhost.AddrsFactory
	DisableNatPortMap bool
	DisableRelay      bool
	EnableRelayHop    bool
	ConnectionManager connmgr.ConnManager
}
