package mycert

import (
	"crypto/tls"
	"fmt"
	"time"
)

func CheckCertExpiry(domain string, port int) (time.Time, error) {
	conn, err := tls.Dial("tcp", domain+fmt.Sprintf(":%d", port), nil)
	if err != nil {
		return time.Time{}, err
	}
	defer func() {
		_ = conn.Close()
	}()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return time.Time{}, fmt.Errorf("no certificates found")
	}

	expiry := certs[0].NotAfter
	return expiry, nil
}
