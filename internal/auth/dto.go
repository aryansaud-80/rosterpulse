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
