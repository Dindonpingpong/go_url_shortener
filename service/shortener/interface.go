package shortener

import (
	"context"

	serviceModel "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/model"
)

type Processor interface {
	SaveURL(ctx context.Context, rawURL string, userId string) (id string, err error)
	GetURL(ctx context.Context, id string) (url string, err error)
	GetURLsByUserID(ctx context.Context, userID string) (urls []serviceModel.FullURL, err error)
}
