package ioc

import (
	"testing"
)

// TestRelativePriority 测试相对优先级的使用
// 演示如何基于其他对象动态确定优先级
func TestRelativePriority(t *testing.T) {
	store := newNamespaceStore("test-relative-priority")

	// 1. 注册基础服务（固定优先级）
	database := &DatabaseTestService{}
	store.Registry(database)

	// 2. 注册缓存服务（相对于数据库的优先级）
	cache := &CacheTestService{
		store:          store,
		databaseName:   "database-test",
		priorityOffset: -10, // 比数据库低10，即在数据库之后初始化
	}
	store.Registry(cache)

	// 3. 注册业务服务（相对于缓存的优先级）
	service := &BusinessTestService{
		store:          store,
		cacheName:      "cache-test",
		priorityOffset: -5, // 比缓存低5
	}
	store.Registry(service)

	// 验证优先级被正确计算
	if len(store.Items()) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(store.Items()))
	}

	// 查找各个对象
	var dbObj, cacheObj, svcObj *ObjectWrapper
	for _, item := range store.Items() {
		switch item.Name {
		case "database-test":
			dbObj = item
		case "cache-test":
			cacheObj = item
		case "business-test":
			svcObj = item
		}
	}

	if dbObj == nil || cacheObj == nil || svcObj == nil {
		t.Fatal("Not all objects were registered")
	}

	// 验证相对优先级
	t.Logf("Database priority: %d", dbObj.Priority)
	t.Logf("Cache priority: %d", cacheObj.Priority)
	t.Logf("Service priority: %d", svcObj.Priority)

	expectedCachePriority := dbObj.Priority - 10
	if cacheObj.Priority != expectedCachePriority {
		t.Errorf("Cache priority expected %d, got %d", expectedCachePriority, cacheObj.Priority)
	}

	expectedServicePriority := cacheObj.Priority - 5
	if svcObj.Priority != expectedServicePriority {
		t.Errorf("Service priority expected %d, got %d", expectedServicePriority, svcObj.Priority)
	}

	// 排序后验证顺序
	store.Sort()

	// 应该按照 Database -> Cache -> Service 的顺序
	if store.Items()[0].Name != "database-test" {
		t.Errorf("Expected database first, got %s", store.Items()[0].Name)
	}
	if store.Items()[1].Name != "cache-test" {
		t.Errorf("Expected cache second, got %s", store.Items()[1].Name)
	}
	if store.Items()[2].Name != "business-test" {
		t.Errorf("Expected service third, got %s", store.Items()[2].Name)
	}

	t.Log("✅ Relative priority works correctly!")
}

// TestConditionalPriority 测试条件优先级
func TestConditionalPriority(t *testing.T) {
	store := newNamespaceStore("test-conditional-priority")

	// 注册一个根据条件确定优先级的服务
	service := &ConditionalPriorityService{
		isHighPriority: true,
	}
	store.Registry(service)

	// 验证优先级
	if len(store.Items()) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(store.Items()))
	}

	if store.Items()[0].Priority != 100 {
		t.Errorf("Expected priority 100, got %d", store.Items()[0].Priority)
	}

	// 注册低优先级版本
	store2 := newNamespaceStore("test-conditional-priority-2")
	service2 := &ConditionalPriorityService{
		isHighPriority: false,
	}
	store2.Registry(service2)

	if store2.Items()[0].Priority != 10 {
		t.Errorf("Expected priority 10, got %d", store2.Items()[0].Priority)
	}

	t.Log("✅ Conditional priority works correctly!")
}

// DatabaseTestService 数据库测试服务（固定优先级）
type DatabaseTestService struct {
	ObjectImpl
}

func (s *DatabaseTestService) Name() string {
	return "database-test"
}

func (s *DatabaseTestService) Priority() int {
	return 100 // 固定的高优先级
}

// CacheTestService 缓存测试服务（相对优先级）
type CacheTestService struct {
	ObjectImpl
	store          *NamespaceStore
	databaseName   string
	priorityOffset int
}

func (s *CacheTestService) Name() string {
	return "cache-test"
}

func (s *CacheTestService) Priority() int {
	// 相对于数据库服务的优先级
	db := s.store.Get(s.databaseName)
	if db != nil {
		return db.Priority() + s.priorityOffset
	}
	return 0
}

// BusinessTestService 业务测试服务（相对优先级）
type BusinessTestService struct {
	ObjectImpl
	store          *NamespaceStore
	cacheName      string
	priorityOffset int
}

func (s *BusinessTestService) Name() string {
	return "business-test"
}

func (s *BusinessTestService) Priority() int {
	// 相对于缓存服务的优先级
	cache := s.store.Get(s.cacheName)
	if cache != nil {
		return cache.Priority() + s.priorityOffset
	}
	return 0
}

// ConditionalPriorityService 条件优先级测试服务
type ConditionalPriorityService struct {
	ObjectImpl
	isHighPriority bool
}

func (s *ConditionalPriorityService) Name() string {
	return "conditional-priority-test"
}

func (s *ConditionalPriorityService) Priority() int {
	// 根据条件返回不同的优先级
	if s.isHighPriority {
		return 100
	}
	return 10
}
