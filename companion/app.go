package main

import (
	"fmt"
	"time"
)

// runApp contain the shared logic for both GUI and CLI modes
func runApp() {
	fmt.Println("Wolite Companion started...")
	// TODO: Implement actual companion logic here (e.g. status reporting, command listening)

	// Keep the goroutine alive if needed, or just run the logic loop
	for {
		time.Sleep(1 * time.Second)
	}
}
