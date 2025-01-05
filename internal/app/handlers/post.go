package handlers

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/alexuryumtsev/go-shortener/internal/app/storage"
)

const baseURL = "http://localhost:8080/"

// PostHandler обрабатывает POST-запросы.
func PostHandler(repo storage.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid method", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		originalURL := strings.TrimSpace(string(body))
		if originalURL == "" {
			http.Error(w, "Empty URL", http.StatusBadRequest)
			return
		}

		id := generateID(originalURL)
		repo.Save(id, originalURL)

		shortenedURL := baseURL + id
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortenedURL))
	}
}

// generateID генерирует случайный идентификатор.
func generateID(url string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(url)))[:8] // Используем MD5 и берём первые 8 символов.
}
