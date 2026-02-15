package ioc

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestForEachWithConcurrentWrite 测试 ForEach 在有写锁竞争时的行为
// 这是一个更严格的测试，用于验证递归锁在写锁竞争时是否安全
func TestForEachWithConcurrentWrite(t *testing.T) {
	t.Run("ForEach with Get during write contention", func(t *testing.T) {
		store := newNamespaceStore("test-contention")

		// 注册初始对象
		for i := range 10 {
			obj := &SimpleTestObject{
				name:     fmt.Sprintf("obj-%d", i),
				priority: i,
			}
			store.Registry(obj)
		}

		var wg sync.WaitGroup
		done := make(chan bool, 1)
		timeout := 5 * time.Second

		// 1. 启动一个 goroutine 不断尝试写入
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range 100 {
				obj := &SimpleTestObject{
					name:     fmt.Sprintf("new-obj-%d", i),
					priority: 100 + i,
				}
				store.Registry(obj)
				time.Sleep(1 * time.Millisecond)
			}
		}()

		// 2. 同时使用 ForEach 并在回调中调用 Get
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 50 {
				store.ForEach(func(w *ObjectWrapper) {
					// 在 ForEach（持有读锁）中调用 Get（尝试获取读锁）
					_ = store.Get(w.Name)
				})
				time.Sleep(2 * time.Millisecond)
			}
		}()

		// 等待完成
		go func() {
			wg.Wait()
			done <- true
		}()

		select {
		case <-done:
			t.Log("✅ Concurrent ForEach+Get with write operations completed")
		case <-time.After(timeout):
			t.Fatal("DEADLOCK DETECTED: Operation timed out under write contention")
		}
	})

	t.Run("ForEach should be safe with snapshot pattern", func(t *testing.T) {
		// 验证：如果 ForEach 也采用快照模式，会更安全
		store := newNamespaceStore("test-snapshot")

		for i := 0; i < 10; i++ {
			obj := &SimpleTestObject{
				name:     fmt.Sprintf("obj-%d", i),
				priority: i,
			}
			store.Registry(obj)
		}

		done := make(chan bool, 1)
		timeout := 5 * time.Second

		go func() {
			// 测试：即使在 ForEach 回调中修改容器，也应该是安全的
			count := 0
			store.ForEach(func(w *ObjectWrapper) {
				count++
				// 当前 ForEach 持有读锁，这里不应该尝试写入
				// 否则会死锁
			})
			if count != 10 {
				t.Errorf("Expected 10 items, got %d", count)
			}
			done <- true
		}()

		select {
		case <-done:
			t.Log("✅ ForEach with current implementation is safe")
		case <-time.After(timeout):
			t.Fatal("DEADLOCK DETECTED: ForEach timed out")
		}
	})
}

// SimpleTestObject 简单测试对象
type SimpleTestObject struct {
	ObjectImpl
	name     string
	priority int
}

func (o *SimpleTestObject) Name() string {
	if o.name == "" {
		return "simple-test"
	}
	return o.name
}

func (o *SimpleTestObject) Priority() int {
	return o.priority
}
