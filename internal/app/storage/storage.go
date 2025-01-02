package storage

import "sync"

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
func (s *Storage) Save(id, url string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[id] = url
}

// Get возвращает оригинальный URL по идентификатору.
func (s *Storage) Get(id string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, exists := s.data[id]
	return url, exists
}
