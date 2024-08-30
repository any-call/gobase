package myvalidator

import (
	"fmt"
	"testing"
)

func TestValidIPV6(t *testing.T) {
	b := ValidIPV6("fe80::ec5e:e3ff:febe:ceb3")
	t.Log("ipv6:", b)

	b = ValidFqdn("abc.mm.")
	t.Log("domain:", b)

	b = ValidSrv("1 1 8686 your-server.l.google.com")
	t.Log("srv :", b)

	b = ValidCAA("12 iodef www.testdns.com")
	t.Log("caa :", b)

	b = ValidSoa("dns.baidu.com sa.baidu.com. 2012146317 300 300 2592000 7200")
	t.Log("soa:", b)
}

func TestCAAInfo_Record(t *testing.T) {
	ret, err := ExtractMethodComments("base_test.go", "MyStruct")
	if err != nil {
		t.Error("extracted err", err)
		return
	}

	t.Log("extracted is :", ret)
}

type MyStruct struct{}

// Greet 定义带有参数的方法
func (m *MyStruct) Greet(name string, times int) {
	for i := 0; i < times; i++ {
		fmt.Printf("Hello, %s!\n", name)
	}
}

func (m *MyStruct) Multiply(a, b int) int {
	return a * b
}
