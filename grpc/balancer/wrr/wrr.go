package wrr

import (
	"fmt"
	"sync/atomic"

	"github.com/infraboard/mcenter/common/rand"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

// Name is the name of weighted_round_robin balancer.
const Name = "weighted_round_robin"

var logger = zap.L().Named(Name)

// newBuilder creates a new weighted roundrobin balancer builder.
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &rrPickerBuilder{}, base.Config{HealthCheck: false})
}

type rrPickerBuilder struct{}

func (*rrPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	logger.Infof("roundrobinPicker: Build called with info: %v", info)
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	scs := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for subConn, subConnInfo := range info.ReadySCs {
		weight := GetWeight(subConnInfo.Address)
		for i := uint32(0); i < weight; i++ {
			scs = append(scs, subConn)
		}
		scs = append(scs, subConn)
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
	fmt.Println(sc)
	logger.Infof("pick conn: %s", sc)
	return balancer.PickResult{SubConn: sc}, nil
}

func init() {
	balancer.Register(newBuilder())
}
