package mybus

// 订阅标准
type BusSubscriber interface {
	Subscribe(key string, fn any) error
	SubscribeAsync(key string, fn any, sequence bool) error
	SubscribeOnce(key string, fn any) error
	SubscribeOnceAsync(key string, fn any, sequence bool) error
	Unsubscribe(key string, fn any)
}

// 发布标准
type BusPublisher interface {
	Publish(key string, args ...any)
	HasCallback(key string) bool
	WaitAsync()
}

type EventBus interface {
	BusSubscriber
	BusPublisher
}
