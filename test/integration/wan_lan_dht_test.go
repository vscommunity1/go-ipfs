package integrationtest

import (
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"testing"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/bootstrap"
	mock "github.com/ipfs/go-ipfs/core/mock"

	corenet "github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	kbucket "github.com/libp2p/go-libp2p-kbucket"
	testutil "github.com/libp2p/go-libp2p-testing/net"
	mocknet "github.com/libp2p/go-libp2p/p2p/net/mock"

	ma "github.com/multiformats/go-multiaddr"
)

func TestDHTConnectivityFast(t *testing.T) {
	conf := testutil.LatencyConfig{
		NetworkLatency:    0,
		RoutingLatency:    0,
		BlockstoreLatency: 0,
	}
	if err := RunDHTConnectivity(conf, 5); err != nil {
		t.Fatal(err)
	}
}

func TestDHTConnectivitySlowNetwork(t *testing.T) {
	SkipUnlessEpic(t)
	conf := testutil.LatencyConfig{NetworkLatency: 400 * time.Millisecond}
	if err := RunDHTConnectivity(conf, 5); err != nil {
		t.Fatal(err)
	}
}

func TestDHTConnectivitySlowRouting(t *testing.T) {
	SkipUnlessEpic(t)
	conf := testutil.LatencyConfig{RoutingLatency: 400 * time.Millisecond}
	if err := RunDHTConnectivity(conf, 5); err != nil {
		t.Fatal(err)
	}
}

var wanPrefix = net.ParseIP("100::")
var lanPrefix = net.ParseIP("fe80::")

func makeAddr(n uint32, wan bool) ma.Multiaddr {
	var ip net.IP
	if wan {
		ip = append(net.IP{}, wanPrefix...)
	} else {
		ip = append(net.IP{}, lanPrefix...)
	}

	binary.LittleEndian.PutUint32(ip[4:], n)
	addr, _ := ma.NewMultiaddr(fmt.Sprintf("/ip6/%s/tcp/4242", ip))
	return addr
}

func RunDHTConnectivity(conf testutil.LatencyConfig, numPeers int) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create network
	mn := mocknet.New(ctx)
	mn.SetLinkDefaults(mocknet.LinkOptions{
		Latency:   conf.NetworkLatency,
		Bandwidth: math.MaxInt32,
	})

	testPeer, err := core.NewNode(ctx, &core.BuildCfg{
		Online: true,
		Host:   mock.MockHostOption(mn),
	})
	if err != nil {
		return err
	}
	defer testPeer.Close()

	wanPeers := []*core.IpfsNode{}
	lanPeers := []*core.IpfsNode{}

	connectionContext, connCtxCancel := context.WithTimeout(ctx, 15*time.Second)
	defer connCtxCancel()
	for i := 0; i < numPeers; i++ {
		wanPeer, err := core.NewNode(ctx, &core.BuildCfg{
			Online: true,
			Host:   mock.MockHostOption(mn),
		})
		if err != nil {
			return err
		}
		defer wanPeer.Close()
		wanAddr := makeAddr(uint32(i), true)
		wanPeer.Peerstore.AddAddr(wanPeer.Identity, wanAddr, peerstore.PermanentAddrTTL)
		for _, p := range wanPeers {
			mn.LinkPeers(p.Identity, wanPeer.Identity)
			wanPeer.PeerHost.Connect(connectionContext, p.Peerstore.PeerInfo(p.Identity))
		}
		wanPeers = append(wanPeers, wanPeer)

		lanPeer, err := core.NewNode(ctx, &core.BuildCfg{
			Online: true,
			Host:   mock.MockHostOption(mn),
		})
		if err != nil {
			return err
		}
		defer lanPeer.Close()
		lanAddr := makeAddr(uint32(i), false)
		lanPeer.Peerstore.AddAddr(lanPeer.Identity, lanAddr, peerstore.PermanentAddrTTL)
		for _, p := range lanPeers {
			mn.LinkPeers(p.Identity, lanPeer.Identity)
			lanPeer.PeerHost.Connect(connectionContext, p.Peerstore.PeerInfo(p.Identity))
		}
		lanPeers = append(lanPeers, lanPeer)
	}
	connCtxCancel()

	// Add interfaces / addresses to test peer.
	wanAddr := makeAddr(0, true)
	testPeer.Peerstore.AddAddr(testPeer.Identity, wanAddr, peerstore.PermanentAddrTTL)
	lanAddr := makeAddr(0, false)
	testPeer.Peerstore.AddAddr(testPeer.Identity, lanAddr, peerstore.PermanentAddrTTL)

	// The test peer is connected to one lan peer.
	for _, p := range lanPeers {
		if _, err := mn.LinkPeers(testPeer.Identity, p.Identity); err != nil {
			return err
		}
	}
	err = testPeer.PeerHost.Connect(ctx, lanPeers[0].Peerstore.PeerInfo(lanPeers[0].Identity))
	if err != nil {
		return err
	}

	startupCtx, startupCancel := context.WithTimeout(ctx, time.Second*15)
	testPeer.DHT.Bootstrap(startupCtx)
