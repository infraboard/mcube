package main

import (
	"fmt"

	"github.com/infraboard/mcube/v2/ioc/config/datasource"
)

func main() {
	db := datasource.DB()
	// 通过db对象进行数据库操作
	fmt.Println(db)
}
