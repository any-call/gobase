package myvalidator

import (
	"net"
)

func ValidIPV4(ipv4 string) bool {
	address := net.ParseIP(ipv4)
	if address == nil {
		return false
	}
	if address.To4() == nil {
		return false
	}

	return true
}
