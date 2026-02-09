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

	standard := []middleware{
		logger,
		recoverer,
		rateLimit,
	}

	authStack := []middleware{
		logger,
		recoverer,
		rateLimit,
		a.bearerTokenAuth,
	}

	handle := func(pattern string, handler func(http.ResponseWriter, *http.Request), middlewares []middleware) {
		var h http.Handler = http.HandlerFunc(handler)
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h)
		}
		mux.Handle(pattern, h)
	}

	handlePublic := func(pattern string, handler func(http.ResponseWriter, *http.Request)) {
		handle(pattern, handler, standard)
	}

	handleAuth := func(pattern string, handler func(http.ResponseWriter, *http.Request)) {
		handle(pattern, handler, authStack)
	}

	// public routes
	handlePublic("GET "+p+"/health", a.handleHealth)
	handlePublic("GET "+p+"/cert-fingerprint", a.handleCertFingerprint)

	// power commands (protected routes)
	handleAuth("POST "+p+"/shutdown", a.handleShutdown)
	handleAuth("POST "+p+"/reboot", a.handleReboot)
	handleAuth("POST "+p+"/sleep", a.handleSleep)
	handleAuth("POST "+p+"/hibernate", a.handleHibernate)
}

func (a *API) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeRespOk(w, "ok", nil)
}
