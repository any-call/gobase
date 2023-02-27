package mynetwork

import (
	"bufio"
	"context"
	"net"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func LookupNS(domain string) (ns []string, err error) {
	if nameServer, err := net.LookupNS(domain); err != nil {
		return nil, err
	} else {
		ns = make([]string, len(nameServer))
		for i, v := range nameServer {
			ns[i] = strings.TrimSuffix(v.Host, ".")
		}
	}
	return
}

func LookupIP(domain string) (ipRec []net.IP, err error) {
	if iprecords, err := net.LookupIP(domain); err != nil {
		return nil, err
	} else {
		ipRec = make([]net.IP, len(iprecords))
		for i, v := range iprecords {
			ipRec[i] = v
		}
	}
	return
}

/*
*PTR记录,是电子邮件系统中的邮件交换记录的一种;另一种邮件交换记录是A记录
（在IPv4协议中）或AAAA记录（在IPv6协议中）.PTR记录常被用于反向地址解析.
根据一个IP值,查找映射的域名值,不一定没有ip地址都回生效,DNS的IP地址可以查到.
*/

func LookupAddr(addr string) (domains []string, err error) {
	if nameServer, err := net.LookupAddr(addr); err != nil {
		return nil, err
	} else {
		domains = make([]string, len(nameServer))
		for i, v := range nameServer {
			domains[i] = v
		}
	}
	return
}

func LookupCName(domain string) (cname string, err error) {
	return net.LookupCNAME(domain)
}

func LookupTXT(domain string) (txt []string, err error) {
	return net.LookupTXT(domain)
}

func DigWithNS(domain string, timeout time.Duration) (nss []string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	//准备参数
	var output []byte
	cmd := exec.CommandContext(ctx, "dig", "ns", domain)
	output, _ = cmd.CombinedOutput()

	f := strings.NewReader(string(output))
	reader := bufio.NewReader(f)
	compileRegex := regexp.MustCompile("\\S+")
	newDomain := domain
	if strings.HasSuffix(domain, ".") == false {
		newDomain = domain + "."
	}

	nss = make([]string, 0)
	for {
		line, _, err1 := reader.ReadLine()
		if err1 != nil {
			break
		}
		strTmp := string(line)
		if strings.Trim(strTmp, " ") == "" {
			continue
		}

		if strings.HasPrefix(strTmp, ";; Received") {
			if len(nss) > 0 {
				return
			}
			continue
		}

		matchArr := compileRegex.FindAllStringSubmatch(strTmp, -1)
		if len(matchArr) == 5 {
			if newDomain == matchArr[0][0] && matchArr[3][0] == "NS" {
				nss = append(nss, strings.TrimSuffix(matchArr[4][0], "."))
			}
		}
	}

	return
}
