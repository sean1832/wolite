package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"wolite/internal/auth"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	OTP      string `json:"otp,omitempty"`
}

type authResponse struct {
	Status string `json:"status"`
	User   string `json:"user,omitempty"`
	HasOTP bool   `json:"has_otp"` // Confirmed 2FA enabled?
}

func (a *API) handleAuthLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeRespErr(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := a.store.FindUser(req.Username)
	if err != nil {
		// Don't reveal if user exists or not, but for now standard 401
		writeRespErr(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !auth.CheckPasswordHash(req.Password, user.Password) {
		writeRespErr(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// If user has OTP enabled, verify it
	if user.OTP != "" {
		if req.OTP == "" {
			writeRespErr(w, "OTP required", http.StatusUnauthorized)
			return
		}
		if !auth.Validate2FA(req.OTP, user.OTP) {
			writeRespErr(w, "Invalid OTP", http.StatusUnauthorized)
			return
		}
	}

	tokenString, expirationTime, err := auth.GenerateJWTToken(user.Username, []byte(a.config.JWTSecret), a.config.JWTExpiry)
	if err != nil {
		writeRespErr(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		// Secure:   true, // TODO: Enable in production
	})

	writeRespOk(w, "authenticated", authResponse{Status: "authenticated", User: user.Username, HasOTP: user.OTP != ""})
}

func (a *API) handleAuthLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
	})
	writeRespOk(w, "logged out", nil)
}

func (a *API) handleAuthStatus(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims == nil {
		slog.Error("claims missing from context", "path", r.URL.Path)
		writeRespErr(w, "internal server error", http.StatusInternalServerError)
		return
	}

	user, err := a.store.FindUser(claims.Username)
	if err != nil {
		slog.Error("User from token not found", "username", claims.Username)
		writeRespErr(w, "User not found", http.StatusUnauthorized)
		return
	}

	writeRespOk(w, "authenticated", authResponse{Status: "authenticated", User: user.Username, HasOTP: user.OTP != ""})
}

// handleAuthInitialized checks if the application has been initialized (has users)
func (a *API) handleAuthInitialized(w http.ResponseWriter, r *http.Request) {
	initialized := a.store.HasUsers()
	writeRespOk(w, "ok", map[string]bool{"initialized": initialized})
}
