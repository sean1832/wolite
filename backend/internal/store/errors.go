package store

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found)")
	ErrDeviceNotFound = errors.New("device not found)")
	ErrUserExists     = errors.New("user already exists")
)
