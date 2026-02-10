package myctrl

import (
	"sync"
	"time"
)

// RateLimiter 限流器
type RateLimiter struct {
	interval time.Duration // 请求间隔
	mutex    sync.Mutex    // 互斥锁
	nextTime time.Time     // 下一次允许请求的时间
}

// NewRateLimiter 创建限流器
func NewRateLimiter(interval time.Duration) *RateLimiter {
	return &RateLimiter{
		interval: interval,
		nextTime: time.Now(),
	}
}

// Wait 等待直到允许请求
func (rl *RateLimiter) Wait() {
	rl.mutex.Lock()
	now := time.Now()
	if now.Before(rl.nextTime) {
		sleep := rl.nextTime.Sub(now)
		rl.mutex.Unlock()
		time.Sleep(sleep)
		rl.mutex.Lock()
	}

	rl.nextTime = rl.nextTime.Add(rl.interval)
	rl.mutex.Unlock()
}
