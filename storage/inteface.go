package storage

import (
	"context"

	serviceModel "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/model"
)

type URLSaver interface {
	SaveShortedURL(ctx context.Context, url string, userID string, shortedURL string) error
}

type URLGetter interface {
	GetURL(ctx context.Context, userID string, shortedURL string) (string, error)
}

type URLsByuserIDGetter interface {
	GetURLsByuserID(ctx context.Context, userID string) (urls []serviceModel.FullURL, err error)
}

type URLsBatchSaver interface {
	SaveBatchShortedURL(ctx context.Context, userID string, urls []serviceModel.FullURL) (err error)
}

type URLsBatchDeleter interface {
	DeleteSoftBatchShortedURL(ctx context.Context, userID string, shortedURLs []string) error
}

type Persister interface {
	PersistStorage() error
}

type Pinger interface {
	Ping() error
}

type Closer interface {
	Close() error
}
type URLStorer interface {
	URLSaver
	URLGetter
	URLsByuserIDGetter
	URLsBatchSaver
	URLsBatchDeleter
	Persister
	Pinger
	Closer
}
