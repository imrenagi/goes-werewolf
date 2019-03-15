package events

import (
	"github.com/imrenagi/goes-werewolf/internal/app/events"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/models"
)

type WolvesVoteStarted struct {
	events.DomainEvent
	GameID    string
	PollID    string
	GameDay   int
	GameState string
	SendTo    models.Player
}

type DoctorTargetRequested struct {
	events.DomainEvent
	GameID    string
	PollID    string
	GameDay   int
	GameState string
	SendTo    models.Player
}

type SeerTargetRequested struct {
	events.DomainEvent
	GameID    string
	PollID    string
	GameDay   int
	GameState string
	SendTo    models.Player
}
