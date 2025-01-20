package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alexuryumtsev/go-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostHandler(t *testing.T) {
	// тестовое хранилище.
	repo := storage.NewMockStorage()
	handler := PostHandler(repo, "http://localhost:8080/")

	type want struct {
		code        int
		body        string
		contentType string
	}

	testCases := []struct {
		name     string
		inputURL string
		want     want
	}{
		{
			name:     "Valid URL",
			inputURL: "https://practicum.yandex.ru/",
			want: want{
				code:        http.StatusCreated,
				body:        "http://localhost:8080/",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:     "Empty URL",
			inputURL: "",
			want: want{
				code:        http.StatusBadRequest,
				body:        "empty URL",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// тестовый HTTP-запрос.
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tc.inputURL))
			rec := httptest.NewRecorder()
			handler(rec, req)

			res := rec.Result()
			defer res.Body.Close()
			assert.Equal(t, tc.want.code, res.StatusCode)

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.True(t, strings.HasPrefix(string(resBody), tc.want.body))
		})
	}
}
