package inmemory

import (
	"context"
	"fmt"

	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage"
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
		return "", fmt.Errorf("URL not found")
	}

	return URL, nil
}

func (s *Storage) SaveShortedURL(ctx context.Context, url string, shortedURL string) error {
	s.DB[shortedURL] = url

	return nil
}
