package user

import (
	"context"
	"log/slog"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo Iuser
	log  *slog.Logger
}

func NewUserService(repo Iuser, log *slog.Logger) *UserService {
	return &UserService{
		repo: repo,
		log:  log.With("component", "user.service"),
	}
}

func (s *UserService) Register(ctx context.Context, req *UserRequest) (*UserResponse, error) {

	cleanEmail := strings.ToLower(strings.TrimSpace(req.Email))

	existing, err := s.repo.FindByEmail(ctx, cleanEmail)
	if err == nil && existing != nil {
		return nil, ErrEmailTaken
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		s.log.Error("failed to hash password", "error", err)
		return nil, ErrFailedToHashPassword
	}

	record := &UserRecord{
		Email:           req.Email,
		NormalizedEmail: cleanEmail,
		PasswordHash:    string(hashedBytes),
		Role:            req.Role,
		Status:          "ACTIVE",
	}

	user, err := s.repo.Create(ctx, record)

	if err != nil {
		s.log.Error("failed to create user", "error", err)
		return nil, ErrFailedToCreateUser
	}

	res := &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return res, nil

}
