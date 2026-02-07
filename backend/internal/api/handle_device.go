package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"wolite/internal/store"
)

type createDeviceRequest struct {
	MACAddress  string `json:"mac_address"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IPAddress   string `json:"ip_address,omitempty"`
}

type updateDeviceRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	IPAddress   string `json:"ip_address,omitempty"`
}

func (r *createDeviceRequest) Validate() error {
	if r.MACAddress == "" {
		return errors.New("mac address is required")
	}
	if r.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

// handleDevicesGetAll returns all devices associated with a username. (jwt protected)
func (a *API) handleDevicesGetAll(w http.ResponseWriter, r *http.Request) {
	claims := a.guard(w, r)
	if claims == nil {
		return
	}

	devices, err := a.store.GetDevicesForUser(claims.Username)
	if err != nil {
		writeRespErr(w, "Failed to retrieve devices", http.StatusInternalServerError)
		slog.Error("failed to retrieve devices", "username", claims.Username)
		return
	}

	writeRespOk(w, "devices retrieved", devices)
	slog.Info("devices retrieved", "username", claims.Username, "devices_count", len(devices))
}

// handleDeviceGet returns a single device by MAC address that is accessible by the user. (jwt protected)
func (a *API) handleDeviceGet(w http.ResponseWriter, r *http.Request) {
	claims := a.guard(w, r)
	if claims == nil {
		return
	}

	// get url id
	id := r.PathValue("id")
	if id == "" {
		writeRespErr(w, "Invalid request", http.StatusBadRequest)
		slog.Error("invalid request", "body", r.Body)
		return
	}

	// Secure Access: Use GetDeviceForUser to ensure ownership
	device, err := a.store.GetDeviceForUser(claims.Username, id)
	if err != nil {
		if err == store.ErrDeviceNotFound {
			writeRespErr(w, "Device not found", http.StatusNotFound)
		} else {
			writeRespErr(w, "Failed to retrieve device", http.StatusInternalServerError)
		}
		slog.Error("failed to retrieve device", "username", claims.Username, "mac_address", id, "error", err)
		return
	}
	writeRespOk(w, "device retrieved", device)
	slog.Info("device retrieved", "username", claims.Username, "mac_address", id)
}

// handleDeviceCreate creates a new device for a user. (jwt protected)
func (a *API) handleDeviceCreate(w http.ResponseWriter, r *http.Request) {
	claims := a.guard(w, r)
	if claims == nil {
		return
	}

	var req createDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeRespErr(w, "Invalid request body", http.StatusBadRequest)
		slog.Error("invalid request body", "username", claims.Username, "error", err)
		return
	}

	if err := req.Validate(); err != nil {
		writeRespErr(w, err.Error(), http.StatusBadRequest)
		slog.Error("validation failed", "username", claims.Username, "error", err)
		return
	}

	device := store.NewDevice(req.MACAddress, req.Name, req.Description, req.IPAddress, store.StatusUnknown)

	// Secure Creation: Use CreateDeviceForUser for atomic creation and assignment
	err := a.store.CreateDeviceForUser(claims.Username, device)
	if err != nil && err == store.ErrDeviceExists {
		writeRespErr(w, "Device already exists", http.StatusBadRequest)
		slog.Error("device already exists", "username", claims.Username, "mac_address", device.MACAddress)
		return
	} else if err != nil {
		writeRespErr(w, "Failed to add device", http.StatusInternalServerError)
		slog.Error("failed to add device", "username", claims.Username, "mac_address", device.MACAddress, "error", err)
		return
	}

	writeRespOk(w, "device added", device)
	slog.Info("device added to user", "username", claims.Username, "mac_address", device.MACAddress)
}

// handleDeviceUpdate updates an existing device. (jwt protected)
func (a *API) handleDeviceUpdate(w http.ResponseWriter, r *http.Request) {
	claims := a.guard(w, r)
	if claims == nil {
		return
	}

	// get url id
	id := r.PathValue("id")
	if id == "" {
		writeRespErr(w, "Invalid request", http.StatusBadRequest)
		slog.Error("invalid request", "body", r.Body)
		return
	}

	// Secure Access: Verify ownership first
	device, err := a.store.GetDeviceForUser(claims.Username, id)
	if err != nil {
		if err == store.ErrDeviceNotFound {
			writeRespErr(w, "Device not found", http.StatusNotFound)
		} else {
			writeRespErr(w, "Failed to retrieve device", http.StatusInternalServerError)
		}
		slog.Error("device not found or access denied", "username", claims.Username, "mac_address", id, "error", err)
		return
	}

	// update device
	var req updateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeRespErr(w, "Invalid request body", http.StatusBadRequest)
		slog.Error("invalid request body", "username", claims.Username, "mac_address", device.MACAddress)
		return
	}
	// Only update fields that are present? The struct pointers are not nil, but strings are values.
	// The current logic updates fields.
	if req.Name != "" {
		device.Name = req.Name
	}
	if req.Description != "" {
		device.Description = req.Description
	}
	if req.IPAddress != "" {
		device.IPAddress = req.IPAddress
	}

	device.Name = req.Name
	device.Description = req.Description
	device.IPAddress = req.IPAddress

	err = a.store.UpdateDevice(device)
	if err != nil {
		writeRespErr(w, "Failed to update device", http.StatusInternalServerError)
		slog.Error("failed to update device", "username", claims.Username, "mac_address", device.MACAddress, "error", err)
		return
	}
	writeRespOk(w, "device updated", device)
	slog.Info("device updated", "username", claims.Username, "mac_address", device.MACAddress)
}

func (a *API) handleDeviceDelete(w http.ResponseWriter, r *http.Request) {
	claims := a.guard(w, r)
	if claims == nil {
		return
	}

	// get url id
	id := r.PathValue("id")
	if id == "" {
		writeRespErr(w, "Invalid request", http.StatusBadRequest)
		slog.Error("invalid id", "username", claims.Username)
		return
	}

	// Secure Access: Verify ownership first
	_, err := a.store.GetDeviceForUser(claims.Username, id)
	if err != nil {
		if err == store.ErrDeviceNotFound {
			writeRespErr(w, "Device not found", http.StatusNotFound)
		} else {
			writeRespErr(w, "Failed to delete device", http.StatusInternalServerError)
		}
		slog.Error("device not found or access denied", "username", claims.Username, "mac_address", id, "error", err)
		return
	}

	err = a.store.DeleteDevice(id)
	if err != nil && err == store.ErrDeviceNotFound {
		writeRespErr(w, "Device not found", http.StatusNotFound)
		slog.Error("device not found", "username", claims.Username, "mac_address", id)
		return
	} else if err != nil {
		writeRespErr(w, "Failed to delete device", http.StatusInternalServerError)
		slog.Error("failed to delete device", "username", claims.Username, "mac_address", id)
		return
	}
	writeRespOk(w, "device deleted", nil)
	slog.Info("device deleted", "username", claims.Username, "mac_address", id)
}

func (a *API) handleDeviceWake(w http.ResponseWriter, r *http.Request) {
}
