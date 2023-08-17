package myconv

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func IPV42Long(ip string) uint32 {
	var long uint32
	parts := strings.Split(ip, ".")
	for i := range parts {
		parts[i] = strings.TrimLeft(parts[i], "0")
		if parts[i] == "" {
			parts[i] = "0"
		}
	}
	newIP := strings.Join(parts, ".")
	binary.Read(bytes.NewBuffer(net.ParseIP(newIP).To4()), binary.BigEndian, &long)
	return long
}

func Long2IPV4(ipInt int64) string {
	b0 := strconv.FormatInt((ipInt>>24)&0xff, 10)
	b1 := strconv.FormatInt((ipInt>>16)&0xff, 10)
	b2 := strconv.FormatInt((ipInt>>8)&0xff, 10)
	b3 := strconv.FormatInt(ipInt&0xff, 10)
	return fmt.Sprintf("%s.%s.%s.%s", b0, b1, b2, b3)
}
