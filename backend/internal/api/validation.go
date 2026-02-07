package api

import (
	"log/slog"
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

// guard checks for valid authentication and returns claims.
// It handles writing error responses for unauthorized or invalid requests.
func (a *API) guard(w http.ResponseWriter, r *http.Request) *auth.Claims {
	claims, err := a.validateCookies(r)
	if err != nil {
		switch err {
		case ErrUnauthorized:
			writeRespErr(w, "unauthenticated", http.StatusUnauthorized)
			slog.Error("unauthenticated", "path", r.URL.Path)
		case ErrInvalidRequest:
			writeRespErr(w, "Invalid request", http.StatusBadRequest)
			slog.Error("invalid request", "path", r.URL.Path, "error", err)
		default:
			writeRespErr(w, "Authentication failed", http.StatusInternalServerError)
			slog.Error("authentication failed", "path", r.URL.Path, "error", err)
		}
		return nil
	}
	return claims
}
