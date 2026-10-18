package user

import "errors"

var (
	ErrEmailTaken           = errors.New("email is already registered")
	ErrUserNotFound         = errors.New("user not found")
	ErrPasswordLength       = errors.New("password must be at least 8 characters long")
	ErrInvalidRole          = errors.New("invalid role")
	ErrInvalidEmail         = errors.New("invalid email")
	ErrFailedToHashPassword = errors.New("failed to hash password")
	ErrFailedToCreateUser   = errors.New("failed to create user")
	ErrRoleIsRequired       = errors.New("role is required")
)
