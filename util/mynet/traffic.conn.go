package mynet

import (
	"net"
	"sync"
)

// TrafficConn 结构体用于包装 net.Conn，增加流量统计功能
type TrafficConn struct {
	net.Conn
	mux        sync.Mutex
	bytesRead  uint64
	bytesWrite uint64
}

func NewTrafficConn(conn net.Conn) *TrafficConn {
	return &TrafficConn{
		Conn:       conn,
		bytesRead:  0,
		bytesWrite: 0,
	}
}

func (c *TrafficConn) ResetTraffic() {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.bytesRead = 0
	c.bytesWrite = 0
}

func (c *TrafficConn) GetTraffic() (rCount, wCount uint64) {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.bytesRead, c.bytesWrite
}

func (c *TrafficConn) TakeTraffic() (rCount, wCount uint64) {
	c.mux.Lock()
	defer c.mux.Unlock()
	tmpR := c.bytesRead
	tmpW := c.bytesWrite
	c.bytesRead = 0
	c.bytesWrite = 0
	return tmpR, tmpW
}

// Read 方法拦截读取操作并统计流量
func (c *TrafficConn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	c.mux.Lock()
	c.bytesRead += uint64(n)
	c.mux.Unlock()
	return
}

// Write 方法拦截写入操作并统计流量
func (c *TrafficConn) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	c.mux.Lock()
	c.bytesWrite += uint64(n)
	c.mux.Unlock()
	return
}
