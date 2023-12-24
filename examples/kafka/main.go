package main

import (
	"fmt"

	"github.com/infraboard/mcube/v2/ioc/config/kafka"
)

func main() {
	// 消息生产者
	producer := kafka.Producer("test")
	fmt.Println(producer)

	// 消息消费者
	consumer := kafka.ConsumerGroup("group id", []string{"topic name"})
	fmt.Println(consumer)
}
