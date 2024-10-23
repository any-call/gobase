package mysocks5

import (
	"io"
	"net"
	"strconv"
	"strings"
)

type (
	Addr []byte
)

func (a Addr) String() string {
	var host, port string

	switch a[0] { // address type
	case ATypeDomain:
		host = string(a[2 : 2+a[1]])
		port = strconv.Itoa((int(a[2+a[1]]) << 8) | int(a[2+a[1]+1]))
	case ATypeIPv4:
		host = net.IP(a[1 : 1+net.IPv4len]).String()
		port = strconv.Itoa((int(a[1+net.IPv4len]) << 8) | int(a[1+net.IPv4len+1]))
	case ATypeIPV6:
		host = net.IP(a[1 : 1+net.IPv6len]).String()
		port = strconv.Itoa((int(a[1+net.IPv6len]) << 8) | int(a[1+net.IPv6len+1]))
	}

	return net.JoinHostPort(host, port)
}

func (a Addr) AddType() (string, int) {
	hostStr, _, _ := net.SplitHostPort(a.String())
	ip := net.ParseIP(hostStr)
	if ip == nil {
		return hostStr, ATypeDomain
	}

	if strings.Contains(hostStr, ".") {
		return hostStr, ATypeIPv4
	}

	if strings.Contains(hostStr, ":") {
		return hostStr, ATypeIPV6
	}

	return hostStr, ATypeDomain
}

// ReadAddr reads just enough bytes from r to get a valid Addr.
func ReadAddr(r io.Reader) (Addr, int, error) {
	b := make([]byte, MaxAddrLen)
	_, err := io.ReadFull(r, b[:1]) // read 1st byte for address type
	if err != nil {
		return nil, 0, err
	}

	switch b[0] {
	case ATypeDomain:
		_, err = io.ReadFull(r, b[1:2]) // read 2nd byte for domain length
		if err != nil {
			return nil, 0, err
		}
		_, err = io.ReadFull(r, b[2:2+b[1]+2])
		newBuff := b[:1+1+b[1]+2]
		return newBuff, int(b[0]), err
	case ATypeIPv4:
		_, err = io.ReadFull(r, b[1:1+net.IPv4len+2])
		return b[:1+net.IPv4len+2], int(b[0]), err
	case ATypeIPV6:
		_, err = io.ReadFull(r, b[1:1+net.IPv6len+2])
		return b[:1+net.IPv6len+2], int(b[0]), err
	}

	return nil, 0, ErrAddressNotSupported
}

// SplitAddr slices a SOCKS address from beginning of b. Returns nil if failed.
func SplitAddr(b []byte) Addr {
	addrLen := 1
	if len(b) < addrLen {
		return nil
	}

	switch b[0] {
	case ATypeDomain:
		if len(b) < 2 {
			return nil
		}
		addrLen = 1 + 1 + int(b[1]) + 2
	case ATypeIPv4:
		addrLen = 1 + net.IPv4len + 2
	case ATypeIPV6:
		addrLen = 1 + net.IPv6len + 2
	default:
		return nil

	}

	if len(b) < addrLen {
		return nil
	}

	return b[:addrLen]
}

// ParseAddr parses the address in string s. Returns nil if failed.
func ParseAddr(s string) Addr {
	var addr Addr
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return nil
	}
	if ip := net.ParseIP(host); ip != nil {
		if ip4 := ip.To4(); ip4 != nil {
			addr = make([]byte, 1+net.IPv4len+2)
			addr[0] = ATypeIPv4
			copy(addr[1:], ip4)
		} else {
			addr = make([]byte, 1+net.IPv6len+2)
			addr[0] = ATypeIPV6
			copy(addr[1:], ip)
		}
	} else {
		if len(host) > 255 {
			return nil
		}
		addr = make([]byte, 1+1+len(host)+2)
		addr[0] = ATypeDomain
		addr[1] = byte(len(host))
		copy(addr[2:], host)
	}

	portnum, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return nil
	}

	addr[len(addr)-2], addr[len(addr)-1] = byte(portnum>>8), byte(portnum)

	return addr
}
