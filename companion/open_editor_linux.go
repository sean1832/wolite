//go:build gui && !windows

package main

import "os/exec"

// openFileInEditor opens the given file in the system's default application.
func openFileInEditor(path string) error {
	return exec.Command("xdg-open", path).Start()
}
