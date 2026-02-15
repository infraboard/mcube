package ioc_test

import (
	"fmt"
	"testing"

	"github.com/infraboard/mcube/v2/ioc"
)

type Database struct {
	ioc.ObjectImpl
}

func (d *Database) Name() string    { return "*ioc_test.Database" }
func (d *Database) Version() string { return "2.0.0" }

type Logger struct {
	ioc.ObjectImpl
}

func (l *Logger) Name() string    { return "*ioc_test.Logger" }
func (l *Logger) Version() string { return "1.5.0" }

type Cache struct {
	ioc.ObjectImpl
}

func (c *Cache) Name() string    { return "*ioc_test.Cache" }
func (c *Cache) Version() string { return "1.2.0" }

type UserRepository struct {
	ioc.ObjectImpl
	DB *Database `ioc:"autowire=true;namespace=configs"`
}

func (r *UserRepository) Name() string    { return "*ioc_test.UserRepository" }
func (r *UserRepository) Version() string { return "3.0.0" }

type UserService struct {
	ioc.ObjectImpl
	Repo   *UserRepository `ioc:"autowire=true;namespace=default_dep"`
	Logger *Logger         `ioc:"autowire=true;namespace=default_dep"`
	Cache  *Cache          `ioc:"autowire=true;namespace=default_dep"`
}

func (s *UserService) Name() string    { return "*ioc_test.UserService" }
func (s *UserService) Version() string { return "4.1.0" }

// TestPrintDependencies 测试依赖关系打印
func TestPrintDependencies(t *testing.T) {
	// 注册对象到不同命名空间
	configs := ioc.DefaultStore.Namespace("configs")
	if configs.Get("*ioc_test.Database") == nil {
		configs.Registry(&Database{})
	}

	defaultNs := ioc.DefaultStore.Namespace("default_dep")
	if defaultNs.Get("*ioc_test.Logger") == nil {
		defaultNs.Registry(&Logger{})
	}
	if defaultNs.Get("*ioc_test.Cache") == nil {
		defaultNs.Registry(&Cache{})
	}
	if defaultNs.Get("*ioc_test.UserRepository") == nil {
		defaultNs.Registry(&UserRepository{})
	}
	if defaultNs.Get("*ioc_test.UserService") == nil {
		defaultNs.Registry(&UserService{})
	}

	// 打印依赖树
	fmt.Println("\n=== Testing Dependency Tree Output ===")
	defaultNs.PrintDependencies()

	// 打印摘要
	defaultNs.PrintDependencySummary()

	// 打印所有命名空间
	ioc.DefaultStore.PrintAllDependencies()
}

type SimpleService struct {
	ioc.ObjectImpl
}

func (s *SimpleService) Name() string    { return "*ioc_test.SimpleService" }
func (s *SimpleService) Version() string { return "1.0.0" }

// TestPrintDependenciesSimple 简单依赖测试
func TestPrintDependenciesSimple(t *testing.T) {
	ns := ioc.DefaultStore.Namespace("simple_test")
	if ns.Get("*ioc_test.SimpleService") == nil {
		ns.Registry(&SimpleService{})
	}

	fmt.Println("\n=== Simple Service (No Dependencies) ===")
	ns.PrintDependencies()
}

type MdDatabase struct {
	ioc.ObjectImpl
}

func (d *MdDatabase) Name() string    { return "*ioc_test.MdDatabase" }
func (d *MdDatabase) Version() string { return "1.0.0" }

type MdService struct {
	ioc.ObjectImpl
	DB *MdDatabase `ioc:"autowire=true;namespace=markdown_test"`
}

func (s *MdService) Name() string    { return "*ioc_test.MdService" }
func (s *MdService) Version() string { return "2.0.0" }

// TestExportMarkdown 测试导出Markdown
func TestExportMarkdown(t *testing.T) {
	ns := ioc.DefaultStore.Namespace("markdown_test")
	if ns.Get("*ioc_test.MdDatabase") == nil {
		ns.Registry(&MdDatabase{})
	}
	if ns.Get("*ioc_test.MdService") == nil {
		ns.Registry(&MdService{})
	}

	// 导出Markdown
	markdown := ns.ExportDependenciesToMarkdown()
	fmt.Println("\n=== Exported Markdown ===")
	fmt.Println(markdown)
}

// === 命令式依赖测试（手动Get()方式） ===

type EmailLogger struct {
	ioc.ObjectImpl
}

func (l *EmailLogger) Name() string    { return "*ioc_test.EmailLogger" }
func (l *EmailLogger) Version() string { return "1.0.0" }

type SMTPClient struct {
	ioc.ObjectImpl
}

func (s *SMTPClient) Name() string    { return "*ioc_test.SMTPClient" }
func (s *SMTPClient) Version() string { return "2.5.0" }

