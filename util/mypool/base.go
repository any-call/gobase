package mypool

import "sync"

type Pool[V any] struct {
	sync.Pool
}

func NewPool[V any]() *Pool[V] {
	return &Pool[V]{}
}

func (p *Pool[V]) GetObj() (V, bool) {
	var value V
	o := p.Get()
	if o != nil {
		return o.(V), true
	}

	return value, false
}

func (p *Pool[V]) SetObj(o V) {
	p.Put(o)
}
