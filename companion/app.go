package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"wolcompanion/internal/api"
	"wolcompanion/internal/auth"
	"wolcompanion/internal/tls"
)

// runApp contain the shared logic for both GUI and CLI modes
func runApp() error {
	// get config directory
	configDir, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("failed to get config dir: %v", err)
	}

	configDir = filepath.Join(configDir, "wolite")

	// create directory if not exist
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	certPath := filepath.Join(configDir, "cert.pem")
	keyPath := filepath.Join(configDir, "key.pem")
	tokenPath := filepath.Join(configDir, "token.sha256")

	// temp token file
	tempTokenFile := filepath.Join(configDir, "token.txt")
	if _, err := os.Stat(tempTokenFile); err == nil {
		// token file exists, meaning it is not cleanup from previous session, delete it
		if err := os.Remove(tempTokenFile); err != nil {
			return fmt.Errorf("failed to remove token file: %v", err)
		}
		slog.Info("token file deleted", "reason", "not cleaned up from previous session", "path", tempTokenFile)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to check token file: %v", err)
	}

	// generate cert if not exist
	if !tls.CertExists(certPath, keyPath) {
		slog.Info("generating self-signed certificate")
		if err := tls.GenerateSelfSignedCert(certPath, keyPath); err != nil {
			return fmt.Errorf("failed to generate cert: %v", err)
		}
		slog.Info("certificate generated", "cert_path", certPath, "key_path", keyPath)
	}

	// generate token if not exist, and save only the hash of the token
	var tokenHash [32]byte
	if !auth.TokenExists(tokenPath) {
		slog.Info("generating authentication token")
		token, err := auth.GenerateToken(32)
		if err != nil {
			return fmt.Errorf("failed to generate token: %v", err)
		}
		if err := os.WriteFile(tempTokenFile, []byte(token), 0600); err != nil {
			return fmt.Errorf("failed to save token: %v", err)
		}
		slog.Info("token saved to file, copy it. (auto-delete after 120 seconds)", "path", tempTokenFile)

		// goroutine: auto delete token file after 120 seconds
		go func() {
			time.Sleep(120 * time.Second)
			if err := os.Remove(tempTokenFile); err != nil {
				slog.Error("failed to remove token file", "error", err)
				return
			}
			slog.Info("token file deleted", "reason", "auto-deleted after 120 seconds", "path", tempTokenFile)
		}()

		// hash
		tokenHash = sha256.Sum256([]byte(token))
		// save only the hash, not the token itself
		if err := os.WriteFile(tokenPath, tokenHash[:], 0600); err != nil {
			return fmt.Errorf("failed to save token hash: %v", err)
		}
	} else {
		// load existing token hash
		hashBytes, err := os.ReadFile(tokenPath)
		if err != nil {
			return fmt.Errorf("failed to read token hash: %v", err)
		}
		if len(hashBytes) != 32 {
			return fmt.Errorf("invalid token hash length: expected 32 bytes, got %d", len(hashBytes))
		}
		copy(tokenHash[:], hashBytes)
	}

	mux := http.NewServeMux()

	apiHandler := api.NewAPI(context.Background(), tokenHash)
	apiHandler.RegisterRoutesV1(mux)

	slog.Info("Server starting", "port", 8081, "protocol", "https")
	err = http.ListenAndServeTLS(":8081", certPath, keyPath, mux)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	return nil
}

// onExit is called when the app is exiting. Gracefully clean up resources.
func onExit() {
	// remove temp token file if it exists
	configDir, err := os.UserCacheDir()
	if err != nil {
		slog.Error("failed to get config dir", "error", err)
		return
	}

	tempTokenFile := filepath.Join(configDir, "wolite", "token.txt")
	if _, err := os.Stat(tempTokenFile); err == nil {
		if err := os.Remove(tempTokenFile); err != nil {
			slog.Error("failed to remove temp token file", "path", tempTokenFile, "error", err)
			return
		}
		slog.Info("temp token file cleaned up", "path", tempTokenFile)
	}
}
