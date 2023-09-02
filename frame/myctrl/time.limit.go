package myctrl

import (
	"github.com/any-call/gobase/util/mycache"
	"sync"
	"sync/atomic"
	"time"
)

const (
	timeDuration = "timeDuration"
)

type goTimeLimiter struct {
	sync.Mutex
	limiter   chan struct{}
	maxNum    int32
	unitNum   *atomic.Int32
	t         time.Duration
	cacheTime mycache.Cache
}

func NewGoTimeLimiter(goNum int, t time.Duration) GoTimelimiter {
	if goNum <= 0 {
		goNum = 1
	}

	if t == 0 {
		t = time.Second
	}

	return &goTimeLimiter{
		limiter:   make(chan struct{}, goNum),
		maxNum:    int32(goNum),
		unitNum:   &atomic.Int32{},
		t:         t,
		cacheTime: mycache.NewCache()}
}

func (self *goTimeLimiter) Begin() {
	self.limiter <- struct{}{}

	if v, b := self.cacheTime.Get(timeDuration); b {
		if intV, ok := v.(int32); ok {
			if intV >= self.maxNum {
				<-self.limiter
				time.Sleep(time.Millisecond * 5)
				self.Begin()
			} else {
				self.Lock()
				defer self.Unlock()
				if err := self.cacheTime.UpdateValue(timeDuration, self.unitNum.Add(1)); err != nil {
					self.unitNum.Store(1)
					self.cacheTime.Set(timeDuration, self.unitNum.Load(), self.t)
				}
			}
		}
	} else {
		self.unitNum.Store(1)
		self.cacheTime.Set(timeDuration, self.unitNum.Load(), self.t)
	}

	return
}

func (self *goTimeLimiter) End() {
	<-self.limiter
	return
}

func (self *goTimeLimiter) Number() int32 {
	if num, ok := self.cacheTime.Get(timeDuration); ok {
		if int32V, ok := num.(int32); ok {
			return int32V
		}
	}

	return 0
}
