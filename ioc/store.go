package ioc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

type defaultStore struct {
	conf  *LoadConfigRequest
	store []*NamespaceStore
}

func (s *defaultStore) Len() int {
	return len(s.store)
}

func (s *defaultStore) ForEatch(fn func(*NamespaceStore)) {
	for _, item := range s.store {
		fn(item)
	}
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

	// 先加载配置文件（多个文件按顺序加载，后面的覆盖前面的）
	if req.ConfigFile.Enabled {
		for _, path := range req.ConfigFile.Paths {
			// 检查文件是否存在
			if !IsFileExists(path) {
				if !req.ConfigFile.SkipIFNotExist {
					return fmt.Errorf("file %s not exist", path)
				}
				debug("skip not exist config file: %s", path)
				continue
			}

			fileType := filepath.Ext(path)
			if err := ValidateFileType(fileType); err != nil {
				return fmt.Errorf("file %s: %w", path, err)
			}

			// 读取文件内容
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", path, err)
			}

			debug("loading config file: %s", path)
			// 为每个namespace加载配置（配置会合并）
			for i := range s.store {
				item := s.store[i]
				err := item.LoadFromFileContent(content, fileType)
				if err != nil {
					errs = append(errs, fmt.Sprintf("file %s: %s", path, err.Error()))
				}
			}
		}
	}

	// 最后加载环境变量（优先级最高，会覆盖文件配置）
	if req.ConfigEnv.Enabled {
		debug("loading config from environment variables with prefix: %s", req.ConfigEnv.Prefix)
		for i := range s.store {
			item := s.store[i]
			err := item.LoadFromEnv(req.ConfigEnv.Prefix)
			if err != nil {
				errs = append(errs, err.Error())
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, ","))
	}

	s.conf = req
	return nil
}

// CallPostConfigHooks 调用所有对象的 PostConfig 钩子
func (s *defaultStore) CallPostConfigHooks() error {
	for i := range s.store {
		namespace := s.store[i]
		err := namespace.CallPostConfigHooks()
		if err != nil {
			return fmt.Errorf("[%s] %s", namespace.Namespace, err)
		}
	}
	return nil
}

// CallPostConfigHooks 调用命名空间内所有对象的 PostConfig 钩子
func (s *NamespaceStore) CallPostConfigHooks() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.Items {
		obj := s.Items[i]
		if hook, ok := obj.Value.(PostConfigHook); ok {
			debug("calling PostConfig hook for %s", obj.Name)
			if err := hook.OnPostConfig(); err != nil {
				return fmt.Errorf("PostConfig hook failed for %s: %w", obj.Name, err)
			}
		}
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

// 倒序遍历s.store, 关闭对象
func (s *defaultStore) Stop(ctx context.Context) {
	for i := len(s.store) - 1; i >= 0; i-- {
		item := s.store[i]
		item.Close(ctx)
	}
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
	// 并发安全保护
	mu sync.RWMutex
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
		Name:     name,
		Version:  version,
		Priority: obj.Priority(),
		Value:    obj,
	}
}

type ObjectWrapper struct {
	Name     string
	Version  string
	Priority int
	Value    Object
}

func (s *NamespaceStore) SetPriority(v int) *NamespaceStore {
	s.Priority = v
	return s
}

func (s *NamespaceStore) Registry(v Object) StoreUser {
	s.mu.Lock()
	defer s.mu.Unlock()

	obj := NewObjectWrapper(v)
	// 根据名字查找（不考虑版本）
	old, index := s.getByNameWithIndex(obj.Name)
	// 没有同名对象, 直接添加
	if old == nil {
		s.Items = append(s.Items, obj)
		return s
	}

	// 已存在同名对象，使用智能版本覆盖策略：高版本可以覆盖低版本
	cmp := CompareVersion(obj.Version, old.Version)
	if cmp > 0 {
		// 新版本更高，允许覆盖
		s.setWithIndex(index, obj)
		return s
	} else if cmp == 0 {
		// 相同版本，不允许重复注册
		panic(fmt.Sprintf("ioc obj %s (version %s) has already registered", obj.Name, obj.Version))
	} else {
		// 新版本更低，不允许覆盖
		panic(fmt.Sprintf("ioc obj %s: cannot register lower version %s (current: %s)", obj.Name, obj.Version, old.Version))
	}
}

