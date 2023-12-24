package main

import (
	"fmt"

	"github.com/infraboard/mcube/v2/ioc/config/redis"
)

func main() {
	client := redis.Client()
	fmt.Println(client)
}
