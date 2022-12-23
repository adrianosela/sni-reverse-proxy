package respmw

import "net/http"

// DoNothingResponseMiddleware is an NO-OP implemenation of the ResponseMiddleware interface
type DoNothingResponseMiddleware struct{}

// ensure DoNothingResponseMiddleware implements ResponseMiddleware at compile-time
var _ ResponseMiddleware = (*DoNothingResponseMiddleware)(nil)

// NewDoNothingResponseMiddleware returns a new DoNothingResponseMiddleware
func NewDoNothingResponseMiddleware() *DoNothingResponseMiddleware {
	return &DoNothingResponseMiddleware{}
}

// Do does nothing to the given response
func (d *DoNothingResponseMiddleware) Do(resp *http.Response) error { return nil }
