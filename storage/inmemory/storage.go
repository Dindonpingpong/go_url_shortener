package inmemory

import (
	"context"
	"sync"

	serviceModel "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/model"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage/errors"
	storageModel "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage/model"
)

var _ storage.URLStorer = (*Storage)(nil)

type Storage struct {
	mu sync.Mutex
	DB map[string]storageModel.URLInDB
}

func NewStorage() *Storage {
	db := make(map[string]storageModel.URLInDB)

	return &Storage{DB: db}
}

func (s *Storage) GetURL(ctx context.Context, shortedURL string) (url string, err error) {
	s.mu.Lock()
	URLInDB, ok := s.DB[shortedURL]

	defer s.mu.Unlock()

	if !ok {
		return "", &errors.StorageEmptyResultError{ID: shortedURL}
	}

	return URLInDB.URL, nil
}

func (s *Storage) SaveShortedURL(ctx context.Context, url string, userId string, shortedURL string) error {
	s.mu.Lock()

	urlInDb := storageModel.URLInDB{
		URL:    url,
		UserID: userId,
	}
	s.DB[shortedURL] = urlInDb
	s.mu.Unlock()

	return nil
}

func (s *Storage) GetURLsByUserID(ctx context.Context, userID string) (urls []serviceModel.FullURL, err error) {
	for shortedURL, url := range s.DB {
		if url.UserID == userID {
			fullURL := serviceModel.FullURL{
				OriginalURL: url.URL,
				ShortURL:    shortedURL,
			}
			urls = append(urls, fullURL)
		}
	}

	return urls, nil
}

func (s *Storage) PersistStorage() error {
	return nil
}
