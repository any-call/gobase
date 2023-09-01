package myctrl

import (
	"github.com/any-call/gobase/util/mycache"
	"sync/atomic"
	"time"
)

const (
	timeDuration = "timeDuration"
)

type goTimeLimiter struct {
	limiter   chan struct{}
	maxNum    int32
	runNum    int32
	t         time.Duration
	cacheTime mycache.Cache
}

func NewGoTimeLimiter(goNum int, t time.Duration) *goTimeLimiter {
	if goNum <= 0 {
		goNum = 1
	}

	if t == 0 {
		t = time.Second
	}

	return &goTimeLimiter{limiter: make(chan struct{}, goNum), maxNum: int32(goNum), t: t, cacheTime: mycache.NewCache()}
}

func (self *goTimeLimiter) Begin() {
	self.limiter <- struct{}{}
	if v, b := self.cacheTime.Get(timeDuration); b {
		if intV, ok := v.(int32); ok {
			if intV > self.maxNum {
			}
		}
		self.cacheTime.UpdateValue(timeDuration, atomic.AddInt32(&self.runNum, 1))
	} else {
		atomic.StoreInt32(&self.runNum, 1)
		self.cacheTime.Set(timeDuration, self.runNum, self.t)
	}

	return
}

func (self *goTimeLimiter) End() {
	<-self.limiter
	if b := self.cacheTime.HasKey(timeDuration); b {
		self.cacheTime.UpdateValue(timeDuration, atomic.AddInt32(&self.runNum, -1))
	} else {
		self.cacheTime.Set(timeDuration, atomic.AddInt32(&self.runNum, -1), self.t)
	}

	return
}

func (self *goTimeLimiter) Number() int32 {
	return atomic.LoadInt32(&self.runNum)
}
