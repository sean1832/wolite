//go:build gui

package main

import (
	"log/slog"
	"wolcompanion/internal/auth"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
)

// Global tokenStore for cleanup on exit
var globalTokenStore *auth.TokenStore

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Wolite Companion")
	systray.SetTooltip("Wake on Lan client side monitoring tool")

	mRegen := systray.AddMenuItem("Regenerate Token", "Generate a new authentication token")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	// Sets the icon of a menu item. Only available on Mac and Windows.
	mQuit.SetIcon(icon.Data)

	// Initialize App
	tokenStore, certPath, keyPath, err := Initialize()
	if err != nil {
		slog.Error("Failed to Initialize app", "error", err)
		systray.Quit()
		return
	}

	// Store globally for cleanup
	globalTokenStore = tokenStore

	// Start the server in a goroutine
	go func() {
		if err := StartServer(tokenStore, certPath, keyPath); err != nil {
			slog.Error("Server crashed", "error", err)
			systray.Quit()
		}
	}()

	// Handle menu clicks
	go func() {
		for {
			select {
			case <-mRegen.ClickedCh:
				newToken, err := tokenStore.Regenerate()
				if err != nil {
					slog.Error("Failed to regenerate token", "error", err)
				} else {
					slog.Info("Token regenerated manually via tray", "token", newToken)
				}
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

// onExit is called when the systray app is exiting.
func onExit() {
	if globalTokenStore != nil {
		globalTokenStore.Cleanup()
	}
}
