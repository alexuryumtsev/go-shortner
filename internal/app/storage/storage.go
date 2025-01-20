package storage

import (
	"sync"

	"github.com/alexuryumtsev/go-shortener/internal/app/models"
)

type URLStorage interface {
	Save(urlModel models.URLModel)
	Get(id string) (models.URLModel, bool)
}

// Storage управляет сохранением и получением данных.
type Storage struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewStorage создаёт новое хранилище.
func NewStorage() *Storage {
	return &Storage{data: make(map[string]string)}
}

// Save сохраняет URL.
func (s *Storage) Save(urlModel models.URLModel) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[urlModel.ID] = urlModel.URL
}

// Get возвращает оригинальный URL по идентификатору.
func (s *Storage) Get(id string) (models.URLModel, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, exists := s.data[id]
	return models.URLModel{ID: id, URL: url}, exists
}
