package config

import (
	"flag"
	"fmt"
	"net/url"
	"regexp"
)

type Config struct {
	ServerAddress string // Адрес запуска HTTP-сервера
	BaseURL       string // Базовый адрес для сокращённых URL
}

func InitConfig() (*Config, error) {
	cfg := &Config{}

	// Определяем флаги
	flag.StringVar(&cfg.ServerAddress, "a", ":8080", "HTTP server address, host:port")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080/", "Base URL for shortened links")

	// Обрабатываем флаги
	flag.Parse()

	// Проверяем значения флагов
	if cfg.ServerAddress != "" {
		// Проверка формата host:port
		hostPortPattern := `^([a-zA-Z0-9.-]+)?(:[0-9]+)$`
		matched, err := regexp.MatchString(hostPortPattern, cfg.ServerAddress)
		if err != nil || !matched {
			return nil, fmt.Errorf("invalid server address format, expected host:port")
		}
	}

	if cfg.BaseURL != "" {
		_, err := url.ParseRequestURI(cfg.BaseURL)
		if err != nil {
			return nil, fmt.Errorf("invalid base URL: %v", err)
		}
	}

	return cfg, nil
}
