//go:build cli

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defer onExit()

	// Create a channel for errors from runApp
	errChan := make(chan error, 1)
	go func() {
		if err := runApp(); err != nil {
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
		// defer onExit() will still run when main() returns underneath.
		log.Printf("Application crashed: %v", err)
	}
}
