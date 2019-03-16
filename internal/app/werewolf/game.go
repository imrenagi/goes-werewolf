package werewolf

import (
	"time"
	"github.com/google/uuid"

	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/events"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/models"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/states"
	"github.com/imrenagi/goes-werewolf/pkg/eventsource/sourcing"
)

type GameDataService interface {
	Save(game *Game) error
	Update(game *Game) error
}

const firstDay = 1

type Game struct {
	sourcing.EventSource
	states.Context

	ID      string
	State   states.State
	Players []models.Player
	day     int

	DS GameDataService
}

func (game *Game) Start() {
	game.Execute()
}

func (g *Game) SetState(s states.State) {
	g.State = s
}

func (g Game) GetDay() int {
	return g.day
}

func (g *Game) Execute() {
	g.State.Execute(g)
}

func (g Game) GameID() string {
	return g.ID
}

func NewGameFromHistory(sourceId sourcing.EventSourceId, history []sourcing.Event) (*Game, error) {
	game := new(Game)
	game.EventSource = sourcing.CreateFromHistory(game, sourceId, history)
	return game, nil
}

func NewGame(channelID, channelName, creatorID, creatorName string, platform string) (*Game, error) {
	game := new(Game)
	game.EventSource = sourcing.CreateNew(game)

	gameID := uuid.New().String()
	game.Apply(events.GameInitialized{
		GameID:      gameID,
		ChannelID:   channelID,
		ChannelName: channelName,
		CreatorID:   creatorID,
		CreatorName: creatorName,
		Platform:    platform,
	})

	return game, nil
}

func (g *Game) HandleGameInitialized(ev events.GameInitialized) {
	g.ID = ev.GameID
	g.day = firstDay
	g.State = states.NewInitialState(10 * time.Millisecond)
	firstPlayer := models.NewPlayer(ev.CreatorName, ev.CreatorID)
	g.Players = []models.Player{firstPlayer}
}

func (g *Game) HandlePlayerJoined(ev events.PlayerJoined) {
	g.Players = append(g.Players, ev.Player)
}
