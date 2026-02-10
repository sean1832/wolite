package auth

import (
	"crypto/sha256"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// TokenStore manages the lifecycle and storage of the authentication token.
type TokenStore struct {
	mu            sync.RWMutex
	tokenHash     [32]byte
	configDir     string
	tokenPath     string
	tempTokenPath string
}

// NewTokenStore initializes a new TokenStore.
// It loads the existing token if present, or generates a new one.
func NewTokenStore(configDir string) (*TokenStore, error) {
	ts := &TokenStore{
		configDir:     configDir,
		tokenPath:     filepath.Join(configDir, "token.sha256"),
		tempTokenPath: filepath.Join(configDir, "token.txt"),
	}

	if err := ts.loadOrGenerate(); err != nil {
		return nil, err
	}

	return ts, nil
}

// loadOrGenerate loads the existing token hash or generates a new one.
func (ts *TokenStore) loadOrGenerate() error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	// Clean up any stale temp token file
	if _, err := os.Stat(ts.tempTokenPath); err == nil {
		if err := os.Remove(ts.tempTokenPath); err != nil {
			return fmt.Errorf("failed to remove stale token file: %v", err)
		}
		slog.Info("stale token file deleted", "path", ts.tempTokenPath)
	}

	if !TokenExists(ts.tokenPath) {
		_, err := ts.generateAndSave()
		return err
	}

	// Load existing token hash
	hashBytes, err := os.ReadFile(ts.tokenPath)
	if err != nil {
		return fmt.Errorf("failed to read token hash: %v", err)
	}
	if len(hashBytes) != 32 {
		return fmt.Errorf("invalid token hash length: expected 32 bytes, got %d", len(hashBytes))
	}
	copy(ts.tokenHash[:], hashBytes)
	return nil
}

// Regenerate generates a new token, saves it, and updates the hash.
// It returns the new plain-text token.
func (ts *TokenStore) Regenerate() (string, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	token, err := ts.generateAndSave()
	if err != nil {
		return "", err
	}

	return token, nil
}

// generateAndSave generates a new token, writes it to disk, and updates the hash.
// It expects the caller to hold the lock. Returns the plain-text token.
func (ts *TokenStore) generateAndSave() (string, error) {
	slog.Info("generating authentication token")
	token, err := GenerateToken(32)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	if err := os.WriteFile(ts.tempTokenPath, token, 0600); err != nil {
		return "", fmt.Errorf("failed to save token: %v", err)
	}
	slog.Info("token saved to file", "info", "will be deleted in 120 seconds", "path", ts.tempTokenPath)

	// Goroutine to delete token file after 120 seconds
	go func() {
		time.Sleep(120 * time.Second)
		if err := os.Remove(ts.tempTokenPath); err != nil {
			if !os.IsNotExist(err) {
				slog.Error("failed to remove token file", "error", err)
			}
			return
		}
		slog.Info("token file deleted", "reason", "auto-deleted after 120 seconds", "path", ts.tempTokenPath)
	}()

	ts.tokenHash = sha256.Sum256(token)

	// Save hash
	if err := os.WriteFile(ts.tokenPath, ts.tokenHash[:], 0600); err != nil {
		return "", fmt.Errorf("failed to save token hash: %v", err)
	}

	return string(token), nil
}

// GetHash returns the current token hash in a thread-safe manner.
func (ts *TokenStore) GetHash() [32]byte {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return ts.tokenHash
}

// Cleanup removes the temporary token file if it exists.
func (ts *TokenStore) Cleanup() {
	if _, err := os.Stat(ts.tempTokenPath); err == nil {
		if err := os.Remove(ts.tempTokenPath); err != nil {
			slog.Error("failed to remove temp token file", "path", ts.tempTokenPath, "error", err)
		} else {
			slog.Info("temp token file cleaned up", "path", ts.tempTokenPath)
		}
	}
}

// GetTempTokenPath returns the path to the temporary token file.
func (ts *TokenStore) GetTempTokenPath() string {
	return ts.tempTokenPath
}
