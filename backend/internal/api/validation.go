package api

import (
	"net/http"
	"wolite/internal/auth"
)

func (a *API) validateCookies(r *http.Request) (*auth.Claims, error) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, ErrUnauthorized
		}
		return nil, ErrInvalidRequest
	}

	claims, err := auth.ValidateJWTToken(c.Value, []byte(a.config.JWTSecret))
	if err != nil {
		return nil, ErrUnauthorized
	}
	return claims, nil
}
