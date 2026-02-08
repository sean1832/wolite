package api

import (
	"context"
	"net/http"
	"wolite/internal/env"
	"wolite/internal/store"
)

type API struct {
	Context context.Context
	store   *store.Store
	config  *env.Config
}

func NewAPI(ctx context.Context, store *store.Store, config *env.Config) *API {
	return &API{
		Context: ctx,
		store:   store,
		config:  config,
	}
}

func (a *API) RegisterRoutesV1(mux *http.ServeMux) {
	const p = "/api/v1"

	// Middleware stacks
	standard := []middleware{
		Logger,
		Recoverer,
	}

	authStack := []middleware{
		Logger,
		Recoverer,
		a.Auth,
	}

	// Add CORS middleware in development mode
	if a.config.DevMode {
		standard = append([]middleware{Cors}, standard...)
		authStack = append([]middleware{Cors}, authStack...)
	}

	// Helper to apply middleware
	handle := func(pattern string, handler func(http.ResponseWriter, *http.Request), middlewares []middleware) {
		h := http.HandlerFunc(handler)
		// Apply in reverse order so the first in list is the outer-most
		for i := len(middlewares) - 1; i >= 0; i-- {
			h = middlewares[i](h).(http.HandlerFunc)
		}
		mux.Handle(pattern, h)
	}

	// Wrapper for standard middleware (no auth)
	handlePublic := func(pattern string, handler func(http.ResponseWriter, *http.Request)) {
		handle(pattern, handler, standard)
	}

	// Wrapper for auth middleware
	handleAuth := func(pattern string, handler func(http.ResponseWriter, *http.Request)) {
		handle(pattern, handler, authStack)
	}

	// User routes
	handlePublic("POST "+p+"/users", a.handleUserCreate) // create a new user (for initial setup)
	handleAuth("PUT "+p+"/users", a.handleUserUpdate)    // update the user (e.g. change password)

	// Device routes
	handleAuth("GET "+p+"/devices", a.handleDevicesGetAll)        // list all devices the user has access to
	handleAuth("POST "+p+"/devices", a.handleDeviceCreate)        // create a new device to the user's account
	handleAuth("GET "+p+"/devices/{id}", a.handleDeviceGet)       // get a specific device by ID
	handleAuth("PUT "+p+"/devices/{id}", a.handleDeviceUpdate)    // update a specific device by ID
	handleAuth("DELETE "+p+"/devices/{id}", a.handleDeviceDelete) // delete a specific device by ID

	// Device Actions:
	handleAuth("POST "+p+"/devices/{id}/wake", a.handleDeviceWake) // wake a specific device by ID

	// Auth routes
	handleAuth("GET "+p+"/auth/status", a.handleAuthStatus)    // check if the user is authenticated
	handlePublic("POST "+p+"/auth/login", a.handleAuthLogin)   // login with username and password (optionally OTP)
	handlePublic("POST "+p+"/auth/logout", a.handleAuthLogout) // logout the user
}
