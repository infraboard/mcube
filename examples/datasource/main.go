package main

import (
	"fmt"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
)

func main() {
	// 开启配置文件读取配置
	ioc.DevelopmentSetupWithPath("etc/application.toml")

	db := datasource.DB()
	// 通过db对象进行数据库操作
	fmt.Println(db)
}
