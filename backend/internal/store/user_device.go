package store

type UserDeviceMapping struct {
	Username   string `json:"username"`
	MACAddress string `json:"mac_address"`
}

// AddUserDeviceMapping adds a new user-device mapping only if the mapping does not exist.
func (s *Store) AddUserDeviceMapping(mapping *UserDeviceMapping) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Ensure inner map exists
	if s.userDeviceMappings[mapping.Username] == nil {
		s.userDeviceMappings[mapping.Username] = make(map[string]UserDeviceMapping)
	}

	// Guard: Check existence
	if _, exists := s.userDeviceMappings[mapping.Username][mapping.MACAddress]; exists {
		return ErrUserDeviceMappingExists
	}

	// Action: Write to map
	s.userDeviceMappings[mapping.Username][mapping.MACAddress] = *mapping

	// Persistence: Flush to disk
	return s.flush()
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

	return devices, nil
}
