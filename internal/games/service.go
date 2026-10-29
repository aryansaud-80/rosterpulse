package games

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type IGameService interface {
	GetGames(ctx context.Context) ([]GameResponse, error)
	GetGameById(ctx context.Context, id string) (*GameResponse, error)
}

type GameService struct {
	repo IGames
	log  *slog.Logger
}

func NewGameService(repo IGames, log *slog.Logger) *GameService {
	return &GameService{
		repo: repo,
		log:  log.With("component", "games.service"),
	}
}

func (s *GameService) GetGames(ctx context.Context) ([]GameResponse, error) {
	games, err := s.repo.FindAll(ctx)
	if err != nil {
		s.log.Error("failed to fetch all games from repository", "error", err)
		return nil, fmt.Errorf("failed to get games: %w", err)
	}

	rgames := make([]GameResponse, 0, len(games))

	for _, game := range games {
		rgames = append(rgames, GameResponse{
			Id:        game.Id,
			Code:      game.Code,
			Name:      game.Name,
			IsActive:  game.IsActive,
			CreatedAt: game.CreatedAt,
		})
	}

	return rgames, nil
}

func (s *GameService) GetGameById(ctx context.Context, id string) (*GameResponse, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidGameID
	}

	game, err := s.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, ErrGameNotFound) {
			return nil, ErrGameNotFound
		}
		s.log.Error("failed to fetch game by id from repository", "game_id", id, "error", err)
		return nil, fmt.Errorf("failed to get game by id: %w", err)
	}

	rgame := &GameResponse{
		Id:        game.Id,
		Code:      game.Code,
		Name:      game.Name,
		IsActive:  game.IsActive,
		CreatedAt: game.CreatedAt,
	}

	return rgame, nil
}
