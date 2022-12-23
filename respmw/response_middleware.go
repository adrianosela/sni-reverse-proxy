package respmw

import "net/http"

// ResponseMiddleware represents the response-processing
// middleware functionality of a multiple-host proxy
type ResponseMiddleware interface {
	Do(resp *http.Response) error
}
