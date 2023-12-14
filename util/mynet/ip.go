package mynet

import (
	"fmt"
	"gitee.com/any-call/gobase/util/myos"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
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

func GetLocalIP() (net.IP, error) {
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
						//return dnsServers, nil
					}
				}
			}
		}
	}

	return nil, fmt.Errorf("unknown err")
}

func GetLocalDNSServer() ([]string, error) {
	//https://juejin.cn/s/%E5%A6%82%E4%BD%95%E8%8E%B7%E5%8F%96%E6%9C%AC%E6%9C%BAIp%E7%BD%91%E5%85%B3DNs
	if myos.IsMac() {
		//准备参数
		//var output []byte
		//cmd := exec.CommandContext(context.Background(), "scutil", "--dns")
		//output, _ = cmd.CombinedOutput()
		//
		//f := strings.NewReader(string(output))
		//reader := bufio.NewReader(f)
		//compileRegex := regexp.MustCompile("\\S+")
		//dns := make([]string, 0)
		//for {
		//	line, _, err1 := reader.ReadLine()
		//	if err1 != nil {
		//		break
		//	}
		//	strTmp := strings.Trim(string(line))
		//	if strings.HasPrefix(strTmp, "nameserver[") {
		//
		//	}
		//
		//	matchArr := compileRegex.FindAllStringSubmatch(strTmp, -1)
		//	if len(matchArr) == 5 {
		//		if newDomain == matchArr[0][0] && matchArr[3][0] == "NS" {
		//			nss = append(nss, strings.TrimSuffix(matchArr[4][0], "."))
		//		}
		//	}
		//}
	}

	return nil, nil
}

func GetPublicIP() (ip string, err error) {
	ip, err = ipinfo()
	if err != nil {
		ip, err = ifconfig()
		if err != nil {
			ip, err = ipify()
		}
	}
	return
}

func ipinfo() (ip string, err error) {
	resp, err := http.Get("https://ipinfo.io/ip")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	sb, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	} else {
		ip = string(sb)
		_, err = checkIP(ip)
	}
	return
}

func ifconfig() (ip string, err error) {
	resp, err := http.Get("https://ifconfig.me/ip")
	if err != nil {
		return
	}

	defer resp.Body.Close()

	sb, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	} else {
		ip = string(sb)
		_, err = checkIP(ip)
	}
	return
}

func ipify() (ip string, err error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	sb, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	} else {
		ip = string(sb)
		_, err = checkIP(ip)
	}
	return
}

var shiftIndex = []int{24, 16, 8, 0}

func checkIP(ip string) (uint32, error) {
	var ps = strings.Split(ip, ".")
	if len(ps) != 4 {
		return 0, fmt.Errorf("invalid ip address `%s`", ip)
	}

	var val = uint32(0)
	for i, s := range ps {
		d, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("the %dth part `%s` is not an integer", i, s)
		}

		if d < 0 || d > 255 {
			return 0, fmt.Errorf("the %dth part `%s` should be an integer bettween 0 and 255", i, s)
		}

		val |= uint32(d) << shiftIndex[i]
	}

	// convert the ip to integer
	return val, nil
}
