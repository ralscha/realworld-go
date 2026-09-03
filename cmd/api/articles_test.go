package main

import (
	"net/http/httptest"
	"testing"
)

func TestPagination(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		wantOffset int
		wantLimit  int
		wantError  bool
	}{
		{name: "defaults", wantOffset: 0, wantLimit: 20},
		{name: "values", query: "?offset=10&limit=50", wantOffset: 10, wantLimit: 50},
		{name: "zero values", query: "?offset=0&limit=0", wantOffset: 0, wantLimit: 0},
		{name: "negative offset", query: "?offset=-1", wantError: true},
		{name: "invalid limit", query: "?limit=many", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/api/articles"+tt.query, nil)
			offset, limit, err := pagination(r)
			if (err != nil) != tt.wantError {
				t.Fatalf("pagination() error = %v, wantError %v", err, tt.wantError)
			}
			if offset != tt.wantOffset || limit != tt.wantLimit {
				t.Fatalf("pagination() = %d, %d, want %d, %d", offset, limit, tt.wantOffset, tt.wantLimit)
			}
		})
	}
}
