package mysocks5

import (
	"errors"
	"fmt"
	"github.com/any-call/gobase/frame/myctrl"
	"github.com/any-call/gobase/util/mylog"
	"io"
	"net"
	"strconv"
	"time"
)

const (
	MethodNoAuth = 0x00 //表示无需认证
	MethodAuth   = 0x02 //表示用户名/密码认证
)

// Handshake fast-tracks SOCKS initialization to get target address to connect.
func Handshake(rw io.ReadWriter, authFn func(username, password string) bool) (Addr, error) {
	// Read RFC 1928 section 4 for request and reply structure and sizes
	buf := make([]byte, MaxReqLen)

	n1, err := rw.Read(buf) // SOCKS version and account methods
	if err != nil {
		return nil, err
	}

	if n1 >= 3 {
		if buf[0] == Version5 { //说明支持socks5
			_, err = rw.Write([]byte{Version5, myctrl.ObjFun(func() byte {
				if authFn != nil {
					return MethodAuth
				}
				return MethodNoAuth
			})}) // SOCKS v5, no account required
			if err != nil {
				return nil, err
			}

			if authFn != nil {
				if !authenticate(rw, authFn) {
					return nil, fmt.Errorf("auth fail")
				}
			}

			n, err := rw.Read(buf) // SOCKS request: VER, CMD, RSV, Addr
			if err != nil {
				return nil, err
			}
			buf = buf[:n]
			if buf[1] == CmdConnect {
				_, err = rw.Write([]byte{Version5, ReplySuccess, Reserved, 1, 0, 0, 0, 0, 0, 0}) // SOCKS v5, reply succeeded
				return buf[3:], err                                                              // skip VER, CMD, RSV fields
			} else if buf[1] == CmdUDPAssociate { //UDP穿透
				return nil, ErrUdpAssociate
			}

			return nil, ErrCommandNotSupported
		} else if buf[0] == Version4 { //说明支持socks4 &socks4a
			if buf[1] != CmdConnect {
				return nil, ErrCommandNotSupported
			}

			_, err = rw.Write([]byte{0, 0x5a, 0, 0, 0, 0, 0, 0}) // SOCKS v4, no account required
			if err != nil {
				return nil, err
			}
			//已建立链接
			if n1 == 9 { //socks4
				rtn := make([]byte, 1+net.IPv4len+2)
				rtn[0] = ATypeIPv4
				copy(rtn[1:], buf[4:8])
				copy(rtn[5:], buf[2:4])
				return rtn, nil
			} else {
				return nil, ErrCommandNotSupported
			}
		}
	}

	return nil, ErrCommandNotSupported

}

//Step 1: 客户端发起握手
//+----+----------+----------+
//|VER | NMETHODS | METHODS  |
//+----+----------+----------+
//| 1  |    1     | 1 to 255 |
//+----+----------+----------+

//服务端响应
//+----+--------+
//|VER | METHOD |
//+----+--------+
//| 1  |   1    |
//+----+--------+

//Step 2: 客户端发起连接请求
//+----+-----+-------+------+----------+----------+
//|VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
//+----+-----+-------+------+----------+----------+
//| 1  |  1  | X'00' |  1   | Variable |    2     |
//+----+-----+-------+------+----------+----------+

//服务器响应
//+----+-----+-------+------+----------+----------+
//|VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
//+----+-----+-------+------+----------+----------+
//| 1  |  1  | X'00' |  1   | Variable |    2     |
//+----+-----+-------+------+----------+----------+

