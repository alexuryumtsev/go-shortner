package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/alexuryumtsev/go-shortener/internal/app/storage"

	"golang.org/x/exp/rand"
)

const baseURL = "http://localhost:8080/"

// PostHandler обрабатывает POST-запросы.
func PostHandler(storage *storage.Storage) http.HandlerFunc {
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

		id := generateID()
		storage.Save(id, originalURL)

		shortenedURL := baseURL + id
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortenedURL))
	}
}

// generateID генерирует случайный идентификатор.
func generateID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const idLength = 8
	b := make([]byte, idLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
