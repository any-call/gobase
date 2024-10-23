package mysocks5

import "fmt"

// SOCKS errors as defined in RFC 1928 section 6.
const (
	ErrGeneralFailure       = Error(1) //General SOCKS server failure
	ErrConnectionNotAllowed = Error(2) //Connection not allowed by ruleset
	ErrNetworkUnreachable   = Error(3) //Network unreachable
	ErrHostUnreachable      = Error(4) //Host unreachable
	ErrConnectionRefused    = Error(5) //Connection refused
	ErrTTLExpired           = Error(6) //TTL expired
	ErrCommandNotSupported  = Error(7) //Command not supported
	ErrAddressNotSupported  = Error(8) //Address type not supported
	ErrUdpAssociate         = Error(9) //udp associate
)

type Error byte

func (self Error) Error() string {
	switch int(self) {
	case 1:
		return fmt.Sprintf("[%d]:General SOCKS server failure", self)
	case 2:
		return fmt.Sprintf("[%d]:Connection not allowed by ruleset", self)
	case 3:
		return fmt.Sprintf("[%d]:Network unreachable", self)
	case 4:
		return fmt.Sprintf("[%d]:Host unreachable", self)
	case 5:
		return fmt.Sprintf("[%d]:Connection refused", self)
	case 6:
		return fmt.Sprintf("[%d]:TTL expired", self)
	case 7:
		return fmt.Sprintf("[%d]:Command not supported", self)
	case 8:
		return fmt.Sprintf("[%d]:Address type not supported", self)
	case 9:
		return fmt.Sprintf("[%d]:udp associate", self)
	}

	return fmt.Sprintf("unknow err:%d", self)
}
