package certmgr

import (
	"crypto/tls"
)

// CertificateManager represents the certificate-management functionality required of a multiple-host proxy
type CertificateManager interface {
	GetCertificate(sni string) (*tls.Certificate, error)
}
