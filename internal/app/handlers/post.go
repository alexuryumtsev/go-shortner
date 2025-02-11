package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/alexuryumtsev/go-shortener/internal/app/models"
	"github.com/alexuryumtsev/go-shortener/internal/app/service"
	"github.com/alexuryumtsev/go-shortener/internal/app/storage"
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

type BatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// PostBatchHandler обрабатывает POST-запросы для создания множества коротких URL.
func PostBatchHandler(repo storage.URLStorage, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baseURL = strings.TrimSuffix(baseURL, "/")

		var batchRequests []BatchRequest
		if err := json.NewDecoder(r.Body).Decode(&batchRequests); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if len(batchRequests) == 0 {
			http.Error(w, "Empty batch", http.StatusBadRequest)
			return
		}

		var batchResponses []BatchResponse
		for _, req := range batchRequests {
			urlModel := models.URLModel{
				ID:  service.GenerateID(req.OriginalURL), // Функция для генерации короткого ID
				URL: req.OriginalURL,
			}

			// Проверяем, существует ли уже оригинальный URL
			existingURLModel, exists := repo.Get(r.Context(), urlModel.ID)
			if exists {
				batchResponses = append(batchResponses, BatchResponse{
					CorrelationID: req.CorrelationID,
					ShortURL:      baseURL + "/" + existingURLModel.ID,
				})
				continue
			}

			if err := repo.Save(r.Context(), urlModel); err != nil {
				http.Error(w, "Failed to save URL", http.StatusInternalServerError)
				return
			}

			batchResponses = append(batchResponses, BatchResponse{
				CorrelationID: req.CorrelationID,
				ShortURL:      baseURL + "/" + urlModel.ID,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(batchResponses); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
