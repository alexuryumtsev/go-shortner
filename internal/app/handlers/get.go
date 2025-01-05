package handlers

import (
	"net/http"

	"github.com/alexuryumtsev/go-shortener/internal/app/storage"
)

// GetHandler обрабатывает GET-запросы с динамическими id.
func GetHandler(repo storage.Repository, id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		originalURL, exists := repo.Get(id)
		if !exists {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		// Ответ с редиректом на оригинальный URL.
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
