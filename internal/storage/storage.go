package storage

type Storage interface {
	SaveURL(id, url string)
	GetURL(id string) (string, bool)
}
