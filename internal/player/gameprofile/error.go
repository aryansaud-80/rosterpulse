package gameprofile

import "errors"

var (
	ErrGameProfileNotFound = errors.New("game profile not found")

	ErrGameProfileAlreadyExists = errors.New("game profile already exists for this game")

	ErrGameInvalidOrInactive = errors.New("selected game does not exist or is inactive")

	ErrForbidden = errors.New("unauthorized access to game profile")
)
