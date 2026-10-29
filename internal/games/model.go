package games

import "time"

type GameRecord struct {
	Id        string
	Code      string
	Name      string
	IsActive  bool
	CreatedAt time.Time
}
