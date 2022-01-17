package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/api/rest/model"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/config"
	shortenerService "github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/shortener"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/shortener/v1"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/storage/inmemory"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type URLHandlerTestSuite struct {
	suite.Suite
	storage          storage.URLStorer
	shortenerService shortenerService.Processor
	urlHandler       *URLHandler
}

func (s *URLHandlerTestSuite) SetupTest() {
	s.storage = inmemory.NewStorage()
	s.shortenerService, _ = shortener.NewShortenerService(s.storage)
	cfg, _ := config.NewDefaultConfiguration()
	s.urlHandler, _ = NewURLHandler(s.shortenerService, cfg.ServerConfig)
}

func TestURLHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(URLHandlerTestSuite))
}

func (s *URLHandlerTestSuite) TestGetURL() {
	ctx := context.Background()
	urlID, _ := s.shortenerService.SaveURL(ctx, "https://yandex.ru")

	r := chi.NewRouter()

	r.Get("/{urlID}", s.urlHandler.HandleGetURL())

	type want struct {
		code int
	}

	tests := []struct {
		name  string
		urlID string
		want  want
	}{
		{
			name:  "Correct id",
			urlID: urlID,
			want: want{
				code: 307,
			},
		},
		{
			name:  "Invalid id",
			urlID: "",
			want: want{
				code: 404,
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(r)
			defer ts.Close()

			req, err := http.NewRequest(http.MethodGet, ts.URL+"/"+tt.urlID, nil)

			if err != nil {
				t.Errorf("Problem with server")
			}

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			res, err := client.Do(req)
			if err != nil {
				t.Errorf(err.Error())
			} else {
				assert.Equal(t, tt.want.code, res.StatusCode)
			}

			defer res.Body.Close()
		})
	}
}

func (s *URLHandlerTestSuite) TestPostURL() {
	r := chi.NewRouter()

	r.Post("/", s.urlHandler.HandlePostURL())

	type want struct {
		code int
	}

	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "Correct url",
			url:  "http://yandex.com",
			want: want{
				code: 201,
			},
		},
		{
			name: "Invalid url",
			url:  "",
			want: want{
				code: 400,
			},
		},
		{
			name: "Invalid url",
			url:  "261341",
			want: want{
				code: 400,
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(r)
			defer ts.Close()

			payload := strings.NewReader(tt.url)
			req, err := http.NewRequest(http.MethodPost, ts.URL+"/", payload)

			if err != nil {
				t.Errorf("Problem with server")
			}

			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			res, err := client.Do(req)

			if err != nil {
				t.Errorf(err.Error())
			} else {
				assert.Equal(t, tt.want.code, res.StatusCode)
			}

			defer res.Body.Close()
		})
	}
}

func (s *URLHandlerTestSuite) TestJSONPostURL() {
	r := chi.NewRouter()

	r.Post("/api/shorten", s.urlHandler.JSONHandlePostURL())

	type want struct {
		code int
	}

	tests := []struct {
		name string
		url  model.RequestURL
		want want
	}{
		{
			name: "Correct url",
			url: model.RequestURL{
				URL: "http://yandex.com",
			},
			want: want{
				code: 201,
			},
		},
		{
			name: "Invalid url",
			url: model.RequestURL{
				URL: "",
			},
			want: want{
				code: 400,
			},
		},
		{
			name: "Invalid url",
			url: model.RequestURL{
				URL: "4123",
			},
			want: want{
				code: 400,
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(r)
			defer ts.Close()

			reqBody, _ := json.Marshal(tt.url)
			payload := strings.NewReader(string(reqBody))
			req, err := http.NewRequest(http.MethodPost, ts.URL+"/api/shorten", payload)
			
			if err != nil {
				t.Errorf("Problem with server")
			}
			
			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			res, err := client.Do(req)

			buf := new(bytes.Buffer)
			buf.ReadFrom(res.Body)
			newStr := buf.String()
			if err != nil {
				t.Errorf(err.Error())
			} else {
				assert.Equal(t, tt.want.code, res.StatusCode, newStr)
			}

			defer res.Body.Close()
		})
	}
}
