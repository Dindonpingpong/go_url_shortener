package rest

import (
	"context"
	"net/http"

	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/api/rest/handlers"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/shortener/v1"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage/inmemory"
	"github.com/go-chi/chi"
)

func InitServer(ctx context.Context) (server *http.Server, err error) {
	storage := inmemory.NewStorage()

	shortenerService, err := shortener.NewShortenerService(storage)

	if err != nil {
		return nil, err
	}

	urlHandler, err := handlers.NewURLHandler(shortenerService)

	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()

	r.Get("/{urlID}", urlHandler.HandleGetURL())
	r.Post("/", urlHandler.HandlePostURL())

	return &http.Server{
		Addr: ":8080",
		Handler: r,
	}, nil
}