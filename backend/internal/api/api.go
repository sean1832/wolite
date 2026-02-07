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
}
