package store

import "sort"

type UserDeviceMapping struct {
	Username   string `json:"username"`
	MACAddress string `json:"mac_address"`
}

// GetDevicesForUser returns all devices associated with a username.
func (s *Store) GetDevicesForUser(username string) ([]Device, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	mappings, ok := s.userDeviceMappings[username]
	if !ok {
		return []Device{}, nil // Return empty list if no mappings found
	}

	devices := make([]Device, 0, len(mappings))
	for mac := range mappings {
		if device, exists := s.devices[mac]; exists {
			devices = append(devices, device)
		}
	}

	sortDevices(devices)

	return devices, nil
}

// sort devices by Order, then Name
func sortDevices(devices []Device) {
	sort.Slice(devices, func(i, j int) bool {
		if devices[i].Order != devices[j].Order {
			return devices[i].Order < devices[j].Order
		}
		return devices[i].Name < devices[j].Name
	})
}

// AddDeviceToUser adds a new device to a user only if the mapping does not exist.
func (s *Store) AddDeviceToUser(username string, device *Device) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Ensure both device and username exists
	if _, exists := s.users[username]; !exists {
		return ErrUserNotFound
	}
	if _, exists := s.devices[device.MACAddress]; !exists {
		return ErrDeviceNotFound
	}

	// Ensure inner map exists
	if s.userDeviceMappings[username] == nil {
		s.userDeviceMappings[username] = make(map[string]UserDeviceMapping)
	}

	// Guard: Check existence
	if _, exists := s.userDeviceMappings[username][device.MACAddress]; exists {
		return ErrUserDeviceMappingExists
	}

	// Action: Write to map
	s.userDeviceMappings[username][device.MACAddress] = UserDeviceMapping{
		Username:   username,
		MACAddress: device.MACAddress,
	}

	// Persistence: Flush to disk
	return s.flush()
}

func (s *Store) RemoveDeviceFromUser(username, macAddress string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Ensure mapping exists
	if _, exists := s.userDeviceMappings[username][macAddress]; !exists {
		return ErrUserDeviceMappingNotFound
	}

	// Action: Remove from map
	delete(s.userDeviceMappings[username], macAddress)

	// Persistence: Flush to disk
	return s.flush()
}

// GetDeviceForUser returns a device only if it is associated with the given username.
func (s *Store) GetDeviceForUser(username, macAddress string) (*Device, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 1. Check if user exists (optional, but good for error clarity)
	if _, exists := s.users[username]; !exists {
		return nil, ErrUserNotFound
	}

	// 2. Check mapping
	mappings, ok := s.userDeviceMappings[username]
	if !ok {
		return nil, ErrDeviceNotFound // Effectively not found for this user
	}

	if _, ok := mappings[macAddress]; !ok {
		return nil, ErrDeviceNotFound
	}

	// 3. Get actual device
	device, ok := s.devices[macAddress]
	if !ok {
		return nil, ErrDeviceNotFound // Should not happen if consistency is maintained
	}

	return &device, nil
}

// CreateDeviceForUser atomically creates a device and assigns it to a user.
func (s *Store) CreateDeviceForUser(username string, device *Device) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. Check user existence
	if _, exists := s.users[username]; !exists {
		return ErrUserNotFound
	}

	// 2. Check device existence
	if _, exists := s.devices[device.MACAddress]; exists {
		return ErrDeviceExists
	}

	// 3. Initialize mapping map if nil
	if s.userDeviceMappings[username] == nil {
		s.userDeviceMappings[username] = make(map[string]UserDeviceMapping)
	}

	// 4. Perform writes
	s.devices[device.MACAddress] = *device
	s.userDeviceMappings[username][device.MACAddress] = UserDeviceMapping{
		Username:   username,
		MACAddress: device.MACAddress,
	}

	// 5. Persist
	return s.flush()
}
