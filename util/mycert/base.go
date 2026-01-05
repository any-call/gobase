package mycert

import (
	"crypto/tls"
	"fmt"
	"os/exec"
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

func CheckCertExpiryEx(domain, servIP string, port int) (time.Time, error) {
	conn, err := tls.Dial("tcp", servIP+fmt.Sprintf(":%d", port), &tls.Config{
		ServerName: domain, // 确保证书验证 SNI
	})
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

// 基于 acme 续期的证书，其文件名不会变，只会变更内容
func RenewCertificate(domain string) error {
	cmd := exec.Command("/root/.acme.sh/acme.sh", "--renew", "-d", domain, "--ecc", "--force")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to renew certificate: %v\n%s", err, string(output))
	}
	fmt.Printf("✅ Certificate renewed for %s:\n%s\n", domain, string(output))
	return nil
}

func RestartService(serviceName string) error {
	cmd := exec.Command("systemctl", "restart", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to restart service: %v\n%s", err, string(output))
	}
	fmt.Printf("🔄 Service %s restarted:\n%s\n", serviceName, string(output))
	return nil
}
