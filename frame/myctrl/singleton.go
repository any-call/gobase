package myctrl

import "sync"

// Singleton 是一个单例接口，包含 GetInstance 方法
type Singleton[T any] interface {
	GetInstance() *T
}

// NewSingleton 创建并返回一个实现 Singleton 接口的单例
func NewSingleton[T any](newInstance func() *T) Singleton[T] {
	return &singleton[T]{
		newInstance: newInstance,
	}
}

// singleton 是实现了 Singleton 接口的具体类型，不对外公开
type singleton[T any] struct {
	instance    *T
	once        sync.Once
	newInstance func() *T
}

// GetInstance 获取单例对象的实例
func (s *singleton[T]) GetInstance() *T {
	s.once.Do(func() {
		s.instance = s.newInstance()
	})
	return s.instance
}
