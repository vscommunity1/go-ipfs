// Package namecache implements background following (resolution and pinning) of names
package namecache

import (
	"context"
<<<<<<< HEAD
	"fmt"
=======
>>>>>>> namecache: ipfs name follow
	"strings"
	"sync"
	"time"

<<<<<<< HEAD
	"github.com/ipfs/go-ipfs/core/coreapi/interface"
	"github.com/ipfs/go-ipfs/dagutils"
	"github.com/ipfs/go-ipfs/namesys"

	"gx/ipfs/QmR8BauakNcBa3RbE4nbQu76PDiJgoQgz8AJdhJuiU4TAw/go-cid"
	ipld "gx/ipfs/QmRL22E4paat7ky7vx9MLpR97JHHbFPrg3ytFQw6qp1y1s/go-ipld-format"
	"gx/ipfs/QmWqh9oob7ZHQRwU5CdTqpnC8ip8BEkFNrwXRxeNo5Y7vA/go-path"
	dag "gx/ipfs/Qmb2UEG2TAeVrEJSjqsZF7Y2he7wRDkrdt6c3bECxwZf4k/go-merkledag"
=======
	namesys "github.com/ipfs/go-ipfs/namesys"
	pin "github.com/ipfs/go-ipfs/pin"

	uio "gx/ipfs/QmUnHNqhSB1JgzVCxL1Kz3yb4bdyB4q1Z9AD5AUBVmt3fZ/go-unixfs/io"
	resolver "gx/ipfs/QmVi2uUygezqaMTqs3Yzt5FcZFHJoYD4B7jQ2BELjj7ZuY/go-path/resolver"
	ipld "gx/ipfs/QmcKKBwfz6FyQdHR2jsXrrF6XeSBXYL86anmWNewpFpoF5/go-ipld-format"
>>>>>>> namecache: ipfs name follow
	logging "gx/ipfs/QmcuXC5cxs79ro2cUuHs4HQ2bkDLJUYokwL8aivcX6HW3C/go-log"
)

const (
<<<<<<< HEAD
	DefaultFollowInterval = 1 * time.Hour
	resolveTimeout        = 1 * time.Minute
=======
	followInterval = 60 * time.Minute
>>>>>>> namecache: ipfs name follow
)

var log = logging.Logger("namecache")

// NameCache represents a following cache of names
type NameCache interface {
<<<<<<< HEAD
	// Follow starts following name
	Follow(name string, prefetch bool, followInterval time.Duration) error
	// Unofollow cancels a follow
	Unfollow(name string) error
	// ListFollows returns a list of followed names
=======
	Follow(name string, pinit bool)
	Unfollow(name string)
>>>>>>> namecache: ipfs name follow
	ListFollows() []string
}

type nameCache struct {
	nsys    namesys.NameSystem
<<<<<<< HEAD
	dag     ipld.NodeGetter
=======
	pinning pin.Pinner
	dag     ipld.DAGService
>>>>>>> namecache: ipfs name follow

	ctx     context.Context
	follows map[string]func()
	mx      sync.Mutex
}

