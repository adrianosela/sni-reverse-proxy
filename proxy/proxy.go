package proxy

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/adrianosela/sni-reverse-proxy/certmgr"
	"github.com/adrianosela/sni-reverse-proxy/hostmgr"
	"github.com/adrianosela/sni-reverse-proxy/reqmw"
	"github.com/adrianosela/sni-reverse-proxy/respmw"
)

const (
	defaultCertPathFmt = "/etc/ssl/certs/%s.cert"
	defaultKeyPathFmt  = "/etc/ssl/certs/%s.key"
)

// Proxy is an multiple-host SNI-based reverse proxy
type Proxy struct {
	requestMiddleware  reqmw.RequestMiddleware
	responseMiddleware respmw.ResponseMiddleware
	hostManager        hostmgr.HostManager
	certManager        certmgr.CertificateManager
}

// New returns a new proxy with default configuration
func New() *Proxy {
	return &Proxy{
		// by default do nothing with requests
		requestMiddleware: reqmw.NewDoNothingRequestMiddleware(),
		// by default do nothing with responses
		responseMiddleware: respmw.NewDoNothingResponseMiddleware(),
		// by default use thread-safe in-memory storage
		hostManager: hostmgr.NewInMemoryHostManager(),
		// by default use file system cert manager
		certManager: certmgr.NewFileSystemCertificateManager(defaultCertPathFmt, defaultKeyPathFmt),
	}
}

// WithRequestMiddleware sets the proxy's request-processing middleware
func (p *Proxy) WithRequestMiddleware(rm reqmw.RequestMiddleware) *Proxy {
	p.requestMiddleware = rm
	return p
}

// WithResponseMiddleware sets the proxy's response-processing middleware
func (p *Proxy) WithResponseMiddleware(rm respmw.ResponseMiddleware) *Proxy {
	p.responseMiddleware = rm
	return p
}

// WithHostManager sets the proxy's host manager
func (p *Proxy) WithHostManager(hm hostmgr.HostManager) *Proxy {
	p.hostManager = hm
	return p
}

// WithCertificateManager sets the proxy's certificate manager
func (p *Proxy) WithCertificateManager(cm certmgr.CertificateManager) *Proxy {
	p.certManager = cm
	return p
}

// AddHost adds a target host to the proxy
func (p *Proxy) AddHost(sni string, addr string) error {
	targetURL, err := url.Parse(addr)
	if err != nil {
		return fmt.Errorf("failed to parse target address \"%s\": %s", addr, err)
	}
	targetRP := httputil.NewSingleHostReverseProxy(targetURL)
	targetRP.ModifyResponse = p.responseMiddleware.Do
	return p.hostManager.PutHost(sni, addr, targetRP)
}

// RemoveHost removes a target host from the proxy
func (p *Proxy) RemoveHost(sni string) error {
	return p.hostManager.RemoveHost(sni)
}

// ListenAndServeTLS starts the proxy server
func (p *Proxy) ListenAndServeTLS(addr string) error {

	wrappedHandler := p.requestMiddleware.Wrap(
		func(rw http.ResponseWriter, req *http.Request) {
			target, found, err := p.hostManager.GetHost(req.TLS.ServerName)
			if err != nil {
				_, _ = rw.Write([]byte(fmt.Sprintf("failed to get target host for SNI \"%s\": %s", req.TLS.ServerName, err)))
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !found {
				_, _ = rw.Write([]byte(fmt.Sprintf("no target host available for SNI \"%s\"", req.TLS.ServerName)))
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			target.ReverseProxy.ServeHTTP(rw, req)
			return
		},
	)

	server := http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(wrappedHandler.Do),
		TLSConfig: &tls.Config{
			GetCertificate: func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
				_, found, err := p.hostManager.GetHost(chi.ServerName)
				if err != nil {
					return nil, fmt.Errorf("failed to get target host for SNI \"%s\": %s", chi.ServerName, err)
				}
				if !found {
					return nil, fmt.Errorf("no target host available for SNI \"%s\"", chi.ServerName)
				}
				cert, err := p.certManager.GetCertificate(chi.ServerName)
				if err != nil {
					return nil, fmt.Errorf("failed to get certificate for target host with SNI \"%s\": %s", chi.ServerName, err)
				}
				return cert, nil
			},
		},
	}
	return server.ListenAndServeTLS("", "")
}
