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
<<<<<<< HEAD
		user.OTP = secret
=======
		user.PendingOTP = secret // Store in pending until verified
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
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
<<<<<<< HEAD
=======
	// Auto-login: Generate JWT token and set cookie
	tokenString, expirationTime, err := auth.GenerateJWTToken(user.Username, []byte(a.config.JWTSecret), a.config.JWTExpiry)
	if err != nil {
		slog.Error("Failed to generate token during auto-login", "error", err)
		// Don't fail the request, just don't auto-login
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			Expires:  expirationTime,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
			// Secure:   true, // TODO: Enable in production
		})
	}

>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
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
<<<<<<< HEAD
	payload := struct {
		Username string `json:"username"`
		Password string `json:"password"`
		OTP      string `json:"otp,omitempty"` // optional OTP for 2FA
=======
	claims := GetUserFromContext(r.Context())
	if claims == nil {
		slog.Error("claims missing from context", "path", r.URL.Path)
		writeRespErr(w, "internal server error", http.StatusInternalServerError)
		return
	}

	payload := struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		OldPassword string `json:"old_password,omitempty"`
		UseOTP      bool   `json:"use_otp,omitempty"` // optional flag to indicate if user wants to use OTP
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
	}{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeRespErr(w, "Invalid request payload", http.StatusBadRequest)
		slog.Error("Failed to decode request payload", "error", err)
		return
	}

<<<<<<< HEAD
	user, err := a.store.FindUser(payload.Username)
	if err != nil {
		writeRespErr(w, "User not found", http.StatusNotFound)
		slog.Error("User not found", "username", payload.Username, "error", err)
=======
	// Determine if username change is requested (not supported yet, but for completeness)
	if payload.Username != "" && payload.Username != claims.Username {
		writeRespErr(w, "Forbidden: Cannot change username", http.StatusForbidden)
		return
	}

	user, err := a.store.FindUser(claims.Username)
	if err != nil {
		writeRespErr(w, "User not found", http.StatusNotFound)
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
		return
	}

	if payload.Password != "" {
<<<<<<< HEAD
		user.Password = payload.Password
	}
	if payload.OTP != "" {
		user.OTP = payload.OTP
=======
		// Secure Password Change Flow
		if payload.OldPassword == "" {
			writeRespErr(w, "Current password is required to set a new password", http.StatusBadRequest)
			return
		}

		if !auth.CheckPasswordHash(payload.OldPassword, user.Password) {
			writeRespErr(w, "Invalid current password", http.StatusUnauthorized)
			return
		}

		hashedPassword, err := auth.HashPassword(payload.Password)
		if err != nil {
			slog.Error("Password hashing failed", "error", err)
			writeRespErr(w, "System error", http.StatusInternalServerError)
			return
		}
		user.Password = hashedPassword
	}
	var otpUrl string
	if payload.UseOTP {
		secret, url, err := auth.GenerateOTPSecret(payload.Username)
		if err != nil {
			slog.Error("OTP secret generation failed", "error", err)
			writeRespErr(w, "System error", http.StatusInternalServerError)
			return
		}
		user.PendingOTP = secret // Store in pending until verified
		otpUrl = url
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
	}

	if err := a.store.UpdateUser(user); err != nil {
		writeRespErr(w, "Update failed", http.StatusInternalServerError)
		slog.Error("Database write failed", "error", err)
		return
	}

<<<<<<< HEAD
	writeRespWithStatus(w, "User updated", nil, http.StatusOK)
	slog.Info("User updated", "username", payload.Username)
}
=======
	if payload.UseOTP {
		writeRespWithStatus(w, "User updated with OTP", map[string]string{"otp_url": otpUrl}, http.StatusOK)
		slog.Info("User updated with OTP", "username", payload.Username)
		return
	}

	writeRespWithStatus(w, "User updated", nil, http.StatusOK)
	slog.Info("User updated", "username", payload.Username)
}

func (a *API) handleUserOTPVerify(w http.ResponseWriter, r *http.Request) {
	claims := GetUserFromContext(r.Context())
	if claims == nil {
		writeRespErr(w, "internal server error", http.StatusInternalServerError)
		return
	}

	payload := struct {
		Code string `json:"code"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeRespErr(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := a.store.FindUser(claims.Username)
	if err != nil {
		writeRespErr(w, "User not found", http.StatusNotFound)
		return
	}

	if user.PendingOTP == "" {
		writeRespErr(w, "No pending OTP setup found", http.StatusBadRequest)
		return
	}

	if !auth.Validate2FA(payload.Code, user.PendingOTP) {
		writeRespErr(w, "Invalid OTP code", http.StatusBadRequest)
		return
	}

	// Code is valid, promote PendingOTP to OTP
	user.OTP = user.PendingOTP
	user.PendingOTP = ""

	if err := a.store.UpdateUser(user); err != nil {
		writeRespErr(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	writeRespOk(w, "2FA enabled successfully", nil)
	slog.Info("2FA verified and enabled", "username", claims.Username)
}
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
