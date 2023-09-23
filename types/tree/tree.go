package tree

func NewArcoDesignTreeSet() *ArcoDesignTreeSet {
	return &ArcoDesignTreeSet{
		Items: []*ArcoDesignTree{},
	}
}

type ArcoDesignTreeSet struct {
	Items []*ArcoDesignTree `json:"items"`
}

func (s *ArcoDesignTreeSet) Add(item *ArcoDesignTree) {
	s.Items = append(s.Items, item)
}

func (s *ArcoDesignTreeSet) ForEatch(fn func(*ArcoDesignTree)) {
	for i := range s.Items {
		fn(s.Items[i])
	}
}

func (s *ArcoDesignTreeSet) GetOrCreateTreeByRootKey(
	key, title string) *ArcoDesignTree {
	for i := range s.Items {
		item := s.Items[i]
		if item.Key == key {
			return item
		}
	}

	item := NewArcoDesignTree(key, title)
	s.Add(item)
	return item
}

func NewArcoDesignTree(key, title string) *ArcoDesignTree {
	return &ArcoDesignTree{
		Key:      key,
		Title:    title,
		Children: []*ArcoDesignTree{},
	}
}

type ArcoDesignTree struct {
	Title    string            `json:"title"`
	Key      string            `json:"key"`
	Children []*ArcoDesignTree `json:"children"`
}

func (t *ArcoDesignTree) Add(item *ArcoDesignTree) {
	t.Children = append(t.Children, item)
}

func (t *ArcoDesignTree) GetOrCreateChildrenByKey(
	key, title string, deep int) *ArcoDesignTree {
	for i := range t.Children {
		c := t.Children[i]
		if c.Key == key {
			return c
		}
	}

	item := NewArcoDesignTree(key, title)
	t.Add(item)
	return item
}
