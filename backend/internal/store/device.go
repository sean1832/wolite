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

func (s *Store) GetDevice(macAddress string) (*Device, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	device, ok := s.devices[macAddress]
	if !ok {
		return nil, ErrDeviceNotFound
	}
	return &device, nil
}

func (s *Store) UpdateDevice(device *Device) error {
	return nil
}

func (s *Store) DeleteDevice(macAddress string) error {
	return nil
}
