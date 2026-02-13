package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"wolcompanion/internal/api"
	"wolcompanion/internal/auth"
	"wolcompanion/internal/config"
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

func InitializeConfig(port int) (*config.ConfigStore, error) {
	// get config directory
	configDir, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config dir: %v", err)
	}

	configDir = filepath.Join(configDir, "wolite")
	configFilePath := filepath.Join(configDir, "config.json")

	// create directory if not exist
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	// create config file if not exist
	var configStore *config.ConfigStore
	if !config.ConfigExists(configFilePath) {
		configStore, err = config.CreateConfig(configFilePath, port)
		if err != nil {
			return nil, fmt.Errorf("failed to create config: %v", err)
		}
		slog.Info("config file created", "path", configFilePath)
	} else {
		configStore, err = config.LoadConfig(configFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %v", err)
		}
		slog.Info("config file loaded", "path", configFilePath)
	}

	return configStore, nil
}

// StartServer starts the HTTP server. It blocks until the server exits.
func StartServer(tokenStore *auth.TokenStore, certPath, keyPath string, port int) error {
	mux := http.NewServeMux()

	apiHandler := api.NewAPI(context.Background(), tokenStore)
	apiHandler.RegisterRoutesV1(mux)

	slog.Info("Server starting", "port", port, "protocol", "https")
	// ListenAndServeTLS blocks
	return http.ListenAndServeTLS(":"+strconv.Itoa(port), certPath, keyPath, mux)
}
