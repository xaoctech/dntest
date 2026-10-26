package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestDownload(t *testing.T) {
	for _, tc := range []struct {
		query string
		size  int
	}{
		{"", 204800},
		{"0", 204800},
		{"-1", 204800},
		{"invalid", 204800},
		{"9223372036854775808", 204800},
		{"1", 1},
		{"1024", 1024},
		{"1025", 1025},
	} {
		t.Run(tc.query, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/?"+tc.query, nil)
			w := httptest.NewRecorder()
			download(w, r)

			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200", w.Code)
			}
			if w.Body.String() != strings.Repeat("0", tc.size) {
				t.Errorf("body differs from %d ASCII zero bytes (got %d bytes)", tc.size, w.Body.Len())
			}
			size := strconv.Itoa(tc.size)
			for name, want := range map[string]string{
				"Content-Type":        "application/octet-stream",
				"Content-Disposition": "attachment; filename=dntest" + size,
				"Content-Length":      size,
				"Pragma":              "no-cache",
				"Expires":             "0",
			} {
				if got := w.Header().Get(name); got != want {
					t.Errorf("%s = %q, want %q", name, got, want)
				}
			}
		})
	}
}

func TestDownloadHead(t *testing.T) {
	w := httptest.NewRecorder()
	download(w, httptest.NewRequest(http.MethodHead, "/?1025", nil))
	if w.Body.Len() != 0 {
		t.Errorf("HEAD body has %d bytes, want 0", w.Body.Len())
	}
	if got := w.Header().Get("Content-Length"); got != "1025" {
		t.Errorf("Content-Length = %q, want 1025", got)
	}
}
