package myctrl

import (
	"sync"
)

// IDLocker 针对相同 Key 串行，不同 Key 并行
type IDLocker[K comparable] struct {
	mu    sync.Mutex
	locks map[K]*sync.Mutex
}

// NewIDLocker 创建泛型版 IDLocker
func NewIDLocker[K comparable]() *IDLocker[K] {
	return &IDLocker[K]{
		locks: make(map[K]*sync.Mutex),
	}
}

// Lock 对对应 Key 加锁
func (l *IDLocker[K]) Lock(key K) {
	l.mu.Lock()
	mtx, ok := l.locks[key]
	if !ok {
		mtx = &sync.Mutex{}
		l.locks[key] = mtx
	}
	l.mu.Unlock()

	mtx.Lock()
}

// Unlock 对对应 Key 解锁
func (l *IDLocker[K]) Unlock(key K) {
	l.mu.Lock()
	mtx, _ := l.locks[key]
	l.mu.Unlock()

	if mtx != nil {
		mtx.Unlock()
	}
}
