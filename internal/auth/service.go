package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/aryansaud-80/rosterpulse/internal/user"
	"golang.org/x/crypto/bcrypt"
)

type RefreshResponse struct {
	accessToken  string
	refreshToken string
}

type IAuthService interface {
	Login(ctx context.Context, req *LoginReq) (accessToken, refreshToken string, userResp *user.UserResponse, err error)
	Me(ctx context.Context, user_id string) (*MeResponse, error)
	Refresh(ctx context.Context, token string) (*RefreshResponse, error)
	Logout(ctx context.Context, token string) error
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

func (s *AuthService) Me(ctx context.Context, user_id string) (*MeResponse, error) {
	users, err := s.userRepo.FindById(ctx, user_id)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		s.log.Error("failed to get user profile", "user_id", user_id, "error", err)
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	userResp := &MeResponse{
		ID:              users.ID,
		Email:           users.Email,
		NormalizedEmail: users.NormalizedEmail,
		Role:            users.Role,
		Status:          users.Status,
		CreatedAt:       users.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       users.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return userResp, nil
}

func (s *AuthService) Refresh(ctx context.Context, token string) (*RefreshResponse, error) {
	hToken := HashRefreshToken(token)
	if strings.TrimSpace(hToken) == "" {
		return nil, ErrEmptyTokenHash
	}

	session, err := s.srepo.FindByHash(ctx, hToken)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			return nil, ErrSessionNotFound
		}
		s.log.Error("failed to find session for refresh", "error", err)
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if time.Now().After(session.ExpiresAt) || session.RevokedAt != nil {
		return nil, ErrInvalidToken
	}

	_, err = s.userRepo.FindById(ctx, session.UserId)
	if err != nil {
		s.log.Warn("user associated with session no longer exists", "user_id", session.UserId, "error", err)
		return nil, ErrUserNotFound
	}

	now := time.Now()

	tx, err := s.srepo.BeginTx(ctx)
	if err != nil {
		s.log.Error("failed to begin transaction for token refresh", "error", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	txRepo := s.srepo.WithTx(tx)

	err = txRepo.RevokeAndUpdateLastUsed(ctx, session.ID, now)
	if err != nil {
		s.log.Error("failed to revoke old session during rotation", "session_id", session.ID, "error", err)
		return nil, fmt.Errorf("failed to revoke session: %w", err)
	}

	newRToken, err := GenerateRefreshToken()
	if err != nil {
		s.log.Error("failed to generate new cryptographically secure refresh token", "error", err)
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	newHashToken := HashRefreshToken(newRToken)
	refreshTokenDuration := s.tokenManager.GetRefreshTokenDuration()
	expireAt := now.Add(refreshTokenDuration)

	newSession := &Session{
		UserId:           session.UserId,
		RefreshTokenHash: newHashToken,
		TokenVersion:     session.TokenVersion + 1,
		ExpiresAt:        expireAt,
	}

	sess, err := txRepo.CreateSession(ctx, newSession)
	if err != nil {
		s.log.Error("failed to create new rotated session record", "user_id", session.UserId, "error", err)
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.log.Error("failed to commit token rotation transaction", "error", err)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	aToken, err := s.tokenManager.GenerateAccessToken(session.UserId, sess.Id, int(newSession.TokenVersion))
	if err != nil {
		s.log.Error("failed to generate access token after rotation", "user_id", session.UserId, "error", err)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &RefreshResponse{
		accessToken:  aToken,
		refreshToken: newRToken,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	hToken := HashRefreshToken(token)
	if strings.TrimSpace(hToken) == "" {
		return ErrEmptyTokenHash
	}

	session, err := s.srepo.FindByHash(ctx, hToken)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			return ErrSessionNotFound
		}
		s.log.Error("failed to find session for logout", "error", err)
		return fmt.Errorf("failed to get session: %w", err)
	}

	now := time.Now()

	err = s.srepo.RevokeAndUpdateLastUsed(ctx, session.ID, now)
	if err != nil {
		s.log.Error("failed to revoke session during logout", "session_id", session.ID, "error", err)
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	return nil
}
