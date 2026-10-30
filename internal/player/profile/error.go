package profile

import "errors"

var (
	ErrProfileNotFound = errors.New("player profile not found")
	ErrProfileAlreadyExists = errors.New("player profile already exists")
	ErrInvalidInput = errors.New("invalid input fields provided")
)
