package mynetwork

import (
	"net"
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
