package certmgr

import (
	"crypto/tls"
	"fmt"
)

// FileSystemCertificateManager is an implementation of the
// CertificateManager iface that gets certificates from disk
type FileSystemCertificateManager struct {
	certPathFmt string
	keyPathFmt  string
}

// ensure FSCertificateManager implements CertificateManager at compile-time
var _ CertificateManager = (*FileSystemCertificateManager)(nil)

// NewFileSystemCertificateManager is the FileSystemCertificateManager constructor
func NewFileSystemCertificateManager(certPathFmt, keyPathFmt string) *FileSystemCertificateManager {
	return &FileSystemCertificateManager{certPathFmt: certPathFmt, keyPathFmt: keyPathFmt}
}

// GetCertificate returns the *tls.Certificate for a given SNI
func (fs *FileSystemCertificateManager) GetCertificate(sni string) (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(
		fmt.Sprintf(fs.certPathFmt, sni),
		fmt.Sprintf(fs.keyPathFmt, sni),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load x509 key pair from file system: %s", err)
	}
	return &cert, err
}
