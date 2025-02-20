package middleware

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alexuryumtsev/go-shortener/internal/app/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// ErrorMiddleware — middleware для обработки ошибок.
func ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// Поймаем панику, если она возникла
			if rec := recover(); rec != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()

		// Выполнение основной логики запроса
		next.ServeHTTP(w, r)
	})
}

// ProcessError — функция для обработки ошибок в контексте работы с БД.
func ProcessError(w http.ResponseWriter, err error, shortenedURL string, responseString bool) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)

		if responseString {
			w.Write([]byte(shortenedURL))
			return
		}

		json.NewEncoder(w).Encode(models.ResponseBody{
			ShortURL: shortenedURL,
		})

		return
	}

	// Обработка других ошибок
	http.Error(w, err.Error(), http.StatusBadRequest)
}
