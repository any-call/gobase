package myctrl

import (
	"sync"
	"time"
)

type goLimiter struct {
	sync.Mutex
	limiter chan struct{}
	goNum   int
}

func NewGolimiter(goNum int) Golimiter {
	if goNum <= 0 {
		goNum = 1
	}

	return &goLimiter{limiter: make(chan struct{}, goNum), goNum: goNum}
}

func (self *goLimiter) Begin() {
	for {
		if self.Number() < self.goNum {
			self.Lock()
			defer self.Unlock()

			self.limiter <- struct{}{}
			return
		}
		time.Sleep(time.Millisecond)
	}
}

func (self *goLimiter) End() {
	self.Lock()
	defer self.Unlock()

	if len(self.limiter) > 0 {
		<-self.limiter
	}
}

func (self *goLimiter) Number() int {
	self.Lock()
	defer self.Unlock()

	return len(self.limiter)
}
