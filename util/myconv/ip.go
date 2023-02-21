package myconv

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

func IPV42Long(ip string) uint32 {
	var long uint32
	binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &long)
	return long
}

func Long2IPV4(ipInt int64) string {
	b0 := strconv.FormatInt((ipInt>>24)&0xff, 10)
	b1 := strconv.FormatInt((ipInt>>16)&0xff, 10)
	b2 := strconv.FormatInt((ipInt>>8)&0xff, 10)
	b3 := strconv.FormatInt(ipInt&0xff, 10)
	return fmt.Sprintf("%s.%s.%s.%s", b0, b1, b2, b3)
}
