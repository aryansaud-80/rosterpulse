package auth

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrIncorrectPassword  = errors.New("incorrect current password")
	ErrInvalidCredentials = errors.New("invalid email or password")

	ErrInvalidToken          = errors.New("invalid token")
	ErrExpiredToken          = errors.New("token has expired")
	ErrInvalidRefreshToken   = errors.New("invalid refresh token")
	ErrTokenGenerationFailed = errors.New("failed to generate token")
	ErrEmptyTokenHash        = errors.New("refresh token hash cannot be empty")

	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session has expired")
	ErrSessionRevoked  = errors.New("session has been revoked")
	ErrSessionCreation = errors.New("failed to create session")
)
