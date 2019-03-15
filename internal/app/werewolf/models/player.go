package models

import (
	"github.com/google/uuid"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/characters"
)

type PlayerLiveness string

const (
	ALIVE PlayerLiveness = "ALIVE"
	DEAD  PlayerLiveness = "DEAD"
)

func NewPlayer(userName, userID string) Player {
	return Player{
		ID:       uuid.New().String(),
		UserName: userName,
		UserID:   userID,
	}
}

type Player struct {
	ID               string //canonical ID
	UserName, UserID string //platform Information specific
	Platform         string

	Character characters.Interface
}
