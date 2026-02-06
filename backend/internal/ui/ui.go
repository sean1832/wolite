package ui

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"
)

// Embed the entire frontend build directory
//
//go:embed all:dist
var dist embed.FS

// GetFileSystem returns the embedded filesystem for the UI
func GetFileSystem() (http.FileSystem, error) {
	fsys, err := fs.Sub(dist, "dist")
	if err != nil {
		return nil, err
	}
	return http.FS(fsys), nil
}

// Handler serves the embedded UI with SPA fallback support
type Handler struct {
	fileServer http.Handler
}

// NewHandler creates a new UI handler
func NewHandler() (*Handler, error) {
	fsys, err := GetFileSystem()
	if err != nil {
		return nil, err
	}

	return &Handler{
		fileServer: http.FileServer(fsys),
	}, nil
}

// ServeHTTP handles HTTP requests for the UI
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// For API routes, don't handle them here
	if strings.HasPrefix(path, "/api/") {
		http.NotFound(w, r)
		return
	}

	// Try to serve the file
	// If it's a directory or doesn't exist, serve index.html for SPA routing
	if strings.HasSuffix(path, "/") || !strings.Contains(path, ".") {
		r.URL.Path = "/"
	}

	slog.Debug("Serving UI", "path", path, "rewritten", r.URL.Path)
	h.fileServer.ServeHTTP(w, r)
}
