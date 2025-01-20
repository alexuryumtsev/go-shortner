package handlers

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/alexuryumtsev/go-shortener/internal/app/service"
	"github.com/alexuryumtsev/go-shortener/internal/app/storage"
)

// PostHandler обрабатывает POST-запросы.
func PostHandler(storage storage.URLStorage, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		originalURL := strings.TrimSpace(string(body))
		shortenedURL, err := service.NewURLService(storage, baseURL).ShortenerURL(originalURL)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortenedURL))
	}
}

// generateID генерирует случайный идентификатор.
func generateID(url string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(url)))[:8] // Используем MD5 и берём первые 8 символов.
}
