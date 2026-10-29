package games

import "errors"

var (
	ErrGameNotFound    = errors.New("game not found")
	ErrInvalidGameID   = errors.New("invalid game identifier format")
	ErrGameInactive    = errors.New("game is currently inactive or unavailable")
	ErrInvalidGameData = errors.New("provided game data is invalid or incomplete")
)
