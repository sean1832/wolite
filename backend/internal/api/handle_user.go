package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"wolite/internal/auth"
	"wolite/internal/store"
)

func (a *API) handleUserGet(w http.ResponseWriter, r *http.Request) {

}

func (a *API) handleUserCreate(w http.ResponseWriter, r *http.Request) {
	// payload should contain username and password
	payload := struct {
		Username string `json:"username"`
		Password string `json:"password"`
		UseOTP   bool   `json:"use_otp,omitempty"` // optional flag to indicate if user wants to use OTP
	}{}

	// Limit request body size to prevent memory exhaustion (e.g. 1MB limit)
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeRespErr(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Input Validation
	// Fail FAST before doing expensive DB or Crypto work
	if len(payload.Username) < 3 || len(payload.Username) > 32 {
		writeRespErr(w, "Username must be between 3 and 32 characters", http.StatusBadRequest)
		return
	}
	if len(payload.Password) < 8 {
		writeRespErr(w, "Password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	// Availability Check (DoS Protection)
	// Check DB before hashing. This prevents attackers from burning your CPU
	// by flooding requests for existing users.
	_, err := a.store.FindUser(payload.Username)
	if err == nil {
		// Error is nil, meaning User WAS found -> Conflict
		writeRespErr(w, "Username taken", http.StatusConflict)
		return
	}

	if err != store.ErrUserNotFound {
		// Any other error is a system failure (disk I/O, permissions, etc.)
		slog.Error("Database check failed", "error", err)
		writeRespErr(w, "System error", http.StatusInternalServerError)
		return
	}

	// Hashing
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		slog.Error("Password hashing failed", "error", err)
		writeRespErr(w, "System error", http.StatusInternalServerError)
		return
	}

	// Construct User
	// generate OTP secret here if needed
	user, err := store.NewUser(payload.Username, hashedPassword)
	var otpUrl string
	if payload.UseOTP {
		secret, url, err := auth.GenerateOTPSecret(payload.Username)
		if err != nil {
			slog.Error("OTP secret generation failed", "error", err)
			writeRespErr(w, "System error", http.StatusInternalServerError)
			return
		}
		otpUrl = url
		user.OTP = secret
	}
	if err != nil {
		// Ensure NewUser doesn't return sensitive system errors
		writeRespErr(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	// Persist
	// We use CreateUser which acts as a final guard
	if err := a.store.CreateUser(*user); err != nil {
		if err == store.ErrUserExists {
			writeRespErr(w, "Username taken", http.StatusConflict)
			return
		}
		slog.Error("Failed to persist user", "error", err)
		writeRespErr(w, "System error", http.StatusInternalServerError)
		return
	}
	if payload.UseOTP {
		// Return OTP provisioning URL for user to set up their authenticator app
		writeRespWithStatus(w, "User created with OTP", map[string]string{"otp_url": otpUrl}, http.StatusCreated)
		slog.Info("User created with OTP", "username", payload.Username)
		return
	}
	// Success
	writeRespWithStatus(w, "User created", nil, http.StatusCreated)
	slog.Info("User created", "username", payload.Username)
}

func (a *API) handleUserUpdate(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Username string `json:"username"`
		Password string `json:"password"`
		OTP      string `json:"otp,omitempty"` // optional OTP for 2FA
	}{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeRespErr(w, "Invalid request payload", http.StatusBadRequest)
		slog.Error("Failed to decode request payload", "error", err)
		return
	}

	user, err := a.store.FindUser(payload.Username)
	if err != nil {
		writeRespErr(w, "User not found", http.StatusNotFound)
		slog.Error("User not found", "username", payload.Username, "error", err)
		return
	}

	if payload.Password != "" {
		user.Password = payload.Password
	}
	if payload.OTP != "" {
		user.OTP = payload.OTP
	}

	if err := a.store.UpdateUser(user); err != nil {
		writeRespErr(w, "Update failed", http.StatusInternalServerError)
		slog.Error("Database write failed", "error", err)
		return
	}

	writeRespWithStatus(w, "User updated", nil, http.StatusOK)
	slog.Info("User updated", "username", payload.Username)
}
