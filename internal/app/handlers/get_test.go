package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexuryumtsev/go-shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
)

func TestGetHandler(t *testing.T) {
	// тестовое хранилище и добавляем тестовые данные.
	id := generateID()
	st := storage.NewStorage()
	st.Save(id, "https://practicum.yandex.ru/")

	type want struct {
		code        int
		header      string
		contentType string
	}

	testCases := []struct {
		name        string
		id          string
		requestPath string
		want        want
	}{
		{
			name:        "Valid ID",
			id:          id,
			requestPath: "/" + id,
			want: want{
				code:        http.StatusTemporaryRedirect,
				header:      "https://practicum.yandex.ru/",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Invalid ID",
			id:          "",
			requestPath: "/",
			want: want{
				code:        http.StatusNotFound,
				header:      "",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// тестовый HTTP-запрос.
			req := httptest.NewRequest(http.MethodGet, tc.requestPath, nil)
			rec := httptest.NewRecorder()

			handler := GetHandler(st, tc.id)
			handler(rec, req)

			res := rec.Result()
			defer res.Body.Close()
			assert.Equal(t, tc.want.code, res.StatusCode)
			assert.Equal(t, tc.want.header, rec.Header().Get("Location"))

			if tc.name != "Valid ID" {
				assert.Equal(t, tc.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
