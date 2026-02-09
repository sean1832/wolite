package api

import (
	"context"
	"net/http"
	"wolcompanion/internal/auth"
)

type API struct {
	Context    context.Context
	tokenStore *auth.TokenStore
}

func NewAPI(ctx context.Context, tokenStore *auth.TokenStore) *API {
	return &API{
		Context:    ctx,
		tokenStore: tokenStore,
	}
}

func (a *API) RegisterRoutesV1(mux *http.ServeMux) {
	const p = "/api/v1"

	standard := []middleware{
		logger,
		securityHeaders,
		recoverer,
		rateLimit,
	}

	authStack := []middleware{
		logger,
		securityHeaders,
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
