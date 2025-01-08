package config

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
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
	hostPortPattern := `^([a-zA-Z0-9.-]+)?(:[0-9]+)$`
	matched, err := regexp.MatchString(hostPortPattern, cfg.ServerAddress)
	if err != nil || !matched {
		return nil, fmt.Errorf("invalid server address format, expected host:port")
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = envBaseURL
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "http://localhost:8080/" // Значение по умолчанию.
	}

	_, err = url.ParseRequestURI(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %v", err)
	}

	return cfg, nil
}
