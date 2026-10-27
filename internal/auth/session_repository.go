package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ISessionRepo interface {
	CreateSession(ctx context.Context, session *Session) (*Session, error)
	GetSessionByTokenHash(ctx context.Context, tokenHash string) (*SessionRecord, error)
}

type SessionRepository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewSessionRepository(db *pgxpool.Pool, log *slog.Logger) *SessionRepository {
	return &SessionRepository{
		db:  db,
		log: log.With("component", "auth.session_repository"),
	}
}

func (r *SessionRepository) CreateSession(ctx context.Context, session *Session) (*Session, error) {
	query := `INSERT INTO sessions (user_id, refresh_token_hash, token_version, expires_at)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

	r.log.Info("Executing session create")

	err := r.db.QueryRow(
		ctx,
		query,
		session.UserId,
		session.RefreshTokenHash,
		session.TokenVersion,
		session.ExpiresAt,
	).Scan(
		&session.Id,
	)

	if err != nil {
		r.log.Error("failed to create session", "error", err)
		return nil, fmt.Errorf("failed to execute session create: %w", err)
	}

	return session, nil
}

func (r *SessionRepository) GetSessionByTokenHash(ctx context.Context, tokenHash string) (*SessionRecord, error) {
	query := `SELECT id, user_id, refresh_token_hash, token_version, expires_at, revoked_at, created_at, last_used_at
	FROM sessions WHERE refresh_token_hash=$1`

	var record SessionRecord

	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&record.ID,
		&record.UserId,
		&record.RefreshTokenHash,
		&record.TokenVersion,
		&record.ExpiresAt,
		&record.RevokedAt,
		&record.CreatedAt,
		&record.LastUsedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
		r.log.Error("failed to query session by hash", "error", err)
		return nil, fmt.Errorf("database query failed: %w", err)
	}

	return &record, nil
}
