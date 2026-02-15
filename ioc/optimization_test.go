package ioc_test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/infraboard/mcube/v2/ioc"
)

// TestConcurrentRegistry 测试并发注册的安全性
func TestConcurrentRegistry(t *testing.T) {
	ns := ioc.Default()

	// 并发注册100个对象
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			obj := &ConcurrentTestObject{ID: idx}
			ns.Registry(obj)
		}(i)
	}
	wg.Wait()

	// 验证所有对象都已注册
	count := 0
	ns.ForEach(func(w *ioc.ObjectWrapper) {
		if _, ok := w.Value.(*ConcurrentTestObject); ok {
			count++
		}
	})

	if count != 100 {
		t.Errorf("Expected 100 objects, got %d", count)
	}
}

type ConcurrentTestObject struct {
	ioc.ObjectImpl
	ID int
}

func (t *ConcurrentTestObject) Name() string {
	return "concurrent_test_" + string(rune(t.ID))
}

// TestGenericGet 测试泛型Get API
func TestGenericGet(t *testing.T) {
	// 注册测试对象
	testObj := &GenericTestObject{Value: "test123"}
	ioc.Default().Registry(testObj)

	// 使用泛型Get获取
	obj, err := ioc.Get[*GenericTestObject](ioc.Default())
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if obj.Value != "test123" {
		t.Errorf("Expected value 'test123', got '%s'", obj.Value)
	}
}

// TestMustGet 测试MustGet API
func TestMustGet(t *testing.T) {
	// 注册测试对象
	testObj := &MustGetTestObject{Data: 42}
	ioc.Default().Registry(testObj)

	// 使用MustGet获取
	obj := ioc.MustGet[*MustGetTestObject](ioc.Default())
	if obj.Data != 42 {
		t.Errorf("Expected data 42, got %d", obj.Data)
	}
}

// TestMustGetPanic 测试MustGet在对象不存在时会panic
func TestMustGetPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic but didn't get one")
		}
	}()

	// 尝试获取不存在的对象
	_ = ioc.MustGet[*NonExistentObject](ioc.Default())
}

type GenericTestObject struct {
	ioc.ObjectImpl
	Value string
}

func (t *GenericTestObject) Name() string {
	return "*ioc_test.GenericTestObject"
}

type MustGetTestObject struct {
	ioc.ObjectImpl
	Data int
}

func (t *MustGetTestObject) Name() string {
	return "*ioc_test.MustGetTestObject"
}

type NonExistentObject struct {
	ioc.ObjectImpl
}

// TestMultipleConfigFiles 测试多配置文件加载
func TestMultipleConfigFiles(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()

	// 创建base配置文件
	baseConfig := filepath.Join(tmpDir, "base.toml")
	baseContent := `
[multi_config_test]
value1 = "base"
value2 = "base"
`
	if err := os.WriteFile(baseConfig, []byte(baseContent), 0644); err != nil {
		t.Fatal(err)
	}

	// 创建override配置文件
	overrideConfig := filepath.Join(tmpDir, "override.toml")
	overrideContent := `
[multi_config_test]
value2 = "override"
value3 = "override"
`
	if err := os.WriteFile(overrideConfig, []byte(overrideContent), 0644); err != nil {
		t.Fatal(err)
	}

	// 注册测试对象
	testObj := &MultiConfigTestObject{}
	ioc.Config().Registry(testObj)

	// 使用多配置文件加载，强制重新加载
	err := ioc.LoadConfig().
		FromFiles(baseConfig, overrideConfig).
		ForceReload().
		Load()

	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// 验证配置合并结果
	if testObj.Value1 != "base" {
		t.Errorf("Expected value1 'base', got '%s'", testObj.Value1)
	}
	if testObj.Value2 != "override" {
		t.Errorf("Expected value2 'override', got '%s'", testObj.Value2)
	}
	if testObj.Value3 != "override" {
		t.Errorf("Expected value3 'override', got '%s'", testObj.Value3)
	}
}

type MultiConfigTestObject struct {
	ioc.ObjectImpl
	Value1 string `toml:"value1"`
	Value2 string `toml:"value2"`
	Value3 string `toml:"value3"`
}

func (t *MultiConfigTestObject) Name() string {
	return "multi_config_test"
}

// TestBuilderAPI 测试Builder风格的API
func TestBuilderAPI(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.toml")

	content := `
[builder_test]
name = "builder"
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// 注册测试对象
	testObj := &BuilderTestObject{}
	ioc.Config().Registry(testObj)

	// 使用Builder API，强制重新加载
	err := ioc.LoadConfig().
		FromFile(configPath).
		SkipIfNotExist().
		ForceReload().
		Load()

	if err != nil {
		t.Fatalf("Builder Load failed: %v", err)
	}

	if testObj.ConfigName != "builder" {
		t.Errorf("Expected name 'builder', got '%s'", testObj.ConfigName)
	}
}

type BuilderTestObject struct {
	ioc.ObjectImpl
	ConfigName string `toml:"name"`
}

func (t *BuilderTestObject) Name() string {
	return "builder_test"
}

func (t *BuilderTestObject) Init() error {
	return nil
}

// TestChainedRegistry 测试链式注册
func TestChainedRegistry(t *testing.T) {
	// 链式注册多个对象
	ioc.Default().
		Registry(&ChainedTestObject1{}).
		Registry(&ChainedTestObject2{}).
		Registry(&ChainedTestObject3{})

	// 验证所有对象都已注册
	obj1 := ioc.Default().Get("*ioc_test.ChainedTestObject1")
	obj2 := ioc.Default().Get("*ioc_test.ChainedTestObject2")
	obj3 := ioc.Default().Get("*ioc_test.ChainedTestObject3")

	if obj1 == nil || obj2 == nil || obj3 == nil {
		t.Error("Chained registry failed")
	}
}

// TestBatchRegistry 测试批量注册
func TestBatchRegistry(t *testing.T) {
	// 批量注册
	ioc.Default().RegistryAll(
		&BatchTestObject1{},
		&BatchTestObject2{},
		&BatchTestObject3{},
	)

	// 验证所有对象都已注册
	obj1 := ioc.Default().Get("*ioc_test.BatchTestObject1")
	obj2 := ioc.Default().Get("*ioc_test.BatchTestObject2")
	obj3 := ioc.Default().Get("*ioc_test.BatchTestObject3")

	if obj1 == nil || obj2 == nil || obj3 == nil {
		t.Error("Batch registry failed")
	}
}

type ChainedTestObject1 struct {
	ioc.ObjectImpl
}

func (t *ChainedTestObject1) Name() string {
	return "*ioc_test.ChainedTestObject1"
}

type ChainedTestObject2 struct {
	ioc.ObjectImpl
}

func (t *ChainedTestObject2) Name() string {
	return "*ioc_test.ChainedTestObject2"
}

type ChainedTestObject3 struct {
	ioc.ObjectImpl
}

func (t *ChainedTestObject3) Name() string {
	return "*ioc_test.ChainedTestObject3"
}

type BatchTestObject1 struct {
	ioc.ObjectImpl
}

func (t *BatchTestObject1) Name() string {
	return "*ioc_test.BatchTestObject1"
}

type BatchTestObject2 struct {
	ioc.ObjectImpl
}

func (t *BatchTestObject2) Name() string {
	return "*ioc_test.BatchTestObject2"
}

type BatchTestObject3 struct {
	ioc.ObjectImpl
}

func (t *BatchTestObject3) Name() string {
	return "*ioc_test.BatchTestObject3"
}
