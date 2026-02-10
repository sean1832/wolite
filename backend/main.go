package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"
	"wolite/internal/api"
	"wolite/internal/env"
	"wolite/internal/store"
	"wolite/internal/ui"
	"wolite/internal/worker"
)

func main() {
	config := env.LoadConfig()
	mux := http.NewServeMux()
	store, err := store.New(config.DatabasePath)
	if err != nil {
		log.Fatalf("failed to initialized JSON database %v", err)
	}

	apiHandler := api.NewAPI(context.Background(), store, config)

	// Start background workers
	statusChecker := worker.NewStatusChecker(store, 30*time.Second)
	go statusChecker.Start(context.Background())

	apiHandler.RegisterRoutesV1(mux)

	// initialize embedded UI handler
	uiHandler, err := ui.NewHandler()
	if err != nil {
		log.Fatalf("failed to initialize UI handler: %v", err)
	}
	mux.Handle("/", uiHandler)

	var handler http.Handler = mux
	if config.DevMode {
		handler = api.Cors(mux)
	}

	slog.Info("Server starting", "port", config.Port)
	err = http.ListenAndServe(":"+config.Port, handler)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
