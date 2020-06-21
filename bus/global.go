package bus

var (
	publisher  Publisher
	subscriber Subscriber
)

// P bus为全局对象
func P() Publisher {
	if publisher == nil {
		panic("publisher not initail")
	}
	return publisher
}

// SetPublisher 设置pub
func SetPublisher(p Publisher) {
	publisher = p
}

// S bus为全局对象
func S() Subscriber {
	if subscriber == nil {
		panic("subscriber not initial")
	}
	return subscriber
}

// SetSubscriber 设置sub
func SetSubscriber(s Subscriber) {
	subscriber = s
}
