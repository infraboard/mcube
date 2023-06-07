package ioc

import (
	"fmt"
	"sort"
)

var (
	store = NewDefaultStore()
)

const (
	DefaultNamespace = "default"
)

// 初始化对象
func InitIocObject() error {
	return store.InitIocObject()
}

// 注册对象到默认空间
func RegistryObject(obj IocObject) {
	RegistryObjectWithNs(DefaultNamespace, obj)
}

// 获取默认空间对象
func GetObject(name string) IocObject {
	return GetObjectWithNs(DefaultNamespace, name)
}

// 注册对象
func RegistryObjectWithNs(namespace string, obj IocObject) {
	store.Namespace(namespace).Add(obj)
}

// 获取对象
func GetObjectWithNs(namespace, name string) IocObject {
	obj := store.Namespace(namespace).Get(name)
	if obj == nil {
		panic(fmt.Sprintf("ioc obj %s not registed", name))
	}

	return obj
}

func NewDefaultStore() *DefaultStore {
	return &DefaultStore{
		store: map[string]*IocObjectSet{},
	}
}

type DefaultStore struct {
	store map[string]*IocObjectSet
}

// 初始化托管的所有对象
func (s *DefaultStore) InitIocObject() error {
	for ns, objects := range s.store {
		objects.Sort()
		for i := range objects.Items {
			obj := objects.Items[i]
			err := obj.Init()
			if err != nil {
				return fmt.Errorf("init object %s.%s error, %s", ns, obj.Name(), err)
			}
		}
	}
	return nil
}

func (s *DefaultStore) Namespace(namespace string) *IocObjectSet {
	if v, ok := s.store[namespace]; ok {
		return v
	}
	ns := NewIocObjectSet()
	s.store[namespace] = ns
	return ns
}

func (s *DefaultStore) ShowRegistryObjectNames() (names []string) {
	for ns, v := range s.store {
		for i := range v.Items {
			obj := v.Items[i]
			names = append(names, ns+"."+obj.Name())
		}
	}
	return
}

func NewIocObjectSet() *IocObjectSet {
	return &IocObjectSet{
		Items: []IocObject{},
	}
}

type IocObjectSet struct {
	Items []IocObject
}

func (s *IocObjectSet) Add(obj IocObject) {
	if s.Exist(obj.Name()) {
		panic(fmt.Sprintf("ioc obj %s has registed", obj.Name()))
	}

	s.Items = append(s.Items, obj)
}

func (s *IocObjectSet) Get(name string) IocObject {
	for i := range s.Items {
		item := s.Items[i]
		if item.Name() == name {
			return item
		}
	}

	return nil
}

func (s *IocObjectSet) Exist(name string) bool {
	obj := s.Get(name)
	return obj != nil
}

func (s *IocObjectSet) ObjectNames() (names []string) {
	for i := range s.Items {
		item := s.Items[i]
		names = append(names, item.Name())
	}
	return
}

func (s *IocObjectSet) Len() int {
	return len(s.Items)
}

func (s *IocObjectSet) Less(i, j int) bool {
	return s.Items[i].Priority() > s.Items[j].Priority()
}

func (s *IocObjectSet) Swap(i, j int) {
	s.Items[i], s.Items[j] = s.Items[j], s.Items[i]
}

// 根据对象的优先级进行排序
func (s *IocObjectSet) Sort() {
	sort.Sort(s)
}
