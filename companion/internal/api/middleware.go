package api

import (
	"crypto/sha256"
	"crypto/subtle"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

type middleware func(http.Handler) http.Handler

// logger logs the request details and execution time
func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		slog.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.statusCode,
			"duration", time.Since(start),
		)
	})
}

// Recoverer recovers from panics, logs the error, and returns a 500
func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				writeRespErr(w, "Internal Server Error", http.StatusInternalServerError)
				slog.Error("panic recovered", "error", err, "stack", string(debug.Stack()))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// Global limiter: 5 requests per second, burst of 10
var limiter = rate.NewLimiter(rate.Every(200*time.Millisecond), 10)

// rateLimit limits the number of requests per second
func rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			writeRespErr(w, "Too Many Requests", http.StatusTooManyRequests)
			slog.Warn("rate limit exceeded", "path", r.URL.Path)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// bearerTokenAuth creates a middleware that enforces authentication.
func (a *API) bearerTokenAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		next.ServeHTTP(w, r)
	})
}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}