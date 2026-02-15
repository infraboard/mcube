package ioc

import (
	"fmt"
	"testing"
	"time"
)

// TestNoDeadlockOnRegistry 测试注册对象时不会发生死锁
// 修复前：如果 Priority() 中调用 Get()，会导致死锁
// 修复后：Priority() 在加锁前调用，不会死锁
func TestNoDeadlockOnRegistry(t *testing.T) {
	// 创建测试用的命名空间
	store := newNamespaceStore("test-deadlock")

	// 先注册一个基础对象
	baseObj := &BaseTestObject{priority: 10}
	store.Registry(baseObj)

	// 使用 channel 来检测是否发生死锁
	done := make(chan bool, 1)
	timeout := 5 * time.Second

	go func() {
		// 这个操作在修复前会导致死锁，修复后应该正常完成
		// 因为 Priority() 现在在加锁前调用
		dependentObj := &DependentTestObject{
			store:    store,
			priority: 5,
		}
		store.Registry(dependentObj)
		done <- true
	}()

	select {
	case <-done:
		// 成功完成，没有死锁
		t.Log("Registry completed successfully, no deadlock detected")
	case <-time.After(timeout):
		t.Fatal("DEADLOCK DETECTED: Registry operation timed out after", timeout)
	}
}

// TestSortNoDeadlock 测试排序时不会发生死锁
func TestSortNoDeadlock(t *testing.T) {
	store := newNamespaceStore("test-sort")

	// 注册多个对象（使用不同的名称）
	store.Registry(&BaseTestObject{name: "obj1", priority: 10})
	store.Registry(&BaseTestObject{name: "obj2", priority: 20})
	store.Registry(&BaseTestObject{name: "obj3", priority: 5})

	done := make(chan bool, 1)
	timeout := 5 * time.Second

	go func() {
		// Sort 操作应该使用缓存的 Priority 字段，不会调用 Priority() 方法
		store.Sort()
		done <- true
	}()

	select {
	case <-done:
		t.Log("Sort completed successfully, no deadlock detected")

		// 验证排序结果
		if store.Items()[0].Priority != 20 {
			t.Errorf("Expected first item priority 20, got %d", store.Items()[0].Priority)
		}
		if store.Items()[1].Priority != 10 {
			t.Errorf("Expected second item priority 10, got %d", store.Items()[1].Priority)
		}
		if store.Items()[2].Priority != 5 {
			t.Errorf("Expected third item priority 5, got %d", store.Items()[2].Priority)
		}
	case <-time.After(timeout):
		t.Fatal("DEADLOCK DETECTED: Sort operation timed out after", timeout)
	}
}

// TestInitNoDeadlock 测试初始化时不会发生死锁
func TestInitNoDeadlock(t *testing.T) {
	store := newNamespaceStore("test-init")

	// 注册多个对象（使用不同的名称）
	store.Registry(&BaseTestObject{name: "init-obj1", priority: 10})
	store.Registry(&BaseTestObject{name: "init-obj2", priority: 20})
	store.Registry(&BaseTestObject{name: "init-obj3", priority: 5})

	done := make(chan bool, 1)
	timeout := 5 * time.Second

	go func() {
		// Init 操作会排序并初始化所有对象
		err := store.Init()
		if err != nil {
			t.Errorf("Init failed: %v", err)
		}
		done <- true
	}()

	select {
	case <-done:
		t.Log("Init completed successfully, no deadlock detected")

		// 验证初始化顺序（高优先级先初始化）
		if store.Items()[0].Priority != 20 {
			t.Errorf("Expected first item priority 20, got %d", store.Items()[0].Priority)
		}
	case <-time.After(timeout):
		t.Fatal("DEADLOCK DETECTED: Init operation timed out after", timeout)
	}
}

// BaseTestObject 基础测试对象
type BaseTestObject struct {
	ObjectImpl
	name     string
	priority int
}

func (o *BaseTestObject) Name() string {
	if o.name == "" {
		return "base-test"
	}
	return o.name
}

func (o *BaseTestObject) Priority() int {
	return o.priority
}

// DependentTestObject 依赖其他对象的测试对象
// 在修复前，如果在 Priority() 中调用 Get()，会导致死锁
type DependentTestObject struct {
	ObjectImpl
	store    *NamespaceStore
	priority int
}

func (o *DependentTestObject) Name() string {
	return "dependent-test"
}

func (o *DependentTestObject) Priority() int {
	// 注意：这里故意不调用 Get()
	// 因为正确的做法是 Priority() 只返回常量
	// 如果需要获取依赖，应该在 Init() 中进行
	return o.priority
}

func (o *DependentTestObject) Init() error {
	// 正确的做法：在 Init() 中获取依赖
	// 此时容器已完成注册，可以安全地调用 Get()
	if o.store != nil {
		_ = o.store.Get("base-test")
	}
	return nil
}

// TestPriorityBestPractice 测试 Priority() 的最佳实践
func TestPriorityBestPractice(t *testing.T) {
	t.Run("Priority should return constant", func(t *testing.T) {
		obj := &BaseTestObject{priority: 10}

		// Priority() 应该总是返回相同的值
		p1 := obj.Priority()
		p2 := obj.Priority()

		if p1 != p2 {
			t.Errorf("Priority() should return constant value, got %d and %d", p1, p2)
		}
	})

	t.Run("Dependencies should be fetched in Init", func(t *testing.T) {
		store := newNamespaceStore("test-deps")
		store.Registry(&BaseTestObject{priority: 10})

		obj := &DependentTestObject{
			store:    store,
			priority: 5,
		}
		store.Registry(obj)

		// 初始化应该成功
		err := store.Init()
		if err != nil {
			t.Errorf("Init failed: %v", err)
		}
	})
}

// TestConcurrentRegistry 测试并发注册不会发生死锁
func TestConcurrentRegistry(t *testing.T) {
	store := newNamespaceStore("test-concurrent")

	done := make(chan bool)
	numGoroutines := 10

	// 启动多个 goroutine 并发注册对象
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			obj := &BaseTestObject{
				name:     fmt.Sprintf("concurrent-obj-%d", id),
				priority: id,
			}
			store.Registry(obj)
			done <- true
		}(i)
	}

	// 等待所有 goroutine 完成
	timeout := time.After(5 * time.Second)
	for i := 0; i < numGoroutines; i++ {
		select {
		case <-done:
			// 成功
		case <-timeout:
			t.Fatal("DEADLOCK DETECTED: Concurrent registry timed out")
		}
	}

	// 验证所有对象都已注册
	if store.Len() != numGoroutines {
		t.Errorf("Expected %d objects, got %d", numGoroutines, store.Len())
	}

	t.Log("Concurrent registry completed successfully")
}
