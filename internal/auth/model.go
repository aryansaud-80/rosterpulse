package auth

import "time"

type SessionRecord struct {
	ID               string
	UserId           string
	RefreshTokenHash string
	TokenVersion     int64
	ExpiresAt        time.Time
	RevokedAt        *time.Time
	CreatedAt        time.Time
	LastUsedAt       *time.Time
}
