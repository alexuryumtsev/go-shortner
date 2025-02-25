package url

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

// URLService интерфейс для сокращения URL
type URLService interface {
	ShortenerURL(originalURL, userID string) (string, error)
	SaveBatchShortenerURL(batchModels []models.URLBatchModel, userID string) ([]string, error)
}

// urlService реализация URLService
type urlService struct {
	ctx     context.Context
	storage storage.URLWriter
	baseURL string
}

// NewURLService конструктор URLService
func NewURLService(ctx context.Context, storage storage.URLWriter, baseURL string) URLService {
	return &urlService{
		ctx:     ctx,
		storage: storage,
		baseURL: strings.TrimSuffix(baseURL, "/"),
	}
}

// ShortenerURL сокращает URL и сохраняет в базе
func (s *urlService) ShortenerURL(originalURL, userID string) (string, error) {
	if originalURL == "" {
		return "", fmt.Errorf("empty URL")
	}

	id := GenerateID(originalURL)
	urlModel := models.URLModel{ID: id, URL: originalURL, UserID: userID}
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

// SaveBatchShortenerURL сокращает пакет URL
func (s *urlService) SaveBatchShortenerURL(batchModels []models.URLBatchModel, userID string) ([]string, error) {
	var urlModels []models.URLModel
	for _, req := range batchModels {
		urlModels = append(urlModels, models.URLModel{
			ID:     GenerateID(req.OriginalURL),
			URL:    req.OriginalURL,
			UserID: userID,
		})
	}

	var shortenedURLs []string
	for _, urlModel := range urlModels {
		shortenedURLs = append(shortenedURLs, s.baseURL + "/" + urlModel.ID)
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

// GenerateID создает короткий идентификатор для URL
func GenerateID(url string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(url)))[:8]
}
