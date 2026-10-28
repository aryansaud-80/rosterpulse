package auth

import (
	"time"

	"github.com/aryansaud-80/rosterpulse/internal/user"
)

type Session struct {
	Id               string
	UserId           string
	RefreshTokenHash string
	TokenVersion     int64
	ExpiresAt        time.Time
}

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	User        *user.UserResponse
}

type MeResponse struct {
	ID              string `json:"id"`
	Email           string `json:"email"`
	NormalizedEmail string `json:"normalized_email"`
	Role            string `json:"role"`
	Status          string `json:"status"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}
