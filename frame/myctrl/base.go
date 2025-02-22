package myctrl

import (
	"fmt"
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

// RetryFunction 是需要执行的函数，返回 error
type RetryFunction func() error

// Retry 执行指定的函数，失败后按指定次数和时间间隔重试
func Retry(fn RetryFunction, retries int, timeout time.Duration) error {
	if fn == nil {
		return nil
	}
	if retries <= 0 {
		retries = 1
	}

	var err error
	for i := 0; i < retries; i++ {
		if err = fn(); err == nil {
			return nil // 成功，直接返回
		}
		fmt.Printf("第 %d 次尝试失败: %v，等待 %v 后重试...\n", i+1, err, timeout)
		time.Sleep(timeout) // 等待后重试
	}
	return fmt.Errorf("执行失败，重试 %d 次后仍然出错: %w", retries, err)
}

func WaitForCondition(checkFunc func() bool, interval time.Duration) {
	if checkFunc == nil {
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if checkFunc() {
				return
			}
		}
	}
}