<<<<<<< HEAD
func NewNameCache(ctx context.Context, nsys namesys.NameSystem, dag ipld.NodeGetter) NameCache {
	return &nameCache{
		ctx:     ctx,
		nsys:    nsys,
=======
func NewNameCache(ctx context.Context, nsys namesys.NameSystem, pinning pin.Pinner, dag ipld.DAGService) NameCache {
	return &nameCache{
		ctx:     ctx,
		nsys:    nsys,
		pinning: pinning,
>>>>>>> namecache: ipfs name follow
		dag:     dag,
		follows: make(map[string]func()),
	}
}

// Follow spawns a goroutine that periodically resolves a name
<<<<<<< HEAD
// and (when dopin is true) pins it in the background
func (nc *nameCache) Follow(name string, prefetch bool, followInterval time.Duration) error {
	nc.mx.Lock()
	defer nc.mx.Unlock()

	if !strings.HasPrefix(name, "/ipns/") {
		name = "/ipns/" + name
	}

	if _, ok := nc.follows[name]; ok {
		return fmt.Errorf("already following %s", name)
	}

	ctx, cancel := context.WithCancel(nc.ctx)
	go nc.followName(ctx, name, prefetch, followInterval)
	nc.follows[name] = cancel

	return nil
}

// Unfollow cancels a follow
func (nc *nameCache) Unfollow(name string) error {
	nc.mx.Lock()
	defer nc.mx.Unlock()

	if !strings.HasPrefix(name, "/ipns/") {
		name = "/ipns/" + name
	}

	cancel, ok := nc.follows[name]
	if !ok {
		return fmt.Errorf("unknown name %s", name)
	}

	cancel()
	delete(nc.follows, name)
	return nil
=======
// and (when pinit is true) pins it in the background
func (nc *nameCache) Follow(name string, pinit bool) {
	nc.mx.Lock()
	defer nc.mx.Unlock()

	if _, ok := nc.follows[name]; ok {
		return
	}

	ctx, cancel := context.WithCancel(nc.ctx)
	go nc.followName(ctx, name, pinit)
	nc.follows[name] = cancel
}

// Unfollow cancels a follow
func (nc *nameCache) Unfollow(name string) {
	nc.mx.Lock()
	defer nc.mx.Unlock()

	cancel, ok := nc.follows[name]
	if ok {
		cancel()
		delete(nc.follows, name)
	}
>>>>>>> namecache: ipfs name follow
}

// ListFollows returns a list of names currently being followed
func (nc *nameCache) ListFollows() []string {
	nc.mx.Lock()
	defer nc.mx.Unlock()

<<<<<<< HEAD
	follows := make([]string, 0, len(nc.follows))
	for name := range nc.follows {
=======
	follows := make([]string, 0)
	for name, _ := range nc.follows {
>>>>>>> namecache: ipfs name follow
		follows = append(follows, name)
	}

	return follows
}

<<<<<<< HEAD
func (nc *nameCache) followName(ctx context.Context, name string, prefetch bool, followInterval time.Duration) {
	emptynode := new(dag.ProtoNode)

	c, err := nc.resolveAndUpdate(ctx, name, prefetch, emptynode.Cid())
	if err != nil {
		log.Errorf("Error following %s: %s", name, err.Error())
	}
=======
func (nc *nameCache) followName(ctx context.Context, name string, pinit bool) {
	nc.resolveAndPin(ctx, name, pinit)
>>>>>>> namecache: ipfs name follow

	ticker := time.NewTicker(followInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
<<<<<<< HEAD
			c, err = nc.resolveAndUpdate(ctx, name, prefetch, c)

			if err != nil {
				log.Errorf("Error following %s: %s", name, err.Error())
			}
=======
			nc.resolveAndPin(ctx, name, pinit)
>>>>>>> namecache: ipfs name follow

		case <-ctx.Done():
			return
		}
	}
}

<<<<<<< HEAD
func (nc *nameCache) resolveAndUpdate(ctx context.Context, name string, prefetch bool, oldcid cid.Cid) (cid.Cid, error) {
	ptr, err := nc.resolve(ctx, name)
	if err != nil {
		return cid.Undef, err
	}

	newcid, err := pathToCid(ptr)
	if err != nil {
		return cid.Undef, err
	}

	if newcid.Equals(oldcid) || !prefetch {
		return newcid, nil
	}

	oldnd, err := nc.dag.Get(ctx, oldcid)
	if err != nil {
		return cid.Undef, err
	}

	newnd, err := nc.dag.Get(ctx, newcid)
	if err != nil {
		return cid.Undef, err
	}

	changes, err := dagutils.Diff(ctx, nc.dag, oldnd, newnd)
	if err != nil {
		return cid.Undef, err
	}

	log.Debugf("fetching changes in %s (%s -> %s)", name, oldcid, newcid)
	for _, change := range changes {
		if change.Type == iface.DiffRemove {
			continue
		}

		toFetch, err := nc.dag.Get(ctx, change.After)
		if err != nil {
			return cid.Undef, err
		}

		// just iterate over all nodes
		walker := ipld.NewWalker(ctx, ipld.NewNavigableIPLDNode(toFetch, nc.dag))
		if err := walker.Iterate(func(node ipld.NavigableNode) error {
			return nil
		}); err != ipld.EndOfDag {
			return cid.Undef, fmt.Errorf("unexpected error when prefetching followed name: %s", err)
		}
	}

	return newcid, err
}

func (nc *nameCache) resolve(ctx context.Context, name string) (path.Path, error) {
	log.Debugf("resolving %s", name)

	rctx, cancel := context.WithTimeout(ctx, resolveTimeout)
	defer cancel()

	p, err := nc.nsys.Resolve(rctx, name)
	if err != nil {
		return "", err
	}

	log.Debugf("resolved %s to %s", name, p)

	return p, nil
}

func pathToCid(p path.Path) (cid.Cid, error) {
	return cid.Decode(p.Segments()[1])
=======
func (nc *nameCache) resolveAndPin(ctx context.Context, name string, pinit bool) {
	log.Debugf("resolving %s", name)

	if !strings.HasPrefix(name, "/ipns/") {
		name = "/ipns/" + name
	}

	p, err := nc.nsys.Resolve(ctx, name)
	if err != nil {
		log.Debugf("error resolving %s: %s", name, err.Error())
		return
	}

	log.Debugf("resolved %s to %s", name, p)

	if !pinit {
		return
	}

	log.Debugf("pinning %s", p)

	r := &resolver.Resolver{
		DAG:         nc.dag,
		ResolveOnce: uio.ResolveUnixfsOnce,
	}

	n, err := r.ResolvePath(ctx, p)
	if err != nil {
		log.Debugf("error resolving path %s to node: %s", p, err.Error())
		return
	}

	err = nc.pinning.Pin(ctx, n, true)
	if err != nil {
		log.Debugf("error pinning path %s: %s", p, err.Error())
		return
	}

	err = nc.pinning.Flush()
	if err != nil {
		log.Debugf("error flushing pin: %s", err.Error())
	}
>>>>>>> namecache: ipfs name follow
}
