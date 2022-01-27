package inmemory

import (
	"context"
	"sync"

	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage/errors"
)

var _ storage.URLStorer = (*Storage)(nil)

type Storage struct {
	mu sync.Mutex
	DB map[string]string
}

func NewStorage() *Storage {
	db := make(map[string]string)

	return &Storage{DB: db}
}

func (s *Storage) GetURL(ctx context.Context, shortedURL string) (url string, err error) {
	s.mu.Lock()
	URL, ok := s.DB[shortedURL]

	defer s.mu.Unlock()

	if !ok {
		return "", &errors.StorageEmptyResultError{ID: shortedURL}
	}

	return URL, nil
}

func (s *Storage) SaveShortedURL(ctx context.Context, url string, shortedURL string) error {
	s.mu.Lock()
	s.DB[shortedURL] = url
	s.mu.Unlock()

	return nil
}
