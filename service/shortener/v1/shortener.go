package shortener

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/pkg/short"
	sertviceErrors "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/errors"
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

func (s *Shortener) SaveURL(ctx context.Context, rawURL string, userId string) (id string, err error) {
	_, err = url.ParseRequestURI(rawURL)

	if err != nil {
		return "", &sertviceErrors.ServiceBusinessError{Msg: "incorrect url"}
	}

	shortURL, err := short.GenereteShortString(rawURL)

	if err != nil {
		return "", &sertviceErrors.ServiceBusinessError{Msg: "incorrect url"}
	}

	err = s.urlStorer.SaveShortedURL(ctx, rawURL, userId, shortURL)

	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (s *Shortener) GetURL(ctx context.Context, id string) (url string, err error) {
	url, err = s.urlStorer.GetURL(ctx, id)

	if err != nil {
		switch err.(type) {
		default:
			return "", err
		case *storageErrors.StorageEmptyResultError:
			return "", &sertviceErrors.ServiceBusinessError{Msg: err.Error()}
		}
	}

	return url, nil
}

func (s *Shortener) GetURLsByUserID(ctx context.Context, userID string) (urls []serviceModel.FullURL, err error) {
	urls, err = s.urlStorer.GetURLsByUserID(ctx, userID)

	if err != nil {
		return nil, err
	}

	return urls, nil
}

func (s *Shortener) PingStorage() error {
	return s.urlStorer.Ping()
}