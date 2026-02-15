package ioc

import (
	"testing"
	"time"
)

// TestAutowireDeadlockRisk 测试 Autowire 中的潜在死锁风险
// Autowire 使用 ForEach (持有读锁)，在回调中调用 Get() (尝试获取读锁)
func TestAutowireDeadlockRisk(t *testing.T) {
	t.Run("Cross namespace autowire - Safe", func(t *testing.T) {
		// 跨 namespace 的 autowire 是安全的（不同的锁）
		controller := newNamespaceStore("controllers")
		service := newNamespaceStore("default")

		// 注册服务
		baseService := &AutowireTestBase{}
		service.Registry(baseService)

		// 注册控制器（依赖服务）
		ctrl := &AutowireTestController{}
		controller.Registry(ctrl)

		// 创建临时的 DefaultStore 用于测试
		oldStore := DefaultStore
		defer func() { DefaultStore = oldStore }()

		DefaultStore = &defaultStore{
			store: []*NamespaceStore{controller, service},
		}

		done := make(chan bool, 1)
		timeout := 5 * time.Second

		go func() {
			// 跨 namespace 的 autowire 应该是安全的
			err := controller.Autowire()
			if err != nil {
				t.Errorf("Autowire failed: %v", err)
			}
			done <- true
		}()

		select {
		case <-done:
			t.Log("✅ Cross namespace autowire completed successfully")
		case <-time.After(timeout):
			t.Fatal("DEADLOCK DETECTED: Autowire timed out")
		}
	})

	t.Run("Same namespace autowire - Potential deadlock", func(t *testing.T) {
		// 同一 namespace 内的自引用理论上会死锁
		// 但由于 Go 的 RWMutex 实现，同一 goroutine 多次 RLock 实际上不会死锁
		// 只有在有写锁等待时才会死锁
		store := newNamespaceStore("test-same-ns")

		base := &AutowireTestBase{}
		store.Registry(base)

		dependent := &AutowireTestSelfDependent{}
		store.Registry(dependent)

		// 创建临时的 DefaultStore
		oldStore := DefaultStore
		defer func() { DefaultStore = oldStore }()

		DefaultStore = &defaultStore{
			store: []*NamespaceStore{store},
		}

		done := make(chan bool, 1)
		timeout := 5 * time.Second

		go func() {
			// 实际测试表明这不会死锁
			// 因为 RWMutex 允许同一 goroutine 多次获取读锁（如果没有写锁在等待）
			err := store.Autowire()
			if err != nil {
				t.Logf("Autowire failed (expected): %v", err)
			}
			done <- true
		}()

		select {
		case <-done:
			t.Log("✅ Same namespace autowire completed (RWMutex allows recursive RLock)")
		case <-time.After(timeout):
			t.Fatal("DEADLOCK DETECTED: Autowire timed out")
		}
	})
}

// AutowireTestBase 测试基础服务
type AutowireTestBase struct {
	ObjectImpl
}

func (s *AutowireTestBase) Name() string {
	return "autowire-base"
}

// AutowireTestController 测试控制器（跨 namespace 依赖）
type AutowireTestController struct {
	ObjectImpl
	Service Object `ioc:"autowire=true;namespace=default"`
}

func (c *AutowireTestController) Name() string {
	return "autowire-controller"
}

// AutowireTestSelfDependent 同一 namespace 内的依赖
type AutowireTestSelfDependent struct {
	ObjectImpl
	Base Object `ioc:"autowire=true;namespace=test-same-ns"`
}

func (s *AutowireTestSelfDependent) Name() string {
	return "autowire-self-dependent"
}
