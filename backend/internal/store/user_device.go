package store

type UserDeviceMapping struct {
	Username   string `json:"username"`
	MACAddress string `json:"mac_address"`
}

// AddUserDeviceMapping adds a new user-device mapping only if the mapping does not exist.
func (s *Store) AddUserDeviceMapping(mapping *UserDeviceMapping) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Guard: Check existence
	if _, exists := s.userDeviceMappings[mapping.Username][mapping.MACAddress]; exists {
		return ErrUserDeviceMappingExists
	}

	// Action: Write to map
	s.userDeviceMappings[mapping.Username][mapping.MACAddress] = *mapping

	// Persistence: Flush to disk
	return s.flush()
}
