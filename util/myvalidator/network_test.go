package myvalidator

import "testing"

func Test_validIPCIDR(t *testing.T) {
	test1 := "192.168.1.2"
	b1 := ValidIPCIDR(test1)
	t.Logf("valid ipcidr :%s result :%v", test1, b1)

	test2 := "192.168.1.2/12"
	b2 := ValidIPCIDR(test2)
	t.Logf("valid ipcidr :%s result :%v", test2, b2)

	testIP3 := "192.168.1.12"
	testIPCidr3 := "192.168.1.12/32"
	b3 := ValidIPBelongIPCidr(testIP3, testIPCidr3)
	t.Logf("valid ip  :%s belong to  :%v is :%v", testIP3, testIPCidr3, b3)
}
