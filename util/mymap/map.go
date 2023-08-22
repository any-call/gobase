package mymap

import (
	"github.com/any-call/gobase/util/mylist"
	"sync"
)

type Map[K comparable, V any] struct {
	lock    sync.RWMutex
	mapList map[K]V
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return new(Map[K, V]).init()
}

func (l *Map[K, V]) init() *Map[K, V] {
	l.mapList = make(map[K]V, 100)
	return l
}

func (l *Map[K, V]) ResetByMap(items map[K]V) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.mapList = make(map[K]V, 100)
	if items != nil {
		for k, v := range items {
			l.mapList[k] = v
		}
	}

	return
}

func (l *Map[K, V]) Insert(k K, v V) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.mapList[k] = v
	return
}

func (l *Map[K, V]) Remove(k K) {
	l.lock.Lock()
	defer l.lock.Unlock()

	delete(l.mapList, k)
	return
}

func (l *Map[K, V]) TakeAt(k K) (v V, b bool) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if v, b = l.mapList[k]; b {
		delete(l.mapList, k)
	}

	return
}

func (l *Map[K, V]) Value(k K) (v V, b bool) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	v, b = l.mapList[k]
	return
}

func (l *Map[K, V]) Values() *mylist.List[V] {
	l.lock.RLock()
	defer l.lock.RUnlock()

	list := mylist.NewList[V]()
	for _, v := range l.mapList {
		list.Append(v)
	}

	return list
}

func (l *Map[K, V]) ToArray() (allKeys []K, allValues []V) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	allKeys = make([]K, len(l.mapList))
	allValues = make([]V, len(l.mapList))
	index := 0
	for k, v := range l.mapList {
		allKeys[index] = k
		allValues[index] = v
		index++
	}
	return
}

func (l *Map[K, V]) ToMap() map[K]V {
	l.lock.RLock()
	defer l.lock.RUnlock()

	ret := make(map[K]V)
	for k, v := range l.mapList {
		ret[k] = v
	}
	return ret
}

func (l *Map[K, V]) Keys() *mylist.List[K] {
	l.lock.RLock()
	defer l.lock.RUnlock()

	list := mylist.NewList[K]()
	for i, _ := range l.mapList {
		list.Append(i)
	}

	return list
}

func (l *Map[K, V]) Range(f func(key K, value V)) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	for i, v := range l.mapList {
		f(i, v)
	}
}

func (l *Map[K, V]) Search(f func(key K, value V) bool) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	for i, v := range l.mapList {
		if b := f(i, v); b {
			goto end
		}
	}
end:
}

func (l *Map[K, V]) Len() int {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return len(l.mapList)
}
