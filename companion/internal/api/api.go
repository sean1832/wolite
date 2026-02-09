package api

import (
	"context"
	"net/http"
)

type API struct {
	Context context.Context
}

func NewAPI(ctx context.Context) *API {
	return &API{
		Context: ctx,
	}
}

func (a *API) RegisterRoutesV1(mux *http.ServeMux) {
	const p = "/api/v1"

	mux.HandleFunc("GET "+p+"/health", a.handleHealth)

	// power commands
	mux.HandleFunc("POST "+p+"/shutdown", a.handleShutdown)
	mux.HandleFunc("POST "+p+"/reboot", a.handleReboot)
	mux.HandleFunc("POST "+p+"/sleep", a.handleSleep)
	mux.HandleFunc("POST "+p+"/hibernate", a.handleHibernate)
}

func (a *API) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeRespOk(w, "ok", nil)
}
