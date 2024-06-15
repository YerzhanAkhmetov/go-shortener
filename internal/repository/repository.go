package repository

import (
	"github.com/YerzhanAkhmetov/go-shortener/internal/domain"
)

type URLRepository interface {
	Save(url domain.URL) error
	FindByID(id string) (domain.URL, bool)
}
