package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Iuser interface {
	Create(ctx context.Context, record *UserRecord) (*UserRecord, error)
	FindByEmail(ctx context.Context, normalizedEmail string) (*UserRecord, error)
}

type UserRepository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewUserRepository(db *pgxpool.Pool, log *slog.Logger) *UserRepository {
	userRepo := &UserRepository{
		db:  db,
		log: log.With("component", "user.repository"),
	}
	return userRepo
}

func (r *UserRepository) Create(ctx context.Context, record *UserRecord) (*UserRecord, error) {
	query := `
		INSERT INTO users (email, normalized_email, password_hash, role, status) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, created_at, updated_at
	`

	r.log.Info("executing user insert", "email", record.Email)

	err := r.db.QueryRow(
		ctx,
		query,
		record.Email,
		record.NormalizedEmail,
		record.PasswordHash,
		record.Role,
		record.Status,
	).Scan(&record.ID, &record.CreatedAt, &record.UpdatedAt)

	if err != nil {
		r.log.Error("failed to insert user", "error", err, "email", record.Email)
		return nil, fmt.Errorf("failed to execute user insert: %w", err)
	}

	return record, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, normalizedEmail string) (*UserRecord, error) {
	query := `SELECT id, email, normalized_email, password_hash, role, status, created_at, updated_at, deleted_at FROM users WHERE normalized_email = $1 AND status = 'ACTIVE'`

	var record UserRecord

	err := r.db.QueryRow(ctx, query, normalizedEmail).Scan(
		&record.ID,
		&record.Email,
		&record.NormalizedEmail,
		&record.PasswordHash,
		&record.Role,
		&record.Status,
		&record.CreatedAt,
		&record.UpdatedAt,
		&record.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		r.log.Error("failed to query user by email", "error", err)
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	return &record, nil
}
