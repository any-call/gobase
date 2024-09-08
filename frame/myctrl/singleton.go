package myctrl

import "sync"

// 单例对象
// Singleton 是一个泛型的单例结构
type Singleton[T any] struct {
	instance *T
	once     sync.Once
}

// GetInstance 获取单例对象的实例，只有在第一次调用时会创建新实例
func (s *Singleton[T]) GetInstance(newInstance func() *T) *T {
	s.once.Do(func() {
		s.instance = newInstance()
	})
	return s.instance
}
