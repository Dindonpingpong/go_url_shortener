package shortener

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/pkg/short"
	serviceErrors "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/errors"
	serviceModel "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/model"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/shortener"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage"
	storageErrors "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage/errors"
)

type Shortener struct {
	urlStorer storage.URLStorer
}

var _ shortener.Processor = (*Shortener)(nil)

func NewShortenerService(st storage.URLStorer) (*Shortener, error) {
	if st == nil {
		return nil, fmt.Errorf("storage.URLStorer: nil")
	}

	return &Shortener{st}, nil
}

func (s *Shortener) SaveURL(ctx context.Context, rawURL string, userID string) (id string, err error) {
	_, err = url.ParseRequestURI(rawURL)

	if err != nil {
		return "", &serviceErrors.ServiceBusinessError{Msg: "incorrect url"}
	}

	shortURL, err := short.GenereteShortString(rawURL)

	if err != nil {
		return "", &serviceErrors.ServiceBusinessError{Msg: "cannot generate short url"}
	}

	err = s.urlStorer.SaveShortedURL(ctx, rawURL, userID, shortURL)

	if err != nil {
		var storageAlreadyExistsError *storageErrors.StorageAlreadyExistsError

		if errors.As(err, &storageAlreadyExistsError) {
			return shortURL, &serviceErrors.ServiceAlreadyExistsError{Msg: storageAlreadyExistsError.Error()}
		}

		return "", &serviceErrors.ServiceBusinessError{Msg: err.Error()}
	}

	return shortURL, nil
}

func (s *Shortener) GetURL(ctx context.Context, id string) (url string, err error) {
	url, err = s.urlStorer.GetURL(ctx, id)

	if err != nil {
		var storageEmptyResultError *storageErrors.StorageEmptyResultError

		if errors.As(err, &storageEmptyResultError) {
			return "", &serviceErrors.ServiceNotFoundByIDError{ID: storageEmptyResultError.Error()}
		}

		var storageDeletedError *storageErrors.StorageDeletedError

		if errors.As(err, &storageDeletedError) {
			return "", &serviceErrors.ServiceEntityDeletedError{Msg: storageDeletedError.Error()}
		}

		return "", err
	}

	return url, nil
}

func (s *Shortener) GetURLsByuserID(ctx context.Context, userID string) (urls []serviceModel.FullURL, err error) {
	urls, err = s.urlStorer.GetURLsByuserID(ctx, userID)

	if err != nil {
		var storageEmptyResultError *storageErrors.StorageEmptyResultError

		if errors.As(err, &storageEmptyResultError) {
			return nil, &serviceErrors.ServiceNotFoundByIDError{ID: storageEmptyResultError.Error()}
		}

		return nil, err
	}

	return urls, nil
}

func (s *Shortener) SaveBatchShortedURL(ctx context.Context, userID string, urls []string) (savedUrls []serviceModel.FullURL, err error) {
	for _, originalURL := range urls {
		_, err = url.ParseRequestURI(originalURL)

		if err != nil {
			return nil, &serviceErrors.ServiceBusinessError{Msg: "incorrect url"}
		}

		shortURL, err := short.GenereteShortString(originalURL)

		if err != nil {
			return nil, &serviceErrors.ServiceBusinessError{Msg: "incorrect url"}
		}

		fullURL := serviceModel.FullURL{
			OriginalURL: originalURL,
			ShortURL:    shortURL,
		}

		savedUrls = append(savedUrls, fullURL)
	}

	err = s.urlStorer.SaveBatchShortedURL(ctx, userID, savedUrls)

	return savedUrls, err
}

func (s *Shortener) DeleteBatchShortedURL(ctx context.Context, userID string, shortedURLs []string) {
	go func() {
		err := s.urlStorer.DeleteSoftBatchShortedURL(ctx, userID, shortedURLs)

		if err != nil {
			log.Fatal(err)
		}
	}()

}

func (s *Shortener) PingStorage() error {
	return s.urlStorer.Ping()
}
