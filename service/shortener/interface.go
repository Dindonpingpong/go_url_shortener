package shortener

import (
	"context"

	serviceModel "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/model"
)

type Processor interface {
	SaveURL(ctx context.Context, rawURL string, userID string) (id string, err error)
	GetURL(ctx context.Context, id string) (url string, err error)
	GetURLsByuserID(ctx context.Context, userID string) (urls []serviceModel.FullURL, err error)
	SaveBatchShortedURL(ctx context.Context, userID string, urls []string) (savedUrls []serviceModel.FullURL, err error)
	DeleteBatchShortedURL(ctx context.Context, userID string, shortedURLs[]string)
	PingStorage() error
}
