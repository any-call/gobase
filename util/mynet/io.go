package mynet

import (
	"context"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

type (
	closeWriter interface {
		CloseWrite() error
	}

	closeReader interface {
		CloseRead() error
	}
)

var (
	// 每次拷贝使用的缓冲区大小（与 io.Copy 默认 32KB 接近）
	defaultBufSize = 32 * 1024
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, defaultBufSize)
	},
}

func copyHalfClose(dst io.Writer, src io.Reader) (int64, error) {
	return copyHalfCloseWithActivity(dst, src, nil)
}

func copyHalfCloseWithActivity(dst io.Writer, src io.Reader, onActivity func(int)) (int64, error) {
	defer func() {
		if c, ok := dst.(closeWriter); ok {
			_ = c.CloseWrite()
		}

		if c, ok := src.(closeReader); ok {
			_ = c.CloseRead()
		}
	}()

	return smartCopy(dst, src, onActivity)
}

// SmartCopy 与 io.Copy 语义一致，并复用内部缓冲区。
// 未启用流量观察时，优先使用 WriterTo/ReaderFrom 优化路径。
func SmartCopy(dst io.Writer, src io.Reader) (int64, error) {
	return smartCopy(dst, src, nil)
}

func smartCopy(dst io.Writer, src io.Reader, onActivity func(int)) (int64, error) {
	if onActivity == nil {
		if wt, ok := src.(io.WriterTo); ok {
			return wt.WriteTo(dst)
		}
		if rf, ok := dst.(io.ReaderFrom); ok {
			return rf.ReadFrom(src)
		}
	}

	buf := bufPool.Get().([]byte)
	defer bufPool.Put(buf)

	var total int64

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nwTotal := 0
			for nwTotal < nr {
				nw, ew := dst.Write(buf[nwTotal:nr])
				if nw > 0 {
					nwTotal += nw
					if onActivity != nil {
						onActivity(nw)
					}
				}
				if ew != nil {
					return total + int64(nwTotal), ew
				}
				if nw == 0 {
					return total + int64(nwTotal), io.ErrShortWrite
				}
			}
			total += int64(nr)
		}

		if er != nil {
			if errors.Is(er, io.EOF) {
				return total, nil
			}
			return total, er
		}
	}
}

func Relay(left, right io.ReadWriter) (int64, int64, error) {
	type res struct {
		N   int64
		Err error
	}

	ch := make(chan res, 1) // 使用缓冲通道
	go func() {
		up_n, err := copyHalfClose(right, left)
		ch <- res{up_n, err}
	}()

	down_n, err := copyHalfClose(left, right)
	rs := <-ch

	if err == nil {
		err = rs.Err
	}

	return down_n, rs.N, err
}

// RelayEx 双向转发数据。
// 正常 EOF 仅执行半关闭；非 EOF 错误会关闭两端。
// 为保证异常时能够及时退出，left 和 right 应实现 io.Closer。
func RelayEx(left, right io.ReadWriter) (downN, upN int64, relayErr error) {
	type result struct {
		isDown bool
		n      int64
		err    error
	}

	results := make(chan result, 2)

	// 客户端 -> 目标服务器
	go func() {
		n, err := copyHalfClose(right, left)
		results <- result{
			isDown: false,
			n:      n,
			err:    err,
		}
	}()

	// 目标服务器 -> 客户端
	go func() {
		n, err := copyHalfClose(left, right)
		results <- result{
			isDown: true,
			n:      n,
			err:    err,
		}
	}()

	for completed := 0; completed < 2; completed++ {
		rs := <-results

		if rs.isDown {
			downN = rs.n
		} else {
			upN = rs.n
		}

		// EOF 是正常半关闭，继续等待另一个方向。
		if rs.err == nil || errors.Is(rs.err, io.EOF) {
			continue
		}

		// 明确异常才关闭两端。
		if relayErr == nil {
			relayErr = rs.err
			closeRelayEndpoints(left, right)
		}
	}

	return downN, upN, relayErr
}

