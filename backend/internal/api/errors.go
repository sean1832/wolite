package api

import "errors"

var (
	ErrUnauthorized   = errors.New("unauthorized")
	ErrInvalidRequest = errors.New("invalid request")
)
