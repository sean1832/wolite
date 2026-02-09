package api

import (
	"context"
	"net/http"
)

type API struct {
	Context   context.Context
	tokenHash [32]byte
}

func NewAPI(ctx context.Context, tokenHash [32]byte) *API {
	return &API{
		Context:   ctx,
		tokenHash: tokenHash,
	}
}

func (a *API) RegisterRoutesV1(mux *http.ServeMux) {
	const p = "/api/v1"

	// public routes
	mux.HandleFunc("GET "+p+"/health", a.handleHealth)

	// power commands (protected routes)
	mux.HandleFunc("POST "+p+"/shutdown", a.requireAuth(a.handleShutdown))
	mux.HandleFunc("POST "+p+"/reboot", a.requireAuth(a.handleReboot))
	mux.HandleFunc("POST "+p+"/sleep", a.requireAuth(a.handleSleep))
	mux.HandleFunc("POST "+p+"/hibernate", a.requireAuth(a.handleHibernate))
}

func (a *API) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeRespOk(w, "ok", nil)
}
