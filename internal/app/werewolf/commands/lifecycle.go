package commands

import "github.com/imrenagi/goes-werewolf/internal/app/werewolf/models"

type InitializeGame struct {
	Platform  string
	ChannelID string
	Initiator models.Player
}

type CancelGame struct {
}

type StartGame struct {
}

type EndGame struct {
}

type Advance struct {
}

type JoinGame struct {
}
