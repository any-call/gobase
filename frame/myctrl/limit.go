package myctrl

import "sync/atomic"

type goLimiter struct {
	limiter chan struct{}
	maxNum  int
	runNum  int32
}

func NewGolimiter(goNum int) *goLimiter {
	if goNum <= 0 {
		goNum = 1
	}

	return &goLimiter{limiter: make(chan struct{}, goNum), maxNum: goNum}
}

func (self *goLimiter) Begin() {
	self.limiter <- struct{}{}
	atomic.AddInt32(&self.runNum, 1)
	return
}

func (self *goLimiter) End() {
	<-self.limiter
	atomic.AddInt32(&self.runNum, -1)
	return
}

func (self *goLimiter) Number() int32 {
	return atomic.LoadInt32(&self.runNum)
}
