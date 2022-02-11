package pgstorage

import (
	"context"

	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/config"
	serviceModel "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/model"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var _ storage.URLStorer = (*Storage)(nil)

type Storage struct {
	Cfg *config.StorageConfig
	db  *sqlx.DB
}

func NewStorage(cfg *config.StorageConfig) (*Storage, error) {
	db, err := sqlx.Open("postgres", cfg.DatabaseDSN)

	if err != nil {
		return nil, err
	}

	err = migrate(db)

	if err != nil {
		return nil, err
	}

	return &Storage{db: db, Cfg: cfg}, nil
}

func (s *Storage) GetURL(ctx context.Context, shortedURL string) (url string, err error) {
	return "", nil
}

func (s *Storage) SaveShortedURL(ctx context.Context, url string, userId string, shortedURL string) error {
	return nil
}

func (s *Storage) GetURLsByUserID(ctx context.Context, userID string) (urls []serviceModel.FullURL, err error) {
	return urls, nil
}

func (s *Storage) PersistStorage() error {
	return nil
}

func (s *Storage) Ping() error {
	return s.db.Ping()
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func migrate(db *sqlx.DB) error {
	query := `CREATE TABLE IF NOT EXISTS urls (
		id bigserial not null,
		user_id uuid not null,
		url text not null unique,
		short_url text not null unique 
	);`

	_, err := db.Exec(query)

	return err
}
