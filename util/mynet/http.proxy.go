package mynet

import (
	"errors"
	"fmt"
	"github.com/any-call/gobase/util/mylog"
	"github.com/any-call/gobase/util/mysocks5"
	"net"
	"net/http"
	"strings"
	"time"
)

type httpProxyUtil struct {
}

func ImpHttpProxy() httpProxyUtil {
	return httpProxyUtil{}
}

func (self httpProxyUtil) GetTargetAddr(req *http.Request) string {
	var ret string
	if req.Method == http.MethodConnect {
		ret = req.Host
	} else {
		ret = req.URL.Host
		if !strings.Contains(ret, ":") {
			// 默认 HTTP 端口
			ret += ":80"
		}
	}

	return ret
}

func (self httpProxyUtil) HandleHttpProxy(w http.ResponseWriter, r *http.Request, specLocalIp string) (int64, int64, error) {
	targetAddr := self.GetTargetAddr(r)
	var dstConn net.Conn
	var err error
	// 建立与目标服务器的 TCP 连接
	if specLocalIp != "" {
		dialer := &net.Dialer{
			LocalAddr: &net.TCPAddr{
				IP: net.ParseIP(specLocalIp), // 本地 IP
			},
		}
		dstConn, err = dialer.Dial("tcp", targetAddr)
	} else {
		dstConn, err = net.Dial("tcp", targetAddr)
	}

	if err != nil {
		http.Error(w, "Unable to connect to target server", http.StatusServiceUnavailable)
		return 0, 0, fmt.Errorf("Unable to connect to target server : %v", err)
	}

	defer func() {
		_ = dstConn.Close()
	}()

	// 将客户端请求写入目标连接
	err = r.Write(dstConn)
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to write request to target: %v", err)
	}

	// 将目标服务器的响应返回给客户端
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return 0, 0, fmt.Errorf("Hijacking not supported")
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to hijack connection: %v", err)
	}
	defer func() {
		_ = clientConn.Close()
	}()

	//// 双向数据转发
	return Relay(clientConn, dstConn)
}

func (self httpProxyUtil) HandleHttpsProxy(w http.ResponseWriter, r *http.Request, specLocalIp string) (int64, int64, error) {
	targetAddr := self.GetTargetAddr(r)
	var dstConn net.Conn
	var err error
	// 建立与目标服务器的 TCP 连接
	if specLocalIp != "" {
		dialer := &net.Dialer{
			LocalAddr: &net.TCPAddr{
				IP: net.ParseIP(specLocalIp), // 本地 IP
			},
		}
		dstConn, err = dialer.Dial("tcp", targetAddr)
	} else {
		dstConn, err = net.Dial("tcp", targetAddr)
	}

	if err != nil {
		http.Error(w, "Unable to connect to target server", http.StatusServiceUnavailable)
		return 0, 0, fmt.Errorf("Unable to connect to target server:%v", err)
	}
	defer func() {
		_ = dstConn.Close()
	}()

	// 通知客户端隧道已建立
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return 0, 0, errors.New("Hijacking not supported")
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to hijack connection: %v", err)
	}
	defer func() {
		_ = clientConn.Close()
	}()

	if _, err = clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n")); err != nil {
		mylog.Debugf("write err:: %v", err)
	}

	// 双向数据转发
	return Relay(clientConn, dstConn)
}
func (self httpProxyUtil) HandleHttpsProxyWithTimeout(w http.ResponseWriter, r *http.Request, specLocalIp string, timeout time.Duration) (int64, int64, error) {
	targetAddr := self.GetTargetAddr(r)
	var dstConn net.Conn
	var err error
	// 建立与目标服务器的 TCP 连接
	if specLocalIp != "" {
		dialer := &net.Dialer{
			LocalAddr: &net.TCPAddr{
				IP: net.ParseIP(specLocalIp), // 本地 IP
			},
		}
		dstConn, err = dialer.Dial("tcp", targetAddr)
	} else {
		dstConn, err = net.Dial("tcp", targetAddr)
	}

	if err != nil {
		http.Error(w, "Unable to connect to target server", http.StatusServiceUnavailable)
		return 0, 0, fmt.Errorf("Unable to connect to target server:%v", err)
	}
	defer func() {
		_ = dstConn.Close()
	}()

	// 通知客户端隧道已建立
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return 0, 0, errors.New("Hijacking not supported")
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to hijack connection: %v", err)
	}
	defer func() {
		_ = clientConn.Close()
	}()

	if _, err = clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n")); err != nil {
		mylog.Debugf("write err:: %v", err)
	}

	// 双向数据转发
	return RelayWithTimeout(clientConn, dstConn, timeout)
}
func (self httpProxyUtil) HandleSocks5Proxy(w http.ResponseWriter, r *http.Request, socks5SrvAddr, socks5Username, socksPwd string) (int64, int64, error) {
	targetStr := self.GetTargetAddr(r)
	targetAddr := mysocks5.ParseAddr(targetStr)
	if targetAddr == nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return 0, 0, fmt.Errorf("invalid target : %s", targetStr)
	}

	var dstConn net.Conn
	var err error
	dstConn, err = mysocks5.ConnToSocks5(targetAddr, socks5SrvAddr, func() (userName, password string) {
		return socks5Username, socksPwd
	})

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return 0, 0, fmt.Errorf("connect socks5 err: %v", err)
	}

	defer func() {
		_ = dstConn.Close()
	}()

	if r.Method != http.MethodConnect { //http enter
		// http 将客户端请求写入目标连接
		if err = r.Write(dstConn); err != nil {
			return 0, 0, fmt.Errorf("Failed to write request to target: %v", err)
		}
	}

	// 将目标服务器的响应返回给客户端
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return 0, 0, errors.New("Hijacking not supported")
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to hijack connection: %v", err)
	}
	defer func() {
		_ = clientConn.Close()
	}()

	if r.Method == http.MethodConnect {
		_, _ = clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	}

	// 双向数据转发
	return Relay(clientConn, dstConn)
}

func (self httpProxyUtil) HandleSocks5ProxyWithTimeout(w http.ResponseWriter, r *http.Request,
	socks5SrvAddr, socks5Username, socksPwd string, timeout time.Duration) (int64, int64, error) {
	targetStr := self.GetTargetAddr(r)
	targetAddr := mysocks5.ParseAddr(targetStr)
	if targetAddr == nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return 0, 0, fmt.Errorf("invalid target : %s", targetStr)
	}

	var dstConn net.Conn
	var err error
	dstConn, err = mysocks5.ConnToSocks5(targetAddr, socks5SrvAddr, func() (userName, password string) {
		return socks5Username, socksPwd
	})

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return 0, 0, fmt.Errorf("connect socks5 err: %v", err)
	}

	defer func() {
		_ = dstConn.Close()
	}()

	if r.Method != http.MethodConnect { //http enter
		// http 将客户端请求写入目标连接
		if err = r.Write(dstConn); err != nil {
			return 0, 0, fmt.Errorf("Failed to write request to target: %v", err)
		}
	}

	// 将目标服务器的响应返回给客户端
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return 0, 0, errors.New("Hijacking not supported")
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to hijack connection: %v", err)
	}
	defer func() {
		_ = clientConn.Close()
	}()

	if r.Method == http.MethodConnect {
		_, _ = clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	}

	// 双向数据转发
	return RelayWithTimeout(clientConn, dstConn, timeout)
}