// EmailService 使用命令式依赖（手动Get()）
type EmailService struct {
	ioc.ObjectImpl
	// 注意：这里没有 ioc 标签，依赖是在 Init() 中手动获取的
	logger     *EmailLogger
	smtpClient *SMTPClient
}

func (s *EmailService) Name() string    { return "*ioc_test.EmailService" }
func (s *EmailService) Version() string { return "3.0.0" }

// Init 手动获取依赖
func (s *EmailService) Init() error {
	// 命令式依赖：手动通过 Get() 获取
	ns := ioc.DefaultStore.Namespace("imperative_test")
	s.logger = ns.Get("*ioc_test.EmailLogger").(*EmailLogger)
	s.smtpClient = ns.Get("*ioc_test.SMTPClient").(*SMTPClient)
	return nil
}

// DeclareDependencies 实现 DependencyDeclarer 接口来声明命令式依赖
// 这样依赖图就能展示 EmailService 依赖 EmailLogger 和 SMTPClient
func (s *EmailService) DeclareDependencies() []ioc.DependencyInfo {
	return []ioc.DependencyInfo{
		{
			Name:      "*ioc_test.EmailLogger",
			Namespace: "imperative_test",
			FieldName: "logger", // 可选：用于文档说明
		},
		{
			Name:      "*ioc_test.SMTPClient",
			Namespace: "imperative_test",
			FieldName: "smtpClient",
		},
	}
}

// TestImperativeDependencies 测试命令式依赖（手动Get()）的可视化
func TestImperativeDependencies(t *testing.T) {
	ns := ioc.DefaultStore.Namespace("imperative_test")

	// 清理之前的注册
	if ns.Get("*ioc_test.EmailLogger") == nil {
		ns.Registry(&EmailLogger{})
	}
	if ns.Get("*ioc_test.SMTPClient") == nil {
		ns.Registry(&SMTPClient{})
	}
	if ns.Get("*ioc_test.EmailService") == nil {
		ns.Registry(&EmailService{})
	}

	// 初始化（触发 Init() 进行手动依赖注入）
	emailService := ns.Get("*ioc_test.EmailService").(*EmailService)
	if err := emailService.Init(); err != nil {
		t.Fatalf("Failed to init EmailService: %v", err)
	}

	// 打印依赖树 - 应该能看到 EmailService 依赖 EmailLogger 和 SMTPClient
	fmt.Println("\n=== Imperative Dependencies (Manual Get()) ===")
	ns.PrintDependencies()

	fmt.Println("\n说明：")
	fmt.Println("  - EmailService 没有使用 ioc 标签")
	fmt.Println("  - 依赖在 Init() 中通过 Get() 手动获取")
	fmt.Println("  - 但实现了 DependencyDeclarer 接口")
	fmt.Println("  - 因此依赖图仍能正确展示依赖关系")
}

// === 混合依赖测试（同时使用标签和接口） ===

type NotificationQueue struct {
	ioc.ObjectImpl
}

func (q *NotificationQueue) Name() string    { return "*ioc_test.NotificationQueue" }
func (q *NotificationQueue) Version() string { return "1.8.0" }

type NotificationService struct {
	ioc.ObjectImpl
	// 声明式依赖：使用 ioc 标签自动注入
	Logger *EmailLogger `ioc:"autowire=true;namespace=mixed_test"`

	// 命令式依赖：手动 Get() 获取（无标签）
	queue *NotificationQueue
}

func (n *NotificationService) Name() string    { return "*ioc_test.NotificationService" }
func (n *NotificationService) Version() string { return "5.0.0" }

func (n *NotificationService) Init() error {
	ns := ioc.DefaultStore.Namespace("mixed_test")
	n.queue = ns.Get("*ioc_test.NotificationQueue").(*NotificationQueue)
	return nil
}

// DeclareDependencies 只需要声明命令式依赖，标签依赖会自动检测
func (n *NotificationService) DeclareDependencies() []ioc.DependencyInfo {
	return []ioc.DependencyInfo{
		{
			Name:      "*ioc_test.NotificationQueue",
			Namespace: "mixed_test",
			FieldName: "queue",
		},
	}
}

