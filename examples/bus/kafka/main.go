package main

import (
	"context"
	"fmt"
	"time"

	"github.com/infraboard/mcube/v2/ioc/config/bus"
	// kafka总线
	_ "github.com/infraboard/mcube/v2/ioc/config/bus/kafka"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/server"
)

const (
	TEST_SUBJECT = "event_bus"
)

func main() {
	ioc.DevelopmentSetup()

	// 消息生产者
	bus.GetService().TopicSubscribe(context.Background(), TEST_SUBJECT, func(e *bus.Event) {
		fmt.Println(string(e.Data))
	})

	// 发布消息
	go func() {
		for {
			time.Sleep(1 * time.Second)
			err := bus.GetService().Publish(context.Background(), &bus.Event{
				Subject: TEST_SUBJECT,
				Data:    []byte("test"),
			})
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	// 消息消费者
	// 启动应用
	err := server.Run(context.Background())
	if err != nil {
		panic(err)
	}
}
