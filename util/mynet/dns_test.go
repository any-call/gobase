package mynet

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"
)

func Test_DigNS(t *testing.T) {
	ns := DigNS("v5dns.xyz", time.Second)
	t.Log(ns)
}

func Test_TT(t *testing.T) {
	dns := "110.242.68.66"
	ns, err := net.LookupAddr(dns)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("ns:", ns)
}

func Test_other(t *testing.T) {
	domain := "baidu.com" // 替换为要查询的域名
	resp, err := http.Head("http://" + domain)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	serverType := resp.Header.Get("Server")
	fmt.Println("Server type:", serverType)
}
