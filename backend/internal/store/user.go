package store

import "fmt"

type User struct {
<<<<<<< HEAD
	Username string `json:"username"`
	Password string `json:"password"`
	OTP      string `json:"otp,omitempty"`
	Devices  []int  `json:"devices,omitempty"` // list of device IDs the user has access to
=======
	Username   string `json:"username"`
	Password   string `json:"password"`
	OTP        string `json:"otp,omitempty"`
	PendingOTP string `json:"pending_otp,omitempty"` // Temp storage for OTP verification
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
}

func NewUser(username, password string) (*User, error) {
	// ensure username and password are not empty
	if username == "" || password == "" {
		return nil, fmt.Errorf("username and password cannot be empty")
	}

	return &User{
		Username: username,
		Password: password,
	}, nil
}

func NewUserWithOTP(username, password, otp string) (*User, error) {
	// ensure username and password are not empty
	if username == "" || password == "" {
		return nil, fmt.Errorf("username and password cannot be empty")
	}

	return &User{
		Username: username,
		Password: password,
		OTP:      otp,
	}, nil
}
<<<<<<< HEAD
=======

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

// HasUsers returns true if at least one user exists in the store.
func (s *Store) HasUsers() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.users) > 0
}
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
