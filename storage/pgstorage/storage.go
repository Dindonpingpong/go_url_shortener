package pgstorage

import (
	"context"
	"errors"

	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/config"
	serviceModel "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/model"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage"
	storageErrors "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage/errors"
	pgModel "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage/pgstorage/model"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func (s *Storage) GetURL(ctx context.Context, userID string, shortedURL string) (url string, err error) {
	var queryResult pgModel.URLInDB

	query := "SELECT * FROM urls WHERE short_url = $1 AND user_id = $2"

	err = s.db.GetContext(ctx, &queryResult, query, shortedURL, userID)

	if queryResult.IsDeleted {
		return "", &storageErrors.StorageDeletedError{ShortURL: shortedURL}
	} 

	return url, err
}

func (s *Storage) SaveShortedURL(ctx context.Context, url string, userID string, shortedURL string) error {
	query := "INSERT INTO urls (user_id, url, short_url) VALUES ($1, $2, $3)"

	var pgErr *pq.Error

	_, err := s.db.ExecContext(ctx, query, userID, url, shortedURL)

	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return &storageErrors.StorageAlreadyExistsError{ShortURL: shortedURL}
	}

	return err
}

func (s *Storage) GetURLsByuserID(ctx context.Context, userID string) (urls []serviceModel.FullURL, err error) {
	var queryResult []pgModel.URLInDB

	query := "SELECT * FROM urls WHERE user_id = $1 AND is_deleted = false"

	err = s.db.SelectContext(ctx, &queryResult, query, userID)

	if err != nil {
		return nil, err
	}

	if len(queryResult) == 0 {
		return nil, &storageErrors.StorageEmptyResultError{ID: userID}
	}

	for _, urlInDB := range queryResult {
		fullURL := serviceModel.FullURL{
			OriginalURL: urlInDB.URL,
			ShortURL:    urlInDB.ShortURL,
		}

		urls = append(urls, fullURL)
	}

	return urls, nil
}

func (s *Storage) SaveBatchShortedURL(ctx context.Context, userID string, urls []serviceModel.FullURL) (err error) {
	var query = "INSERT INTO urls (user_id, url, short_url) VALUES ($1, $2, $3)"

	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	for _, url := range urls {
		_, err = tx.ExecContext(
			ctx,
			query,
			userID,
			url.OriginalURL,
			url.ShortURL,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Storage) DeleteSoftBatchShortedURL(ctx context.Context, userID string, shortedURLs []string) error {
	var query = "UPDATE urls SET is_deleted = true WHERE user_id = $1 AND short_url = $2"

	tx, err := s.db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	for _, url := range shortedURLs {
		_, err = tx.ExecContext(
			ctx,
			query,
			userID,
			url,
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
		is_deleted boolean not null DEFAULT false,
		short_url text not null unique 
	);`

	_, err := db.Exec(query)

	return err
}
