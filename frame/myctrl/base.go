package myctrl

import (
	"time"
)

type (
	Golimiter interface {
		Begin()
		End()
		Number() int32
	}
)

func DelayExec(t time.Duration, fn func()) {
	if fn != nil {
		timer := time.NewTimer(t)
		go func(t *time.Timer) {
			<-t.C
			fn()
			timer.Stop()
		}(timer)
	}
}
