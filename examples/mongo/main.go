package main

import (
	"fmt"

	"github.com/infraboard/mcube/v2/ioc/config/mongo"
)

func main() {
	// 获取mongodb 客户端对象
	client := mongo.Client()
	fmt.Println(client)

	// 获取DB对象
	db := mongo.DB()
	fmt.Println(db)
}
