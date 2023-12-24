package main

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc/config/cache"
)

type TestStruct struct {
	FiledA string
}

func main() {
	ctx := context.Background()

	c := cache.C()

	key := "test"
	obj := &TestStruct{FiledA: "test"}

	// 设置缓存
	err := c.Set(ctx, key, obj, cache.WithExpiration(300))
	if err != nil {
		panic(err)
	}

	// 后期缓存
	var v *TestStruct
	err = c.Get(ctx, key, v)
	if err != nil {
		panic(err)
	}
}
