//go:build cli

package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Start the shared application logic
	go runApp()

	// Block until interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
