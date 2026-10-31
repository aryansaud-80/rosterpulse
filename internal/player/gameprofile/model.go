package gameprofile

import "time"

type GameProfileRecord struct {
	Id              string
	PlayerProfileId string
	GameId          string
	InGameName      string
	InGameId        string
	ServerId        *string
	CompetitiveRole *string
	Rank            *string
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}