StartupWait:
	for {
		select {
		case err, done := <-testPeer.DHT.LAN.RefreshRoutingTable():
			if err.Error() == kbucket.ErrLookupFailure.Error() ||
				testPeer.DHT.LAN.RoutingTable() == nil ||
				testPeer.DHT.LAN.RoutingTable().Size() == 0 {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if err != nil || !done {
				if !done {
					err = fmt.Errorf("expected refresh routing table to close")
				}
				fmt.Fprintf(os.Stderr, "how odd. that was lookupfailure.\n")
				startupCancel()
				return err
			}
			break StartupWait
		case <-startupCtx.Done():
			startupCancel()
			return fmt.Errorf("expected faster dht bootstrap")
		}
	}
	startupCancel()

	fmt.Fprintf(os.Stderr, "finding provider\n")
	// choose a lan peer and validate lan DHT is functioning.
	i := rand.Intn(len(lanPeers))
	if testPeer.PeerHost.Network().Connectedness(lanPeers[i].Identity) == corenet.Connected {
		i = (i + 1) % len(lanPeers)
		if testPeer.PeerHost.Network().Connectedness(lanPeers[i].Identity) == corenet.Connected {
			testPeer.PeerHost.Network().ClosePeer(lanPeers[i].Identity)
			testPeer.PeerHost.Peerstore().ClearAddrs(lanPeers[i].Identity)
		}
	}
	// That peer will provide a new CID, and we'll validate the test node can find it.
	provideCid := cid.NewCidV1(cid.Raw, []byte("Lan Provide Record"))
	provideCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if err := lanPeers[i].DHT.Provide(provideCtx, provideCid, true); err != nil {
		return err
	}
	provChan := testPeer.DHT.FindProvidersAsync(provideCtx, provideCid, 0)
	prov, ok := <-provChan
	if !ok || prov.ID == "" {
		return fmt.Errorf("Expected provider. stream closed early")
	}
	if prov.ID != lanPeers[i].Identity {
		return fmt.Errorf("Unexpected lan peer provided record")
	}

	// Now, bootstrap from a wan peer.
	for _, p := range wanPeers {
		if _, err := mn.LinkPeers(testPeer.Identity, p.Identity); err != nil {
			return err
		}
	}
	bis := wanPeers[0].Peerstore.PeerInfo(wanPeers[0].PeerHost.ID())
	bcfg := bootstrap.BootstrapConfigWithPeers([]peer.AddrInfo{bis})
	if err := testPeer.Bootstrap(bcfg); err != nil {
		return err
	}

	err, done := <-testPeer.DHT.WAN.RefreshRoutingTable()
	if err != nil || !done {
		if !done {
			err = fmt.Errorf("expected refresh routing table to close")
		}
		return err
	}

	// choose a wan peer and validate wan DHT is functioning.
	i = rand.Intn(len(wanPeers))
	if testPeer.PeerHost.Network().Connectedness(wanPeers[i].Identity) == corenet.Connected {
		testPeer.PeerHost.Network().ClosePeer(wanPeers[i].Identity)
		testPeer.PeerHost.Peerstore().ClearAddrs(wanPeers[i].Identity)
	}
	// That peer will provide a new CID, and we'll validate the test node can find it.
	wanCid := cid.NewCidV1(cid.Raw, []byte("Wan Provide Record"))
	wanProvideCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if err := wanPeers[i].DHT.Provide(wanProvideCtx, wanCid, true); err != nil {
		return err
	}
	provChan = testPeer.DHT.FindProvidersAsync(wanProvideCtx, wanCid, 0)
	prov, ok = <-provChan
	if !ok || prov.ID == "" {
		return fmt.Errorf("Expected one provider, closed early")
	}
	if prov.ID != wanPeers[i].Identity {
		return fmt.Errorf("Unexpected lan peer provided record")
	}

	// Finally, re-share the lan provided cid from a wan peer and expect a merged result.
	i = rand.Intn(len(wanPeers))
	if testPeer.PeerHost.Network().Connectedness(wanPeers[i].Identity) == corenet.Connected {
		testPeer.PeerHost.Network().ClosePeer(wanPeers[i].Identity)
		testPeer.PeerHost.Peerstore().ClearAddrs(wanPeers[i].Identity)
	}

	provideCtx, cancel = context.WithTimeout(ctx, time.Second)
	defer cancel()
	if err := wanPeers[i].DHT.Provide(provideCtx, provideCid, true); err != nil {
		return err
	}
	provChan = testPeer.DHT.FindProvidersAsync(provideCtx, provideCid, 0)
	prov, ok = <-provChan
	if !ok {
		return fmt.Errorf("Expected two providers, got 0")
	}
	prov, ok = <-provChan
	if !ok {
		return fmt.Errorf("Expected two providers, got 1")
	}

	return nil
}
