package types

import (
	"errors"
	"strings"
	"testing"
)

func TestAPIError_Error_WithDetail(t *testing.T) {
	e := APIError{Code: 400, Message: "bad request", Detail: "field required"}
	got := e.Error()
	if !strings.Contains(got, "bad request") {
		t.Errorf("missing message in %q", got)
	}
	if !strings.Contains(got, "field required") {
		t.Errorf("missing detail in %q", got)
	}
}

func TestAPIError_Error_WithoutDetail(t *testing.T) {
	e := APIError{Code: 400, Message: "bad request"}
	got := e.Error()
	if !strings.Contains(got, "bad request") {
		t.Errorf("missing message in %q", got)
	}
	if strings.Contains(got, "(") {
		t.Errorf("unexpected parentheses in %q (no detail expected)", got)
	}
}

func TestAPIError_Error_EmptyMessage(t *testing.T) {
	e := APIError{Code: 404}
	got := e.Error()
	if got == "" {
		t.Error("Error() returned empty string")
	}
}

func TestResponseError_Error_WithErrors(t *testing.T) {
	e := &ResponseError{
		StatusCode: 404,
		Errors:     []APIError{{Message: "not found", Detail: "resource missing"}},
	}
	got := e.Error()
	if !strings.Contains(got, "404") {
		t.Errorf("missing status code in %q", got)
	}
	if !strings.Contains(got, "not found") {
		t.Errorf("missing error message in %q", got)
	}
}

func TestResponseError_Error_MultipleErrors(t *testing.T) {
	e := &ResponseError{
		StatusCode: 422,
		Errors: []APIError{
			{Message: "first error"},
			{Message: "second error"},
		},
	}
	got := e.Error()
	if !strings.Contains(got, "first error") {
		t.Errorf("missing first error in %q", got)
	}
	if !strings.Contains(got, "second error") {
		t.Errorf("missing second error in %q", got)
	}
}

func TestResponseError_Error_WithoutErrors(t *testing.T) {
	e := &ResponseError{StatusCode: 500}
	got := e.Error()
	if !strings.Contains(got, "500") {
		t.Errorf("missing status code in %q", got)
	}
}

func TestResponseError_Error_ZeroErrors(t *testing.T) {
	e := &ResponseError{StatusCode: 503, Errors: []APIError{}}
	got := e.Error()
	if !strings.Contains(got, "503") {
		t.Errorf("missing status code in %q", got)
	}
}

func TestAsResponseError_Success(t *testing.T) {
	orig := &ResponseError{StatusCode: 401, Errors: []APIError{{Message: "unauthorized"}}}
	got, ok := AsResponseError(orig)
	if !ok {
		t.Fatal("expected ok=true for *ResponseError")
	}
	if got != orig {
		t.Error("expected same pointer returned")
	}
}

func TestAsResponseError_WrappedError(t *testing.T) {
	orig := &ResponseError{StatusCode: 403}
	wrapped := errors.Join(orig)
	got, ok := AsResponseError(wrapped)
	if !ok {
		t.Fatal("expected ok=true for wrapped *ResponseError")
	}
	if got.StatusCode != 403 {
		t.Errorf("expected 403, got %d", got.StatusCode)
	}
}

func TestAsResponseError_NonResponseError(t *testing.T) {
	_, ok := AsResponseError(errors.New("plain error"))
	if ok {
		t.Fatal("expected ok=false for plain error")
	}
}

func TestAsResponseError_Nil(t *testing.T) {
	_, ok := AsResponseError(nil)
	if ok {
		t.Fatal("expected ok=false for nil error")
	}
}
