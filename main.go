package main

import (
	"log"

	"github.com/adrianosela/sni-reverse-proxy/certmgr"
	"github.com/adrianosela/sni-reverse-proxy/hostmgr"
	"github.com/adrianosela/sni-reverse-proxy/proxy"
	"github.com/adrianosela/sni-reverse-proxy/reqmw"
	"github.com/adrianosela/sni-reverse-proxy/respmw"
)

func main() {
	rp := proxy.New().
		// do nothing with requests for now (same as the default behavior)
		WithRequestMiddleware(reqmw.NewDoNothingRequestMiddleware()).
		// do nothing with responses for now (same as the default behavior)
		WithResponseMiddleware(respmw.NewDoNothingResponseMiddleware()).
		// store hosts configuration in-memory (same as the default behavior)
		WithHostManager(hostmgr.NewInMemoryHostManager()).
		// using file-system certificate manager with test certs directory
		WithCertificateManager(certmgr.NewFileSystemCertificateManager(
			"./_test_certs_/%s-cert.pem",
			"./_test_certs_/%s-key.pem",
		))

	// add target hosts to proxy
	for sni, target := range map[string]string{
		"backend_a.adrianosela.com": "http://localhost:8085",
		"backend_b.adrianosela.com": "http://localhost:8086",
		"backend_c.adrianosela.com": "http://localhost:8087",
	} {
		if err := rp.AddHost(sni, target); err != nil {
			log.Fatalf("failed to add proxy target host: %s", err)
		}
	}

	// start the proxy server
	if err := rp.ListenAndServeTLS(":8443"); err != nil {
		log.Fatalf("failed to listen and serve TLS: %s", err)
	}
}
