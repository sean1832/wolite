//go:build gui

package main

import (
	_ "embed"
	"log/slog"
	"wolcompanion/internal/auth"

	"github.com/getlantern/systray"
)

// Global tokenStore for cleanup on exit
var globalTokenStore *auth.TokenStore

//go:embed assets/logo.ico
var iconData []byte

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("Wolite Companion")
	systray.SetTooltip("Wake on Lan client side monitoring tool")

	mRegen := systray.AddMenuItem("Generate Token", "Generate a new authentication token")
	mConfig := systray.AddMenuItem("Configure", "Configure the app")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	// Initialize App
	tokenStore, certPath, keyPath, err := Initialize()
	if err != nil {
		slog.Error("Failed to Initialize app", "error", err)
		systray.Quit()
		return
	}
	configStore, err := InitializeConfig(8081)
	if err != nil {
		slog.Error("Failed to Initialize config", "error", err)
		systray.Quit()
		return
	}

	// Store globally for cleanup
	globalTokenStore = tokenStore

	// Start the server in a goroutine
	go func() {
		if err := StartServer(tokenStore, certPath, keyPath, configStore.GetPort()); err != nil {
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
					slog.Info("Token regenerated", "token", newToken)
					// Open the token file in the default text editor
					if err := openFileInEditor(tokenStore.GetTempTokenPath()); err != nil {
						slog.Error("Failed to open token file", "error", err)
					}
				}
			case <-mConfig.ClickedCh:
				// Open the config file in the default text editor
				if err := openFileInEditor(configStore.GetConfigPath()); err != nil {
					slog.Error("Failed to open config file", "error", err)
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
