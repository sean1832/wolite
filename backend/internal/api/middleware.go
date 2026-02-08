package api

import (
	"context"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
	"wolite/internal/auth"
)

type middleware func(http.Handler) http.Handler

// CreateStack wraps a handler with a list of middleware
func CreateStack(xs ...middleware) middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}

// Logger logs the request details and execution time
func Logger(next http.Handler) http.Handler {
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
func Recoverer(next http.Handler) http.Handler {
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

type contextKey string

const userContextKey contextKey = "user"

// Auth validates the JWT token and adds the claims to the context
func (a *API) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := a.validateCookies(r)
		if err != nil {
			switch err {
			case ErrUnauthorized:
				writeRespErr(w, "unauthenticated", http.StatusUnauthorized)
			case ErrInvalidRequest:
				writeRespErr(w, "Invalid request", http.StatusBadRequest)
			default:
				writeRespErr(w, "Authentication failed", http.StatusInternalServerError)
			}
			slog.Warn("authentication failed", "path", r.URL.Path, "error", err)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext retrieves the user claims from the context
func GetUserFromContext(ctx context.Context) *auth.Claims {
	claims, ok := ctx.Value(userContextKey).(*auth.Claims)
	if !ok {
		return nil
	}
	return claims
}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

// Cors handles Cross-Origin Resource Sharing for development mode
func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from the frontend dev server
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
