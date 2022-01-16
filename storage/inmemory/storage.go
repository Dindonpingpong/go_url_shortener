package inmemory

import (
	"context"

	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage/errors"
)

var _ storage.URLStorer = (*Storage)(nil)

type Storage struct {
	DB map[string]string
}

func NewStorage() *Storage {
	db := make(map[string]string)

	return &Storage{DB: db}
}

func (s *Storage) GetURL(ctx context.Context, shortedURL string) (url string, err error) {
	URL, ok := s.DB[shortedURL]

	if !ok {
		return "", &errors.StorageEmptyResultError{ID: shortedURL}
	}

	return URL, nil
}

func (s *Storage) SaveShortedURL(ctx context.Context, url string, shortedURL string) error {
	s.DB[shortedURL] = url

	return nil
}
// Errors types, generete service