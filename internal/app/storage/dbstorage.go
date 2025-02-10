package storage

import (
	"context"

	"github.com/alexuryumtsev/go-shortener/internal/app/db"
	"github.com/alexuryumtsev/go-shortener/internal/app/models"
)

// DatabaseStorage управляет сохранением и получением данных в базе данных.
type DatabaseStorage struct {
	db *db.Database
}

// NewDatabaseStorage создаёт новое хранилище для базы данных.
func NewDatabaseStorage(db *db.Database) *DatabaseStorage {
	return &DatabaseStorage{db: db}
}

// Save сохраняет URL в базе данных.
func (s *DatabaseStorage) Save(ctx context.Context, urlModel models.URLModel) error {
	return nil
}

// Get возвращает оригинальный URL по идентификатору из базы данных.
func (s *DatabaseStorage) Get(ctx context.Context, id string) (models.URLModel, bool) {
	query := `SELECT original_url FROM urls WHERE short_url = $1`
	row := s.db.Pool.QueryRow(ctx, query, id)

	var urlModel models.URLModel
	urlModel.ID = id
	err := row.Scan(&urlModel.URL)
	if err != nil {
		return models.URLModel{}, false
	}
	return urlModel, true
}

// LoadFromFile загружает данные из базы данных (не требуется для базы данных).
func (s *DatabaseStorage) LoadFromFile() error {
	return nil
}

// Ping проверяет соединение с базой данных.
func (s *DatabaseStorage) Ping(ctx context.Context) error {
	return s.db.Ping(ctx)
}
