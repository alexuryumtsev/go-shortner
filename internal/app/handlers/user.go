package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/alexuryumtsev/go-shortener/internal/app/service/user"
	"github.com/alexuryumtsev/go-shortener/internal/app/storage"
)

type UserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func GetUserURLsHandler(repo storage.URLStorage, userService user.UserService, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := userService.GetUserIDFromCookie(r)
		urls, err := repo.GetUserURLs(r.Context(), userID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if len(urls) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		baseURL := strings.TrimSuffix(baseURL, "/")
		var userURLs []UserURL
		for _, url := range urls {
			if url.Deleted {
				continue
			}
			userURLs = append(userURLs, UserURL{
				ShortURL:    baseURL + "/" + url.ID,
				OriginalURL: url.URL,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userURLs)
	}
}

func DeleteUserURLsHandler(repo storage.URLStorage, userService user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var shortURLs []string
		if err := json.NewDecoder(r.Body).Decode(&shortURLs); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		userID := userService.GetUserIDFromCookie(r)

		// Создаем каналы для передачи идентификаторов URL и завершения работы
		idChan := make(chan string, len(shortURLs))
		done := make(chan struct{})
		quit := make(chan struct{})

		// Запускаем горутину для заполнения канала идентификаторами
		go func() {
			defer close(idChan)
			for _, shortURL := range shortURLs {
				select {
				case idChan <- shortURL:
				case <-quit:
					return
				}
			}
		}()

		// Запускаем горутину для удаления URL
		go func() {
			defer close(done)
			var batch []string
			for {
				select {
				case id, ok := <-idChan:
					if !ok {
						// Канал закрыт, обновляем оставшиеся идентификаторы
						if len(batch) > 0 {
							repo.DeleteUserURLs(r.Context(), userID, batch)
						}
						return
					}
					batch = append(batch, id)
					// Если буфер заполнен, выполняем обновление
					if len(batch) >= 10 {
						repo.DeleteUserURLs(r.Context(), userID, batch)
						batch = batch[:0]
					}
				case <-quit:
					return
				}
			}
		}()

		// Ожидаем завершения удаления
		<-done

		// Закрываем канал quit, чтобы завершить все горутины
		close(quit)

		w.WriteHeader(http.StatusAccepted)
	}
}
