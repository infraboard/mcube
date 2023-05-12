package static

import (
	"github.com/infraboard/mcube/grpc/balancer/wrr"
	"google.golang.org/grpc/resolver"
)

const (
	Scheme = "static"
)

// Following is an example name resolver implementation. Read the name
// resolution example to learn more about it.

type staticResolverBuilder struct{}

func (*staticResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &staticResolver{
		target: target,
		cc:     cc,
		store:  store,
	}
	r.start()
	return r, nil
}

func (*staticResolverBuilder) Scheme() string { return Scheme }

type staticResolver struct {
	target resolver.Target
	cc     resolver.ClientConn

	store *Store
}

func (r *staticResolver) start() {
	ts := r.store.Get(r.target.URL.Host)
	addrs := []resolver.Address{}
	for i := range ts.Items {
		item := ts.Items[i]
		addr := resolver.Address{Addr: item.Address}
		wrr.SetWeight(&addr, item.weight)
		addrs = append(addrs, addr)
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (*staticResolver) ResolveNow(o resolver.ResolveNowOptions) {}

func (*staticResolver) Close() {}

func init() {
	resolver.Register(&staticResolverBuilder{})
}
