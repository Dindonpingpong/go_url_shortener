package middlewares

import (
	"fmt"
	"net/http"

	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/config"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/secretary"
	"github.com/google/uuid"
)

type CookieHandler struct {
	svc secretary.Secretary
	cfg *config.SecretConfig
}

const (
	UserCookieKey = "user"
)

func NewCookieHandler(svc secretary.Secretary, cfg *config.SecretConfig) (*CookieHandler, error) {
	if svc == nil {
		return nil, fmt.Errorf("secretary.Secretary: nil")
	}

	return &CookieHandler{
		svc: svc,
		cfg: cfg,
	}, nil
}

func (c *CookieHandler) AuthCookieHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(UserCookieKey)

		if err != nil {
			if err == http.ErrNoCookie {
				userID := uuid.New().String()

				token := c.svc.Encode(userID)

				generatedCookie := &http.Cookie{
					Name:  UserCookieKey,
					Value: token,
					Path:  "/",
				}

				http.SetCookie(w, generatedCookie)
				r.AddCookie(generatedCookie)
			} else {
				http.Error(w, "Something wrong with cookie", http.StatusInternalServerError)
			}
		} else {
			_, err := c.svc.Decode(cookie.Value)

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
		}

		next.ServeHTTP(w, r)
	})
}
