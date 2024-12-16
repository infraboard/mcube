package lock_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/lock"
	"github.com/infraboard/mcube/v2/tools/file"
)

var (
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
)

func TestRedisLockUnLock(t *testing.T) {
	os.Setenv("LOCK_PROVIDER", lock.PROVIDER_REDIS)
	ioc.DevelopmentSetup()

	m := lock.L().New("test", 10*time.Second)
	if err := m.Lock(ctx); err != nil {
		t.Fatal(err)
	}
	if err := m.UnLock(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestRedisLock(t *testing.T) {
	os.Setenv("LOCK_PROVIDER", lock.PROVIDER_REDIS)
	ioc.DevelopmentSetup()
	g := &sync.WaitGroup{}
	for i := range 9 {
		go LockTest(i, g)
	}
	g.Wait()
	time.Sleep(10 * time.Second)
}

func TestRedisTryLock(t *testing.T) {
	os.Setenv("LOCK_PROVIDER", lock.PROVIDER_REDIS)
	ioc.DevelopmentSetup()
	g := &sync.WaitGroup{}
	for i := range 9 {
		go TryLockTest(i, g)
	}
	g.Wait()
	time.Sleep(10 * time.Second)
}

func TestGoCacheRedisLock(t *testing.T) {
	ioc.DevelopmentSetup()
	g := &sync.WaitGroup{}
	for i := range 9 {
		go LockTest(i, g)
	}
	g.Wait()
	time.Sleep(10 * time.Second)
}

func LockTest(number int, g *sync.WaitGroup) {
	fmt.Println(number, "start")
	g.Add(1)
	defer g.Done()
	m := lock.L().New("test", 1*time.Second)
	if err := m.Lock(ctx); err != nil {
		fmt.Println(err)
	}
	fmt.Println(number, "down")
}

func TryLockTest(number int, g *sync.WaitGroup) {
	fmt.Println(number, "start")
	g.Add(1)
	defer g.Done()
	m := lock.L().New("test", 1*time.Second)
	if err := m.TryLock(ctx); err != nil {
		fmt.Println(number, err)
		return
	}
	defer m.UnLock(ctx)
	fmt.Println(number, "obtained lock")
}

func TestDefaultConfig(t *testing.T) {
	file.MustToToml(
		lock.AppName,
		ioc.Config().Get(lock.AppName),
		"test/default.toml",
	)
}
