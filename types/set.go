package types

import "github.com/infraboard/mcube/v2/tools/pretty"

// 构造函数
func New[T any]() *Set[T] {
	return &Set[T]{
		Items: []T{},
	}
}

type Set[T any] struct {
	Total int64 `json:"total"`
	Items []T   `json:"items"`
}

func (s *Set[T]) String() string {
	return pretty.ToJSON(s)
}

func (s *Set[T]) Add(items ...T) {
	s.Items = append(s.Items, items...)
	s.Total += int64(len(items))
}

func (s *Set[T]) Len() int {
	return len(s.Items)
}

type ItemHandleFunc[T any] func(t T)

func (s *Set[T]) ForEach(h ItemHandleFunc[T]) {
	for i := range s.Items {
		h(s.Items[i])
	}
}

type ItemFilterFunc[T any] func(t T) bool

func (s Set[T]) Filter(f ItemFilterFunc[T]) Set[T] {
	var filteredItems []T
	for _, item := range s.Items {
		if f(item) {
			filteredItems = append(filteredItems, item)
		}
	}
	return Set[T]{
		Total: int64(len(filteredItems)),
		Items: filteredItems,
	}
}
