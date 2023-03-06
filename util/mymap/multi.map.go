package mymap

import (
	"github.com/any-call/gobase/util/mylist"
	"sync"
)

type MultiMap[K, V comparable] struct {
	lock    sync.RWMutex
	mapList map[K][]V
}

func NewMultiMap[K, V comparable]() *MultiMap[K, V] {
	return new(MultiMap[K, V]).init()
}

func (l *MultiMap[K, V]) init() *MultiMap[K, V] {
	l.mapList = make(map[K][]V, 100)
	return l
}

func (l *MultiMap[K, V]) Insert(key K, value V) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if values, ok := l.mapList[key]; ok {
		values = append(values, value)
		l.mapList[key] = values
	} else {
		l.mapList[key] = []V{value}
	}

	return
}

func (l *MultiMap[K, V]) Removes(key K) {
	l.lock.Lock()
	defer l.lock.Unlock()

	delete(l.mapList, key)
	return
}

func (l *MultiMap[K, V]) Remove(key K) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if values, ok := l.mapList[key]; ok {
		values = values[:len(values)-1]
		if len(values) == 0 {
			delete(l.mapList, key)
		} else {
			l.mapList[key] = values
		}
	}

	return
}

func (l *MultiMap[K, V]) TakeAt(k K) (v V, b bool) {
	l.lock.Lock()
	defer l.lock.Unlock()

	var values []V
	if values, b = l.mapList[k]; b {
		if len(values) == 0 {
			b = false
			delete(l.mapList, k)
		} else {
			v = values[len(values)-1]
			values = values[:len(values)-1]
			l.mapList[k] = values
		}
	}

	return
}

func (l *MultiMap[K, V]) Value(k K) (v V, b bool) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	var values []V
	if values, b = l.mapList[k]; b {
		if len(values) > 0 {
			v = values[len(values)-1]
		} else {
			b = false
		}
	}

	return
}

func (l *MultiMap[K, V]) Values(key K) (values []V, b bool) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	values, b = l.mapList[key]
	return
}

func (l *MultiMap[K, V]) Keys() *mylist.List[K] {
	l.lock.RLock()
	defer l.lock.RUnlock()

	list := mylist.NewList[K]()
	for i, _ := range l.mapList {
		list.Append(i)
	}

	return list
}

func (l *MultiMap[K, V]) Range(f func(key K, value V)) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	for i, values := range l.mapList {
		for _, v := range values {
			f(i, v)
		}
	}
}

func (l *MultiMap[K, V]) Len() int {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return len(l.mapList)
}
