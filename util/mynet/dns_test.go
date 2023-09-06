package mynet

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func Test_DigNS(t *testing.T) {
	ns := DigNS("v5dns.xyz", time.Second)
	t.Log(ns)
}

func Test_TT(t *testing.T) {
	if ns, err := LookupNSEx("baidu.com", time.Millisecond*100); err != nil {
		t.Error(err)
		return
	} else {
		t.Log("ns is :", ns)
	}
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
