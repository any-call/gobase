package mybus

// 订阅标准
type BusSubscriber interface {
	Subscribe(key string, fn any) error
	SubscribeAsync(key string, fn any) error
	SubscribeOnce(key string, fn any) error
	SubscribeOnceAsync(key string, fn any) error
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

type RpcBus interface {
	Start() error
	Stop()
	ServerAddr() string
	ServerPath() string
	Bus() EventBus
}

type ClientBus interface {
	RpcBus
	Subscribe(key string, fn any, serverAddr, serverPath string) error
	SubscribeOnce(key string, fn any, serverAddr, serverPath string) error
	PushEvent(arg *ClientArg, reply *bool) error
}
