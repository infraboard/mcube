package ioc

import (
	"context"
	"testing"
	"time"
)

// TestNoDeadlockInCallbacks 测试在用户回调中调用Get()不会死锁
// 这是修复后的关键测试：确保在Init、Close等回调中可以安全地访问容器
func TestNoDeadlockInCallbacks(t *testing.T) {
	t.Run("Get in Init callback", func(t *testing.T) {
		store := newNamespaceStore("test-init-callback")

		// 先注册一个被依赖的对象
		base := &CallbackTestBase{}
		store.Registry(base)

		// 注册一个在Init中调用Get的对象
		dependent := &CallbackTestDependent{
			store:      store,
			dependency: "callback-test-base",
		}
		store.Registry(dependent)

		done := make(chan bool, 1)
		timeout := 5 * time.Second

		go func() {
			// Init会在内部调用Get()
			err := store.Init()
			if err != nil {
				t.Errorf("Init failed: %v", err)
			}
			done <- true
		}()

		select {
		case <-done:
			t.Log("Init with Get() callback completed successfully")
			if !dependent.initialized {
				t.Error("Dependent object was not initialized")
			}
			if dependent.fetchedDep == nil {
				t.Error("Failed to fetch dependency in Init")
			}
		case <-time.After(timeout):
			t.Fatal("DEADLOCK DETECTED: Init with Get() callback timed out")
		}
	})

	t.Run("Get in PostConfig callback", func(t *testing.T) {
		store := newNamespaceStore("test-postconfig-callback")

		// 先注册一个被依赖的对象
		base := &CallbackTestBase{}
		store.Registry(base)

		// 注册一个在PostConfig中调用Get的对象
		dependent := &CallbackTestPostConfig{
			store:      store,
			dependency: "callback-test-base",
		}
		store.Registry(dependent)

		done := make(chan bool, 1)
		timeout := 5 * time.Second

		go func() {
			err := store.CallPostConfigHooks()
			if err != nil {
				t.Errorf("CallPostConfigHooks failed: %v", err)
			}
			done <- true
		}()

		select {
		case <-done:
			t.Log("PostConfig with Get() callback completed successfully")
			if dependent.fetchedDep == nil {
				t.Error("Failed to fetch dependency in PostConfig")
			}
		case <-time.After(timeout):
			t.Fatal("DEADLOCK DETECTED: PostConfig with Get() callback timed out")
		}
	})

	t.Run("Get in Close callback", func(t *testing.T) {
		store := newNamespaceStore("test-close-callback")

		// 先注册一个被依赖的对象
		base := &CallbackTestBase{}
		store.Registry(base)

		// 注册一个在Close中调用Get的对象
		dependent := &CallbackTestClose{
			store:      store,
			dependency: "callback-test-base",
		}
		store.Registry(dependent)

		done := make(chan bool, 1)
		timeout := 5 * time.Second

		go func() {
			store.Close(context.Background())
			done <- true
		}()

		select {
		case <-done:
			t.Log("Close with Get() callback completed successfully")
			if dependent.fetchedDep == nil {
				t.Error("Failed to fetch dependency in Close")
			}
		case <-time.After(timeout):
			t.Fatal("DEADLOCK DETECTED: Close with Get() callback timed out")
		}
	})
}

// CallbackTestBase 基础测试对象
type CallbackTestBase struct {
	ObjectImpl
}

func (o *CallbackTestBase) Name() string {
	return "callback-test-base"
}

// CallbackTestDependent 在Init中调用Get的测试对象
type CallbackTestDependent struct {
	ObjectImpl
	store       *NamespaceStore
	dependency  string
	initialized bool
	fetchedDep  Object
}

func (o *CallbackTestDependent) Name() string {
	return "callback-test-dependent"
}

func (o *CallbackTestDependent) Init() error {
	// 在Init中调用Get - 这在修复前会导致死锁
	o.fetchedDep = o.store.Get(o.dependency)
	o.initialized = true
	return nil
}

// CallbackTestPostConfig 在PostConfig中调用Get的测试对象
type CallbackTestPostConfig struct {
	ObjectImpl
	store      *NamespaceStore
	dependency string
	fetchedDep Object
}

func (o *CallbackTestPostConfig) Name() string {
	return "callback-test-postconfig"
}

func (o *CallbackTestPostConfig) OnPostConfig() error {
	// 在PostConfig中调用Get - 这在修复前会导致死锁
	o.fetchedDep = o.store.Get(o.dependency)
	return nil
}

// CallbackTestClose 在Close中调用Get的测试对象
type CallbackTestClose struct {
	ObjectImpl
	store      *NamespaceStore
	dependency string
	fetchedDep Object
}

func (o *CallbackTestClose) Name() string {
	return "callback-test-close"
}

func (o *CallbackTestClose) Close(ctx context.Context) {
	// 在Close中调用Get - 这在修复前会导致死锁
	o.fetchedDep = o.store.Get(o.dependency)
}
