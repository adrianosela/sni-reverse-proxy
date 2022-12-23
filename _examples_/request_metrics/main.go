package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/adrianosela/sni-reverse-proxy/certmgr"
	"github.com/adrianosela/sni-reverse-proxy/hostmgr"
	"github.com/adrianosela/sni-reverse-proxy/proxy"
	"github.com/adrianosela/sni-reverse-proxy/reqmw"
)

func writeResponse(rec *httptest.ResponseRecorder, rw http.ResponseWriter) {
	// copy headers
	for key, values := range rec.Header() {
		for _, value := range values {
			rw.Header().Add(key, value)
		}
	}

	// copy status code
	rw.WriteHeader(rec.Code)

	// copy body
	rw.Write(rec.Body.Bytes())
}

func requestMetrics(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	rec := httptest.NewRecorder()

	start := time.Now()
	next(rec, req)
	end := time.Now()

	writeResponse(rec, rw)

	go func() {
		log.Printf(
			"{\"method\":\"%s\",\"url\":\"https://%s%s\",\"status\":%d,\"content_length\":%d,\"content_type\":\"%s\",\"duration\":\"%d μs\"}",
			req.Method,
			req.TLS.ServerName,
			req.URL.String(),
			rec.Result().StatusCode,
			rec.Result().ContentLength,
			rec.HeaderMap.Get("Content-Type"),
			end.Sub(start).Microseconds(),
		)
	}()
}

func main() {
	rp := proxy.New().
		WithRequestMiddleware(reqmw.NewRequestMiddlewareChain(
			requestMetrics,
		)).
		WithHostManager(hostmgr.NewInMemoryHostManager()).
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
