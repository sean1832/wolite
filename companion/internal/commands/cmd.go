package commands

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

// PrepareShutdown validates requirements and returns a closure to execute Shutdown.
func PrepareShutdown(delaySeconds int) (func() error, error) {
	var cmd *exec.Cmd

	if delaySeconds < 0 {
		return nil, fmt.Errorf("delaySeconds must be non-negative")
	}

	switch runtime.GOOS {
	case "windows":
		// Check for executable existence
		if _, err := exec.LookPath("shutdown"); err != nil {
			return nil, fmt.Errorf("shutdown executable not found")
		}
		// /s = shutdown, /t 0 = immediate (we handle delay in Go)
		cmd = exec.Command("shutdown", "/s", "/t", "0")
	case "linux", "darwin":
		if _, err := exec.LookPath("shutdown"); err != nil {
			return nil, fmt.Errorf("shutdown executable not found")
		}
		// -h = halt/shutdown, now = immediate
		cmd = exec.Command("shutdown", "-h", "now")
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return createOp(delaySeconds, cmd), nil
}

// PrepareReboot validates requirements and returns a closure to execute Reboot.
func PrepareReboot(delaySeconds int) (func() error, error) {
	var cmd *exec.Cmd

	if delaySeconds < 0 {
		return nil, fmt.Errorf("delaySeconds must be non-negative")
	}

	switch runtime.GOOS {
	case "windows":
		if _, err := exec.LookPath("shutdown"); err != nil {
			return nil, fmt.Errorf("shutdown executable not found")
		}
		// /r = reboot
		cmd = exec.Command("shutdown", "/r", "/t", "0")
	case "linux", "darwin":
		// 'reboot' is standard, but good to check
		if _, err := exec.LookPath("reboot"); err != nil {
			// Fallback to shutdown -r if reboot missing
			if _, err := exec.LookPath("shutdown"); err != nil {
				return nil, fmt.Errorf("neither reboot nor shutdown executables found")
			}
			cmd = exec.Command("shutdown", "-r", "now")
		} else {
			cmd = exec.Command("reboot")
		}
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return createOp(delaySeconds, cmd), nil
}

// PrepareHibernate validates requirements and returns a closure to execute Hibernate.
func PrepareHibernate(delaySeconds int) (func() error, error) {
	var cmd *exec.Cmd

	if delaySeconds < 0 {
		return nil, fmt.Errorf("delaySeconds must be non-negative")
	}

	switch runtime.GOOS {
	case "windows":
		if _, err := exec.LookPath("shutdown"); err != nil {
			return nil, fmt.Errorf("shutdown executable not found")
		}
		cmd = exec.Command("shutdown", "/h")
	case "linux":
		if _, err := exec.LookPath("systemctl"); err != nil {
			return nil, fmt.Errorf("systemctl not found, cannot hibernate")
		}
		cmd = exec.Command("systemctl", "hibernate")
	case "darwin":
		return nil, fmt.Errorf("manual hibernation not supported on macOS")
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return createOp(delaySeconds, cmd), nil
}

// PrepareSleep validates requirements and returns a closure to execute Sleep.
func PrepareSleep(delaySeconds int) (func() error, error) {
	var cmd *exec.Cmd

	if delaySeconds < 0 {
		return nil, fmt.Errorf("delaySeconds must be non-negative")
	}

	switch runtime.GOOS {
	case "windows":
		if _, err := exec.LookPath("powershell"); err != nil {
			return nil, fmt.Errorf("powershell not found")
		}
		psCmd := "Add-Type -AssemblyName System.Windows.Forms; [System.Windows.Forms.Application]::SetSuspendState('Suspend', $false, $false)"
		cmd = exec.Command("powershell", "-Command", psCmd)
	case "linux":
		if _, err := exec.LookPath("systemctl"); err != nil {
			return nil, fmt.Errorf("systemctl not found")
		}
		cmd = exec.Command("systemctl", "suspend")
	case "darwin":
		if _, err := exec.LookPath("pmset"); err != nil {
			return nil, fmt.Errorf("pmset not found")
		}
		cmd = exec.Command("pmset", "sleepnow")
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return createOp(delaySeconds, cmd), nil
}

// createOp acts as a factory for the unified execution logic.
// It captures the delay and the command in a closure.
func createOp(delaySeconds int, cmd *exec.Cmd) func() error {
	return func() error {
		if delaySeconds > 0 {
			time.Sleep(time.Duration(delaySeconds) * time.Second)
		}
		// cmd.Run() blocks until the command finishes
		return cmd.Run()
	}
}
