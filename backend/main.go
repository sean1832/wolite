package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"wolite/internal/api"
<<<<<<< HEAD
=======
	"wolite/internal/env"
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
	"wolite/internal/store"
	"wolite/internal/ui"
)

func main() {
<<<<<<< HEAD
	mux := http.NewServeMux()

	store, err := store.New("db.json")
	if err != nil {
		log.Fatalf("failed to initialized JSON database %v", err)
	}
	apiHandler := api.NewAPI(context.Background(), store)
=======
	config := env.LoadConfig()
	mux := http.NewServeMux()
	store, err := store.New(config.DatabasePath)
	if err != nil {
		log.Fatalf("failed to initialized JSON database %v", err)
	}

	apiHandler := api.NewAPI(context.Background(), store, config)
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de

	apiHandler.RegisterRoutesV1(mux)

	// initialize embedded UI handler
	uiHandler, err := ui.NewHandler()
	if err != nil {
		log.Fatalf("failed to initialize UI handler: %v", err)
	}
	mux.Handle("/", uiHandler)

<<<<<<< HEAD
	slog.Info("Server starting", "url", "http://localhost:8080")
	err = http.ListenAndServe(":8080", mux)
=======
	var handler http.Handler = mux
	if config.DevMode {
		handler = api.Cors(mux)
	}

	slog.Info("Server starting", "url", "http://localhost:8080")
	err = http.ListenAndServe(":8080", handler)
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
