package handlers

import (
	"net/http"

	"github.com/alexuryumtsev/go-shortener/internal/app/storage"
	"github.com/alexuryumtsev/go-shortener/internal/app/storage/pg"
)

// PingHandler проверяет соединение с базой данных.
func PingHandler(repo storage.URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, является ли хранилище экземпляром DatabaseStorage
		if dbRepo, ok := repo.(*pg.DatabaseStorage); ok {
			if err := dbRepo.Ping(r.Context()); err != nil {
				http.Error(w, "Database connection error", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		}
	}
}
