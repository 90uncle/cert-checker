package checker

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

type CertInfo struct {
	Domain     string
	ExpiresIn  int       // 剩余天数
	ExpiryDate time.Time // 过期时间
	Issuer     string
	CommonName string
	IsExpired  bool
	IsWarning  bool
}

func CheckCert(domain string, alertThreshold int) (*CertInfo, error) {
	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: 10 * time.Second},
		"tcp",
		domain,
		&tls.Config{InsecureSkipVerify: true},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %v", domain, err)
	}
	defer conn.Close()
	cert := conn.ConnectionState().PeerCertificates[0]
	now := time.Now()
	expiresIn := int(cert.NotAfter.Sub(now).Hours() / 24)
	info := &CertInfo{
		Domain:     domain,
		ExpiresIn:  expiresIn,
		ExpiryDate: cert.NotAfter,
		Issuer:     cert.Issuer.String(),
		CommonName: cert.Subject.CommonName,
		IsExpired:  expiresIn < 0,
		IsWarning:  expiresIn <= alertThreshold && expiresIn >= 0,
	}
	return info, nil
}
