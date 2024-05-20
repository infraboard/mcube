package main

import (
	"fmt"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/ip2region"
)

func main() {
	ioc.DevelopmentSetup()

	resp, err := ip2region.Get().LookupIP("117.136.38.42")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
	// 中国|0|北京|北京市|移动
}
