package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alexuryumtsev/go-shortener/config"
	"github.com/alexuryumtsev/go-shortener/internal/app/handlers"
	"github.com/alexuryumtsev/go-shortener/internal/app/storage"

	"github.com/go-chi/chi/v5"
)

func ShortenerRouter(cfg *config.Config) chi.Router {
	// Инициализация хранилища.
	var repo storage.Repository = storage.NewStorage()

	// Регистрация маршрутов.
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", handlers.PostHandler(repo, cfg.BaseURL))
		r.Get("/{id}", handlers.GetHandler(repo))
	})

	return r
}

func main() {
	// Инициализируем конфигурацию
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// Запуск сервера.
	fmt.Println("Server started at", cfg.ServerAddress)
	http.ListenAndServe(cfg.ServerAddress, ShortenerRouter(cfg))
}