// RelayWithIdleTimeout 双向转发数据，并在整条连接持续无流量时退出。
// 任一方向成功转发数据都会刷新空闲计时，适用于长时间单向下载场景。
func RelayWithIdleTimeout(left, right io.ReadWriter, idleTimeout time.Duration) (int64, int64, error) {
	if idleTimeout <= 0 {
		return Relay(left, right)
	}

	type res struct {
		down bool
		n    int64
		err  error
	}

	results := make(chan res, 2)
	activity := make(chan struct{}, 1)
	var downN, upN int64
	var lastActivity int64
	atomic.StoreInt64(&lastActivity, time.Now().UnixNano())

	notify := func(counter *int64) func(int) {
		return func(n int) {
			atomic.AddInt64(counter, int64(n))
			atomic.StoreInt64(&lastActivity, time.Now().UnixNano())
			select {
			case activity <- struct{}{}:
			default:
			}
		}
	}

	go func() {
		n, err := copyHalfCloseWithActivity(right, left, notify(&upN))
		results <- res{n: n, err: err}
	}()
	go func() {
		n, err := copyHalfCloseWithActivity(left, right, notify(&downN))
		results <- res{down: true, n: n, err: err}
	}()

	timer := time.NewTimer(idleTimeout)
	defer timer.Stop()
	remainingIdle := func() time.Duration {
		return idleTimeout - time.Since(time.Unix(0, atomic.LoadInt64(&lastActivity)))
	}

	completed := 0
	for completed < 2 {
		select {
		case rs := <-results:
			completed++
			if rs.down {
				atomic.StoreInt64(&downN, rs.n)
			} else {
				atomic.StoreInt64(&upN, rs.n)
			}
			if rs.err != nil {
				closeRelayEndpoints(left, right)
				return atomic.LoadInt64(&downN), atomic.LoadInt64(&upN), rs.err
			}
		case <-activity:
			remaining := remainingIdle()
			if remaining <= 0 {
				closeRelayEndpoints(left, right)
				return atomic.LoadInt64(&downN), atomic.LoadInt64(&upN), context.DeadlineExceeded
			}
			resetTimer(timer, remaining)
		case <-timer.C:
			remaining := remainingIdle()
			if remaining > 0 {
				resetTimer(timer, remaining)
				continue
			}
			closeRelayEndpoints(left, right)
			return atomic.LoadInt64(&downN), atomic.LoadInt64(&upN), context.DeadlineExceeded
		}
	}

	return atomic.LoadInt64(&downN), atomic.LoadInt64(&upN), nil
}

func resetTimer(timer *time.Timer, timeout time.Duration) {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
	timer.Reset(timeout)
}

func closeRelayEndpoints(left, right io.ReadWriter) {
	if closer, ok := left.(io.Closer); ok {
		_ = closer.Close()
	}
	if closer, ok := right.(io.Closer); ok {
		_ = closer.Close()
	}
}

func RelayWithTimeout(left, right io.ReadWriter, timeout time.Duration) (int64, int64, error) {
	type res struct {
		N   int64
		Err error
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	upCh := make(chan res, 1)
	downCh := make(chan res, 1)

	// 开始双向拷贝
	go func() {
		up_n, err := copyHalfClose(right, left)
		upCh <- res{up_n, err}
	}()

	go func() {
		down_n, err := copyHalfClose(left, right)
		downCh <- res{down_n, err}
	}()

	var upRes, downRes res
	for i := 0; i < 2; i++ {
		select {
		case upRes = <-upCh:
		case downRes = <-downCh:
		case <-ctx.Done():
			closeRelayEndpoints(left, right)
			return downRes.N, upRes.N, ctx.Err()
		}
	}

	if upRes.Err != nil {
		return downRes.N, upRes.N, upRes.Err
	}
	if downRes.Err != nil {
		return downRes.N, upRes.N, downRes.Err
	}

	return downRes.N, upRes.N, nil
}
