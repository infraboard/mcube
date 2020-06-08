package enum

import (
	"regexp"
	"strings"
)

// NewEnumSet todo
func NewEnumSet() *Set {
	return &Set{
		Items: []*Enum{},
	}
}

// Set 枚举集合
type Set struct {
	Items []*Enum
}

// Length 长度
func (s *Set) Length() int {
	return len(s.Items)
}

// Get 获取一个枚举类
func (s *Set) Get(name string) *Enum {
	for _, e := range s.Items {
		if e.Name == name {
			return e
		}
	}

	enum := NewEnum(name)
	s.Add(enum)
	return enum
}

// Add 添加一个类型
func (s *Set) Add(i *Enum) {
	s.Items = append(s.Items, i)
}

// GetLatest 获取最新一个
func (s *Set) GetLatest() *Enum {
	if s.Length() == 0 {
		return s.Get("default")
	}

	return s.Items[s.Length()-1]
}

// NewEnum todo
func NewEnum(name string) *Enum {
	return &Enum{
		Name:  name,
		Items: []*Item{},
	}
}

// Enum 枚举类型
type Enum struct {
	Name  string
	Doc   string
	Items []*Item
}

// Add todo
func (e *Enum) Add(i *Item) {
	e.Items = append(e.Items, i)
}

// NewItem todo
func NewItem(name, doc string) *Item {
	return &Item{
		Name: name,
		Doc:  doc,
	}
}

var (
	itemParamRe = regexp.MustCompile(`.*?\((.*?)\).*`)
)

// Item 枚举项
// 通过获取注释里()中的内容作为自定义参数
type Item struct {
	Name string // 枚举的名称, 对应常量的名称
	Doc  string // 文档, 常量的文档
}

// Show 枚举项显示
func (i *Item) Show() string {
	matchs := itemParamRe.FindStringSubmatch(i.Doc)
	if len(matchs) < 2 {
		return i.defaultShow()
	}

	return matchs[1]
}

func (i *Item) defaultShow() string {
	return strings.ToLower(i.Name)
}
