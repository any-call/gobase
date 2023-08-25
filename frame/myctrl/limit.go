package myctrl

type goLimiter struct {
	limiter chan struct{}
	goNum   int
}

func NewGolimiter(goNum int) *goLimiter {
	if goNum <= 0 {
		goNum = 1
	}

	return &goLimiter{limiter: make(chan struct{}, goNum), goNum: goNum}
}

func (self *goLimiter) Begin() {
	self.limiter <- struct{}{}
	return
}

func (self *goLimiter) End() {
	<-self.limiter
	return
}

func (self *goLimiter) Number() int {
	return len(self.limiter)
}
