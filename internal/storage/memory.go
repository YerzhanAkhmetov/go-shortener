package storage

type MemoryStorage struct {
	urlStore map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{urlStore: make(map[string]string)}
}

func (m *MemoryStorage) SaveURL(id, url string) {
	m.urlStore[id] = url
}

func (m *MemoryStorage) GetURL(id string) (string, bool) {
	url, exists := m.urlStore[id]
	return url, exists
}
