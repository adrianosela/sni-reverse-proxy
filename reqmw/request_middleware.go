package reqmw

import "net/http"

// RequestMiddleware represents the request-processing
// middleware functionality of a multiple-host proxy
type RequestMiddleware interface {
	Do(http.ResponseWriter, *http.Request)
	Wrap(http.HandlerFunc) RequestMiddleware
}
