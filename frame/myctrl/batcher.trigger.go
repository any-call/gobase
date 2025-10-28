package myctrl

import (
	"sync"
	"time"
)

type BatchTrigger[T any] struct {
	mu       sync.Mutex
	buffer   []T
	maxSize  int
	interval time.Duration
	reportFn func([]T)

	timer    *time.Ticker
	inputCh  chan T
	stopCh   chan struct{}
	flushing bool // 防止重复上报
}

func NewBatchTrigger[T any](maxSize int, interval time.Duration, reportFn func([]T)) *BatchTrigger[T] {
	if maxSize <= 0 {
		maxSize = 100
	}
	if interval <= 0 {
		interval = time.Minute
	}

	bt := &BatchTrigger[T]{
		maxSize:  maxSize,
		interval: interval,
		reportFn: reportFn,
		// 缓冲区自动计算：2~5倍 maxSize，根据实际压力可调整
		inputCh: make(chan T, maxSize*4),
		stopCh:  make(chan struct{}),
		timer:   time.NewTicker(interval),
	}
	go bt.loop()
	return bt
}

func (b *BatchTrigger[T]) loop() {
	for {
		select {
		case <-b.stopCh:
			b.flushNow() // 退出前强制 flush
			return

		case item := <-b.inputCh:
			b.mu.Lock()
			b.buffer = append(b.buffer, item)
			full := len(b.buffer) >= b.maxSize
			b.mu.Unlock()

			if full {
				b.flushNow()
			}

		case <-b.timer.C:
			b.flushNow()
		}
	}
}

func (b *BatchTrigger[T]) Add(item T) {
	select {
	case b.inputCh <- item:
	default:
		// 防止通道阻塞（可选：丢弃或扩容）
	}
}

func (b *BatchTrigger[T]) flushNow() {
	b.mu.Lock()
	// 防止并发 flush
	if b.flushing || len(b.buffer) == 0 {
		b.mu.Unlock()
		return
	}
	b.flushing = true

	// 提取当前批次
	data := b.buffer
	b.buffer = nil
	b.mu.Unlock()

	// 异步执行上报
	go func(batch []T) {
		defer func() {
			b.mu.Lock()
			b.flushing = false
			b.timer.Reset(b.interval) //重置定时器
			b.mu.Unlock()
		}()
		b.reportFn(batch)
	}(data)
}

func (b *BatchTrigger[T]) Stop() {
	close(b.stopCh)
	b.timer.Stop()
}
