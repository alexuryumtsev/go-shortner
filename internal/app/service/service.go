package service

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/alexuryumtsev/go-shortener/internal/app/storage"
)

type URLService struct {
	storage storage.URLStorage
	baseURL string
}

func NewURLService(storage storage.URLStorage, baseURL string) *URLService {
	return &URLService{
		storage: storage,
		baseURL: strings.TrimSuffix(baseURL, "/"),
	}
}

func (s *URLService) ShortenerURL(originalURL string) (string, error) {
	if originalURL == "" {
		return "", fmt.Errorf("empty URL")
	}

	id := generateID(originalURL)
	s.storage.Save(id, originalURL)

	shortenedURL := s.baseURL + "/" + id
	return shortenedURL, nil

}

func generateID(url string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(url)))[:8] // Используем MD5 и берём первые 8 символов.
}
