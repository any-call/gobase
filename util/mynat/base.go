package mynat

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/any-call/gobase/util/mylog"
	"github.com/any-call/gobase/util/mysocks5"
)

type Mode int

const (
	TargetServer Mode = iota
	SSNodeServer
	Socks5Server
)

const udpBufSize = 64 * 1024

// ==============================
// NAT Key
// ==============================

type NatKey struct {
	Client string
	Role   Mode
}

// ==============================
// NAT Entry
// ==============================

type natEntry struct {
	pc   net.PacketConn
	role Mode

	// 只有 socks5 才有
	socksCtrl net.Conn // TCP 控制连接
	udpRelay  net.Addr // socks5 UDP relay 地址
}

func NewNatEntry(pc net.PacketConn, r Mode) *natEntry {
	return &natEntry{
		pc:   pc,
		role: r,
	}
}

func (e *natEntry) WriteTo(payload []byte, tgt net.Addr) (int, error) {
	switch e.role {

	case TargetServer:
		return e.pc.WriteTo(payload, tgt)

	case Socks5Server:
		// 构造 socks5 UDP 包
		socksPacket := buildSocks5UDPPacket(tgt, payload)
		return e.pc.WriteTo(socksPacket, e.udpRelay)
	}

	return 0, fmt.Errorf("unknown role")
}

func (self *natEntry) SetSocks5Server(ctrlConn net.Conn, udpRalay net.Addr) {
	self.socksCtrl = ctrlConn
	self.udpRelay = udpRalay
}

func (e *natEntry) CloseAll() {
	if e.pc != nil {
		_ = e.pc.Close()
		e.pc = nil
	}

	if e.socksCtrl != nil {
		_ = e.socksCtrl.Close()
		e.socksCtrl = nil
	}
}

func (e *natEntry) watchSocksCtrl(onClose func()) {
	if e.socksCtrl == nil {
		return
	}

	go func() {
		buf := make([]byte, 1)
		_, err := e.socksCtrl.Read(buf)
		if err != nil {
			// ctrl 断开
			onClose()
		}
	}()
}

// ==============================
// NAT Map
// ==============================

type natmap struct {
	sync.RWMutex
	m       map[NatKey]*natEntry
	timeout time.Duration
}

func NewNATmap(timeout time.Duration) *natmap {
	return &natmap{
		m:       make(map[NatKey]*natEntry),
		timeout: timeout,
	}
}

func (m *natmap) Get(key NatKey) *natEntry {
	m.RLock()
	defer m.RUnlock()
	return m.m[key]
}

func (m *natmap) Set(key NatKey, entry *natEntry) {
	m.Lock()
	defer m.Unlock()
	m.m[key] = entry
}

func (m *natmap) Del(key NatKey) *natEntry {
	m.Lock()
	defer m.Unlock()

	entry, ok := m.m[key]
	if ok {
		delete(m.m, key)
		return entry
	}
	return nil
}

func (m *natmap) Count() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.m)
}

func (m *natmap) StartCopy(
	key NatKey,
	entry *natEntry,
	dst net.PacketConn,
	peer net.Addr,
	trace *mylog.Trace,
) {

	// 如果是 socks5，监听 ctrl 断开
	if entry.role == Socks5Server && entry.socksCtrl != nil {
		entry.watchSocksCtrl(func() {
			trace.Info("socks5 ctrl closed")

			if e := m.Del(key); e != nil {
				e.CloseAll()
			}
		})
	}

	go func() {
		timedCopy(dst, peer, entry.pc, m.timeout, entry.role, trace)

		if e := m.Del(key); e != nil {
			e.CloseAll()
		}
	}()
}

// ==============================
// Add NAT Entry
// ==============================

func (m *natmap) Add(
	clientAddr string,
	role Mode,
	peer net.Addr,
	dst net.PacketConn,
	src net.PacketConn,
	trace *mylog.Trace,
) *natEntry {

	key := NatKey{
		Client: clientAddr,
		Role:   role,
	}

	entry := &natEntry{
		pc:   src,
		role: role,
	}

	m.Set(key, entry)

	go func(tlog *mylog.Trace) {
		timedCopy(dst, peer, src, m.timeout, role, tlog)

		if e := m.Del(key); e != nil {
			e.CloseAll()
		}
	}(trace)

	return entry
}

// ==============================
// timedCopy
// ==============================

func timedCopy(
	dst net.PacketConn,
	target net.Addr,
	src net.PacketConn,
	timeout time.Duration,
	role Mode,
	trace *mylog.Trace,
) error {

	buf := make([]byte, udpBufSize)

	for {
		src.SetReadDeadline(time.Now().Add(timeout))

		n, raddr, err := src.ReadFrom(buf)
		if err != nil {
			return err
		}

		switch role {

		case TargetServer: // 目标 -> ss1 : add original packet source
			trace.Info("targetServer")
			srcAddr := mysocks5.ParseAddr(raddr.String())
			copy(buf[len(srcAddr):], buf[:n])
			copy(buf, srcAddr)

			_, err = dst.WriteTo(buf[:len(srcAddr)+n], target)

		case SSNodeServer: // client -> user : strip original packet source
			trace.Info("ssNodeServer")
			srcAddr := mysocks5.SplitAddr(buf[:n])
			_, err = dst.WriteTo(buf[len(srcAddr):n], target)

		case Socks5Server:
			// 至少 RSV(2) + FRAG(1) + ATYP(1)
			if n < 4 {
				trace.Info("socks5Server", "len err", buf[:n])
				continue
			}
			// 1️⃣ 检查 RSV
			if buf[0] != 0 || buf[1] != 0 {
				trace.Info("socks5Server", "rsv err", buf[:2])
				continue
			}

			// 2️⃣ 检查 FRAG
			if buf[2] != 0 {
				trace.Info("socks5Server", "flag err", buf[:3])
				// 不支持分片，直接丢弃
				continue
			}

			trace.Info("socks5Server")
			// 1️⃣ 去掉 RSV + FRAG
			packet := buf[3:n]

			// 2️⃣ 解析 socks5 addr
			addr := mysocks5.SplitAddr(packet)
			addrLen := len(addr)
			if addrLen == 0 || addrLen > len(packet) {
				continue
			}

			// 3️⃣ 取 payload
			payload := packet[addrLen:]

			// 4️⃣ ⚠️ 重新按你自己的协议格式封装
			// 和 TargetServer 保持一致
			reply := make([]byte, addrLen+len(payload))
			copy(reply, addr)
			copy(reply[addrLen:], payload)

			_, err = dst.WriteTo(reply, target)
		}

		if err != nil {
			return err
		}
	}
}

func buildSocks5UDPPacket(tgt net.Addr, payload []byte) []byte {
	udpAddr := tgt.(*net.UDPAddr)
	addr := mysocks5.ParseAddr(udpAddr.String())
	packet := make([]byte, 3+len(addr)+len(payload))

	// RSV
	packet[0] = 0
	packet[1] = 0

	// FRAG
	packet[2] = 0

	// ATYP + DST + PORT
	copy(packet[3:], addr)

	// DATA
	copy(packet[3+len(addr):], payload)
	return packet
}
