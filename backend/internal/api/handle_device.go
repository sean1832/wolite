package api

import (
	"net/http"
	"wolite/internal/auth"
)

// handleDevicesGetAll returns all devices associated with a username. (jwt protected)
func (a *API) handleDevicesGetAll(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			writeRespErr(w, "unauthenticated", http.StatusUnauthorized)
			return
		}
		writeRespErr(w, "Invalid request", http.StatusBadRequest)
		return
	}

	claims, err := auth.ValidateJWTToken(c.Value, []byte(a.config.JWTSecret))
	if err != nil {
		writeRespErr(w, "unauthenticated", http.StatusUnauthorized)
		return
	}

	devices, err := a.store.GetDevicesForUser(claims.Username)
	if err != nil {
		writeRespErr(w, "Failed to retrieve devices", http.StatusInternalServerError)
		return
	}

	writeRespOk(w, "devices retrieved", devices)
}

func (a *API) handleDeviceGet(w http.ResponseWriter, r *http.Request) {
}

func (a *API) handleDeviceCreate(w http.ResponseWriter, r *http.Request) {
}

func (a *API) handleDeviceUpdate(w http.ResponseWriter, r *http.Request) {
}

func (a *API) handleDeviceDelete(w http.ResponseWriter, r *http.Request) {
}

func (a *API) handleDeviceWake(w http.ResponseWriter, r *http.Request) {
}
