package filestorage

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/config"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage/errors"
	storageModel "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage/model"
)

var _ storage.URLPersistanceStorer = (*Storage)(nil)

type Storage struct {
	mu  sync.Mutex
	Cfg *config.StorageConfig
	DB  map[string]string
}

func NewStorage(Cfg *config.StorageConfig) (*Storage, error) {
	db := make(map[string]string)

	st := Storage{
		DB:  db,
		Cfg: Cfg,
	}

	err := st.restoreFromFile()

	return &st, err
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

	for ID, URL := range s.DB {
		rowToEncode := storageModel.RowInURLStorage{
			ID:  ID,
			URL: URL,
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
		s.DB[row.ID] = row.URL
	}

	return nil
}
