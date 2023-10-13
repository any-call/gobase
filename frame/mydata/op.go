package mydata

import (
	"sync"
)

type dataOP[DATA any] struct {
	sync.Mutex
	data DATA
}

func NewData[DATA any](obj DATA) Data[DATA] {
	return &dataOP[DATA]{data: obj}
}

func (self *dataOP[DATA]) Get() DATA {
	self.Lock()
	defer self.Unlock()

	return self.data
}

func (self *dataOP[DATA]) Set(d DATA) {
	self.Lock()
	defer self.Unlock()

	self.data = d
}

func (self *dataOP[DATA]) SetItem(fn func(d DATA)) {
	self.Lock()
	defer self.Unlock()

	if fn != nil {
		fn(self.data)
	}
}
