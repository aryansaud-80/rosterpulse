package profile

import (
	"time"

	"github.com/google/uuid"
)

type ProfileRecord struct {
	Id                  string
	UserId              string
	DisplayName         string
	Bio                 string
	Region              string
	ProfileImageAssetId *uuid.UUID
	Status              string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
}
