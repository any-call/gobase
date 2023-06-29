package myssl

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"time"
)

func GetSSLCert(domain string) (*x509.Certificate, error) {
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:443", domain), &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	return conn.ConnectionState().PeerCertificates[0], nil
}

func GetStartTime(cert *x509.Certificate) time.Time {
	if cert == nil {
		return time.Time{}
	}

	return cert.NotBefore
}

func GetEndTime(cert *x509.Certificate) time.Time {
	if cert == nil {
		return time.Time{}
	}

	return cert.NotAfter
}

func GetRemainDays(cert *x509.Certificate) int {
	if cert == nil {
		return 0
	}

	return int(cert.NotAfter.Sub(time.Now()).Hours() / 24)
}

func GetMethod(cert *x509.Certificate) string {
	if cert == nil {
		return ""
	}

	return cert.SignatureAlgorithm.String()
}

func GetType(cert *x509.Certificate) string {
	if cert == nil {
		return ""
	}

	//判定证书类型
	if len(cert.Subject.Organization) > 0 {
		if len(cert.Subject.Organization[0]) > 0 {
			return "OV"
		}
	}

	for _, usage := range cert.ExtKeyUsage {
		if usage == x509.ExtKeyUsageServerAuth || usage == x509.ExtKeyUsageClientAuth {
			return "EV"
		}
	}

	return "DV"
}

func GetSHA1(cert *x509.Certificate) string {
	if cert == nil {
		return ""
	}

	sha1Fingerprint := sha1.Sum(cert.Raw)
	return hex.EncodeToString(sha1Fingerprint[:])
}

func GetSHA256(cert *x509.Certificate) string {
	if cert == nil {
		return ""
	}

	sha256Fingerprint := sha256.Sum256(cert.Raw)
	return hex.EncodeToString(sha256Fingerprint[:])
}

func GetDomain(cert *x509.Certificate) string {
	if cert == nil {
		return ""
	}

	return cert.Subject.CommonName
}

func GetDnsNames(cert *x509.Certificate) []string {
	if cert == nil {
		return nil
	}

	return cert.DNSNames
}

func GetIssuer(cert *x509.Certificate) string {
	if cert == nil {
		return ""
	}

	return cert.Issuer.CommonName
}
