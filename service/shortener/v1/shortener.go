package shortener

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"

	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/shortener"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage"
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

func (s *Shortener) SaveURL(ctx context.Context, rawURL string) (id string, err error) {
	_, err = url.ParseRequestURI(rawURL)

	if err != nil {
		return "", fmt.Errorf("incorrect url")
	}

	id = generateShortURL(rawURL)

	err = s.urlStorer.SaveShortedURL(ctx, rawURL, id)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Shortener) GetURL(ctx context.Context, id string) (url string, err error) {
	url, err = s.urlStorer.GetURL(ctx, id)

	if err != nil {
		return "", err
	}

	return url, nil
}

func generateShortURL(url string) string {
	h := md5.New()

	h.Write([]byte(url))

	return hex.EncodeToString(h.Sum(nil))
}
