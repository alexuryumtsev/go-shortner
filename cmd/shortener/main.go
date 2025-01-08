package main

import (
	"fmt"
	"net/http"

	"github.com/alexuryumtsev/go-shortener/internal/app/handlers"
	"github.com/alexuryumtsev/go-shortener/internal/app/storage"

	"github.com/go-chi/chi/v5"
)

func ShortenerRouter() chi.Router {
	// Инициализация хранилища.
	var repo storage.Repository = storage.NewStorage()

	// Регистрация маршрутов.
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.PostHandler(repo))
		r.Get("/{id}", handlers.GetHandler(repo))
	})

	return r
}

func main() {
	// Запуск сервера.
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", ShortenerRouter())
}
