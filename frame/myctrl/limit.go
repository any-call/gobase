package myctrl

import "sync/atomic"

type goLimiter struct {
	limiter chan struct{}
	maxNum  int
	runNum  *atomic.Int32
}

func NewGolimiter(goNum int) Golimiter {
	if goNum <= 0 {
		goNum = 1
	}

	ret := &goLimiter{
		limiter: make(chan struct{}, goNum),
		maxNum:  goNum,
		runNum:  &atomic.Int32{},
	}
	ret.runNum.Add(0)
	return ret
}

func (self *goLimiter) Begin() {
	self.limiter <- struct{}{}
	self.runNum.Add(1)
	return
}

func (self *goLimiter) End() {
	<-self.limiter
	self.runNum.Add(-1)
	return
}

func (self *goLimiter) Number() int32 {
	return self.runNum.Load()
}

func (self *goLimiter) MaxNumber() int {
	return self.maxNum
}
