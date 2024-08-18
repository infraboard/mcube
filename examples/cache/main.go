package main

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc/config/cache"
)

type TestStruct struct {
	Id     string
	FiledA string
}

func main() {
	ctx := context.Background()

	var v *TestStruct

	// objectId --->  cached Objected
	// 获取objectId对应的对象, 如果缓存中有则之间从缓存中获取, 如果没有 则通过提供的ObjectFinder直接获取
	err := cache.NewGetter(ctx, func(ctx context.Context, objectId string) (any, error) {
		return &TestStruct{Id: objectId, FiledA: "test"}, nil
	}).Get("objectId", v)
	if err != nil {
		panic(err)
	}
}
