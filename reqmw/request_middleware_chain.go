package reqmw

import (
	"net/http"
)

// MiddlewareNodeFn represents the functionality of a single middleware
// in the chain. Each of these functions is responsible for invoking
// the next() middleware in the chain and handling the error.
type MiddlewareNodeFn func(
	rw http.ResponseWriter,
	req *http.Request,
	next http.HandlerFunc,
)

// RequestMiddlewareChain represents a linked-list-like data
// structure where each node is a middleware to be executed
type RequestMiddlewareChain struct {
	current MiddlewareNodeFn
	next    *RequestMiddlewareChain
}

// ensure RequestMiddlewareChain implements RequestMiddleware at compile-time
var _ RequestMiddleware = (*RequestMiddlewareChain)(nil)

// NewRequestMiddlewareChain returns a new RequestMiddlewareChain initialized with
// the given list of MiddlewareNodeFn. They will be executed in the order in which they are given.
func NewRequestMiddlewareChain(layers ...MiddlewareNodeFn) *RequestMiddlewareChain {
	if len(layers) == 0 {
		return &RequestMiddlewareChain{}
	}
	return &RequestMiddlewareChain{
		current: layers[0],
		next:    NewRequestMiddlewareChain(layers[1:]...),
	}
}

// Do executes the current MiddlewareNodeFn, with the next MiddlewareNodeFn as the next argument
func (c *RequestMiddlewareChain) Do(rw http.ResponseWriter, req *http.Request) {
	if c.current == nil {
		return
	}

	next := func(passedRW http.ResponseWriter, passedReq *http.Request) {
		if c.next != nil {
			c.next.Do(passedRW, passedReq)
		}
	}
	c.current(rw, req, next)
}

// Wrap adds a function to the end of a RequestMiddlewareChain
func (c *RequestMiddlewareChain) Wrap(fn http.HandlerFunc) RequestMiddleware {
	// The last middleware in the chain will have a nil
	// pointer for the current field, the next field is
	// never nil.
	last := c
	for last.current != nil {
		last = last.next
	}

	// we set the given function to be the last
	last.current = func(
		rw http.ResponseWriter,
		req *http.Request,
		next http.HandlerFunc,
	) {
		fn(rw, req) // run the passed handlerfunc
		next(rw, req)
	}

	// and maintain the behavior that the next field must not be nil
	last.next = NewRequestMiddlewareChain()

	return c
}
