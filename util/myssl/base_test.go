package myssl

import (
	"fmt"
	"testing"
)

func TestSSL(t *testing.T) {
	domain := "baidu.com"
	cert, err := GetSSLCert(domain)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("ca name:", GetDomain(cert))
	fmt.Println("ca days:", GetRemainDays(cert))
	fmt.Println("ca all name:", GetDnsNames(cert))
	fmt.Println("ca type:", GetType(cert))
	fmt.Println("ca method:", GetMethod(cert))
	fmt.Println("ca sha1:", GetSHA1(cert))
	fmt.Println("ca sha256:", GetSHA256(cert))
	fmt.Println("ca issuer:", GetIssuer(cert))
}
