package store

type Status string

const (
	StatusOnline  Status = "online"
	StatusOffline Status = "offline"
	StatusUnknown Status = "unknown"
	StatusError   Status = "error"
)

type Device struct {
	MACAddress  string `json:"mac_address"`           // unique identifier for the device
	Name        string `json:"name"`                  // human-readable name for the device
	Description string `json:"description,omitempty"` // optional description of the device
	IPAddress   string `json:"ip_address,omitempty"`  // optional IP address of the device

	Status Status `json:"status"` // current status of the device
}

func NewDevice(macAddress, name, description, ipAddress string, status Status) *Device {
	return &Device{
		MACAddress:  macAddress,
		Name:        name,
		Description: description,
		IPAddress:   ipAddress,
		Status:      status,
	}
}

// AddDevice adds a new device only if the MAC address is unique.
func (s *Store) AddDevice(device *Device) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Guard: Check existence
	if _, exists := s.devices[device.MACAddress]; exists {
		return ErrDeviceExists
	}

	// Action: Write to map
	s.devices[device.MACAddress] = *device

	// Persistence: Flush to disk
	return s.flush()
}

// GetDeviceByMacAddress returns a copy of the device. It returns a value (Device), not a pointer, ensuring immutability of the internal cache.
func (s *Store) GetDeviceByMacAddress(macAddress string) (*Device, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	device, ok := s.devices[macAddress]
	if !ok {
		return nil, ErrDeviceNotFound
	}
	return &device, nil
}

// UpdateDevice updates an existing device. It requires the MAC address to be unchanged.
func (s *Store) UpdateDevice(device *Device) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Guard: Check existence
	if _, exists := s.devices[device.MACAddress]; !exists {
		return ErrDeviceNotFound
	}

	// Action: Write to map
	s.devices[device.MACAddress] = *device

	// Persistence: Flush to disk
	return s.flush()
}

// DeleteDevice removes a device from the store and all user-device mappings.
func (s *Store) DeleteDevice(macAddress string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Guard: Check existence
	if _, exists := s.devices[macAddress]; !exists {
		return ErrDeviceNotFound
	}

	// Action: Delete from map
	delete(s.devices, macAddress)

	// Clean up: Remove from all user mappings to ensure consistency
	for _, mappings := range s.userDeviceMappings {
		delete(mappings, macAddress)
	}

	// Persistence: Flush to disk
	return s.flush()
}
