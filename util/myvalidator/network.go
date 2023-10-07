package myvalidator

import (
	"net"
	"regexp"
	"strings"
)

func ValidIP(ip string) bool {
	address := net.ParseIP(ip)
	if address == nil {
		return false
	}

	if address.To4() != nil ||
		address.To16() != nil {
		return true
	}

	return false
}

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

func ValidIPV6(ipv6 string) bool {
	address := net.ParseIP(ipv6)
	if address == nil {
		return false
	}
	if address.To16() == nil {
		return false
	}

	return true
}

// 完整FQDN（Fully Qualified Domain Name）是完全合格的域名，
// 它是一个完整的域名，包括完整的子域名和主域名部分。
// 例如，www.example.com是一个FQDN，其中“www”是子域名，“example”是主域名，“com”是顶级域名。
func ValidFqdn(domain string) bool {
	if !strings.HasSuffix(domain, ".") {
		return false
	} else {
		domain = strings.TrimSuffix(domain, ".")
	}

	pattern := `^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6}$`
	// 使用正则表达式进行匹配验证
	match, _ := regexp.MatchString(pattern, domain)
	return match
}

// 验证域名
func ValidDomain(domain string) bool {
	pattern := `^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6}$`
	// 使用正则表达式进行匹配验证
	match, _ := regexp.MatchString(pattern, domain)
	return match
}

func ValidIPCIDR(ipcidr string) bool {
	if _, _, err := net.ParseCIDR(ipcidr); err != nil {
		return false
	}

	return true
}

func ValidIPBelongIPCidr(ipinfo, ipcidr string) bool {
	ip := net.ParseIP(ipinfo)
	if ip == nil {
		return false
	}

	_, cidr, _ := net.ParseCIDR(ipcidr)
	if cidr == nil {
		return false
	}

	if cidr.Contains(ip) {
		return true
	}

	return false
}
