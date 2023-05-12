package static

import (
	"github.com/infraboard/mcube/tools/pretty"
)

var (
	store = NewStore()
)

func GetStore() *Store {
	return store
}

func NewStore() *Store {
	return &Store{
		targets: map[string]*TargetSet{},
	}
}

type Store struct {
	targets map[string]*TargetSet
}

func (s *Store) Get(service string) *TargetSet {
	v, ok := s.targets[service]
	if ok {
		return v
	}

	v = NewTargetSet()
	s.targets[service] = v
	return v
}

func (s *Store) Add(service string, targets ...*Target) {
	v := s.Get(service)
	v.Add(targets...)
}

func NewTargetSet() *TargetSet {
	return &TargetSet{
		Items: []*Target{},
	}
}

type TargetSet struct {
	Items []*Target
}

func (s *TargetSet) Add(ts ...*Target) {
	s.Items = append(s.Items, ts...)
}

func NewTarget(address string) *Target {
	return &Target{
		Address: address,
		weight:  1,
	}
}

type Target struct {
	Address string
	weight  uint32
}

func (t *Target) String() string {
	return pretty.ToJSON(t)
}

func (t *Target) SetWeight(weight uint32) *Target {
	t.weight = weight
	return t
}