// RegistryAll 批量注册对象
func (s *NamespaceStore) RegistryAll(objs ...Object) StoreUser {
	for _, obj := range objs {
		s.Registry(obj)
	}
	return s
}

func (s *NamespaceStore) Get(name string, opts ...GetOption) Object {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 智能版本覆盖策略下，Get 只根据名字查找，返回当前版本的对象
	obj, _ := s.getByNameWithIndex(name)
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

// getByNameWithIndex 根据名字查找对象（不考虑版本）
func (s *NamespaceStore) getByNameWithIndex(name string) (*ObjectWrapper, int) {
	for i := range s.Items {
		obj := s.Items[i]
		if obj.Name == name {
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
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.Items) == 0 {
		return nil
	}

	return s.Items[0].Value
}

// 最后一个
func (s *NamespaceStore) Last() Object {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.Items) == 0 {
		return nil
	}

	return s.Items[len(s.Items)-1].Value
}

func (s *NamespaceStore) ForEach(fn func(*ObjectWrapper)) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.Items {
		item := s.Items[i]
		fn(item)
	}
}

// 寻找实现了接口的对象
func (s *NamespaceStore) ImplementInterface(objType reflect.Type, opts ...GetOption) (objs []Object) {
	s.mu.RLock()
	defer s.mu.RUnlock()

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
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.Items {
		item := s.Items[i]
		uids = append(uids, ObjectUid(item))
	}
	return
}

func (s *NamespaceStore) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
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
	s.mu.Lock()
	defer s.mu.Unlock()
	sort.Sort(s)
}

func (s *NamespaceStore) Init() error {
	s.mu.Lock()
	// 使用标准库的sort.Slice，避免调用Len()方法造成死锁
	sort.Slice(s.Items, func(i, j int) bool {
		return s.Items[i].Value.Priority() > s.Items[j].Value.Priority()
	})
	s.mu.Unlock()

	debug("init namespace [%s] app ...", s.Namespace)
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.Items {
		obj := s.Items[i]

		// 1. PreInit 钩子
		if hook, ok := obj.Value.(PreInitHook); ok {
			debug("calling PreInit hook for %s", obj.Name)
			if err := hook.OnPreInit(); err != nil {
				return fmt.Errorf("PreInit hook failed for %s: %w", obj.Name, err)
			}
		}

		// 2. 主要初始化
		err := obj.Value.Init()
		if err != nil {
			return fmt.Errorf("init object %s error, %s", obj.Name, err)
		}
		debug("init app %s[priority: %d] ok.", obj.Value.Name(), obj.Value.Priority())

		// 3. PostInit 钩子
		if hook, ok := obj.Value.(PostInitHook); ok {
			debug("calling PostInit hook for %s", obj.Name)
			if err := hook.OnPostInit(); err != nil {
				// PostInit 错误只记录，不中断流程
				debug("PostInit hook failed for %s: %v", obj.Name, err)
			}
		}
	}
	return nil
}

// 倒序关闭
func (s *NamespaceStore) Close(ctx context.Context) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	debug("closing namespace [%s] app ...", s.Namespace)

	for i := len(s.Items) - 1; i >= 0; i-- {
		obj := s.Items[i]

		// 1. PreStop 钩子
		if hook, ok := obj.Value.(PreStopHook); ok {
			debug("calling PreStop hook for %s", obj.Name)
			if err := hook.OnPreStop(ctx); err != nil {
				// PreStop 错误只记录，不中断关闭流程
				debug("PreStop hook failed for %s: %v", obj.Name, err)
			}
		}

		// 2. 主要清理
		obj.Value.Close(ctx)
		debug("closed app %s", obj.Value.Name())

		// 3. PostStop 钩子
		if hook, ok := obj.Value.(PostStopHook); ok {
			debug("calling PostStop hook for %s", obj.Name)
			if err := hook.OnPostStop(ctx); err != nil {
				// PostStop 错误只记录，不中断流程
				debug("PostStop hook failed for %s: %v", obj.Name, err)
			}
		}
	}
}

