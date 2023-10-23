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

	//if ns, err := LookupAddrEx("39.156.66.10", time.Millisecond*1000); err != nil {
	//	t.Error(err)
	//	return
	//} else {
	//	t.Log("addr is :", ns)
	//}

	if ns, err := LookupIPEx("baidu.com", time.Millisecond*100); err != nil {
		t.Error(err)
		return
	} else {
		t.Log("baidu.com ip is :", ns)
	}

	if ns, err := LookupTXTEx("baidu.com", time.Millisecond*100); err != nil {
		t.Error(err)
		return
	} else {
		t.Log("txt is :", ns)
	}

	if ns, err := LookupCNameEx("baidu.com", time.Millisecond*100); err != nil {
		t.Error(err)
		return
	} else {
		t.Log("cname is :", ns)
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

func Test_lookNS(t *testing.T) {
	list, err := LookupNSWithSer("rwscode.com", time.Second, "114.114.114.114")
	if err != nil {
		t.Error("look ns err:", err)
		return
	}

	t.Log("look ns :", list)
}
