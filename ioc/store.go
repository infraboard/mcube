package ioc

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/infraboard/mcube/tools/file"
)

var (
	store = newDefaultStore()
)

func ConfigIocObject(req *LoadConfigRequest) error {
	// 加载对象的配置
	err := store.LoadConfig(req)
	if err != nil {
		return err
	}

	// 初始化对象
	err = store.InitIocObject()
	return err
}

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

func newDefaultStore() *defaultStore {
	return &defaultStore{
		store: []*NamespaceStore{},
	}
}

type defaultStore struct {
	store []*NamespaceStore
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

// 获取一个对象存储空间
func (s *defaultStore) Namespace(namespace string) *NamespaceStore {
	for i := range s.store {
		item := s.store[i]
		if item.Namespace == namespace {
			return item
		}
	}

	ns := newNamespaceStore(namespace)
	s.store = append(s.store, ns)
	return ns
}

// 加载对象配置
func (s *defaultStore) LoadConfig(req *LoadConfigRequest) error {
	errs := []string{}

	// 优先加载环境变量
	if req.ConfigEnv.Enabled {
		for i := range s.store {
			item := s.store[i]
			err := item.LoadFromEnv(req.ConfigEnv.Prefix)
			if err != nil {
				errs = append(errs, err.Error())
			}
		}
	}

	// 再加载配置文件
	if req.ConfigFile.Enabled {
		for i := range s.store {
			item := s.store[i]
			err := item.LoadFromFile(req.ConfigFile.Path)
			if err != nil {
				errs = append(errs, err.Error())
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, ","))
	}

	return nil
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

func newNamespaceStore(namespace string) *NamespaceStore {
	return &NamespaceStore{
		Namespace: namespace,
		Items:     []Object{},
	}
}

type NamespaceStore struct {
	// 空间名称
	Namespace string
	// 空间优先级
	Priority int
	// 空间对象列表
	Items []Object
}

func (s *NamespaceStore) SetPriority(v int) *NamespaceStore {
	s.Priority = v
	return s
}

func (s *NamespaceStore) Registry(obj Object) {
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

func (s *NamespaceStore) Get(name string, opts ...GetOption) Object {
	opt := defaultOption().Apply(opts...)
	obj, _ := s.getWithIndex(name, opt.version)
	return obj
}

func (s *NamespaceStore) setWithIndex(index int, obj Object) {
	s.Items[index] = obj
}

func (s *NamespaceStore) getWithIndex(name, version string) (Object, int) {
	for i := range s.Items {
		obj := s.Items[i]
		if obj.Name() == name && obj.Version() == version {
			return obj, i
		}
	}

	return nil, -1
}

// 第一个
func (s *NamespaceStore) First() Object {
	if s.Len() == 0 {
		return nil
	}

	return s.Items[0]
}

// 最后一个
func (s *NamespaceStore) Last() Object {
	if s.Len() == 0 {
		return nil
	}

	return s.Items[s.Len()-1]
}

func (s *NamespaceStore) ForEach(fn func(Object)) {
	for i := range s.Items {
		item := s.Items[i]
		fn(item)
	}
}

func (s *NamespaceStore) ObjectUids() (uids []string) {
	for i := range s.Items {
		item := s.Items[i]
		uids = append(uids, ObjectUid(item))
	}
	return
}

func (s *NamespaceStore) Len() int {
	return len(s.Items)
}

func (s *NamespaceStore) Less(i, j int) bool {
	return s.Items[i].Priority() > s.Items[j].Priority()
}

func (s *NamespaceStore) Swap(i, j int) {
	s.Items[i], s.Items[j] = s.Items[j], s.Items[i]
}

// 根据对象的优先级进行排序
func (s *NamespaceStore) Sort() {
	sort.Sort(s)
}

func (s *NamespaceStore) Init() error {
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

func (s *NamespaceStore) Close() error {
	for i := range s.Items {
		obj := s.Items[i]
		obj.Destory()
	}
	return nil
}

// 从环境变量中加载对象配置
func (i *NamespaceStore) LoadFromEnv(prefix string) error {
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

func ValidateFileType(ext string) error {
	exist := false
	validateFileType := []string{".toml", ".yml", ".yaml", ".json"}
	for _, ft := range validateFileType {
		if ext == ft {
			exist = true
		}
	}
	if !exist {
		return fmt.Errorf("not support format: %s", ext)
	}

	return nil
}

func (i *NamespaceStore) LoadFromFile(filename string) error {
	if filename == "" {
		return nil
	}

	fileType := filepath.Ext(filename)
	if err := ValidateFileType(fileType); err != nil {
		return err
	}

	errs := []string{}
	i.ForEach(func(o Object) {
		cfg := map[string]Object{
			o.Name(): o,
		}
		var err error
		switch fileType {
		case ".toml":
			_, err = toml.DecodeFile(filename, cfg)
		case ".yml", ".yaml":
			err = file.ReadYamlFile(filename, cfg)
		case ".json":
			err = file.ReadJsonFile(filename, cfg)
		}
		if err != nil {
			errs = append(errs, err.Error())
		}
	})

	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, ","))
	}
	return nil
}
