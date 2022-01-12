package shortener

import "context"

type Processor interface {
	SaveURL(ctx context.Context, url string) (id string ,err error)
	GetURL(ctx context.Context, id string) (url string, err error)
}