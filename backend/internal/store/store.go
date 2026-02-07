package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Store manages the JSON persistence.
type Store struct {
	mu   sync.RWMutex
	path string
	// Internal cache
	users   map[string]User   // map for O(1) lookup
	devices map[string]Device // map for O(1) lookup
}

// New initializes the store.
func New(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	s := &Store{
		path:    path,
		users:   make(map[string]User),
		devices: make(map[string]Device),
	}

	// Load existing data if file exists
	if _, err := os.Stat(path); err == nil {
		if err := s.load(); err != nil {
			return nil, err
		}
	} else {
		// Initialize empty file
		if err := s.flush(); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// FindUser returns a copy of the user.
// It returns a value (User), not a pointer, ensuring immutability of the internal cache.
func (s *Store) FindUser(username string) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.users[username]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return u, nil
}

// CreateUser adds a new user only if the username is unique.
func (s *Store) CreateUser(u User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Guard: Check existence
	if _, exists := s.users[u.Username]; exists {
		return ErrUserExists
	}

	// Action: Write to map
	s.users[u.Username] = u

	// Persistence: Flush to disk
	return s.flush()
}

// UpdateUser replaces an existing user's data.
// It fails if the user does not exist.
func (s *Store) UpdateUser(u User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Guard: Ensure target exists
	if _, exists := s.users[u.Username]; !exists {
		return ErrUserNotFound
	}

	// Action: Overwrite entry
	s.users[u.Username] = u

	// Persistence: Flush to disk
	return s.flush()
}

// flush writes the memory state to disk atomically.
func (s *Store) flush() error {
	// Convert maps to slices for JSON marshaling
	data := struct {
		Users   []User   `json:"users"`
		Devices []Device `json:"devices"`
	}{
		Users:   make([]User, 0, len(s.users)),
		Devices: make([]Device, 0, len(s.devices)),
	}

	for _, u := range s.users {
		data.Users = append(data.Users, u)
	}
	for _, d := range s.devices {
		data.Devices = append(data.Devices, d)
	}

	// Atomic Write Pattern
	tmp, err := os.CreateTemp(filepath.Dir(s.path), "db-tmp-*.json")
	if err != nil {
		return err
	}

	enc := json.NewEncoder(tmp)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}

	if err := tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}

	return os.Rename(tmp.Name(), s.path)
}

// load reads from disk into the maps.
func (s *Store) load() error {
	f, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Temp struct for decoding
	var data struct {
		Users   []User   `json:"users"`
		Devices []Device `json:"devices"`
	}

	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return err
	}

	// Hydrate maps
	s.users = make(map[string]User, len(data.Users))
	for _, u := range data.Users {
		s.users[u.Username] = u
	}

	s.devices = make(map[string]Device, len(data.Devices))
	for _, d := range data.Devices {
		s.devices[d.MACAddress] = d
	}

	return nil
}
