package myctrl

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

	go func() {
		defer func() {
			<-self.limiter
		}()
		fn()
	}()
}

func (self *goLimiter) MaxNumber() int {
	return self.maxNum
}

func (g *goLimiter) Number() int {
	return len(g.limiter)
}
