package rest

import (
	"context"
	"net/http"

	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/api/rest/handlers"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/api/rest/middlewares"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/config"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/secretary/v1"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/shortener/v1"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage"
	"github.com/go-chi/chi"
)

func InitServer(ctx context.Context, cfg *config.Config, st storage.URLStorer) (server *http.Server, err error) {
	shortenerService, err := shortener.NewShortenerService(st)

	if err != nil {
		return nil, err
	}

	urlHandler, err := handlers.NewURLHandler(shortenerService, cfg.ServerConfig)

	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()

	secretaryService, err := secretary.NewSecretaryService(cfg.SecretConfig)

	if err != nil {
		return nil, err
	}

	cookieHandler, err := middlewares.NewCookieHandler(secretaryService, cfg.SecretConfig)

	if err != nil {
		return nil, err
	}

	r.Use(cookieHandler.AuthCookieHandle)
	r.Use(middlewares.CompressHandle)
	r.Use(middlewares.DecompressHandle)
	r.Get("/{urlID}", urlHandler.HandleGetURL())
	r.Get("/user/urls", urlHandler.HandleGetURLsByUserID())
	r.Post("/api/shorten", urlHandler.JSONHandlePostURL())
	r.Post("/", urlHandler.HandlePostURL())

	return &http.Server{
		Addr: cfg.ServerConfig.ServerAddress,
		Handler: r,
	}, nil
}
