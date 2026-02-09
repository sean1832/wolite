package api

import (
	"log/slog"
	"net/http"
	"wolcompanion/internal/commands"
)

func (a *API) handleShutdown(w http.ResponseWriter, r *http.Request) {
	// command validation
	// catch unsupported platforms or missing binaries IMMEDIATELY.
	execute, err := commands.PrepareShutdown(3)
	if err != nil {
		writeRespErr(w, "command validation failed", http.StatusInternalServerError)
		slog.Error("SHUTDOWN validation failed", "error", err, "ip", r.RemoteAddr)
		return
	}

	writeRespOk(w, "SHUTDOWN command validated, executing in 3 seconds", nil)

	go func() {
		// wait 3 seconds, then execute
		if err := execute(); err != nil {
			// at this point, the HTTP connection is closed. Only log the execution failure.
			slog.Error("failed to execute delayed SHUTDOWN", "error", err)
		}
		slog.Info("command executed", "command", "SHUTDOWN", "ip", r.RemoteAddr)
	}()
}

func (a *API) handleReboot(w http.ResponseWriter, r *http.Request) {
	// command validation
	// catch unsupported platforms or missing binaries IMMEDIATELY.
	execute, err := commands.PrepareReboot(3)
	if err != nil {
		writeRespErr(w, "command validation failed", http.StatusInternalServerError)
		slog.Error("REBOOT validation failed", "error", err, "ip", r.RemoteAddr)
		return
	}

	writeRespOk(w, "REBOOT command validated, executing in 3 seconds", nil)

	go func() {
		// wait 3 seconds, then execute
		if err := execute(); err != nil {
			// at this point, the HTTP connection is closed. Only log the execution failure.
			slog.Error("failed to execute delayed REBOOT", "error", err, "ip", r.RemoteAddr)
		}
		slog.Info("command executed", "command", "REBOOT", "ip", r.RemoteAddr)
	}()
}

func (a *API) handleSleep(w http.ResponseWriter, r *http.Request) {

	// command validation
	// catch unsupported platforms or missing binaries IMMEDIATELY.
	execute, err := commands.PrepareSleep(3)
	if err != nil {
		writeRespErr(w, "command validation failed", http.StatusInternalServerError)
		slog.Error("SLEEP validation failed", "error", err, "ip", r.RemoteAddr)
		return
	}

	writeRespOk(w, "SLEEP command validated, executing in 3 seconds", nil)

	go func() {
		// wait 3 seconds, then execute
		if err := execute(); err != nil {
			// at this point, the HTTP connection is closed. Only log the execution failure.
			slog.Error("failed to execute delayed SLEEP", "error", err, "ip", r.RemoteAddr)
		}
		slog.Info("command executed", "command", "SLEEP", "ip", r.RemoteAddr)
	}()
}

func (a *API) handleHibernate(w http.ResponseWriter, r *http.Request) {
	// command validation
	// catch unsupported platforms or missing binaries IMMEDIATELY.
	execute, err := commands.PrepareHibernate(3)
	if err != nil {
		writeRespErr(w, "command validation failed", http.StatusInternalServerError)
		slog.Error("HIBERNATE validation failed", "error", err, "ip", r.RemoteAddr)
		return
	}

	writeRespOk(w, "HIBERNATE command validated, executing in 3 seconds", nil)

	go func() {
		// wait 3 seconds, then execute
		if err := execute(); err != nil {
			// at this point, the HTTP connection is closed. Only log the execution failure.
			slog.Error("failed to execute delayed HIBERNATE", "error", err, "ip", r.RemoteAddr)
		}
		slog.Info("command executed", "command", "HIBERNATE", "ip", r.RemoteAddr)
	}()
}
