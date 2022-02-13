package filestorage

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/config"
	serviceModel "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/model"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage/errors"
	storageModel "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage/model"
)

var _ storage.URLStorer = (*Storage)(nil)

type Storage struct {
	mu  sync.Mutex
	Cfg *config.StorageConfig
	DB  map[string]storageModel.URLInDB
}

func NewStorage(Cfg *config.StorageConfig) (*Storage, error) {
	db := make(map[string]storageModel.URLInDB)

	st := Storage{
		DB:  db,
		Cfg: Cfg,
	}

	err := st.restoreFromFile()

	return &st, err
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

	if len(urls) == 0 {
		return nil, &errors.StorageEmptyResultError{ID: userID}
	}

	return urls, nil
}

func (s *Storage) SaveBatchShortedURL(ctx context.Context, userID string,urls []serviceModel.FullURL) (err error) {
	s.mu.Lock()

	for _, url := range urls {
		urlInDb := storageModel.URLInDB{
			URL:    url.OriginalURL,
			UserID: userID,
		}

		s.DB[url.ShortURL] = urlInDb
	}

	s.mu.Unlock()

	return nil
}

func (s *Storage) PersistStorage() error {
	if len(s.DB) == 0 {
		return nil
	}

	file, err := os.OpenFile(s.Cfg.FileStoragePath, os.O_RDWR|os.O_CREATE, 0777)

	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	var rows []storageModel.RowInURLStorage

	for ID, URLInDB := range s.DB {
		rowToEncode := storageModel.RowInURLStorage{
			ID:     ID,
			URL:    URLInDB.URL,
			UserID: URLInDB.UserID,
		}

		rows = append(rows, rowToEncode)
	}

	err = encoder.Encode(rows)

	if err != nil {
		return err
	}

	log.Print("Storage persisted")

	return nil
}

func (s *Storage) Ping() error {
	return nil
}

func (s *Storage) Close() error {
	return nil
}

func (s *Storage) restoreFromFile() error {
	var rows []storageModel.RowInURLStorage

	file, err := os.OpenFile(s.Cfg.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0777)

	if err != nil {
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	decoder.Decode(&rows)

	log.Print("Restored from file")

	for _, row := range rows {
		urlInDb := storageModel.URLInDB{
			URL:    row.URL,
			UserID: row.UserID,
		}

		s.DB[row.ID] = urlInDb
	}

	return nil
}
