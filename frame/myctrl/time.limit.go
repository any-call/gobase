package myctrl

import (
	"github.com/any-call/gobase/util/mycache"
	"github.com/any-call/gobase/util/mylog"
	"sync/atomic"
	"time"
)

const (
	timeDuration = "timeDuration"
)

type goTimeLimiter struct {
	limiter   chan struct{}
	maxNum    int32
	unitNum   int32
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

	return &goTimeLimiter{limiter: make(chan struct{}, goNum), maxNum: int32(goNum), t: t, cacheTime: mycache.NewCache()}
}

func (self *goTimeLimiter) Begin() {
	self.limiter <- struct{}{}
	if v, b := self.cacheTime.Get(timeDuration); b {
		if intV, ok := v.(int32); ok {
			if intV >= self.maxNum {
				//等待
				<-self.limiter
				time.Sleep(time.Millisecond * 5)
				self.Begin()
			} else {
				self.cacheTime.UpdateValue(timeDuration, atomic.AddInt32(&self.unitNum, 1))
			}
		} else {
			mylog.Debug("will send not int32:", v)
		}
	} else {
		atomic.StoreInt32(&self.unitNum, 1)
		self.cacheTime.Set(timeDuration, int32(1), self.t)
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
