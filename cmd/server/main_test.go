package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithCORSSetsHeaders(t *testing.T) {
	handler := withCORS(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/state", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("expected CORS origin header to be set, got %q", got)
	}
}

func TestWithCORSHandlesPreflight(t *testing.T) {
	called := false
	handler := withCORS(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})

	req := httptest.NewRequest(http.MethodOptions, "/api/move", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if called {
		t.Fatal("OPTIONS preflight request should not reach the wrapped handler")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for preflight, got %d", w.Code)
	}
}
