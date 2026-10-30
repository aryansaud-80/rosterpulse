package profile

import (
	"time"

	"github.com/google/uuid"
)

type CreateProfileRequest struct {
	DisplayName         string     `json:"display_name" validate:"required,min=3,max=50"`
	Bio                 string     `json:"bio" validate:"max=500"`
	Region              string     `json:"region" validate:"required"`
	ProfileImageAssetId *uuid.UUID `json:"profile_image_asset_id" validate:"omitempty"`
	Status              string     `json:"status" validate:"omitempty"`
}

type UpdateProfileParams struct {
	DisplayName         string     `json:"display_name" validate:"omitempty,min=3,max=50"`
	Bio                 string     `json:"bio" validate:"omitempty,max=500"`
	Region              string     `json:"region" validate:"omitempty"`
	ProfileImageAssetId *uuid.UUID `json:"profile_image_asset_id" validate:"omitempty"`
	Status              string     `json:"status" validate:"omitempty,oneof=ACTIVE HIDDEN DISABLED"`
}

type ProfileResponse struct {
	Id                  string     `json:"id"`
	UserId              string     `json:"user_id"`
	DisplayName         string     `json:"display_name"`
	Bio                 string     `json:"bio"`
	Region              string     `json:"region"`
	ProfileImageAssetId *uuid.UUID `json:"profile_image_asset_id"`
	Status              string     `json:"status"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at"`
}

type PublicProfileResponse struct {
	Id                  string     `json:"id"`
	DisplayName         string     `json:"display_name"`
	Bio                 string     `json:"bio"`
	Region              string     `json:"region"`
	ProfileImageAssetId *uuid.UUID `json:"profile_image_asset_id"`
	Status              string     `json:"status"`
	CreatedAt           time.Time  `json:"created_at"`
}
