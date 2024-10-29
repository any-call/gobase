package mynet

import (
	"github.com/any-call/gobase/util/mylog"
	"io"
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
