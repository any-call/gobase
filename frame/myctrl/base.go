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

	GoTimelimiter interface {
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

func ObjFun[T any](f func() T) T {
	return f()
}

func TimerExec(t time.Duration, fn func()) {
	if fn != nil {
		timer := time.NewTimer(t)

		for {
			select {
			case <-timer.C:
				fn()
				timer.Reset(t)
				break
			}
		}
	}
}

func WaitForSignal[T any](timeout time.Duration, signal <-chan T) (T, bool) {
	select {
	case data := <-signal:
		// 收到信号
		return data, true
	case <-time.After(timeout):
		// 超时
		var zero T
		return zero, false
	}
}
