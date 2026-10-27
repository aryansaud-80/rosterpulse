package auth

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/aryansaud-80/rosterpulse/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Login(ctx context.Context, req *LoginReq) (accessToken, refreshToken string, userResp *user.UserResponse, err error)
}

type AuthService struct {
	srepo        ISessionRepo
	userRepo     user.Iuser
	log          *slog.Logger
	tokenManager *TokenManager
}

func NewAuthService(srepo ISessionRepo, urepo user.Iuser, log *slog.Logger, tokenManager *TokenManager) *AuthService {
	return &AuthService{
		srepo:        srepo,
		userRepo:     urepo,
		log:          log.With("component", "auth.service"),
		tokenManager: tokenManager,
	}
}

func (s *AuthService) Login(ctx context.Context, req *LoginReq) (accessToken, refreshToken string, userResp *user.UserResponse, err error) {
	cleanEmail := strings.ToLower(strings.TrimSpace(req.Email))

	existing, err := s.userRepo.FindByEmail(ctx, cleanEmail)
	if err != nil {
		return "", "", nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(existing.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", "", nil, ErrInvalidCredentials
	}

	rToken, err := GenerateRefreshToken()
	if err != nil {
		return "", "", nil, ErrTokenGenerationFailed
	}
	hrToken := HashRefreshToken(rToken)

	if strings.TrimSpace(hrToken) == "" {
		return "", "", nil, ErrEmptyTokenHash
	}

	refreshTokenDuration := s.tokenManager.GetRefreshTokenDuration()
	expireAt := time.Now().Add(refreshTokenDuration)

	session := &Session{
		UserId:           existing.ID,
		RefreshTokenHash: hrToken,
		TokenVersion:     1,
		ExpiresAt:        expireAt,
	}

	ses, err := s.srepo.CreateSession(ctx, session)
	if err != nil {
		return "", "", nil, ErrSessionCreation
	}

	aToken, err := s.tokenManager.GenerateAccessToken(existing.ID, ses.Id, 1)
	if err != nil {
		return "", "", nil, ErrTokenGenerationFailed
	}

	uRep := &user.UserResponse{
		ID:        existing.ID,
		Email:     existing.Email,
		Role:      existing.Role,
		Status:    existing.Status,
		CreatedAt: existing.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return aToken, rToken, uRep, nil
}
