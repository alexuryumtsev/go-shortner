package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/alexuryumtsev/go-shortener/internal/app/models"
	"github.com/alexuryumtsev/go-shortener/internal/app/service"
	"github.com/alexuryumtsev/go-shortener/internal/app/storage"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// PostHandler обрабатывает POST-запросы.
func PostHandler(storage storage.URLWriter, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		ctx := r.Context()
		originalURL := strings.TrimSpace(string(body))
		shortenedURL, err := service.NewURLService(ctx, storage, baseURL).ShortenerURL(originalURL)

		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(shortenedURL))
				return
			}

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(shortenedURL))
	}
}

// RequestBody определяет структуру входных данных.
type RequestBody struct {
	URL string `json:"url"`
}

// ResponseBody определяет структуру ответа.
type ResponseBody struct {
	ShortURL string `json:"result"`
}

// PostHandler обрабатывает POST-запросы для создания коротких URL.
func PostJSONHandler(storage storage.URLWriter, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestBody
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		ctx := r.Context()
		shortenedURL, err := service.NewURLService(ctx, storage, baseURL).ShortenerURL(req.URL)

		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				resp := ResponseBody{
					ShortURL: shortenedURL,
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(resp)
				return
			}

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp := ResponseBody{
			ShortURL: shortenedURL,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}

// PostBatchHandler обрабатывает POST-запросы для создания множества коротких URL.
func PostBatchHandler(repo storage.URLStorage, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baseURL = strings.TrimSuffix(baseURL, "/")

		var batchModels []models.URLBatchModel
		if err := json.NewDecoder(r.Body).Decode(&batchModels); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if len(batchModels) == 0 {
			http.Error(w, "Empty batch", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		urlService := service.NewURLService(ctx, repo, baseURL)

		shortenedURLs, err := urlService.SaveBatchShortenerURL(batchModels)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				var batchResponseModels []models.BatchResponseModel
				for i, shortenedURL := range shortenedURLs {
					batchResponseModels = append(batchResponseModels, models.BatchResponseModel{
						CorrelationID: batchModels[i].CorrelationID,
						ShortURL:      shortenedURL,
					})
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict)
				if err := json.NewEncoder(w).Encode(batchResponseModels); err != nil {
					http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				}
				return
			}

			http.Error(w, "Failed to save URL", http.StatusInternalServerError)
			return
		}

		var batchResponseModels []models.BatchResponseModel
		for i, shortenedURL := range shortenedURLs {
			batchResponseModels = append(batchResponseModels, models.BatchResponseModel{
				CorrelationID: batchModels[i].CorrelationID,
				ShortURL:      shortenedURL,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(batchResponseModels); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
