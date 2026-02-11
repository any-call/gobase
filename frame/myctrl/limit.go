package myctrl

import "github.com/any-call/gobase/util/mylog"

type goLimiter struct {
	limiter chan struct{}
	maxNum  int
}

func NewGolimiter(goNum int) Golimiter {
	if goNum <= 0 {
		goNum = 1
	}

	ret := &goLimiter{
		limiter: make(chan struct{}, goNum),
		maxNum:  goNum,
	}
	return ret
}

func (self *goLimiter) Do(fn func()) {
	self.limiter <- struct{}{}
	self.runAsync(fn)
}

func (self *goLimiter) TryDo(fn func()) bool {
	select {
	case self.limiter <- struct{}{}:
		self.runAsync(fn)
		return true
	default:
		return false
	}
}

func (self *goLimiter) DoAndWait(fn func()) {
	self.limiter <- struct{}{}
	self.runSync(fn)
}

func (self *goLimiter) TryDoAndWait(fn func()) bool {
	select {
	case self.limiter <- struct{}{}:
		self.runSync(fn)
		return true
	default:
		return false
	}
}

func (self *goLimiter) MaxNumber() int {
	return self.maxNum
}

func (self *goLimiter) Number() int {
	return len(self.limiter)
}

// / 内部函数
func (g *goLimiter) runAsync(fn func()) {
	go func() {
		defer g.finish()
		if fn != nil {
			fn()
		}
	}()
}

func (g *goLimiter) runSync(fn func()) {
	defer g.finish()
	if fn != nil {
		fn()
	}
}

func (g *goLimiter) finish() {
	if r := recover(); r != nil {
		mylog.Debug("panic:", r)
	}
	<-g.limiter
}
