package service

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"strings"

	"github.com/alexuryumtsev/go-shortener/internal/app/models"
	"github.com/alexuryumtsev/go-shortener/internal/app/storage"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type URLService struct {
	ctx     context.Context
	storage storage.URLWriter
	baseURL string
}

func NewURLService(ctx context.Context, storage storage.URLWriter, baseURL string) *URLService {
	return &URLService{
		ctx:     ctx,
		storage: storage,
		baseURL: strings.TrimSuffix(baseURL, "/"),
	}
}

func (s *URLService) ShortenerURL(originalURL string) (string, error) {
	if originalURL == "" {
		return "", fmt.Errorf("empty URL")
	}

	id := GenerateID(originalURL)
	urlModel := models.URLModel{ID: id, URL: originalURL}
	shortenedURL := s.baseURL + "/" + id

	err := s.storage.Save(s.ctx, urlModel)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return shortenedURL, err
		}
	}

	return shortenedURL, nil
}

func (s *URLService) SaveBatchShortenerURL(batchModels []models.URLBatchModel) ([]string, error) {
	var urlModels []models.URLModel
	for _, req := range batchModels {
		urlModels = append(urlModels, models.URLModel{
			ID:  GenerateID(req.OriginalURL), // Функция для генерации короткого ID
			URL: req.OriginalURL,
		})
	}

	var shortenedURLs []string
	for _, urlModel := range urlModels {
		shortenedURLs = append(shortenedURLs, s.baseURL+"/"+urlModel.ID)
	}

	err := s.storage.SaveBatch(s.ctx, urlModels)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return shortenedURLs, err
		}
		return nil, err
	}

	return shortenedURLs, nil
}

func GenerateID(url string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(url)))[:8] // Используем MD5 и берём первые 8 символов.
}
