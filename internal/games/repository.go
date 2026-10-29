package games

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/aryansaud-80/rosterpulse/internal/platform/exec"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IGames interface {
	FindAll(ctx context.Context) ([]GameRecord, error)
	FindById(ctx context.Context, id string) (*GameRecord, error)
}

type GameRepository struct {
	db  exec.DBExecutor
	log *slog.Logger
}

func NewGameRepository(db *pgxpool.Pool, log *slog.Logger) *GameRepository {
	return &GameRepository{
		db:  db,
		log: log.With("component", "games.repository"),
	}
}

func (r *GameRepository) FindAll(ctx context.Context) ([]GameRecord, error) {
	var games []GameRecord
	query := `SELECT id, code, name, is_active, created_at FROM games WHERE is_active = $1`

	rows, err := r.db.Query(ctx, query, true)
	if err != nil {
		r.log.Error("failed to query all games", "error", err)
		return nil, fmt.Errorf("failed to query games: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var g GameRecord
		if err := rows.Scan(&g.Id, &g.Code, &g.Name, &g.IsActive, &g.CreatedAt); err != nil {
			r.log.Error("failed to scan game record", "error", err)
			return nil, fmt.Errorf("failed to scan game record: %w", err)
		}
		games = append(games, g)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("error iterating game rows", "error", err)
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return games, nil
}

func (r *GameRepository) FindById(ctx context.Context, id string) (*GameRecord, error) {
	var g GameRecord
	query := `SELECT id, code, name, is_active, created_at FROM games WHERE id = $1 AND is_active = $2`

	err := r.db.QueryRow(ctx, query, id, true).Scan(&g.Id, &g.Code, &g.Name, &g.IsActive, &g.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrGameNotFound
		}
		r.log.Error("failed to query game by id", "game_id", id, "error", err)
		return nil, fmt.Errorf("failed to get game by id: %w", err)
	}

	return &g, nil
}
