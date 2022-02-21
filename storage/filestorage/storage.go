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

func (s *Storage) GetURL(ctx context.Context, userID string, shortedURL string) (string, error)  {
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

	for ID, urlInDB := range s.DB {
		rowToEncode := storageModel.RowInURLStorage{
			ID:     ID,
			URL:    urlInDB.URL,
			UserID: urlInDB.UserID,
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
		urlInDB := storageModel.URLInDB{
			URL:    row.URL,
			UserID: row.UserID,
		}

		s.DB[row.ID] = urlInDB
	}

	return nil
}
