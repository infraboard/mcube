package tree

import "fmt"

func NewArcoDesignTreeSet() *ArcoDesignTreeSet {
	return &ArcoDesignTreeSet{
		Items: []*ArcoDesignTree{},
	}
}

type ArcoDesignTreeSet struct {
	Items []*ArcoDesignTree `json:"items"`
}

func (s *ArcoDesignTreeSet) GetOrCreateTreeByRootKey(key, title string) *ArcoDesignTree {
	return nil
}

func NewArcoDesignTree() *ArcoDesignTree {
	return &ArcoDesignTree{
		Children: []*ArcoDesignTree{},
	}
}

type ArcoDesignTree struct {
	Title    string            `json:"title"`
	Key      string            `json:"key"`
	Children []*ArcoDesignTree `json:"children"`
}

func (t *ArcoDesignTree) GetOrCreateChildrenByKey(
	key, title string, deep int) *ArcoDesignTree {
	for i := range t.Children {
		c := t.Children[i]
		fmt.Print(c)
	}
	return nil
}
