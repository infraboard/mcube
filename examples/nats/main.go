package main

import (
	"context"
	"fmt"
	"time"

	"github.com/infraboard/mcube/v2/ioc"
	ioc_nats "github.com/infraboard/mcube/v2/ioc/config/nats"
	"github.com/infraboard/mcube/v2/ioc/server"
	"github.com/nats-io/nats.go"
)

const (
	TEST_SUBJECT = "event_bus"
)

func main() {
	ioc.DevelopmentSetup()

	// 订阅消息
	_, err := ioc_nats.Get().Subscribe(TEST_SUBJECT, func(msg *nats.Msg) {
		fmt.Println(string(msg.Data))
	})
	if err != nil {
		fmt.Println(err)
	}

	// 发布消息
	go func() {
		for {
			time.Sleep(1 * time.Second)
			err = ioc_nats.Get().Publish(TEST_SUBJECT, []byte("test"))
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	// 启动应用
	err = server.Run(context.Background())
	if err != nil {
		panic(err)
	}
}
