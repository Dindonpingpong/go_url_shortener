package storage

import "context"

type URLSaver interface{
	SaveShortedURL(ctx context.Context, url string, shortedURL string) error
}

type URLGetter interface {
	GetURL(ctx context.Context, shortedURL string) (url string, err error)
}

type URLStorer interface {
	URLSaver
	URLGetter
}

type Persister interface {
	PersistStorage() error
}

type URLPersistanceStorer interface {
	URLStorer
	Persister
}