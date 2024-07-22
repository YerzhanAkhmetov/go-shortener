package repository

import (
	"github.com/YerzhanAkhmetov/go-shortener/internal/domain"
	"github.com/YerzhanAkhmetov/go-shortener/internal/storage"
)

type MemoryURLRepository struct {
	storage storage.Storage
}

func NewURLRepository(storage storage.Storage) URLRepository {
	return &MemoryURLRepository{storage: storage}
}

func (r *MemoryURLRepository) Save(url domain.URL) error {
	r.storage.SaveURL(url.ID, url.OriginalURL)
	return nil
}

func (r *MemoryURLRepository) FindByID(id string) (domain.URL, bool) {
	originalURL, exists := r.storage.GetURL(id)
	if !exists {
		return domain.URL{}, false
	}
	return domain.URL{ID: id, OriginalURL: originalURL}, true
}
