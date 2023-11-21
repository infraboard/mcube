package ioc

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
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
	if err != nil {
		return err
	}

	// 依赖自动注入
	return store.Autowire()
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
		store: []*NamespaceStore{
			newNamespaceStore(CONFIG_NAMESPACE).SetPriority(99),
			newNamespaceStore(CONTROLLER_NAMESPACE).SetPriority(0),
			newNamespaceStore(DEFAULT_NAMESPACE).SetPriority(9),
			newNamespaceStore(API_NAMESPACE).SetPriority(-99),
		},
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

// 初始化托管的所有对象
func (s *defaultStore) Autowire() error {
	for i := range s.store {
		item := s.store[i]
		err := item.Autowire()
		if err != nil {
			return fmt.Errorf("[%s] %s", item.Namespace, err)
		}
	}
	return nil
}

func newNamespaceStore(namespace string) *NamespaceStore {
	return &NamespaceStore{
		Namespace: namespace,
		Items:     []*ObjectWrapper{},
	}
}

type NamespaceStore struct {
	// 空间名称
	Namespace string
	// 空间优先级
	Priority int
	// 空间对象列表
	Items []*ObjectWrapper
}

func NewObjectWrapper(obj Object) *ObjectWrapper {
	name, version := GetIocObjectUid(obj)
	return &ObjectWrapper{
		Name:           name,
		Version:        version,
		Priority:       obj.Priority(),
		AllowOverwrite: obj.AllowOverwrite(),
		Value:          obj,
	}
}

type ObjectWrapper struct {
	Name           string
	Version        string
	AllowOverwrite bool
	Priority       int
	Value          Object
}

func (s *NamespaceStore) SetPriority(v int) *NamespaceStore {
	s.Priority = v
	return s
}

func (s *NamespaceStore) Registry(v Object) {
	obj := NewObjectWrapper(v)
	old, index := s.getWithIndex(obj.Name, obj.Version)
	// 没有, 直接添加
	if old == nil {
		s.Items = append(s.Items, obj)
		return
	}

	// 有, 允许盖写则直接修改
	if obj.AllowOverwrite {
		zap.L().Infof("%s object overwrite", obj.Name)
		s.setWithIndex(index, obj)
		return
	}

	// 有, 不允许修改
	panic(fmt.Sprintf("ioc obj %s has registed", obj.Name))
}

func (s *NamespaceStore) Get(name string, opts ...GetOption) Object {
	opt := defaultOption().Apply(opts...)
	obj, _ := s.getWithIndex(name, opt.version)
	if obj == nil {
		return nil
	}
	return obj.Value
}

// 根据对象对象加载对象
func (s *NamespaceStore) Load(target any, opts ...GetOption) error {
	t := reflect.TypeOf(target)
	v := reflect.ValueOf(target)

	var obj Object
	switch t.Kind() {
	case reflect.Interface:
		objs := s.ImplementInterface(t)
		if len(objs) > 0 {
			obj = objs[0]
		}
	default:
		obj = s.Get(t.String(), opts...)
	}

	// 注入值
	if obj != nil {
		objValue := reflect.ValueOf(obj)
		if !(v.Kind() == reflect.Ptr && objValue.Kind() == reflect.Ptr) {
			return fmt.Errorf("target and object must both be pointers or non-pointers")
		}
		v.Elem().Set(objValue.Elem())
	}
	return nil
}

func (s *NamespaceStore) setWithIndex(index int, obj *ObjectWrapper) {
	s.Items[index] = obj
}

func (s *NamespaceStore) getWithIndex(name, version string) (*ObjectWrapper, int) {
	for i := range s.Items {
		obj := s.Items[i]
		if obj.Name == name && obj.Version == version {
			return obj, i
		}
	}

	return nil, -1
}

// 对象个数统计
func (s *NamespaceStore) Count() int {
	return s.Len()
}

// 第一个
func (s *NamespaceStore) First() Object {
	if s.Len() == 0 {
		return nil
	}

	return s.Items[0].Value
}

