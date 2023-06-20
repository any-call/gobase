package myconv

import (
	"net"
	"testing"
)

func TestIPV42Long(t *testing.T) {
	cidr := "192.168.0.0/24"
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		t.Error("Error parsing CIDR:", err)
		return
	}

	var startIP net.IP
	var endIP net.IP

	for i, b := range ip {
		startIP = startIP[:i]
		endIP = endIP[:i]
		startIP = append(startIP, b&ipnet.Mask[i]|ipnet.IP[i]&^ipnet.Mask[i])
		endIP = append(endIP, b&ipnet.Mask[i]|ipnet.IP[i]&^ipnet.Mask[i]|^ipnet.Mask[i])
	}

	t.Log("Start IP:", startIP.String())
	t.Log("End IP:", endIP.String())
}
