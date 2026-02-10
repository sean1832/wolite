package worker

import (
	"context"
	"log/slog"
	"time"
	"wolite/internal/companion"
	"wolite/internal/store"
)

type StatusChecker struct {
	store    *store.Store
	interval time.Duration
}

func NewStatusChecker(store *store.Store, interval time.Duration) *StatusChecker {
	return &StatusChecker{
		store:    store,
		interval: interval,
	}
}

func (s *StatusChecker) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.checkAll()
		}
	}
}

func (s *StatusChecker) checkAll() {
	devices, err := s.store.GetAllDevices()
	if err != nil {
		slog.Error("failed to get devices for status check", "error", err)
		return
	}

	for _, device := range devices {
		if device.CompanionURL == "" {
			continue
		}

		go s.checkDevice(device)
	}
}

func (s *StatusChecker) checkDevice(device store.Device) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := companion.NewClient(device.CompanionURL, device.CompanionToken, device.CompanionAuthFingerprint)
	if err != nil {
		slog.Warn("invalid companion config during check", "mac", device.MACAddress, "error", err)
		return
	}

	newStatus := store.StatusOffline
	if err := client.Ping(ctx); err == nil {
		newStatus = store.StatusOnline
	}

	// Only update if status changed
	if device.Status != newStatus {
		device.Status = newStatus
		if err := s.store.UpdateDevice(&device); err != nil {
			slog.Error("failed to update device status", "mac", device.MACAddress, "status", newStatus, "error", err)
		} else {
			slog.Info("device status updated", "mac", device.MACAddress, "status", newStatus)
		}
	}
}
