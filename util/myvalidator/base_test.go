package myvalidator

import "testing"

func TestValidIPV6(t *testing.T) {
	b := ValidIPV6("fe80::ec5e:e3ff:febe:ceb3")
	t.Log("ipv6:", b)

	b = ValidFqdn("abc.mm.")
	t.Log("domain:", b)

	b = ValidSrv("1 1 8686 your-server.l.google.com")
	t.Log("srv :", b)

	b = ValidCAA("12 iodef www.testdns.com")
	t.Log("caa :", b)
}
