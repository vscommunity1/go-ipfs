package namecache

import (
	"encoding/json"
	"time"

	ds "gx/ipfs/Qmf4xQhNomPNhrtZc67qSnfJSjxjXs9LWvknJtSXwimPrM/go-datastore"
	nsds "gx/ipfs/Qmf4xQhNomPNhrtZc67qSnfJSjxjXs9LWvknJtSXwimPrM/go-datastore/namespace"
	dsq "gx/ipfs/Qmf4xQhNomPNhrtZc67qSnfJSjxjXs9LWvknJtSXwimPrM/go-datastore/query"
)

var dsPrefix = ds.NewKey("/namecache")

// persistent is a cache layer which persists followed names between node
// restarts
type persistent struct {
	NameCache

	ds ds.Datastore
}

type follow struct {
<<<<<<< HEAD
	Prefetch bool
=======
	Pin bool
>>>>>>> namecache: basic persistent version
	Deadline time.Time
}

func NewPersistentCache(base NameCache, d ds.Datastore) (NameCache, error) {
	d = nsds.Wrap(d, dsPrefix)

	q ,err := d.Query(dsq.Query{})
	if err != nil {
		return nil, err
	}
	defer q.Close()
	for e := range q.Next() {
		var f follow
		if err := json.Unmarshal(e.Value, &f); err != nil {
			return nil, err
		}
<<<<<<< HEAD
		if err := base.Follow(e.Key, f.Prefetch, time.Now().Sub(f.Deadline)); err != nil {
=======
		if err := base.Follow(e.Key, f.Pin, time.Now().Sub(f.Deadline)); err != nil {
>>>>>>> namecache: basic persistent version
			return nil, err
		}
	}


	return &persistent{
		NameCache: base,
		ds: d,
	}, nil
}

<<<<<<< HEAD
func (p *persistent) Follow(name string, prefetch bool, followInterval time.Duration) error {
	b, err := json.Marshal(&follow{
		Prefetch: prefetch,
=======
func (p *persistent) Follow(name string, dopin bool, followInterval time.Duration) error {
	b, err := json.Marshal(&follow{
		Pin: dopin,
>>>>>>> namecache: basic persistent version
		Deadline: time.Now().Add(followInterval),
	})
	if err != nil {
		return err
	}

<<<<<<< HEAD
	if err := p.NameCache.Follow(name, prefetch, followInterval); err != nil {
=======
	if err := p.NameCache.Follow(name, dopin, followInterval); err != nil {
>>>>>>> namecache: basic persistent version
		return err
	}
	return p.ds.Put(ds.NewKey(name), b)
}

func (p *persistent) Unfollow(name string) error {
	if err := p.NameCache.Unfollow(name); err != nil {
		return err
	}
	return p.ds.Delete(ds.NewKey(name))
}