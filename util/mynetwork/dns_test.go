package mynetwork

import (
	"testing"
	"time"
)

func Test_DigNS(t *testing.T) {
	ns := DigNS("v5dns.xyz", time.Second)
	t.Log(ns)
}
