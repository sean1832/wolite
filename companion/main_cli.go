//go:build cli

package main

import (
	"flag"
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

	// Normal run mode
	tokenStore, certPath, keyPath, err := Initialize()
	if err != nil {
		log.Fatalf("Failed to Initialize app: %v", err)
	}
	defer onExit()

	// Store globally for cleanup
	globalTokenStore = tokenStore

	// flags
	// -p --port
	defaultPort := 8443
	apiPort := flag.Int("port", defaultPort, "Port to run the server on")
	flag.IntVar(apiPort, "p", defaultPort, "Port to run the server on (shorthand)")

	var showHelp bool
	flag.BoolVar(&showHelp, "help", false, "Show this help message")
	flag.BoolVar(&showHelp, "h", false, "Show this help message (shorthand)")

	flag.Parse()

	if showHelp {
		printHelp()
		return
	}

	// Create a channel for errors from StartServer
	errChan := make(chan error, 1)
	go func() {
		if err := StartServer(tokenStore, certPath, keyPath, *apiPort); err != nil {
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

func printHelp() bool {
	println("Usage: wolite-companion-cli [options]")
	println("Commands:")
	println("  regen-token        Generate a new authentication token")
	println("Options:")
	println("  -p, --port <port>  Port to run the server on")
	println("  -h, --help         Show this help message")
	return true
}

// onExit is called when the CLI app is exiting.
func onExit() {
	if globalTokenStore != nil {
		globalTokenStore.Cleanup()
	}
}
