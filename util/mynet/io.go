package mynet

import (
	"context"
	"github.com/any-call/gobase/util/mylog"
	"io"
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

func copyHalfClose(dst io.Writer, src io.Reader) (int64, error) {
	defer func() {
		if c, ok := dst.(closeWriter); ok {
			mylog.Info("enter close writer 1")
			_ = c.CloseWrite()
		}

		if c, ok := src.(closeReader); ok {
			mylog.Info("enter close writer 2")
			_ = c.CloseRead()
		}
	}()

	return io.Copy(dst, src)
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
		close(ch) // 完成后关闭通道
	}()

	down_n, err := copyHalfClose(left, right)
	rs := <-ch

	if err == nil {
		err = rs.Err
	}

	return down_n, rs.N, err
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
		close(upCh)
	}()

	go func() {
		down_n, err := copyHalfClose(left, right)
		downCh <- res{down_n, err}
		close(downCh)
	}()

	var upRes, downRes res
	for i := 0; i < 2; i++ {
		select {
		case upRes = <-upCh:
		case downRes = <-downCh:
		case <-ctx.Done():
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
