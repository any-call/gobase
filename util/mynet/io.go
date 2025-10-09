package mynet

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
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
	// 每次 read 的超时，用于可以设置 read deadline 的 src（例如 net.Conn）
	// 设为 0 表示不设置 deadline
	readDeadline = 30 * time.Second

	// 每次拷贝使用的缓冲区大小（与 io.Copy 默认 32KB 接近）
	defaultBufSize = 32 * 1024
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, defaultBufSize)
	},
}

func copyHalfClose(dst io.Writer, src io.Reader) (int64, error) {
	defer func() {
		if c, ok := dst.(closeWriter); ok {
			//mylog.Info("enter close writer 1")
			_ = c.CloseWrite()
		}

		if c, ok := src.(closeReader); ok {
			//mylog.Info("enter close writer 2")
			_ = c.CloseRead()
		}
	}()

	//return io.Copy(dst, src)
	return SmartCopy(dst, src)
}

// SmartCopy：像 io.Copy 一样的签名。
// - 优先走 src.(io.WriterTo) / dst.(io.ReaderFrom) 优化路径（零拷贝）
// - 否则使用自定义的缓冲循环，复用缓冲池、支持读超时、半关闭等增强行为
func SmartCopy(dst io.Writer, src io.Reader) (int64, error) {
	// 优化路径 1：src 提供 WriteTo
	if wt, ok := src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}
	// 优化路径 2：dst 提供 ReadFrom
	if rf, ok := dst.(io.ReaderFrom); ok {
		return rf.ReadFrom(src)
	}

	// 回退路径：我们自己做循环复制
	buf := bufPool.Get().([]byte)
	defer bufPool.Put(buf)

	var total int64

	// 检测 src 是否支持 SetReadDeadline（通常为 net.Conn）
	type setReadDeadline interface {
		SetReadDeadline(t time.Time) error
	}
	var srcDeadline setReadDeadline
	if sd, ok := src.(setReadDeadline); ok && readDeadline > 0 {
		srcDeadline = sd
	}

	for {
		// 如果支持 read-deadline，则提前设置（用于避免长时间阻塞）
		if srcDeadline != nil {
			_ = srcDeadline.SetReadDeadline(time.Now().Add(readDeadline))
		}

		nr, er := src.Read(buf)
		if nr > 0 {
			nwTotal := 0
			for nwTotal < nr {
				nw, ew := dst.Write(buf[nwTotal:nr])
				if nw > 0 {
					nwTotal += nw
				}
				if ew != nil {
					// 写错误直接返回（返回已写入的字节数）
					return total + int64(nwTotal), ew
				}
			}
			total += int64(nr)
		}

		if er != nil {
			// 如果是超时且可重试（net.Error.Timeout），则继续循环读取
			if ne, ok := er.(net.Error); ok && ne.Timeout() {
				// 继续尝试读取
				continue
			}

			if errors.Is(er, io.EOF) {
				// 按 io.Copy 的语义，返回总字节数和 io.EOF（或 nil）。
				// io.Copy 在遇 EOF 时会返回写入总字节数和 nil（注意：Read 返回 EOF，io.Copy 会将 EOF 作为结束信号并返回 nil）
				// 为更贴近 io.Copy 行为，我们在 EOF 时返回 (total, nil)
				return total, nil
			}

			// 其他读错误则返回
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
