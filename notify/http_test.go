package notify

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// English text doJSON English text HTTP English text。
//
// English text：English text（English text errcode、
// English text code、Telegram ok English text），English text**HTTP English text**English text doJSON English text，
// English text。English text，English text 503 English text、
// English text；English text 403 English text、English text。

func replyServer(t *testing.T, status int, body string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestDoJSONClassifiesHTTPStatus(t *testing.T) {
	cases := []struct {
		name      string
		status    int
		permanent bool
	}{
		{"200 English text", 200, false},
		{"429 English text", 429, false},
		{"408 English text", 408, false},
		{"500 English text", 500, false},
		{"502 English text", 502, false},
		{"503 English text", 503, false},
		{"400 English text", 400, true},
		{"401 English text", 401, true},
		{"403 English text", 403, true},
		{"404 English text", 404, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := replyServer(t, tc.status, `{"detail":"upstream says no"}`)
			_, err := doJSON(context.Background(), "GET", srv.URL, nil, nil)
			if tc.status < 300 {
				if err != nil {
					t.Fatalf("2xx English text: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("English text 2xx English text")
			}
			if got := IsPermanent(err); got != tc.permanent {
				t.Fatalf("HTTP %d English text permanent English text：English text %v English text %v (%v)",
					tc.status, tc.permanent, got, err)
			}
			// English text，English text。
			// English text Go English text StatusText：English text
			// （English text），English text、English text。
			if !strings.Contains(err.Error(), strconv.Itoa(tc.status)) {
				t.Errorf("English text HTTP English text %d，English text %v", tc.status, err)
			}
		})
	}
}

// TestDoJSONIncludesResponseSnippet English text snippet：English text，
// English text「English text」，English text。
func TestDoJSONIncludesResponseSnippet(t *testing.T) {
	srv := replyServer(t, 400, `{"error":"invalid webhook token"}`)
	_, err := doJSON(context.Background(), "GET", srv.URL, nil, nil)
	if err == nil {
		t.Fatal("English text")
	}
	if !strings.Contains(err.Error(), "invalid webhook token") {
		t.Errorf("English text，English text %v", err)
	}
}

// TestDoJSONSnippetIsSingleLineAndBounded English text snippet English text：
// English text last_error English text，English text/English text。
func TestDoJSONSnippetIsSingleLineAndBounded(t *testing.T) {
	// English text、English text 5000 English text。
	long := strings.Repeat("x", 5000)
	srv := replyServer(t, 500, "line1\nline2\r\n\tline3 "+long)
	_, err := doJSON(context.Background(), "GET", srv.URL, nil, nil)
	if err == nil {
		t.Fatal("English text")
	}
	msg := err.Error()
	if strings.ContainsAny(msg, "\r\n\t") {
		t.Errorf("English text，English text %q", msg)
	}
	// snippet English text 200 English text + English text，English text。
	if len(msg) > 400 {
		t.Errorf("English text（%d English text），English text snippet English text: %q", len(msg), msg)
	}
}

// TestDoJSONRejectsOversizedResponse English text：English text
// English text（English text last_error）。
func TestDoJSONRejectsOversizedResponse(t *testing.T) {
	huge := strings.Repeat("A", 1<<20) // 1 MiB
	srv := replyServer(t, 400, huge)
	_, err := doJSON(context.Background(), "GET", srv.URL, nil, nil)
	if err == nil {
		t.Fatal("English text")
	}
	if len(err.Error()) > 400 {
		t.Errorf("English text，English text %d", len(err.Error()))
	}
}
