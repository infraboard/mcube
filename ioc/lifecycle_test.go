package ioc_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/infraboard/mcube/v2/ioc"
)

// TestLifecycleHooks 测试完整的生命周期钩子流程
func TestLifecycleHooks(t *testing.T) {
	// 创建一个跟踪对象来记录钩子调用顺序
	tracker := &LifecycleTracker{
		callOrder: []string{},
	}

	// 注册测试对象
	testObj := &MockObject{
		name:    "lifecycle_full_test",
		tracker: tracker,
	}
	ioc.Default().Registry(testObj)

	// 初始化
	err := ioc.LoadConfig().ForceReload().Load()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// 验证初始化阶段的钩子调用顺序
	expectedOrder := []string{
		"OnPostConfig",
		"OnPreInit",
		"Init",
		"OnPostInit",
	}

	if len(tracker.callOrder) < len(expectedOrder) {
		t.Fatalf("Expected at least %d calls, got %d: %v",
			len(expectedOrder), len(tracker.callOrder), tracker.callOrder)
	}

	for i, expected := range expectedOrder {
		if tracker.callOrder[i] != expected {
			t.Errorf("Step %d: expected %s, got %s", i, expected, tracker.callOrder[i])
		}
	}

	// 停止服务
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ioc.DefaultStore.Stop(ctx)

	// 验证停止阶段的钩子调用顺序
	expectedStopOrder := []string{
		"OnPreStop",
		"Close",
		"OnPostStop",
	}

	totalExpected := len(expectedOrder) + len(expectedStopOrder)
	if len(tracker.callOrder) != totalExpected {
		t.Fatalf("Expected %d total calls, got %d: %v",
			totalExpected, len(tracker.callOrder), tracker.callOrder)
	}

	for i, expected := range expectedStopOrder {
		idx := len(expectedOrder) + i
		if tracker.callOrder[idx] != expected {
			t.Errorf("Stop step %d: expected %s, got %s", i, expected, tracker.callOrder[idx])
		}
	}

	t.Logf("Full lifecycle order: %v", tracker.callOrder)
}

// TestPostConfigHook_ValidationError 测试配置验证失败
func TestPostConfigHook_ValidationError(t *testing.T) {
	callCount := 0
	testObj := &MockObject{
		name: "validation_controlled_test",
		OnPostConfigFunc: func() error {
			callCount++
			// 只在第一次调用时返回错误
			if callCount == 1 {
				return fmt.Errorf("validation failed: invalid configuration")
			}
			return nil
		},
	}
	ioc.Default().Registry(testObj)

	err := ioc.LoadConfig().ForceReload().Load()
	if err == nil {
		t.Fatal("Expected validation error, got nil")
	}

	t.Logf("Got expected error: %v", err)
}

// TestPostConfigHook_ValidationSuccess 测试配置验证成功
func TestPostConfigHook_ValidationSuccess(t *testing.T) {
	testObj := &MockObject{
		name: "validation_success_unique_test",
		OnPostConfigFunc: func() error {
			return nil // 验证通过
		},
	}
	ioc.Default().Registry(testObj)

	err := ioc.LoadConfig().ForceReload().Load()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

// TestPreInitHook_Error 测试 PreInit 错误中断初始化
func TestPreInitHook_Error(t *testing.T) {
	preInitCalled := false
	initCalled := false
	preInitCallCount := 0

	testObj := &MockObject{
		name: "preinit_controlled_test",
		OnPreInitFunc: func() error {
			preInitCallCount++
			// 只在第一次调用时返回错误
			if preInitCallCount == 1 {
				preInitCalled = true
				return fmt.Errorf("PreInit failed")
			}
			return nil
		},
	}
	ioc.Default().Registry(testObj)

	err := ioc.LoadConfig().ForceReload().Load()
	if err == nil {
		t.Fatal("Expected PreInit error, got nil")
	}

	if !preInitCalled {
		t.Error("PreInit should have been called")
	}

	if initCalled {
		t.Error("Init should not be called after PreInit error")
	}

	t.Logf("Got expected error: %v", err)
}

// TestPostInitHook_NonBlocking 测试 PostInit 错误不阻止流程
func TestPostInitHook_NonBlocking(t *testing.T) {
	postInitCalled := false

	testObj := &MockObject{
		name: "postinit_error_unique_test",
		OnPostInitFunc: func() error {
			postInitCalled = true
			return fmt.Errorf("PostInit failed (should not block)")
		},
	}
	ioc.Default().Registry(testObj)

	err := ioc.LoadConfig().ForceReload().Load()
	if err != nil {
		t.Fatalf("PostInit error should not block: %v", err)
	}

	if !postInitCalled {
		t.Error("PostInit should have been called")
	}
}

// TestPreStopHook_GracefulShutdown 测试优雅停机
func TestPreStopHook_GracefulShutdown(t *testing.T) {
	var mu sync.Mutex
	activeRequests := 3
	preStopCalled := false

	testObj := &MockObject{
		name: "graceful_shutdown_unique_test",
		OnPreStopFunc: func(ctx context.Context) error {
			preStopCalled = true
			// 等待所有请求完成
			ticker := time.NewTicker(10 * time.Millisecond)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return fmt.Errorf("timeout waiting for requests to complete")
				case <-ticker.C:
					mu.Lock()
					active := activeRequests
					mu.Unlock()

					if active == 0 {
						return nil
					}
				}
			}
		},
	}
	ioc.Default().Registry(testObj)

	err := ioc.LoadConfig().ForceReload().Load()
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// 启动异步请求处理
	go func() {
		time.Sleep(100 * time.Millisecond)
		mu.Lock()
		activeRequests--
		activeRequests--
		activeRequests--
		mu.Unlock()
	}()

	// 停止服务
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ioc.DefaultStore.Stop(ctx)

	if !preStopCalled {
		t.Error("PreStop should have been called")
	}

	mu.Lock()
	finalActive := activeRequests
	mu.Unlock()

	if finalActive != 0 {
		t.Errorf("Expected 0 active requests, got %d", finalActive)
	}
}

