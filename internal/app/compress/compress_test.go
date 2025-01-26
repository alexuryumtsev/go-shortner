package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGzipMiddleware(t *testing.T) {
	handler := GzipMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "hello, world"}`))
	}))

	tests := []struct {
		name               string
		acceptEncoding     string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "Client supports gzip",
			acceptEncoding:     "gzip",
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message": "hello, world"}`,
		},
		{
			name:               "Client does not support gzip",
			acceptEncoding:     "",
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message": "hello, world"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Accept-Encoding", tt.acceptEncoding)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedStatusCode, res.StatusCode)
			}

			var body []byte
			var err error

			if strings.Contains(res.Header.Get("Content-Encoding"), "gzip") {
				gz, err := gzip.NewReader(res.Body)
				body, err = io.ReadAll(gz)
				if err != nil {
					t.Fatalf("failed to create gzip reader: %v", err)
				}
				defer gz.Close()
			} else {
				body, err = io.ReadAll(res.Body)
			}

			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}

			if string(body) != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, body)
			}
		})
	}
}
