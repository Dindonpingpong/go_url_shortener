package storage

import (
	"context"

	serviceModel "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/model"
)

type URLSaver interface {
	SaveShortedURL(ctx context.Context, url string, userId string, shortedURL string) error
}

type URLGetter interface {
	GetURL(ctx context.Context, shortedURL string) (url string, err error)
}

type URLsByUserIDGetter interface {
	GetURLsByUserID(ctx context.Context, userID string) (urls []serviceModel.FullURL, err error)
}

type URLsBatchSaver interface {
	SaveBatchShortedURL(ctx context.Context, userID string, urls []serviceModel.FullURL) (err error)
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
	URLsByUserIDGetter
	URLsBatchSaver
	Persister
	Pinger
	Closer
}


