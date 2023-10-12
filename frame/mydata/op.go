package mydata

import (
	"sync"
)

type dataOP[DATA any] struct {
	sync.RWMutex
	data DATA
}

func NewData[DATA any](obj DATA) Data[DATA] {
	return &dataOP[DATA]{data: obj}
}

func (self *dataOP[DATA]) Get() DATA {
	self.RLock()
	defer self.RUnlock()

	return self.data
}

func (self *dataOP[DATA]) Set(fn func(d DATA)) {
	self.Lock()
	defer self.Unlock()

	if fn != nil {
		fn(self.data)
	}
}

func (self *dataOP[DATA]) TrySet(fn func(d DATA)) bool {
	b := self.TryLock()
	if b {
		if fn != nil {
			fn(self.data)
		}
		return b
	}

	return b
}
