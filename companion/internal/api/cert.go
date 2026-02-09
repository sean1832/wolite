package api

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"wolcompanion/internal/tls"
)

func (a *API) handleCertFingerprint(w http.ResponseWriter, r *http.Request) {
	// get config directory
	configDir, err := os.UserCacheDir()
	if err != nil {
		writeRespErr(w, "failed to get config dir", http.StatusInternalServerError)
		slog.Error("failed to get config dir", "error", err, "ip", r.RemoteAddr)
		return
	}
	certPath := filepath.Join(configDir, "cert.pem")
	fingerprint, err := tls.GetCertFingerprint(certPath)
	if err != nil {
		writeRespErr(w, "failed to retrieve certificate fingerprint", http.StatusInternalServerError)
		slog.Error("failed to get certificate fingerprint", "error", err, "ip", r.RemoteAddr)
		return
	}
	writeRespOk(w, "certificate fingerprint retrieved successfully", map[string]string{"fingerprint": fingerprint})
	slog.Info("certificate fingerprint retrieved", "fingerprint", fingerprint, "ip", r.RemoteAddr)
}