// TestMixedDependencies 测试混合依赖（标签 + 接口）
func TestMixedDependencies(t *testing.T) {
	ns := ioc.DefaultStore.Namespace("mixed_test")

	if ns.Get("*ioc_test.EmailLogger") == nil {
		ns.Registry(&EmailLogger{})
	}
	if ns.Get("*ioc_test.NotificationQueue") == nil {
		ns.Registry(&NotificationQueue{})
	}
	if ns.Get("*ioc_test.NotificationService") == nil {
		ns.Registry(&NotificationService{})
	}

	// 初始化
	notifService := ns.Get("*ioc_test.NotificationService").(*NotificationService)
	if err := notifService.Init(); err != nil {
		t.Fatalf("Failed to init NotificationService: %v", err)
	}

	// 打印依赖树 - 应该能看到两种依赖
	fmt.Println("\n=== Mixed Dependencies (Tag + Interface) ===")
	ns.PrintDependencies()

	fmt.Println("\n说明：")
	fmt.Println("  - NotificationService.Logger: 使用 ioc 标签（自动检测）")
	fmt.Println("  - NotificationService.queue: 手动 Get()（通过接口声明）")
	fmt.Println("  - 两种依赖方式可以共存，都会在依赖图中显示")
}

// === 接口依赖测试（字段类型是接口 + 指定具体实现名称） ===

// Storage 存储接口
type Storage interface {
	Save(key, value string) error
	Load(key string) (string, error)
}

// RedisStorage Redis 实现
type RedisStorage struct {
	ioc.ObjectImpl
	data map[string]string
}

func (r *RedisStorage) Name() string    { return "*ioc_test.RedisStorage" }
func (r *RedisStorage) Version() string { return "2.0.0" }

func (r *RedisStorage) Save(key, value string) error {
	if r.data == nil {
		r.data = make(map[string]string)
	}
	r.data[key] = value
	return nil
}

func (r *RedisStorage) Load(key string) (string, error) {
	return r.data[key], nil
}

// MemoryStorage 内存实现
type MemoryStorage struct {
	ioc.ObjectImpl
	cache map[string]string
}

func (m *MemoryStorage) Name() string    { return "*ioc_test.MemoryStorage" }
func (m *MemoryStorage) Version() string { return "1.5.0" }

func (m *MemoryStorage) Save(key, value string) error {
	if m.cache == nil {
		m.cache = make(map[string]string)
	}
	m.cache[key] = value
	return nil
}

func (m *MemoryStorage) Load(key string) (string, error) {
	return m.cache[key], nil
}

// CacheService 使用接口类型字段，通过 name 指定具体实现
type CacheService struct {
	ioc.ObjectImpl
	// 字段类型是接口，通过 name 指定使用 RedisStorage 实现
	Storage Storage `ioc:"autowire=true;namespace=interface_test;name=*ioc_test.RedisStorage"`
}

func (c *CacheService) Name() string    { return "*ioc_test.CacheService" }
func (c *CacheService) Version() string { return "3.0.0" }

// TestInterfaceDependency 测试接口依赖：字段类型是接口 + name 指定具体实现
func TestInterfaceDependency(t *testing.T) {
	ns := ioc.DefaultStore.Namespace("interface_test")

	// 注册两个 Storage 接口的实现
	if ns.Get("*ioc_test.RedisStorage") == nil {
		ns.Registry(&RedisStorage{})
	}
	if ns.Get("*ioc_test.MemoryStorage") == nil {
		ns.Registry(&MemoryStorage{})
	}
	if ns.Get("*ioc_test.CacheService") == nil {
		ns.Registry(&CacheService{})
	}

	// 执行自动注入
	if err := ns.Autowire(); err != nil {
		t.Fatalf("Autowire failed: %v", err)
	}

	// 验证注入是否正确
	cacheService := ns.Get("*ioc_test.CacheService").(*CacheService)
	if cacheService.Storage == nil {
		t.Fatal("Storage should be injected")
	}

	// 验证注入的是 RedisStorage（通过 name 指定）
	if _, ok := cacheService.Storage.(*RedisStorage); !ok {
		t.Fatalf("Expected RedisStorage, got %T", cacheService.Storage)
	}

	// 功能测试
	if err := cacheService.Storage.Save("key1", "value1"); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	value, err := cacheService.Storage.Load("key1")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if value != "value1" {
		t.Fatalf("Expected value1, got %s", value)
	}

	fmt.Println("\n=== Interface Dependency (Field Type: Interface + Name Specified) ===")
	fmt.Println("✅ 功能验证通过：")
	fmt.Println("  - 字段类型是接口（Storage interface）")
	fmt.Println("  - 通过 name 指定具体实现（*ioc_test.RedisStorage）")
	fmt.Println("  - 依赖注入成功，功能正常")

	// 打印依赖树 - 验证依赖可视化
	fmt.Println("\n📊 依赖图可视化：")
	ns.PrintDependencies()

	fmt.Println("\n说明：")
	fmt.Println("  - CacheService.Storage 字段类型是 interface")
	fmt.Println("  - 通过 ioc 标签的 name 参数指定使用 RedisStorage")
	fmt.Println("  - 依赖图应该展示 CacheService → RedisStorage 的关系")
}
