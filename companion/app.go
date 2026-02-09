package main

import (
	"context"
	"crypto/sha256"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"wolcompanion/internal/api"
	"wolcompanion/internal/auth"
	"wolcompanion/internal/tls"
)

// runApp contain the shared logic for both GUI and CLI modes
func runApp() {
	// get config directory
	configDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatalf("failed to get config dir: %v", err)
	}

	configDir = filepath.Join(configDir, "wolite")

	// create directory if not exist
	if err := os.MkdirAll(configDir, 0700); err != nil {
		log.Fatalf("failed to create directory: %v", err)
	}

	certPath := filepath.Join(configDir, "cert.pem")
	keyPath := filepath.Join(configDir, "key.pem")
	tokenPath := filepath.Join(configDir, "token.sha256")

	// generate cert if not exist
	if !tls.CertExists(certPath, keyPath) {
		slog.Info("generating self-signed certificate")
		if err := tls.GenerateSelfSignedCert(certPath, keyPath); err != nil {
			log.Fatalf("failed to generate cert: %v", err)
		}
		slog.Info("certificate generated", "cert_path", certPath, "key_path", keyPath)
	}

	// generate token if not exist, and save only the hash of the token
	var tokenHash [32]byte
	if !auth.TokenExists(tokenPath) {
		slog.Info("generating authentication token")
		token, err := auth.GenerateToken(32)
		if err != nil {
			log.Fatalf("failed to generate token: %v", err)
		}
		slog.Info("authentication token generated", "token", string(token))

		// hash
		tokenHash = sha256.Sum256([]byte(token))
		// save only the hash, not the token itself
		if err := os.WriteFile(tokenPath, tokenHash[:], 0600); err != nil {
			log.Fatalf("failed to save token hash: %v", err)
		}
	} else {
		// load existing token hash
		hashBytes, err := os.ReadFile(tokenPath)
		if err != nil {
			log.Fatalf("failed to read token hash: %v", err)
		}
		if len(hashBytes) != 32 {
			log.Fatalf("invalid token hash length: expected 32 bytes, got %d", len(hashBytes))
		}
		copy(tokenHash[:], hashBytes)
	}

	mux := http.NewServeMux()

	apiHandler := api.NewAPI(context.Background(), tokenHash)
	apiHandler.RegisterRoutesV1(mux)

	slog.Info("Server starting", "port", 8081, "protocol", "https")
	err = http.ListenAndServeTLS(":8081", certPath, keyPath, mux)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
