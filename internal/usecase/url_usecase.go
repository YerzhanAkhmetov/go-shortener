package usecase

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/YerzhanAkhmetov/go-shortener/internal/domain"
	"github.com/YerzhanAkhmetov/go-shortener/internal/repository"
)

// URLUsecase представляет интерфейс для работы с сокращением URL.
type URLUsecase interface {
	// Create создает новую запись сокращенного URL на основе оригинального URL.
	// Возвращает созданный объект domain.URL и ошибку, если возникла.
	Create(originalURL string) (domain.URL, error)

	// GetByID возвращает сокращенный URL по его идентификатору.
	// Возвращает найденный объект domain.URL и флаг, указывающий на наличие записи.
	GetByID(id string) (domain.URL, bool)
}

// urlUsecase представляет реализацию интерфейса URLUsecase.
type urlUsecase struct {
	repo repository.URLRepository // Репозиторий для работы с хранилищем URL
}

// NewURLUsecase создает новый экземпляр URLUsecase.
func NewURLUsecase(repo repository.URLRepository) URLUsecase {
	return &urlUsecase{repo: repo}
}

// Create создает новую запись сокращенного URL на основе оригинального URL.
// Генерирует уникальный идентификатор, сохраняет URL в репозитории.
// Возвращает созданный объект domain.URL и ошибку, если возникла.
func (u *urlUsecase) Create(originalURL string) (domain.URL, error) {
	id, err := generateID() // Генерация уникального идентификатора
	if err != nil {
		return domain.URL{}, err
	}
	url := domain.URL{
		ID:          id,
		OriginalURL: originalURL,
	}
	err = u.repo.Save(url) // Сохранение URL в репозитории
	if err != nil {
		return domain.URL{}, err
	}
	return url, nil
}

// GetByID возвращает сокращенный URL по его идентификатору.
// Использует репозиторий для поиска URL по заданному идентификатору.
// Возвращает найденный объект domain.URL и флаг, указывающий на наличие записи.
func (u *urlUsecase) GetByID(id string) (domain.URL, bool) {
	return u.repo.FindByID(id)
}

// generateID генерирует случайный уникальный идентификатор длиной 6 байт.
// Использует криптографический рандом для генерации.
// Возвращает сгенерированный идентификатор в виде строки и ошибку, если возникла.
func generateID() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b) // Генерация случайных байт с использованием криптографического рандома
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil // Кодирование в base64 для представления в виде строки
}
