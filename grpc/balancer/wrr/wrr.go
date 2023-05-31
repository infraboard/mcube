package wrr

import (
	"math/rand"
	"sync/atomic"

	"github.com/infraboard/mcube/tools/pretty"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
)

// Name is the name of weighted_round_robin balancer.
const Name = "weighted_round_robin"

var logger = grpclog.Component("weighted_round_robin")

// newBuilder creates a new weighted roundrobin balancer builder.
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &rrPickerBuilder{}, base.Config{HealthCheck: true})
}

type rrPickerBuilder struct{}

func (*rrPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	logger.Infof("WeightedRoundrobinPicker: Build called with info: %v", pretty.ToJSON(info))
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	scs := []balancer.SubConn{}
	for subConn, subConnInfo := range info.ReadySCs {
		weight := GetWeight(subConnInfo.Address)
		logger.Infof("WeightedRoundrobinPicker: address: %s, weighted: %d", subConnInfo.Address.Addr, weight)
		for i := uint32(0); i < weight; i++ {
			scs = append(scs, subConn)
		}
	}
	return &rrPicker{
		subConns: scs,
		// Start at a random index, as the same RR balancer rebuilds a new
		// picker when SubConn states change, and we don't want to apply excess
		// load to the first server in the list.
		next: uint32(rand.Intn(len(scs))),
	}
}

type rrPicker struct {
	// subConns is the snapshot of the roundrobin balancer when this picker was
	// created. The slice is immutable. Each Get() will do a round robin
	// selection from it and return the selected SubConn.
	subConns []balancer.SubConn
	next     uint32
}

func (p *rrPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	subConnsLen := uint32(len(p.subConns))
	nextIndex := atomic.AddUint32(&p.next, 1)

	sc := p.subConns[nextIndex%subConnsLen]
	return balancer.PickResult{SubConn: sc}, nil
}

func init() {
	balancer.Register(newBuilder())
}
