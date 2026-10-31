package gameprofile

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

type IGameProfileRepo interface {
	Create(ctx context.Context, profile *GameProfileRecord) error
	FindByID(ctx context.Context, id string) (*GameProfileRecord, error)
	FindByPlayerID(ctx context.Context, playerProfileId string) ([]GameProfileRecord, error)
	Update(ctx context.Context, profile *GameProfileRecord) error
	Delete(ctx context.Context, id string) error
}

type GameProfileRepo struct {
	db  exec.DBExecutor
	log *slog.Logger
}

func NewGameProfileRepo(db *pgxpool.Pool, log *slog.Logger) *GameProfileRepo {
	return &GameProfileRepo{
		db:  db,
		log: log.With("component", "gameprofile.repository"),
	}
}

func (r *GameProfileRepo) Create(ctx context.Context, profile *GameProfileRecord) error {
	query := `INSERT INTO player_game_profiles (player_profile_id, game_id, in_game_name, in_game_id, server_id, competitive_role, rank)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		ctx,
		query,
		profile.PlayerProfileId,
		profile.GameId,
		profile.InGameName,
		profile.InGameId,
		profile.ServerId,
		profile.CompetitiveRole,
		profile.Rank,
	).Scan(&profile.Id, &profile.CreatedAt, &profile.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrGameProfileAlreadyExists
		}
		r.log.Error("failed to create player game profile", "err", err, "player_profile_id", profile.PlayerProfileId, "game_id", profile.GameId)
		return fmt.Errorf("failed to create game profile: %w", err)
	}

	return nil
}

func (r *GameProfileRepo) FindByID(ctx context.Context, id string) (*GameProfileRecord, error) {
	query := `SELECT id, player_profile_id, game_id, in_game_name, in_game_id, server_id, competitive_role, rank, created_at, updated_at 
	FROM player_game_profiles WHERE id = $1`

	var gameProfile GameProfileRecord

	err := r.db.QueryRow(ctx, query, id).Scan(
		&gameProfile.Id,
		&gameProfile.PlayerProfileId,
		&gameProfile.GameId,
		&gameProfile.InGameName,
		&gameProfile.InGameId,
		&gameProfile.ServerId,
		&gameProfile.CompetitiveRole,
		&gameProfile.Rank,
		&gameProfile.CreatedAt,
		&gameProfile.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrGameProfileNotFound
		}
		r.log.Error("failed to find game profile by id", "err", err, "id", id)
		return nil, fmt.Errorf("failed to find game profile: %w", err)
	}

	return &gameProfile, nil
}

func (r *GameProfileRepo) FindByPlayerID(ctx context.Context, playerProfileId string) ([]GameProfileRecord, error) {
	query := `SELECT id, player_profile_id, game_id, in_game_name, in_game_id, server_id, competitive_role, rank, created_at, updated_at 
	FROM player_game_profiles WHERE player_profile_id = $1`

	rows, err := r.db.Query(ctx, query, playerProfileId)
	if err != nil {
		r.log.Error("failed to query game profiles by player id", "err", err, "player_profile_id", playerProfileId)
		return nil, fmt.Errorf("failed to query game profiles: %w", err)
	}
	defer rows.Close()

	var gameProfiles []GameProfileRecord

	for rows.Next() {
		var gprofile GameProfileRecord
		err := rows.Scan(
			&gprofile.Id,
			&gprofile.PlayerProfileId,
			&gprofile.GameId,
			&gprofile.InGameName,
			&gprofile.InGameId,
			&gprofile.ServerId,
			&gprofile.CompetitiveRole,
			&gprofile.Rank,
			&gprofile.CreatedAt,
			&gprofile.UpdatedAt,
		)
		if err != nil {
			r.log.Error("failed to scan game profile row", "err", err, "player_profile_id", playerProfileId)
			return nil, fmt.Errorf("failed to scan game profile: %w", err)
		}
		gameProfiles = append(gameProfiles, gprofile)
	}

	return gameProfiles, nil
}

func (r *GameProfileRepo) Update(ctx context.Context, profile *GameProfileRecord) error {
	query := `UPDATE player_game_profiles 
	SET in_game_name = $1, in_game_id = $2, server_id = $3, competitive_role = $4, rank = $5, updated_at = NOW() 
	WHERE id = $6`

	tag, err := r.db.Exec(
		ctx,
		query,
		profile.InGameName,
		profile.InGameId,
		profile.ServerId,
		profile.CompetitiveRole,
		profile.Rank,
		profile.Id,
	)

	if err != nil {
		r.log.Error("failed to update game profile", "err", err, "id", profile.Id)
		return fmt.Errorf("failed to update game profile: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrGameProfileNotFound
	}

	return nil
}

func (r *GameProfileRepo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM player_game_profiles WHERE id = $1`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.log.Error("failed to delete game profile", "err", err, "id", id)
		return fmt.Errorf("failed to delete game profile: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrGameProfileNotFound
	}

	return nil
}
