package storage

import "github.com/alexuryumtsev/go-shortener/internal/app/models"

type MockStorage struct {
	data map[string]models.URLModel
}

func NewMockStorage() *MockStorage {
	return &MockStorage{data: make(map[string]models.URLModel)}
}

func (m *MockStorage) Save(urlModel models.URLModel) {
	m.data[urlModel.ID] = urlModel
}

func (m *MockStorage) Get(id string) (models.URLModel, bool) {
	urlModel, exists := m.data[id]
	return urlModel, exists
}
