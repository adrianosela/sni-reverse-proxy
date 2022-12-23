package reqmw

import "net/http"

// DoNothingRequestMiddleware an NO-OP implemenation of the RequestMiddleware interface
type DoNothingRequestMiddleware struct{}

// ensure DoNothingRequestMiddleware implements RequestMiddleware at compile-time
var _ RequestMiddleware = (*DoNothingRequestMiddleware)(nil)

// NewDoNothingRequestMiddleware returns an DoNothingRequestMiddleware
func NewDoNothingRequestMiddleware() *DoNothingRequestMiddleware {
	return &DoNothingRequestMiddleware{}
}

// Do does nothing to the given request
func (a *DoNothingRequestMiddleware) Do(rw http.ResponseWriter, req *http.Request) {
	return
}

// Wrap adds functionality to a middleware
func (a *DoNothingRequestMiddleware) Wrap(fn http.HandlerFunc) RequestMiddleware {
	return NewRequestMiddlewareChain(
		func(
			rw http.ResponseWriter,
			req *http.Request,
			next http.HandlerFunc,
		) {
			fn(rw, req)
			next(rw, req)
		},
	)
}
