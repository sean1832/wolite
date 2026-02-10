package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"wolite/internal/companion"
	"wolite/internal/store"
)

type pairCompanionRequest struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

type companionActionRequest struct {
	Action string `json:"action"` // shutdown, reboot, sleep, hibernate
}

// handleDeviceCompanionPair pairs a device with a companion app.
// It fetches the companion's certificate fingerprint and verifies connectivity before saving.
func (a *API) handleDeviceCompanionPair(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims == nil {
		slog.Error("claims missing from context", "path", r.URL.Path)
		writeRespErr(w, "internal server error", http.StatusInternalServerError)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		writeRespErr(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var req pairCompanionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeRespErr(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" || req.Token == "" {
		writeRespErr(w, "URL and Token are required", http.StatusBadRequest)
		return
	}

	// 1. Get Device to ensure ownership
	device, err := a.store.GetDeviceForUser(claims.Username, id)
	if err != nil {
		if err == store.ErrDeviceNotFound {
			writeRespErr(w, "Device not found", http.StatusNotFound)
		} else {
			writeRespErr(w, "Failed to retrieve device", http.StatusInternalServerError)
		}
		return
	}

	// 2. Fetch Fingerprint (TOFU)
	fingerprint, err := companion.GetFingerprint(r.Context(), req.URL)
	if err != nil {
		writeRespErr(w, "Failed to connect to companion: "+err.Error(), http.StatusBadGateway)
		slog.Error("failed to get fingerprint", "url", req.URL, "error", err)
		return
	}

	// 3. Verify connection with the fetched fingerprint
	client, err := companion.NewClient(req.URL, req.Token, fingerprint)
	if err != nil {
		writeRespErr(w, "Failed to initialize client: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := client.Ping(r.Context()); err != nil {
		writeRespErr(w, "Failed to verify connection after pairing: "+err.Error(), http.StatusBadGateway)
		slog.Error("pairing verification failed", "url", req.URL, "fingerprint", fingerprint, "error", err)
		return
	}

	// 4. Update Device
	device.CompanionURL = req.URL
	device.CompanionToken = req.Token
	device.CompanionAuthFingerprint = fingerprint

	if err := a.store.UpdateDevice(device); err != nil {
		writeRespErr(w, "Failed to save device", http.StatusInternalServerError)
		return
	}

	writeRespOk(w, "Companion paired successfully", device)
	slog.Info("companion paired", "mac", device.MACAddress, "url", req.URL, "fingerprint", fingerprint)
}

// handleDeviceCompanionUnpair removes companion details from a device.
func (a *API) handleDeviceCompanionUnpair(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims == nil {
		writeRespErr(w, "internal server error", http.StatusInternalServerError)
		return
	}

	id := r.PathValue("id")

	device, err := a.store.GetDeviceForUser(claims.Username, id)
	if err != nil {
		if err == store.ErrDeviceNotFound {
			writeRespErr(w, "Device not found", http.StatusNotFound)
		} else {
			writeRespErr(w, "Failed to retrieve device", http.StatusInternalServerError)
		}
		return
	}

	device.CompanionURL = ""
	device.CompanionToken = ""
	device.CompanionAuthFingerprint = ""

	if err := a.store.UpdateDevice(device); err != nil {
		writeRespErr(w, "Failed to update device", http.StatusInternalServerError)
		return
	}

	writeRespOk(w, "Companion unpaired", device)
	slog.Info("companion unpaired", "mac", device.MACAddress)
}

// handleDeviceCompanionAction sends a power command to the paired companion.
func (a *API) handleDeviceCompanionAction(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims == nil {
		writeRespErr(w, "internal server error", http.StatusInternalServerError)
		return
	}

	id := r.PathValue("id")
	var req companionActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeRespErr(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	device, err := a.store.GetDeviceForUser(claims.Username, id)
	if err != nil {
		if err == store.ErrDeviceNotFound {
			writeRespErr(w, "Device not found", http.StatusNotFound)
		} else {
			writeRespErr(w, "Failed to retrieve device", http.StatusInternalServerError)
		}
		return
	}

	if device.CompanionURL == "" || device.CompanionToken == "" {
		writeRespErr(w, "Companion not paired", http.StatusBadRequest)
		return
	}

	client, err := companion.NewClient(device.CompanionURL, device.CompanionToken, device.CompanionAuthFingerprint)
	if err != nil {
		writeRespErr(w, "Invalid companion configuration", http.StatusInternalServerError)
		return
	}

	// Map generic action string to companion package type if needed, or just pass string if package handles it.
	// The companion package defines PowerAction constants.
	var action companion.PowerAction
	switch req.Action {
	case "shutdown":
		action = companion.ActionShutdown
	case "reboot":
		action = companion.ActionReboot
	case "sleep":
		action = companion.ActionSleep
	case "hibernate":
		action = companion.ActionHibernate
	default:
		writeRespErr(w, "Invalid action", http.StatusBadRequest)
		return
	}

	if err := client.Power(r.Context(), action); err != nil {
		writeRespErr(w, "Failed to execute command: "+err.Error(), http.StatusBadGateway)
		slog.Error("companion command failed", "mac", device.MACAddress, "action", action, "error", err)
		return
	}

	writeRespOk(w, "Command executed successfully", nil)
	slog.Info("companion command executed", "mac", device.MACAddress, "action", action)
}
