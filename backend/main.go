package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"wolite/internal/api"
	"wolite/internal/store"
	"wolite/internal/ui"
)

func main() {
	mux := http.NewServeMux()

	store, err := store.New("db.json")
	if err != nil {
		log.Fatalf("failed to initialized JSON database %v", err)
	}
	apiHandler := api.NewAPI(context.Background(), store)

	apiHandler.RegisterRoutesV1(mux)

	// initialize embedded UI handler
	uiHandler, err := ui.NewHandler()
	if err != nil {
		log.Fatalf("failed to initialize UI handler: %v", err)
	}
	mux.Handle("/", uiHandler)

	slog.Info("Server starting", "url", "http://localhost:8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
