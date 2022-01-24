package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/api/rest/model"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/config"
	sertviceErrors "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/errors"
	shortenerService "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/shortener"
	"github.com/go-chi/chi"
)

type URLHandler struct {
	svc          shortenerService.Processor
	serverConfig *config.ServerConfig
}

func NewURLHandler(svc shortenerService.Processor, serverConfig *config.ServerConfig) (*URLHandler, error) {
	if svc == nil {
		return nil, fmt.Errorf("shortenerService: nil")
	}

	return &URLHandler{svc: svc, serverConfig: serverConfig}, nil
}

func (h *URLHandler) HandleGetURL() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		urlID := chi.URLParam(r, "urlID")

		ctx := context.Background()
		url, err := h.svc.GetURL(ctx, urlID)

		if err != nil {
			switch err.(type) {
			default:
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			case *sertviceErrors.ServiceBusinessError:
				http.Error(rw, err.Error(), http.StatusNotFound)
			}
			return
		}

		rw.Header().Set("Location", url)
		rw.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func (h *URLHandler) HandlePostURL() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		id, err := h.svc.SaveURL(ctx, string(b))

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		rw.WriteHeader(http.StatusCreated)

		u := &url.URL{
			Scheme: "http",
			Host:   h.serverConfig.BaseURL,
			Path:   id,
		}

		rw.Write([]byte(u.String()))
	}
}

func (h *URLHandler) JSONHandlePostURL() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var post model.RequestURL

		rContentType := r.Header.Get("Content-Type")

		if rContentType != "application/json" {
			http.Error(rw, "Invalid Content-Type", http.StatusBadRequest)
		}

		b, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		json.Unmarshal(b, &post)

		ctx := context.Background()
		id, err := h.svc.SaveURL(ctx, post.URL)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		u := &url.URL{
			Scheme: "h.serverConfig.BaseURL",
			Path:   id,
		}

		resData := model.ResponseURL{
			ShortURL: u.String(),
		}

		resBody, err := json.Marshal(resData)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)
		rw.Write(resBody)
	}
}
