package main

import (
	"context"
	"time"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/lock"
)

func main() {
	ioc.DevelopmentSetup()

	// 创建一个key为test, 超时时间为10秒的锁
	locker := lock.L().New("test", 10*time.Second)

	ctx := context.Background()

	// 获取锁
	if err := locker.Lock(ctx); err != nil {
		panic(err)
	}

	// 释放锁
	defer locker.UnLock(ctx)
}
