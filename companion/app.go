package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"wolcompanion/internal/api"
	"wolcompanion/internal/auth"
	"wolcompanion/internal/tls"
)

// Initialize prepares the application environment.
// It creates directories, generates certificates, and initializes the token store.
// Returns the TokenStore, certPath, keyPath, and any error.
func Initialize() (*auth.TokenStore, string, string, error) {
	// get config directory
	configDir, err := os.UserCacheDir()
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to get config dir: %v", err)
	}

	configDir = filepath.Join(configDir, "wolite")

	// create directory if not exist
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, "", "", fmt.Errorf("failed to create directory: %v", err)
	}

	certPath := filepath.Join(configDir, "cert.pem")
	keyPath := filepath.Join(configDir, "key.pem")

	// generate cert if not exist
	if !tls.CertExists(certPath, keyPath) {
		slog.Info("generating self-signed certificate")
		if err := tls.GenerateSelfSignedCert(certPath, keyPath); err != nil {
			return nil, "", "", fmt.Errorf("failed to generate cert: %v", err)
		}
		slog.Info("certificate generated", "cert_path", certPath, "key_path", keyPath)
	}

	// Initialize TokenStore (handles token generation, loading, and file cleanup)
	tokenStore, err := auth.NewTokenStore(configDir)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to initialize token store: %v", err)
	}

	return tokenStore, certPath, keyPath, nil
}

// StartServer starts the HTTP server. It blocks until the server exits.
func StartServer(tokenStore *auth.TokenStore, certPath, keyPath string) error {
	mux := http.NewServeMux()

	apiHandler := api.NewAPI(context.Background(), tokenStore)
	apiHandler.RegisterRoutesV1(mux)

	slog.Info("Server starting", "port", 8081, "protocol", "https")
	// ListenAndServeTLS blocks
	return http.ListenAndServeTLS(":8081", certPath, keyPath, mux)
}
