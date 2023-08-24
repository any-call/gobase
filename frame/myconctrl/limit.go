package myconctrl

type goLimiter struct {
	limiter chan struct{}
}

func NewGolimiter(goNum uint) Golimiter {
	if goNum == 0 {
		goNum = 1
	}

	return &goLimiter{limiter: make(chan struct{}, goNum)}
}

func (self *goLimiter) Begin() {
	self.limiter <- struct{}{}
}

func (self *goLimiter) End() {
	<-self.limiter
}
