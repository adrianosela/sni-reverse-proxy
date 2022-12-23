package respmw

import "net/http"

// MiddlewareNodeFn represents the functionality of a single middleware
// in the chain. Each of these functions is responsible for invoking
// the next() middleware in the chain and handling the error.
type MiddlewareNodeFn func(
	resp *http.Response,
	next func(passed *http.Response) error,
) error

// ResponseMiddlewareChain represents a linked-list-like data
// structure where each node is a middleware to be executed
type ResponseMiddlewareChain struct {
	current MiddlewareNodeFn
	next    *ResponseMiddlewareChain
}

// ensure LayeredResponseMiddleware implements ResponseMiddleware at compile-time
var _ ResponseMiddleware = (*ResponseMiddlewareChain)(nil)

// NewResponseMiddlewareChain returns a new ResponseMiddlewareChain initialized with
// the given list of MiddlewareNodeFn. They will be executed in the order in which they are given.
func NewResponseMiddlewareChain(layers ...MiddlewareNodeFn) *ResponseMiddlewareChain {
	if len(layers) == 0 {
		return &ResponseMiddlewareChain{}
	}
	return &ResponseMiddlewareChain{
		current: layers[0],
		next:    NewResponseMiddlewareChain(layers[1:]...),
	}
}

// Do executes the current MiddlewareNodeFn, with the next MiddlewareNodeFn as the next argument
func (c *ResponseMiddlewareChain) Do(resp *http.Response) error {
	if c.current == nil {
		return nil
	}

	next := func(passed *http.Response) error {
		if c.next != nil {
			return c.next.Do(passed)
		}
		return nil
	}
	return c.current(resp, next)
}
