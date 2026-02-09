package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"wolcompanion/internal/api"
	"wolcompanion/internal/tls"
)

// runApp contain the shared logic for both GUI and CLI modes
func runApp() {
	// get config directory
	configDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatalf("failed to get config dir: %v", err)
	}

	configDir = filepath.Join(configDir, "wolite", "certs")

	// create directory if not exist
	if err := os.MkdirAll(configDir, 0700); err != nil {
		log.Fatalf("failed to create directory: %v", err)
	}

	certPath := filepath.Join(configDir, "cert.pem")
	keyPath := filepath.Join(configDir, "key.pem")

	// generate cert if not exist
	if !tls.CertExists(certPath, keyPath) {
		slog.Info("generating self-signed certificate")
		if err := tls.GenerateSelfSignedCert(certPath, keyPath); err != nil {
			log.Fatalf("failed to generate cert: %v", err)
		}
		slog.Info("certificate generated", "cert_path", certPath, "key_path", keyPath)
	}

	// TODO: Implement actual companion logic here (e.g. status reporting, command listening)
	mux := http.NewServeMux()

	apiHandler := api.NewAPI(context.Background())
	apiHandler.RegisterRoutesV1(mux)

	slog.Info("Server starting", "port", 8081, "protocol", "https")
	err = http.ListenAndServeTLS(":8081", certPath, keyPath, mux)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
