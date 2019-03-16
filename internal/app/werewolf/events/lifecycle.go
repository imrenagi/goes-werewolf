package events

import "github.com/imrenagi/goes-werewolf/internal/app/werewolf/models"

type GameInitialized struct {
	GameID string
	CreatorID, CreatorName string
	ChannelID, ChannelName string
	Platform string
}

type GameCanceled struct {
}

type GameStarted struct {
}

type GameEnded struct {
}

type PlayerJoined struct {
	Player models.Player
	GameID string
}

type DayStarted struct {
}

type NightStarted struct {
}
