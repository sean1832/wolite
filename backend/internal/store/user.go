package store

import "fmt"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	OTP      string `json:"otp,omitempty"`
	Devices  []int  `json:"devices,omitempty"` // list of device IDs the user has access to
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
