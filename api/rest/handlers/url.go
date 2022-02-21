package handlers

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/api/rest/middlewares"
	restModel "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/api/rest/model"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/config"
	serviceErrors "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/errors"
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
		userID, err := getuserID(r)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			log.Printf("cannot get userid")
			log.Fatal(err)
			return
		}
		url, err := h.svc.GetURL(ctx, userID, urlID)

		if err != nil {
			var serviceNotFound *serviceErrors.ServiceNotFoundByIDError

			if errors.As(err, &serviceNotFound) {
				http.Error(rw, serviceNotFound.Error(), http.StatusNotFound)
				return
			}

			var serviceURLDeleted *serviceErrors.ServiceEntityDeletedError

			if errors.As(err, &serviceURLDeleted) {
				http.Error(rw, serviceURLDeleted.Error(), http.StatusGone)
				return
			}

			http.Error(rw, err.Error(), http.StatusInternalServerError)
			log.Printf("unknown error from service")
			log.Fatal(err)
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

		userID, err := getuserID(r)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := h.svc.SaveURL(ctx, string(b), userID)

		if err != nil {
			var serviceAlreadyExistsError *serviceErrors.ServiceAlreadyExistsError

			if errors.As(err, &serviceAlreadyExistsError) {
				u, err := createFullURL(h.serverConfig.BaseURL, id)

				if err != nil {
					http.Error(rw, serviceAlreadyExistsError.Error(), http.StatusInternalServerError)
					return
				}

				rw.WriteHeader(http.StatusConflict)
				rw.Write([]byte(u))
				return
			}

			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		rw.WriteHeader(http.StatusCreated)

		u, err := createFullURL(h.serverConfig.BaseURL, id)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Write([]byte(u))
	}
}

func (h *URLHandler) JSONHandlePostURL() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var requestURL restModel.RequestURL

		rContentType := r.Header.Get("Content-Type")

		if rContentType != "application/json" {
			http.Error(rw, "Invalid Content-Type", http.StatusBadRequest)
			return
		}

		b, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		json.Unmarshal(b, &requestURL)

		ctx := context.Background()

		userID, err := getuserID(r)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := h.svc.SaveURL(ctx, requestURL.URL, userID)

		if err != nil {
			var serviceAlreadyExistsError *serviceErrors.ServiceAlreadyExistsError

			if errors.As(err, &serviceAlreadyExistsError) {
				u, err := createFullURL(h.serverConfig.BaseURL, id)

				if err != nil {
					http.Error(rw, serviceAlreadyExistsError.Error(), http.StatusInternalServerError)
					return
				}

				resData := restModel.ResponseURL{
					ShortURL: u,
				}

				resBody, err := json.Marshal(resData)

				if err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}

				rw.Header().Set("Content-Type", "application/json")
				rw.WriteHeader(http.StatusConflict)
				rw.Write([]byte(resBody))
				return
			}

			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		u, err := createFullURL(h.serverConfig.BaseURL, id)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		resData := restModel.ResponseURL{
			ShortURL: u,
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

func (h *URLHandler) HandleGetURLsByuserID() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		userID, err := getuserID(r)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		urls, err := h.svc.GetURLsByuserID(ctx, userID)

		if err != nil {
			var serviceNotFound *serviceErrors.ServiceNotFoundByIDError

			if errors.As(err, &serviceNotFound) {
				http.Error(rw, serviceNotFound.Error(), http.StatusNoContent)
				return
			}

			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		responseURLs := make([]restModel.ResponseFullURL, 0)

		for _, fullURL := range urls {
			u, err := createFullURL(h.serverConfig.BaseURL, fullURL.ShortURL)

			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			responseURL := restModel.ResponseFullURL{
				URL:      fullURL.OriginalURL,
				ShortURL: u,
			}

			responseURLs = append(responseURLs, responseURL)
		}

		resBody, err := json.Marshal(responseURLs)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.Write(resBody)
	}
}

func (h *URLHandler) HandleBatchPostURLs() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rContentType := r.Header.Get("Content-Type")

		if rContentType != "application/json" {
			http.Error(rw, "Invalid Content-Type", http.StatusBadRequest)
			return
		}

		b, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		requestURLs := make([]restModel.RequestBatchItem, 0)

		json.Unmarshal(b, &requestURLs)

		ctx := context.Background()
		userID, err := getuserID(r)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		urlsToSave := make([]string, 0)

		for _, url := range requestURLs {
			urlsToSave = append(urlsToSave, url.OriginalURL)
		}

		savedUrls, err := h.svc.SaveBatchShortedURL(ctx, userID, urlsToSave)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		responseURLs := make([]restModel.ResponseBatchItem, 0)

		for k, item := range requestURLs {
			url := savedUrls[k]

			u, err := createFullURL(h.serverConfig.BaseURL, url.ShortURL)

			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			responseURL := restModel.ResponseBatchItem{
				CorrelationID: item.CorrelationID,
				ShortURL:      u,
			}

			responseURLs = append(responseURLs, responseURL)
		}

		resBody, err := json.Marshal(responseURLs)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)
		rw.Write(resBody)
	}
}

func (h *URLHandler) HandleBatchDeleteURLs() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rContentType := r.Header.Get("Content-Type")

		if rContentType != "application/json" {
			http.Error(rw, "Invalid Content-Type", http.StatusBadRequest)
			return
		}

		b, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		requestURLs := make([]string, 0)

		json.Unmarshal(b, &requestURLs)

		ctx := context.Background()
		userID, err := getuserID(r)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		h.svc.DeleteBatchShortedURL(ctx, userID, requestURLs)

		rw.WriteHeader(http.StatusAccepted)
	}
}

func (h *URLHandler) HandlePing() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := h.svc.PingStorage()

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func getuserID(r *http.Request) (string, error) {
	userCookie, err := r.Cookie(middlewares.UserCookieKey)

	if err != nil {
		return "", err
	}

	token := userCookie.Value

	data, err := hex.DecodeString(token)

	if err != nil {
		return "", err
	}

	if len(data) < 16 {
		return "", errors.New("decoded string is not correct")
	}

	userID := data[:16]

	return hex.EncodeToString(userID), nil
}

func createFullURL(baseURL string, path string) (string, error) {
	u, err := url.Parse(baseURL)

	if err != nil {
		return "", err
	}

	u.Path = path

	return u.String(), nil
}