func ConnToSocks5(addr Addr, dialTimeoutSec int, remoteAddr string, authfn func() (userName, password string), dialCtrl myctrl.Golimiter) (net.Conn, error) {
	if dialTimeoutSec < 0 {
		dialTimeoutSec = 0
	}
	d := net.Dialer{
		Timeout: time.Duration(dialTimeoutSec) * time.Second, // 5 秒超时
	}

	var conn net.Conn
	var err error
	if dialCtrl != nil {
		dialCtrl.Begin()
		conn, err = d.Dial("tcp", remoteAddr)
		dialCtrl.End()
	} else {
		conn, err = d.Dial("tcp", remoteAddr)
	}

	if err != nil {
		return nil, fmt.Errorf("连接 SOCKS5 代理服务器失败:%v", err)
	}

	var username, password string
	if authfn != nil { //存在
		username, password = authfn()
	}

	//建立链接
	_, err = conn.Write(myctrl.ObjFun(func() []byte {
		if username == "" && password == "" {
			return []byte{5, 1, 0}
		}

		return []byte{5, 2, 0, 2}
	}))
	if err != nil {
		defer func() {
			_ = conn.Close()
		}()
		return nil, fmt.Errorf("发送握手请求失败:%v", err)
	}

	//读服务端响应
	response := make([]byte, 2)
	_, err = conn.Read(response) // SOCKS request: VER, CMD, RSV, Addr
	if err != nil {
		defer func() {
			_ = conn.Close()
		}()
		return nil, errors.New("读取握手响应失败")
	}

	if response[0] != 5 || (response[1] != 0 && response[1] != 2) {
		defer func() {
			_ = conn.Close()
		}()
		return nil, fmt.Errorf("SOCKS5 握手失败，代理服务器应答:%v", response)
	}

	if response[1] == 2 { //服务端需求用户名与密码认证
		if username == "" && password == "" {
			defer func() {
				_ = conn.Close()
			}()
			return nil, fmt.Errorf("SOCKS5 需求用户密码认证")
		}

		//mylog.Debug("用户名:", username, ";password:", password)
		// 发送用户名/密码认证信息
		auth := make([]byte, 3+len(username)+len(password))
		auth[0] = 0x01                              // 认证版本
		auth[1] = byte(len(username))               // 用户名长度
		copy(auth[2:], username)                    // 用户名
		auth[2+len(username)] = byte(len(password)) // 密码长度
		copy(auth[3+len(username):], password)      // 密码

		if _, err := conn.Write(auth); err != nil {
			defer func() {
				_ = conn.Close()
			}()
			return nil, fmt.Errorf("发送认证信息失败:%v", err)
		}

		// 读取认证结果
		authResponse := make([]byte, 2)
		if _, err := conn.Read(authResponse); err != nil {
			defer func() {
				_ = conn.Close()
			}()
			return nil, fmt.Errorf("读取认证响应失败:%v", err)
		}

		// 检查认证结果
		if authResponse[1] != 0x00 {
			defer func() {
				_ = conn.Close()
			}()
			return nil, errors.New("用户名/密码认证失败")
		}
	}

	buf := make([]byte, MaxReqLen)
	//然后跟目标服务建立连接
	copy(buf, []byte{5, 1, 0})
	copy(buf[3:], addr)
	_, err = conn.Write(buf[:3+len(addr)])
	if err != nil {
		defer func() {
			_ = conn.Close()
		}()
		return nil, fmt.Errorf("SOCKS5 建立链接时出错：%v", err)
	}

	n, err := conn.Read(buf)
	if err != nil {
		defer func() {
			_ = conn.Close()
		}()
		return nil, fmt.Errorf("SOCKS5 读取出错:%v to %s", err, addr.String())
	}

	if buf[0] != 5 || buf[1] != 0 {
		defer func() {
			_ = conn.Close()
		}()
		return nil, fmt.Errorf("SOCKS5 响应错误：%v to %s", buf[:n], addr.String())
	}

	return conn, nil
}

// 处理客户端认证请求
func authenticate(rw io.ReadWriter, validfn func(username, password string) bool) bool {
	buf := make([]byte, MaxAddrLen)
	n1, err := rw.Read(buf)
	if err != nil {
		mylog.Debug("authenticate read err:", err, n1)
		return false
	}

	//mylog.Info("authenticate read data :", buf[:n1])
	if n1 < 4 { //长度不足
		if n1 == 3 && buf[0] == 0x01 && buf[1] == 0x00 && buf[2] == 0x00 { //没有用户/密码
			if validfn == nil {
				_, _ = rw.Write([]byte{0x01, 0x00})
				return true
			}

			if validfn("", "") {
				_, _ = rw.Write([]byte{0x01, 0x00})
				return true
			} else {
				_, _ = rw.Write([]byte{0x01, 0x01})
				return false
			}
		}

		return false
	}

	if buf[0] != 0x01 { //这属于 SOCKS5 协议内的一个子协议（用于处理用户名/密码的认证机制）。
		return false
	}

	//用户名长度
	uLen := int(buf[1])
	if n1 < (uLen + 3) {
		return false //用户名长度不足
	}
	pLen := int(buf[uLen+2])
	if n1 < (uLen + pLen + 3) {
		return false //密码不足
	}

	uname := string(buf[2 : uLen+2])
	passwd := string(buf[uLen+3 : uLen+3+pLen])

	if validfn == nil || validfn(string(uname), string(passwd)) {
		_, _ = rw.Write([]byte{0x01, 0x00}) // 验证成功
		return true
	}

	_, _ = rw.Write([]byte{0x01, 0x01}) // 验证失败
	return false
}

//客户端发起认证请求
//+----+------+----------+----------+
//|VER | ULEN |  UNAME   | PLEN     |
//+----+------+----------+----------+
//| 1  |  1   | 1 to 255 | 1 to 255 |
//+----+------+----------+----------+
//VER：认证协议的版本号
//ULEN：用户名长度（1字节）。
//UNAME：用户名（ULEN 长度）。
//PLEN：密码长度（1字节）。
//PASSWD：密码（PLEN 长度）。

//服务端响应认证
//+----+--------+
//|VER | STATUS |
//+----+--------+
//| 1  |   1    |
//+----+--------+
//STATUS：
//•	0x00 表示成功。
//•	0x01 表示失败。

func HandleUdpAssociate(conn net.Conn, LitenUdpAddr string) error {
	host, portStr, _ := net.SplitHostPort(LitenUdpAddr)
	pport, _ := strconv.Atoi(portStr)
	ip := net.ParseIP(host).To4()
	if ip == nil {
		return fmt.Errorf("invalid ip address from:%s", LitenUdpAddr)
	}
	resp := []byte{
		0x05, // VER
		0x00, // REP = succeeded
		0x00, // RSV
		0x01, // ATYP = IPv4
		ip[0], ip[1], ip[2], ip[3],
		byte(pport >> 8), byte(pport),
	}

	_, err := conn.Write(resp)
	if err != nil {
		return err
	}
	return nil
}
