package wol

import (
	"fmt"
	"log/slog"
	"net"
)

// SendMagicPacket sends a Magic Packet to the specified broadcast address.
//
// targetMAC:       The MAC address of the machine to wake (e.g., "00:11:22:33:44:55").
// broadcastAddress: The broadcast address of the network (e.g., "192.168.1.255:9").
//
// Note: Do NOT use the target machine's IP for broadcastAddress; use the network's broadcast IP.
//
// Broadcast address is usually the network address with the last octet set to 255.
// For example: target ip 192.168.50.100 -> broadcast 192.168.50.255
func SendMagicPacket(macAddress, broadcastAddr string) error {
	// Parse the MAC address
	mac, err := net.ParseMAC(macAddress)
	if err != nil {
		return fmt.Errorf("invalid mac: %w", err)
	}

	if len(mac) != 6 {
		return fmt.Errorf("invalid mac length: %d bytes (expected 6)", len(mac))
	}

	// construct payload: 6 bytes of 0xFF followed by MAC repeated 16 times
	// 102 bytes = 6 bytes header + (16 * 6 bytes MAC)
	var packet [102]byte
	// Copy header (6x 0xFF)
	copy(packet[:], "\xff\xff\xff\xff\xff\xff")

	// fill payload: 16 repetition of the mac addr
	for i := 6; i < 102; i += 6 {
		copy(packet[i:], mac)
	}

	slog.Debug("sending magic packet", "mac", macAddress, "broadcast", broadcastAddr, "packet_size", len(packet))

	// Resolve the UDP address
	addr, err := net.ResolveUDPAddr("udp4", broadcastAddr)
	if err != nil {
		return fmt.Errorf("invalid broadcast address: %w", err)
	}

	// Use DialUDP instead of Dial to access UDP-specific methods
	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to dial UDP: %w", err)
	}
	defer conn.Close()

	// Enable broadcast
	// Linux/Unix requires this socket option to be set to send to broadcast addresses (e.g. 255.255.255.255).
	// Without this, the call to Write() may fail with "permission denied" or "invalid argument".
	if err := setBroadcast(conn); err != nil {
		return fmt.Errorf("failed to set broadcast: %w", err)
	}

	_, err = conn.Write(packet[:])
	if err != nil {
		return fmt.Errorf("failed to send packet: %w", err)
	}
	return nil
}
