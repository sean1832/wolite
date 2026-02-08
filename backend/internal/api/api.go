package api

import (
	"context"
	"net/http"
<<<<<<< HEAD
=======
	"wolite/internal/env"
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
	"wolite/internal/store"
)

type API struct {
	Context context.Context
	store   *store.Store
<<<<<<< HEAD
}

func NewAPI(ctx context.Context, store *store.Store) *API {
	return &API{
		Context: ctx,
		store:   store,
=======
	config  *env.Config
}

func NewAPI(ctx context.Context, store *store.Store, config *env.Config) *API {
	return &API{
		Context: ctx,
		store:   store,
		config:  config,
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
	}
}

func (a *API) RegisterRoutesV1(mux *http.ServeMux) {
	const p = "/api/v1"
<<<<<<< HEAD
	// User routes
	mux.HandleFunc("POST "+p+"/users", a.handleUserCreate) // create a new user (for initial setup)
	mux.HandleFunc("PUT "+p+"/users", a.handleUserUpdate)  // update the user (e.g. change password)

	// Device routes
	mux.HandleFunc("GET "+p+"/devices", a.handleDevices)              // list all devices the user has access to
	mux.HandleFunc("POST "+p+"/devices", a.handleDeviceCreate)        // create a new device to the user's account
	mux.HandleFunc("GET "+p+"/devices/{id}", a.handleDeviceGet)       // get a specific device by ID
	mux.HandleFunc("PUT "+p+"/devices/{id}", a.handleDeviceUpdate)    // update a specific device by ID
	mux.HandleFunc("DELETE "+p+"/devices/{id}", a.handleDeviceDelete) // delete a specific device by ID
	// Device Actions:
	mux.HandleFunc("POST "+p+"/devices/{id}/wake", a.handleDeviceWake) // wake a specific device by ID
	// TODO: add more actions like shutdown, restart, etc.
	// mux.HandleFunc("POST "+p+"/devices/{id}/shutdown", a.handleDeviceShutdown) // shutdown a specific device by ID
	// mux.HandleFunc("POST "+p+"/devices/{id}/restart", a.handleDeviceRestart)   // restart a specific device by ID
	// mux.HandleFunc("POST "+p+"/devices/{id}/sleep", a.handleDeviceSleep)   // put a specific device to sleep by ID
	// mux.HandleFunc("POST "+p+"/devices/{id}/status", a.handleDeviceStatus)     // get the status of a specific device by ID

	// Auth routes
	mux.HandleFunc("GET "+p+"/auth/status", a.handleAuthStatus)  // check if the user is authenticated
	mux.HandleFunc("POST "+p+"/auth/login", a.handleAuthLogin)   // login with username and password (optionally OTP)
	mux.HandleFunc("POST "+p+"/auth/logout", a.handleAuthLogout) // logout the user
=======

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
	handlePublic("POST "+p+"/users", a.handleUserCreate)             // create a new user (for initial setup)
	handleAuth("PUT "+p+"/users", a.handleUserUpdate)                // update the user (e.g. change password)
	handleAuth("POST "+p+"/users/otp/verify", a.handleUserOTPVerify) // verify and enable OTP

	// Device routes
	handleAuth("GET "+p+"/devices", a.handleDevicesGetAll)        // list all devices the user has access to
	handleAuth("POST "+p+"/devices", a.handleDeviceCreate)        // create a new device to the user's account
	handleAuth("GET "+p+"/devices/{id}", a.handleDeviceGet)       // get a specific device by ID
	handleAuth("PUT "+p+"/devices/{id}", a.handleDeviceUpdate)    // update a specific device by ID
	handleAuth("DELETE "+p+"/devices/{id}", a.handleDeviceDelete) // delete a specific device by ID

	// Device Actions:
	handleAuth("POST "+p+"/devices/{id}/wake", a.handleDeviceWake) // wake a specific device by ID

	// Auth routes
	handleAuth("GET "+p+"/auth/status", a.handleAuthStatus)             // check if the user is authenticated
	handlePublic("POST "+p+"/auth/login", a.handleAuthLogin)            // login with username and password (optionally OTP)
	handlePublic("GET "+p+"/auth/initialized", a.handleAuthInitialized) // check if app has users
	handlePublic("POST "+p+"/auth/logout", a.handleAuthLogout)          // logout the user
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
}
