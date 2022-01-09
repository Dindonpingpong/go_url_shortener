package handlers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	shortenerService "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/shortener"
	"github.com/go-chi/chi"
)

type URLHandler struct {
	svc shortenerService.Processor
}

func NewURLHandler(svc shortenerService.Processor) (*URLHandler, error) {
	if svc == nil {
		return nil, fmt.Errorf("shortenerService: nil")
	}

	return &URLHandler{svc: svc}, nil
}

func (h *URLHandler) HandleGetURL() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		urlID := chi.URLParam(r, "urlID")

		ctx := context.Background()
		url, err := h.svc.GetURL(ctx, urlID)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusNotFound)
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
		rw.Write([]byte("http://" + r.Host + "/" + id))
	}
}
