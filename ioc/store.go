package ioc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

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
		if !req.ConfigFile.SkipIFNotExist && !IsFileExists(req.ConfigFile.Path) {
			return fmt.Errorf("file %s not exist", req.ConfigFile.Path)
		}

		fileType := filepath.Ext(req.ConfigFile.Path)
		if err := ValidateFileType(fileType); err != nil {
			return err
		}

		// 读取文件内容
		content, err := os.ReadFile(req.ConfigFile.Path)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		for i := range s.store {
			item := s.store[i]
			err := item.LoadFromFileContent(content, fileType)

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
		log.Printf("init app %s[priority: %d] ok.", obj.Value.Name(), obj.Value.Priority())
	}
	return nil
}

// 倒序关闭
func (s *NamespaceStore) Close(ctx context.Context) {
	for i := len(s.Items) - 1; i >= 0; i-- {
		obj := s.Items[i]
		obj.Value.Close(ctx)
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
					objs := DefaultStore.Namespace(tag.Namespace).ImplementInterface(fieldType)
					if len(objs) > 0 {
						obj = objs[0]
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
