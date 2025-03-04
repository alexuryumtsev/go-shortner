package pg

import (
	"context"
	"fmt"

	"github.com/alexuryumtsev/go-shortener/internal/app/db"
	"github.com/alexuryumtsev/go-shortener/internal/app/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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
	query := `INSERT INTO urls (short_url, original_url) VALUES ($1, $2)`
	_, err := s.db.Pool.Exec(ctx, query, urlModel.ID, urlModel.URL)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return fmt.Errorf("conflict: %w", err)
		}
		return fmt.Errorf("failed to save URL: %w", err)
	}
	return nil
}

// SaveBatch сохраняет множество URL в базе данных.
func (s *DatabaseStorage) SaveBatch(ctx context.Context, urlModels []models.URLModel) error {
	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, urlModel := range urlModels {
		query := `INSERT INTO urls (short_url, original_url) VALUES ($1, $2) ON CONFLICT (short_url) DO NOTHING`
		_, err := tx.Exec(ctx, query, urlModel.ID, urlModel.URL)
		if err != nil {
			return fmt.Errorf("failed to save URL: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
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
