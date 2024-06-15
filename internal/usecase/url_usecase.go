package usecase

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/YerzhanAkhmetov/go-shortener/internal/domain"
	"github.com/YerzhanAkhmetov/go-shortener/internal/repository"
)

type URLUsecase interface {
	Create(originalURL string) (domain.URL, error)
	GetByID(id string) (domain.URL, bool)
}

type urlUsecase struct {
	repo repository.URLRepository
}

func NewURLUsecase(repo repository.URLRepository) URLUsecase {
	return &urlUsecase{repo: repo}
}

func (u *urlUsecase) Create(originalURL string) (domain.URL, error) {
	id, err := generateID()
	if err != nil {
		return domain.URL{}, err
	}
	url := domain.URL{
		ID:          id,
		OriginalURL: originalURL,
	}
	err = u.repo.Save(url)
	if err != nil {
		return domain.URL{}, err
	}
	return url, nil
}

func (u *urlUsecase) GetByID(id string) (domain.URL, bool) {
	return u.repo.FindByID(id)
}

func generateID() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