// 最后一个
func (s *NamespaceStore) Last() Object {
	if s.Len() == 0 {
		return nil
	}

	return s.Items[s.Len()-1].Value
}

func (s *NamespaceStore) ForEach(fn func(*ObjectWrapper)) {
	for i := range s.Items {
		item := s.Items[i]
		fn(item)
	}
}

// 寻找实现了接口的对象
func (s *NamespaceStore) ImplementInterface(objType reflect.Type, opts ...GetOption) (objs []Object) {
	opt := defaultOption().Apply(opts...)

	for i := range s.Items {
		o := s.Items[i]
		// 断言获取的对象是否满足接口类型
		if o != nil && reflect.TypeOf(o.Value).Implements(objType) {
			if o.Version == opt.version {
				objs = append(objs, o.Value)
			}
		}
	}
	return
}

func (s *NamespaceStore) List() (uids []string) {
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
	return s.Items[i].Value.Priority() > s.Items[j].Value.Priority()
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
		err := obj.Value.Init()
		if err != nil {
			return fmt.Errorf("init object %s error, %s", obj.Name, err)
		}
	}
	return nil
}

func (s *NamespaceStore) Close(ctx context.Context) error {
	errs := []string{}
	for i := range s.Items {
		obj := s.Items[i]
		if err := obj.Value.Close(ctx); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("close error, %s", strings.Join(errs, ","))
	}
	return nil
}

// 从环境变量中加载对象配置
func (i *NamespaceStore) LoadFromEnv(prefix string) error {
	errs := []string{}
	i.ForEach(func(w *ObjectWrapper) {
		err := env.Parse(w.Value, env.Options{
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

// 从环境变量中加载对象配置
func (i *NamespaceStore) Autowire() error {
	i.ForEach(func(w *ObjectWrapper) {
		pt := reflect.TypeOf(w.Value).Elem()
		// go语言所有函数传的都是值，所以要想修改原来的值就需要传指
		// 通过Elem()返回指针指向的对象
		v := reflect.ValueOf(w.Value).Elem()

		for i := 0; i < pt.NumField(); i++ {
			tag := ParseInjectTag(pt.Field(i).Tag.Get("ioc"))
			if tag.Autowire {
				fieldType := v.Field(i).Type()
				var obj Object
				// 根据字段的类型获取值
				switch fieldType.Kind() {
				case reflect.Interface:
					// 为接口类型注入值
					objs := store.Namespace(tag.Namespace).ImplementInterface(fieldType)
					if len(objs) > 0 {
						obj = objs[0]
					}
				default:
					// 为结构体变量注入值
					if tag.Name == "" {
						tag.Name = fieldType.String()
					}
					fmt.Println(tag.Version)
					obj = store.Namespace(tag.Namespace).Get(tag.Name, WithVersion(tag.Version))
				}
				// 注入值
				if obj != nil {
					v.Field(i).Set(reflect.ValueOf(obj))
				}
			}
		}
	})
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

	// 准备一个map读取配置
	cfg := map[string]any{}
	i.ForEach(func(w *ObjectWrapper) {
		cfg[w.Value.Name()] = w.Value
	})

	var err error
	switch fileType {
	case ".toml":
		_, err = toml.DecodeFile(filename, &cfg)
	case ".yml", ".yaml":
		err = file.ReadYamlFile(filename, &cfg)
	case ".json":
		err = file.ReadJsonFile(filename, &cfg)
	default:
		err = fmt.Errorf("unspport format: %s", fileType)
	}
	if err != nil {
		return err
	}

	// 加载到对象中
	errs := []string{}
	i.ForEach(func(w *ObjectWrapper) {
		dj, err := json.Marshal(cfg[w.Value.Name()])
		if err != nil {
			errs = append(errs, err.Error())
		}

		err = json.Unmarshal(dj, w.Value)
		if err != nil {
			errs = append(errs, err.Error())
		}
	})
	if len(errs) > 0 {
		return fmt.Errorf("load config error, %s", strings.Join(errs, ","))
	}

	return nil
}
