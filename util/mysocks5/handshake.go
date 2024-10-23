package mysocks5

import (
	"io"
	"net"
)

const (
	MethodNoAuth = 0x00 //表示无需认证
	MethodAuth   = 0x02 //表示用户名/密码认证
)

// Handshake fast-tracks SOCKS initialization to get target address to connect.
func Handshake(rw io.ReadWriter) (Addr, error) {
	// Read RFC 1928 section 4 for request and reply structure and sizes
	buf := make([]byte, MaxReqLen)

	n1, err := rw.Read(buf) // SOCKS version and account methods
	if err != nil {
		return nil, err
	}

	if n1 >= 3 {
		if buf[0] == Version5 { //说明支持socks5
			_, err = rw.Write([]byte{Version5, MethodNoAuth}) // SOCKS v5, no account required
			if err != nil {
				return nil, err
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
