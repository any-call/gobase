package myvalidator

import "testing"

func TestValidIPV6(t *testing.T) {
	b := ValidIPV6("fe80::ec5e:e3ff:febe:ceb3")
	t.Log("ipv6:", b)

	b = ValidFqdn("abc.mm.")
	t.Log("domain:", b)
}
