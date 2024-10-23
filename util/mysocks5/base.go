package mysocks5

const (
	Version5     byte = 0x05 //SOCKS5 的版本号，固定为 0x05。
	Version4     byte = 0x04 //SOCKS5 的版本号，固定为 0x05。
	Reserved     byte = 0x00 //保留字段，固定为 0x00。
	ReplySuccess byte = 0x00 //成功应答
)

// socks5 request commands
const (
	CmdConnect      = 0x01 //客户端通过代理服务器与目标服务器建立一个 TCP 连接
	CmdBind         = 0x02 //用于服务器模式下的连接。这个命令允许客户端告诉代理服务器监听特定端口
	CmdUDPAssociate = 0x03 //此命令用于通过 SOCKS5 代理发送和接收 UDP 数据报
)

const (
	ATypeIPv4   = 0x01 //IPv4 地址（4 字节）
	ATypeDomain = 0x03 //域名（第一个字节是域名长度，后跟域名字符串）
	ATypeIPV6   = 0x04 //IPv6 地址（16 字节）
)

const (
	MaxAddrLen = 1 + 1 + 255 + 2
	MaxReqLen  = 1 + 1 + 1 + MaxAddrLen
)

//请求格式 说明
// +--------+----------+-----------+---------------+---------------------+----------------+
// |version | command  | reserved  |  address type | Destination Address | Destination Port
// +--------+----------+-----------+---------------+---------------------+----------------+
// | 1 byte |  1 byte  |   1 byte  |    1 byte     | Destination Address |     2 byte
// +--------+----------+-----------+---------------+---------------------+----------------+
//	0x01：IPv4 地址（4 字节）。
//	0x03：域名（第一个字节是域名长度，后跟域名字符串）。
//	0x04：IPv6 地址（16 字节）。

//响应格式
// +--------+---------+
// |version |   Reply
// +--------+---------+
// | 1 byte |  1 byte
// +--------+---------+
//	Version (1 byte)：SOCKS5 版本号，固定为 0x05。
//	Reply (1 byte)：响应状态码：
//	0x00：成功。
