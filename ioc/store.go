package ioc

import (
	"fmt"
	"sort"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/infraboard/mcube/logger/zap"
)

var (
	store = newDefaultsapceStore()
)

// 初始化对象
func InitIocObject() error {
	return store.InitIocObject()
}

// 注册对象
func RegistryObjectWithNs(namespace string, obj Object) {
	store.Namespace(namespace).Registry(obj)
}

// 获取对象
func GetObjectWithNs(namespace, name string) Object {
	obj := store.Namespace(namespace).Get(name)
	if obj == nil {
		panic(fmt.Sprintf("ioc obj %s not registed", name))
	}

	return obj
}

func newDefaultsapceStore() *defaultStore {
	return &defaultStore{
		store: []*ObjectSet{},
	}
}

type defaultStore struct {
	store []*ObjectSet
}

func (s *defaultStore) Len() int {
	return len(s.store)
}

func (s *defaultStore) Less(i, j int) bool {
	return s.store[i].Priority > s.store[j].Priority
}

func (s *defaultStore) Swap(i, j int) {
	s.store[i], s.store[j] = s.store[j], s.store[i]
}

// 根据空间优先级进行排序
func (s *defaultStore) Sort() {
	sort.Sort(s)
}

// 初始化托管的所有对象
func (s *defaultStore) InitIocObject() error {
	s.Sort()

	for i := range s.store {
		item := s.store[i]
		err := item.Init()
		if err != nil {
			return fmt.Errorf("[%s] %s", item.Namespace, err)
		}
	}
	return nil
}

func (s *defaultStore) Namespace(namespace string) *ObjectSet {
	for i := range s.store {
		item := s.store[i]
		if item.Namespace == namespace {
			return item
		}
	}

	ns := newIocObjectSet(namespace)
	s.store = append(s.store, ns)
	return ns
}

func newIocObjectSet(namespace string) *ObjectSet {
	return &ObjectSet{
		Namespace: namespace,
		Items:     []Object{},
	}
}

type ObjectSet struct {
	// 空间名称
	Namespace string
	// 优先级
	Priority int
	// 对象列表
	Items []Object
}

func (s *ObjectSet) SetPriority(v int) *ObjectSet {
	s.Priority = v
	return s
}

func (s *ObjectSet) Registry(obj Object) {
	err := ValidateIocObject(obj)
	if err != nil {
		panic(err)
	}

	old, index := s.getWithIndex(obj.Name(), obj.Version())
	// 没有, 直接添加
	if old == nil {
		s.Items = append(s.Items, obj)
		return
	}

	// 有, 允许盖写则直接修改
	if obj.AllowOverwrite() {
		zap.L().Infof("%s object overwrite", obj.Name())
		s.setWithIndex(index, obj)
		return
	}

	// 有, 不允许修改
	panic(fmt.Sprintf("ioc obj %s has registed", obj.Name()))
}

func (s *ObjectSet) Get(name string, opts ...GetOption) Object {
	opt := defaultOption().Apply(opts...)
	obj, _ := s.getWithIndex(name, opt.version)
	return obj
}

func (s *ObjectSet) setWithIndex(index int, obj Object) {
	s.Items[index] = obj
}

func (s *ObjectSet) getWithIndex(name, version string) (Object, int) {
	for i := range s.Items {
		obj := s.Items[i]
		if obj.Name() == name && obj.Version() == version {
			return obj, i
		}
	}

	return nil, -1
}

// 第一个
func (s *ObjectSet) First() Object {
	if s.Len() == 0 {
		return nil
	}

	return s.Items[0]
}

// 最后一个
func (s *ObjectSet) Last() Object {
	if s.Len() == 0 {
		return nil
	}

	return s.Items[s.Len()-1]
}

func (s *ObjectSet) ForEach(fn func(Object)) {
	for i := range s.Items {
		item := s.Items[i]
		fn(item)
	}
}

func (s *ObjectSet) ObjectUids() (uids []string) {
	for i := range s.Items {
		item := s.Items[i]
		uids = append(uids, ObjectUid(item))
	}
	return
}

func (s *ObjectSet) Len() int {
	return len(s.Items)
}

func (s *ObjectSet) Less(i, j int) bool {
	return s.Items[i].Priority() > s.Items[j].Priority()
}

func (s *ObjectSet) Swap(i, j int) {
	s.Items[i], s.Items[j] = s.Items[j], s.Items[i]
}

// 根据对象的优先级进行排序
func (s *ObjectSet) Sort() {
	sort.Sort(s)
}

func (s *ObjectSet) Init() error {
	s.Sort()
	for i := range s.Items {
		obj := s.Items[i]
		err := obj.Init()
		if err != nil {
			return fmt.Errorf("init object %s error, %s", obj.Name(), err)
		}
	}
	return nil
}

func (s *ObjectSet) Close() error {
	for i := range s.Items {
		obj := s.Items[i]
		obj.Destory()
	}
	return nil
}

// 从环境变量中加载对象配置
func (i *ObjectSet) LoadFromEnv(prefix string) error {
	errs := []string{}
	i.ForEach(func(o Object) {
		err := env.Parse(o, env.Options{
			Prefix: prefix,
		})
		if err != nil {
			errs = append(errs, err.Error())
		}
	})
	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, ","))
	}

	return nil
}

func (i *ObjectSet) LoadFromFile() error {
	return nil
}
