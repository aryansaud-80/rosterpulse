package profile

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type IProfileService interface {
	CreateProfile(ctx context.Context, userId string, req CreateProfileRequest) error
	GetProfileByUserID(ctx context.Context, userId string) (*ProfileResponse, error)
	UpdateProfile(ctx context.Context, userId string, req *UpdateProfileParams) error
	GetPublicProfileByID(ctx context.Context, id string) (*PublicProfileResponse, error)
}

type ProfileService struct {
	repo IProfileRepo
	log  *slog.Logger
}

func NewProfileService(profileRepo IProfileRepo, log *slog.Logger) *ProfileService {
	return &ProfileService{
		repo: profileRepo,
		log:  log.With("component", "profile.service"),
	}
}

func (s *ProfileService) CreateProfile(ctx context.Context, userId string, req CreateProfileRequest) error {
	profile := &ProfileRecord{
		UserId:              userId,
		DisplayName:         req.DisplayName,
		Bio:                 req.Bio,
		Region:              req.Region,
		ProfileImageAssetId: req.ProfileImageAssetId,
		Status:              req.Status,
	}

	err := s.repo.Create(ctx, profile)
	if err != nil {
		if errors.Is(err, ErrProfileAlreadyExists) {
			return ErrProfileAlreadyExists
		}
		s.log.Error("failed to create player profile in repo", "user_id", userId, "err", err)
		return fmt.Errorf("failed to insert player profile: %w", err)
	}

	return nil
}

func (s *ProfileService) GetProfileByUserID(ctx context.Context, userId string) (*ProfileResponse, error) {
	profile, err := s.repo.FindByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			return nil, ErrProfileNotFound
		}
		s.log.Error("failed to fetch profile by user id", "user_id", userId, "err", err)
		return nil, fmt.Errorf("failed to fetch profile: %w", err)
	}

	return &ProfileResponse{
		Id:                  profile.Id,
		UserId:              profile.UserId,
		DisplayName:         profile.DisplayName,
		Bio:                 profile.Bio,
		Region:              profile.Region,
		ProfileImageAssetId: profile.ProfileImageAssetId,
		Status:              profile.Status,
		CreatedAt:           profile.CreatedAt,
		UpdatedAt:           profile.UpdatedAt,
		DeletedAt:           profile.DeletedAt,
	}, nil
}

func (s *ProfileService) GetPublicProfileByID(ctx context.Context, id string) (*PublicProfileResponse, error) {
	profile, err := s.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			return nil, ErrProfileNotFound
		}
		s.log.Error("failed to fetch public profile by id", "profile_id", id, "err", err)
		return nil, fmt.Errorf("failed to fetch profile: %w", err)
	}

	return &PublicProfileResponse{
		Id:                  profile.Id,
		DisplayName:         profile.DisplayName,
		Bio:                 profile.Bio,
		Region:              profile.Region,
		ProfileImageAssetId: profile.ProfileImageAssetId,
		Status:              profile.Status,
		CreatedAt:           profile.CreatedAt,
	}, nil
}

func (s *ProfileService) UpdateProfile(ctx context.Context, userId string, req *UpdateProfileParams) error {
	err := s.repo.Update(ctx, userId, req)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			return ErrProfileNotFound
		}
		s.log.Error("failed to update player profile in repo", "user_id", userId, "err", err)
		return fmt.Errorf("failed to update player profile: %w", err)
	}

	return nil
}
