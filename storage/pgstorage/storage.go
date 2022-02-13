package pgstorage

import (
	"context"

	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/config"
	serviceModel "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/model"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage/errors"
	pgModel "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage/pgstorage/model"
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
	query := "SELECT url FROM urls WHERE short_url = $1"

	err = s.db.GetContext(ctx, &url, query, shortedURL)

	return url, err
}

func (s *Storage) SaveShortedURL(ctx context.Context, url string, userId string, shortedURL string) error {
	query := "INSERT INTO urls (user_id, url, short_url) VALUES ($1, $2, $3)"

	_, err := s.db.ExecContext(ctx, query, userId, url, shortedURL)

	return err
}

func (s *Storage) GetURLsByUserID(ctx context.Context, userID string) (urls []serviceModel.FullURL, err error) {
	var queryResult []pgModel.URLInDB

	query := "SELECT * FROM urls WHERE user_id = $1"

	err = s.db.SelectContext(ctx, &queryResult, query, userID)

	if err != nil {
		return nil, err
	}

	if len(queryResult) == 0 {
		return nil, &errors.StorageEmptyResultError{ID: userID}
	}

	for _, urlInDb := range queryResult {
		fullURL := serviceModel.FullURL{
			OriginalURL: urlInDb.Url,
			ShortURL:    urlInDb.ShortURL,
		}

		urls = append(urls, fullURL)
	}

	return urls, nil
}

func (s *Storage) SaveBatchShortedURL(ctx context.Context, userID string, urls []serviceModel.FullURL) (err error) {
	var query = "INSERT INT urls (user_id, url, short_url) 	VALUES ($1, $2, $3)"

	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	for _, url := range urls {
		_, err = tx.ExecContext(
			ctx,
			query,
			userID,
			url.ShortURL,
			url.OriginalURL,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
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