// 从环境变量中加载对象配置
func (i *NamespaceStore) LoadFromEnv(prefix string) error {
	errs := []string{}
	i.ForEach(func(w *ObjectWrapper) {
		prefixList := strings.ToUpper(w.Name) + "_"
		if prefix != "" {
			prefixList = fmt.Sprintf("%s_%s", strings.ToUpper(prefix), prefixList)
		}
		err := env.Parse(w.Value, env.Options{
			Prefix: prefixList,
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
	var errs []string

	i.ForEach(func(w *ObjectWrapper) {
		pt := reflect.TypeOf(w.Value).Elem()
		// go语言所有函数传的都是值，所以要想修改原来的值就需要传指
		// 通过Elem()返回指针指向的对象
		v := reflect.ValueOf(w.Value).Elem()

		for i := 0; i < pt.NumField(); i++ {
			fieldTag := pt.Field(i).Tag.Get("ioc")
			if fieldTag == "" {
				continue
			}

			tag, err := ParseInjectTagWithError(fieldTag)
			if err != nil {
				errs = append(errs, fmt.Sprintf("parse tag for %s.%s: %v",
					pt.Name(), pt.Field(i).Name, err))
				continue
			}

			if tag.Autowire {
				fieldType := v.Field(i).Type()
				var obj Object
				// 根据字段的类型获取值
				switch fieldType.Kind() {
				case reflect.Interface:
					// 为接口类型注入值
					if tag.Name != "" {
						// 如果指定了 name，直接获取指定的对象（支持接口字段指定具体实现）
						obj = DefaultStore.Namespace(tag.Namespace).Get(tag.Name, WithVersion(tag.Version))
					} else {
						// 否则自动查找实现该接口的对象（取第一个）
						objs := DefaultStore.Namespace(tag.Namespace).ImplementInterface(fieldType)
						if len(objs) > 0 {
							obj = objs[0]
						}
					}
				default:
					// 为结构体变量注入值
					if tag.Name == "" {
						tag.Name = fieldType.String()
					}
					fmt.Println(tag.Version)
					obj = DefaultStore.Namespace(tag.Namespace).Get(tag.Name, WithVersion(tag.Version))
				}
				// 注入值
				if obj != nil {
					v.Field(i).Set(reflect.ValueOf(obj))
				}
			}
		}
	})

	if len(errs) > 0 {
		return fmt.Errorf("autowire errors:\n  %s", strings.Join(errs, "\n  "))
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

func (i *NamespaceStore) LoadFromFileContent(fileContent []byte, fileType string) error {
	// 准备一个临时结构体来解析文件内容
	fileData := make(map[string]any)

	tagName := "toml"
	// 根据文件类型解码
	switch fileType {
	case ".toml":
		if _, err := toml.Decode(string(fileContent), &fileData); err != nil {
			return fmt.Errorf("toml decode error: %w", err)
		}
	case ".yml", ".yaml":
		tagName = "yaml"
		if err := yaml.Unmarshal(fileContent, &fileData); err != nil {
			return fmt.Errorf("yaml decode error: %w", err)
		}
	case ".json":
		tagName = "json"
		if err := json.Unmarshal(fileContent, &fileData); err != nil {
			return fmt.Errorf("json decode error: %w", err)
		}
	default:
		return fmt.Errorf("unsupported format: %s", fileType)
	}

	// 使用mapstructure直接加载到目标对象
	var errs []error
	i.ForEach(func(w *ObjectWrapper) {
		if configData, exists := fileData[w.Value.Name()]; exists {
			decoderConfig := &mapstructure.DecoderConfig{
				Result:           w.Value,
				TagName:          tagName,
				Squash:           true,              // 嵌套结构体
				WeaklyTypedInput: true,              // 允许弱类型转换
				MatchName:        strings.EqualFold, // 大小写不敏感匹配
			}

			decoder, err := mapstructure.NewDecoder(decoderConfig)
			if err != nil {
				errs = append(errs, fmt.Errorf("create decoder for %s error: %w", w.Value.Name(), err))
				return
			}

			if err := decoder.Decode(configData); err != nil {
				errs = append(errs, fmt.Errorf("decode %s error: %w", w.Value.Name(), err))
			}
		}
	})

	if len(errs) > 0 {
		return fmt.Errorf("load config errors: %v", errors.Join(errs...))
	}

	return nil
}
