package types

import (
	"errors"
	"fmt"
	"strings"
)

const errTag = "sentinelone: "

// Pagination holds cursor-based pagination metadata returned by list endpoints.
//
// When NextCursor is non-nil, pass its value as ListParams.Cursor on the next
// call to retrieve the following page.  NextCursor is nil on the last page.
// TotalItems reports the total number of matching records across all pages
// (unless the request set SkipCount, in which case it may be zero).
type Pagination struct {
	NextCursor *string `json:"nextCursor,omitempty"`
	TotalItems int     `json:"totalItems"`
}

// APIError represents a single structured error object returned by the API
// inside the top-level "errors" array of the response envelope.
type APIError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

func (e APIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s%s (%s)", errTag, e.Message, e.Detail)
	}

	return fmt.Sprintf("%s%s", errTag, e.Message)
}

// ResponseError is returned when the server responds with a non-2xx status
// code.  Use [AsResponseError] to unwrap it without importing the errors
// package directly:
//
//	if respErr, ok := types.AsResponseError(err); ok {
//	    fmt.Println(respErr.StatusCode, respErr.Errors)
//	}
type ResponseError struct {
	StatusCode int
	Errors     []APIError
}

func (e *ResponseError) Error() string {
	if len(e.Errors) > 0 {
		msgs := make([]string, len(e.Errors))
		for i, err := range e.Errors {
			msgs[i] = err.Error()
		}

		return fmt.Sprintf("%sHTTP %d: %s", errTag, e.StatusCode, strings.Join(msgs, "; "))
	}

	return fmt.Sprintf("%sHTTP %d", errTag, e.StatusCode)
}

// AsResponseError unwraps err into a *ResponseError and reports whether the
// conversion succeeded.
func AsResponseError(err error) (*ResponseError, bool) {
	var e *ResponseError
	if errors.As(err, &e) {
		return e, true
	}

	return nil, false
}
