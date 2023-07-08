package mynet

import (
	"fmt"
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

func GetLocalDnsIP() (net.IP, error) {
	// 获取本机的网络接口列表
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	// 遍历每个网络接口
	for _, iface := range interfaces {
		// 排除非活动接口和回环接口
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			// 获取接口的地址列表
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}
			// 遍历每个地址
			for _, addr := range addrs {
				// 检查地址是否为IP地址
				ip, ok := addr.(*net.IPNet)
				if ok && !ip.IP.IsLoopback() && ip.IP.To4() != nil {
					// 获取DNS服务器地址
					dnsServers := net.ParseIP(ip.IP.String()).To4()
					if dnsServers != nil {
						return dnsServers, nil
					}
				}
			}
		}
	}

	return nil, fmt.Errorf("unknown err")
}