// ===== 测试对象定义 =====

type LifecycleTracker struct {
	mu        sync.Mutex
	callOrder []string
}

func (t *LifecycleTracker) Record(method string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.callOrder = append(t.callOrder, method)
}

// MockObject 基于 ObjectImpl 的通用测试对象
type MockObject struct {
	ioc.ObjectImpl
	name string

	// 钩子函数字段，可按需配置
	OnPostConfigFunc func() error
	OnPreInitFunc    func() error
	OnPostInitFunc   func() error
	OnPreStopFunc    func(context.Context) error
	OnPostStopFunc   func(context.Context) error
	InitFunc         func() error
	CloseFunc        func(context.Context)

	// 状态追踪
	tracker *LifecycleTracker
}

func (o *MockObject) Name() string {
	return o.name
}

func (o *MockObject) OnPostConfig() error {
	if o.tracker != nil {
		o.tracker.Record("OnPostConfig")
	}
	if o.OnPostConfigFunc != nil {
		return o.OnPostConfigFunc()
	}
	// 调用父类默认实现
	return o.ObjectImpl.OnPostConfig()
}

func (o *MockObject) OnPreInit() error {
	if o.tracker != nil {
		o.tracker.Record("OnPreInit")
	}
	if o.OnPreInitFunc != nil {
		return o.OnPreInitFunc()
	}
	return o.ObjectImpl.OnPreInit()
}

func (o *MockObject) Init() error {
	if o.tracker != nil {
		o.tracker.Record("Init")
	}
	if o.InitFunc != nil {
		return o.InitFunc()
	}
	return o.ObjectImpl.Init()
}

func (o *MockObject) OnPostInit() error {
	if o.tracker != nil {
		o.tracker.Record("OnPostInit")
	}
	if o.OnPostInitFunc != nil {
		return o.OnPostInitFunc()
	}
	return o.ObjectImpl.OnPostInit()
}

func (o *MockObject) OnPreStop(ctx context.Context) error {
	if o.tracker != nil {
		o.tracker.Record("OnPreStop")
	}
	if o.OnPreStopFunc != nil {
		return o.OnPreStopFunc(ctx)
	}
	return o.ObjectImpl.OnPreStop(ctx)
}

func (o *MockObject) Close(ctx context.Context) {
	if o.tracker != nil {
		o.tracker.Record("Close")
	}
	if o.CloseFunc != nil {
		o.CloseFunc(ctx)
		return
	}
	o.ObjectImpl.Close(ctx)
}

func (o *MockObject) OnPostStop(ctx context.Context) error {
	if o.tracker != nil {
		o.tracker.Record("OnPostStop")
	}
	if o.OnPostStopFunc != nil {
		return o.OnPostStopFunc(ctx)
	}
	return o.ObjectImpl.OnPostStop(ctx)
}
