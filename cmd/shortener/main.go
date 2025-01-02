package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alexuryumtsev/go-shortener/internal/app/handlers"
	"github.com/alexuryumtsev/go-shortener/internal/app/storage"
)

func main() {
	// Инициализация хранилища.
	st := storage.NewStorage()

	// Регистрация маршрутов.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Проверка на метод POST для главной страницы.
		if r.Method == http.MethodPost {
			handlers.PostHandler(st)(w, r)
			return
		}

		// Парсим id из пути запроса для GET-запросов.
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path != "" && r.Method == http.MethodGet {
			handlers.GetHandler(st, path)(w, r)
			return
		}

		// Обработка ошибок для некорректных запросов.
		http.Error(w, "Invalid request", http.StatusBadRequest)
	})

	// Запуск сервера.
	fmt.Println("Server started at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
