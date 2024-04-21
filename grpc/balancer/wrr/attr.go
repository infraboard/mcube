package wrr

import (
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

const (
	WEIGHT_ATTRIBUTE_KEY = "weight"
)

func SetWeight(addr *resolver.Address, weight uint32) {
	if addr.Attributes == nil {
		addr.Attributes = attributes.New(WEIGHT_ATTRIBUTE_KEY, weight)
	} else {
		addr.Attributes.WithValue(WEIGHT_ATTRIBUTE_KEY, weight)
	}
}

func GetWeight(addr resolver.Address) uint32 {
	v := addr.Attributes.Value(WEIGHT_ATTRIBUTE_KEY)
	ai, _ := v.(uint32)
	return ai
}
