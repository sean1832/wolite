//go:build cli

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"wolcompanion/internal/auth"
)

// Global tokenStore for cleanup on exit
var globalTokenStore *auth.TokenStore

func main() {
	// Handle CLI commands
	if len(os.Args) > 1 && os.Args[1] == "regen-token" {
		tokenStore, _, _, err := Initialize()
		if err != nil {
			log.Fatalf("Failed to Initialize: %v", err)
		}
		_, err = tokenStore.Regenerate()
		if err != nil {
			log.Fatalf("Failed to regenerate token: %v", err)
		}
		return
	}

	defer onExit()

	// Normal run mode
	tokenStore, certPath, keyPath, err := Initialize()
	if err != nil {
		log.Fatalf("Failed to Initialize app: %v", err)
	}

	// Store globally for cleanup
	globalTokenStore = tokenStore

	// Create a channel for errors from StartServer
	errChan := make(chan error, 1)
	go func() {
		if err := StartServer(tokenStore, certPath, keyPath); err != nil {
			errChan <- err
		}
	}()

	// Block until interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	select {
	case <-c:
		// Triggered by User (Ctrl+C)
	case err := <-errChan:
		// Triggered by App Failure.
		log.Printf("Application crashed: %v", err)
	}
}

// onExit is called when the CLI app is exiting.
func onExit() {
	if globalTokenStore != nil {
		globalTokenStore.Cleanup()
	}
}
