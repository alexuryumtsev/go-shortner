package middleware

import (
	"errors"
	"net/http"

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
func ProcessError(w http.ResponseWriter, err error, shortenedURL string) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		// Если ошибка уникальности
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(shortenedURL))
		return
	}

	// Обработка других ошибок
	http.Error(w, err.Error(), http.StatusBadRequest)
}
