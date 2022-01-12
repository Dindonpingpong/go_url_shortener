package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
	s.urlHandler, _ = NewURLHandler(s.shortenerService)
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
