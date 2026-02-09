package api

import (
	"crypto/sha256"
	"crypto/subtle"
	"log/slog"
	"net/http"
	"strings"
)

// requireAuth creates a closure that enforces authentication.
func (a *API) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeRespErr(w, "Unauthorized", http.StatusUnauthorized)
			slog.Error("Unauthorized", "error", "missing Authorization header")
			return
		}

		// Expecting "Bearer <token>"
		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			writeRespErr(w, "Unauthorized", http.StatusUnauthorized)
			slog.Error("Unauthorized", "error", "invalid Authorization header format")
			return
		}

		inputToken := strings.TrimPrefix(authHeader, prefix)

		// Hash the input token to compare against the stored hash
		inputHash := sha256.Sum256([]byte(inputToken))

		// ConstantTimeCompare requires byte slices.
		// It returns 1 if equal, 0 otherwise.
		if subtle.ConstantTimeCompare(inputHash[:], a.tokenHash[:]) != 1 {
			writeRespErr(w, "Unauthorized", http.StatusUnauthorized)
			slog.Error("Unauthorized", "error", "invalid token")
			return
		}

		next(w, r)
	}
}
