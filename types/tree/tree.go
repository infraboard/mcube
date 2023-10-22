package tree

func NewArcoDesignTree() *ArcoDesignTree {
	return &ArcoDesignTree{
		Items: []*ArcoDesignTreeNode{},
	}
}

type ArcoDesignTree struct {
	Items []*ArcoDesignTreeNode `json:"items"`
}

func (s *ArcoDesignTree) Add(item *ArcoDesignTreeNode) {
	s.Items = append(s.Items, item)
}

func (s *ArcoDesignTree) ForEatch(fn func(*ArcoDesignTreeNode)) {
	for i := range s.Items {
		fn(s.Items[i])
	}
}

func (s *ArcoDesignTree) GetOrCreateTreeByRootKey(
	key, title, nodeType string) *ArcoDesignTreeNode {
	for i := range s.Items {
		item := s.Items[i]
		if item.Key == key {
			return item
		}
	}

	item := NewArcoDesignTreeNode(key, title, nodeType)
	s.Add(item)
	return item
}

func NewArcoDesignTreeNode(key, title, nodeType string) *ArcoDesignTreeNode {
	if title == "" {
		title = key
	}
	return &ArcoDesignTreeNode{
		Key:      key,
		Title:    title,
		Type:     nodeType,
		Extra:    map[string]string{},
		Children: []*ArcoDesignTreeNode{},
	}
}

// https://arco.design/vue/component/tree#API
type ArcoDesignTreeNode struct {
	// 该节点显示的标题
	Title string `json:"title"`
	// 唯一标示
	Key string `json:"key"`
	// 是否禁用节点
	Disabled bool `json:"disabled"`
	// 是否是叶子节点。动态加载时有效
	IsLeaf bool `json:"is_leaf"`
	// 节点类型
	Type string `json:"type"`
	// 其他扩展属性
	Extra map[string]string `json:"extra"`
	// 子节点
	Children []*ArcoDesignTreeNode `json:"children"`
}

func (t *ArcoDesignTreeNode) SetTitle(title string) {
	if title == "" {
		return
	}

	t.Title = title
}

func (t *ArcoDesignTreeNode) Add(item *ArcoDesignTreeNode) {
	t.Children = append(t.Children, item)
}

func (t *ArcoDesignTreeNode) Walk(fn func(*ArcoDesignTreeNode)) {
	walk(t, fn)
}

func walk(t *ArcoDesignTreeNode, fn func(*ArcoDesignTreeNode)) {
	for i := range t.Children {
		fn(t.Children[i])
		walk(t.Children[i], fn)
	}
}

func (t *ArcoDesignTreeNode) GetOrCreateChildrenByKey(
	key, title, nodeType string) *ArcoDesignTreeNode {
	var item *ArcoDesignTreeNode
	t.Walk(func(adt *ArcoDesignTreeNode) {
		if adt.Key == key {
			item = adt
		}
	})
	if item == nil {
		item = NewArcoDesignTreeNode(key, title, nodeType)
		t.Add(item)
	}

	return item
}
