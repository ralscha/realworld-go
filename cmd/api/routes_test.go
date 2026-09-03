package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
)

func TestClientIPRateLimitKeyUsesProxyAppendedXForwardedForEntry(t *testing.T) {
	var got string
	handler := middleware.ClientIPFromXFF()(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		got, _ = clientIPRateLimitKey(r)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.RemoteAddr = "10.0.0.2:1234"
	r.Header.Set("X-Forwarded-For", "203.0.113.10, 198.51.100.25")
	handler.ServeHTTP(httptest.NewRecorder(), r)

	if got != "198.51.100.25" {
		t.Fatalf("rate-limit key = %q, want proxy-appended client address", got)
	}
}
