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

func TestDefaultConfig(t *testing.T) {
	file.MustToToml(
		lock.AppName,
		ioc.Config().Get(lock.AppName),
		"test/default.toml",
	)
}
