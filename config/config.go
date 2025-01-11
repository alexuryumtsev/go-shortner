package config

import (
	"flag"
	"os"

	"github.com/alexuryumtsev/go-shortener/internal/app/validator"
)

type Config struct {
	ServerAddress string // Адрес запуска HTTP-сервера
	BaseURL       string // Базовый адрес для сокращённых URL
}

func InitConfig() (*Config, error) {
	cfg := &Config{}

	// Получаем значения из переменных окружения.
	envServerAddress := os.Getenv("SERVER_ADDRESS")
	envBaseURL := os.Getenv("BASE_URL")

	// Определяем флаги
	flag.StringVar(&cfg.ServerAddress, "a", "", "HTTP server address, host:port")
	flag.StringVar(&cfg.BaseURL, "b", "", "Base URL for shortened links")

	// Обрабатываем флаги
	flag.Parse()

	// Проверяем значения флагов и переменных окружения
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = envServerAddress
	}
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = ":8080" // Значение по умолчанию.
	}

	// Проверка формата host:port
	err := validator.ValidateServerAddress(cfg.ServerAddress)

	if cfg.BaseURL == "" {
		cfg.BaseURL = envBaseURL
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "http://localhost:8080/" // Значение по умолчанию.
	}

	// Проверка корректности URL
	err = validator.ValidateBaseURL(cfg.BaseURL)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
