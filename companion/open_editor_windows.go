//go:build gui && windows

package main

import (
	"syscall"
	"unsafe"
)

// openFileInEditor opens the given file in the system's default application.
func openFileInEditor(path string) error {
	shell32 := syscall.NewLazyDLL("shell32.dll")
	proc := shell32.NewProc("ShellExecuteW")

	verb, _ := syscall.UTF16PtrFromString("open")
	file, _ := syscall.UTF16PtrFromString(path)

	// ShellExecuteW(hwnd, lpOperation, lpFile, lpParameters, lpDirectory, nShowCmd)
	ret, _, err := proc.Call(
		0,
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(file)),
		0,
		0,
		1, // SW_SHOWNORMAL
	)

	// ShellExecute returns a value > 32 on success.
	if ret <= 32 {
		return err
	}
	return nil
}
