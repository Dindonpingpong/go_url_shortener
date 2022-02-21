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

func (s *Storage) GetURL(ctx context.Context, shortedURL string) (string, error) {
	s.mu.Lock()
	urlInDB, ok := s.DB[shortedURL]

	defer s.mu.Unlock()

	if !ok {
		return "", &errors.StorageEmptyResultError{ID: shortedURL}
	}

	return urlInDB.URL, nil
}

func (s *Storage) SaveShortedURL(ctx context.Context, url string, userID string, shortedURL string) error {
	s.mu.Lock()

	urlInDB := storageModel.URLInDB{
		URL:    url,
		UserID: userID,
	}
	s.DB[shortedURL] = urlInDB
	s.mu.Unlock()

	return nil
}

func (s *Storage) GetURLsByuserID(ctx context.Context, userID string) (urls []serviceModel.FullURL, err error) {
	for shortedURL, url := range s.DB {
		if url.UserID == userID {
			fullURL := serviceModel.FullURL{
				OriginalURL: url.URL,
				ShortURL:    shortedURL,
			}
			urls = append(urls, fullURL)
		}
	}

	if len(urls) == 0 {
		return nil, &errors.StorageEmptyResultError{ID: userID}
	}

	return urls, nil
}

func (s *Storage) SaveBatchShortedURL(ctx context.Context, userID string, urls []serviceModel.FullURL) (err error) {
	s.mu.Lock()

	for _, url := range urls {
		urlInDB := storageModel.URLInDB{
			URL:    url.OriginalURL,
			UserID: userID,
		}

		s.DB[url.ShortURL] = urlInDB
	}

	s.mu.Unlock()

	return nil
}

func (s *Storage) DeleteSoftBatchShortedURL(ctx context.Context, userID string, shortedURLs []string) error {
	return nil
}

func (s *Storage) PersistStorage() error {
	return nil
}

func (s *Storage) Ping() error {
	return nil
}

func (s *Storage) Close() error {
	return nil
}
