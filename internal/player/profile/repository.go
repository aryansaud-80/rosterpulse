package profile

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/aryansaud-80/rosterpulse/internal/platform/exec"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IProfileRepo interface {
	Create(ctx context.Context, profile *ProfileRecord) error
	FindByUserId(ctx context.Context, userId string) (*ProfileRecord, error)
	FindById(ctx context.Context, id string) (*ProfileRecord, error)
	Update(ctx context.Context, userId string, params *UpdateProfileParams) error
}

type ProfileRepository struct {
	db  exec.DBExecutor
	log *slog.Logger
}

func NewProfileRepository(db *pgxpool.Pool, log *slog.Logger) *ProfileRepository {
	return &ProfileRepository{
		db:  db,
		log: log.With("component", "profile.repository"),
	}
}

func (r *ProfileRepository) Create(ctx context.Context, profile *ProfileRecord) error {
	query := `
		INSERT INTO player_profiles (user_id, display_name, bio, region, profile_image_asset_id, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		profile.UserId,
		profile.DisplayName,
		profile.Bio,
		profile.Region,
		profile.ProfileImageAssetId,
		profile.Status,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrProfileAlreadyExists
		}
		r.log.Error("failed to create player profile", "err", err, "user_id", profile.UserId)
		return fmt.Errorf("failed to insert player profile: %w", err)
	}

	return nil
}

func (r *ProfileRepository) FindByUserId(ctx context.Context, userid string) (*ProfileRecord, error) {
	query := `SELECT id, user_id, display_name, bio, region, profile_image_asset_id, status, created_at, updated_at, deleted_at FROM player_profiles WHERE user_id = $1`

	var profile ProfileRecord
	err := r.db.QueryRow(ctx, query, userid).Scan(
		&profile.Id, &profile.UserId, &profile.DisplayName,
		&profile.Bio, &profile.Region, &profile.ProfileImageAssetId,
		&profile.Status, &profile.CreatedAt, &profile.UpdatedAt, &profile.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		r.log.Error("failed to query profile by user_id", "err", err, "user_id", userid)
		return nil, fmt.Errorf("database query failed: %w", err)
	}

	return &profile, nil
}

func (r *ProfileRepository) FindById(ctx context.Context, id string) (*ProfileRecord, error) {
	query := `SELECT id, user_id, display_name, bio, region, profile_image_asset_id, status, created_at, updated_at, deleted_at FROM player_profiles WHERE id = $1`

	var profile ProfileRecord
	err := r.db.QueryRow(ctx, query, id).Scan(
		&profile.Id, &profile.UserId, &profile.DisplayName,
		&profile.Bio, &profile.Region, &profile.ProfileImageAssetId,
		&profile.Status, &profile.CreatedAt, &profile.UpdatedAt, &profile.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		r.log.Error("failed to query profile by id", "err", err, "id", id)
		return nil, fmt.Errorf("database query failed: %w", err)
	}

	return &profile, nil
}

func (r *ProfileRepository) Update(ctx context.Context, userId string, params *UpdateProfileParams) error {
	query := "UPDATE player_profiles SET updated_at = NOW()"
	var args []interface{}
	argIdx := 1

	if params.DisplayName != "" {
		query += fmt.Sprintf(", display_name = $%d", argIdx)
		args = append(args, params.DisplayName)
		argIdx++
	}
	if params.Bio != "" {
		query += fmt.Sprintf(", bio = $%d", argIdx)
		args = append(args, params.Bio)
		argIdx++
	}
	if params.Region != "" {
		query += fmt.Sprintf(", region = $%d", argIdx)
		args = append(args, params.Region)
		argIdx++
	}
	if params.ProfileImageAssetId != nil {
		query += fmt.Sprintf(", profile_image_asset_id = $%d", argIdx)
		args = append(args, params.ProfileImageAssetId)
		argIdx++
	}
	if params.Status != "" {
		query += fmt.Sprintf(", status = $%d", argIdx)
		args = append(args, params.Status)
		argIdx++
	}

	if len(args) == 0 {
		return nil
	}

	query += fmt.Sprintf(" WHERE user_id = $%d", argIdx)
	args = append(args, userId)

	tag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		r.log.Error("failed to update player profile", "err", err, "user_id", userId)
		return fmt.Errorf("failed to update player profile: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrProfileNotFound
	}

	return nil
}
