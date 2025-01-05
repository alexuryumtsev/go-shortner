package storage

type MockStorage struct {
	data map[string]string
}

func NewMockStorage() *MockStorage {
	return &MockStorage{data: make(map[string]string)}
}

func (m *MockStorage) Save(id, url string) {
	m.data[id] = url
}

func (m *MockStorage) Get(id string) (string, bool) {
	url, exists := m.data[id]
	return url, exists
}
