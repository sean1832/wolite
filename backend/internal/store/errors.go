package store

import "errors"

var (
<<<<<<< HEAD
	ErrUserNotFound   = errors.New("user not found)")
	ErrDeviceNotFound = errors.New("device not found)")
	ErrUserExists     = errors.New("user already exists")
=======
	ErrUserNotFound              = errors.New("user not found")
	ErrDeviceNotFound            = errors.New("device not found")
	ErrUserExists                = errors.New("user already exists")
	ErrDeviceExists              = errors.New("device already exists")
	ErrUserDeviceMappingExists   = errors.New("user-device mapping already exists")
	ErrUserDeviceMappingNotFound = errors.New("user-device mapping not found")
>>>>>>> 783f6b3d4350d11bfa0b962a4329534f17ed71de
)
