package user

import (
	"time"
)

type UserRecord struct {
	ID              string
	Email           string
	NormalizedEmail string
	PasswordHash    string
	Role            string
	Status          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}
