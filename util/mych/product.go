package mych

type product[DATA any] struct {
	ch chan DATA
}

func NewProduct[DATA any](chLen uint) Product[DATA] {
	return &product[DATA]{
		ch: make(chan DATA, chLen),
	}
}

func (self *product[DATA]) Producter() chan<- DATA {
	return self.ch
}

func (self *product[DATA]) Consumer() <-chan DATA {
	return self.ch
}

func (self *product[DATA]) Send(data DATA) {
	self.ch <- data
}

func (self *product[DATA]) SendBy(valid ValidFunc[DATA], data DATA) {
	if valid != nil {
		if b := valid(data); b {
			self.Send(data)
		}
	}
}

func (self *product[DATA]) Receive() DATA {
	return <-self.Consumer()
}

func (self *product[DATA]) ReceiveBy(handler HandleFunc[DATA]) {
	if handler != nil {
		for handler(self.Receive()) {
		}
	}
}

func (self *product[DATA]) DataLen() int { //返回 channel 中已缓存的数据数量
	return len(self.ch)
}
