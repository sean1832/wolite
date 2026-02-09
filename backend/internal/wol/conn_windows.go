//go:build windows

package wol

import (
	"net"
	"syscall"
)

func setBroadcast(conn *net.UDPConn) error {
	c, err := conn.SyscallConn()
	if err != nil {
		return err
	}
	var err2 error
	err = c.Control(func(fd uintptr) {
		err2 = syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
	})
	if err != nil {
		return err
	}
	return err2
}
