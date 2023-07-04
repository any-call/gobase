package mynet

import (
	"net"
	"net/http"
)

//A类：10段，后三位自由分配，也就是 10.0.0.0 - 10.255.255.255；
//B类：172.16段，后两位自由分配，也就是 172.16.0.0 - 172.31.255.255；
//C类：192.168段，后两位自由分配，也就是 192.168.0.0 - 192.168.255.255；

// RemoteIP 通过 RemoteAddr 获取 IP 地址， 只是一个快速解析方法。
func RemoteIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

func IsLocalIP(ip string) bool { return !IsPublicIP(net.ParseIP(ip)) }
