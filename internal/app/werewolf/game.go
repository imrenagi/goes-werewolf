package werewolf

import (
	"time"

	"github.com/google/uuid"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/events"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/models"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/states"
	"github.com/imrenagi/goes-werewolf/pkg/eventbus"
)

type GameDataService interface {
	Save(game *Game) error
	Update(game *Game) error
}

type Game struct {
	ID string
	states.Context
	State   states.State
	Players []models.Player
	Day     int

	DS GameDataService
}

func (game *Game) Start() {
	game.Execute()
}

func (g *Game) SetState(s states.State) {
	g.State = s
}

func (g Game) GetDay() int {
	return g.Day
}

func (g *Game) Execute() {
	g.State.Execute(g)
}

func (g Game) GameID() string {
	return g.ID
}

func NewGame() *Game {
	return &Game{
		State: states.NewInitialState(10 * time.Millisecond),
	}
}

func (g *Game) ApplyInitializeGame(ev events.GameInitialized) {

	g.ID = uuid.New().String()
	g.Day = 1
	g.Players := []models.Player{ev.}

	eventbus.Publish(events.GameInitialized{

	})

	eventbus.Publish(events.PlayerJoined{

	})


	// g.DS = what need to be assigned to this?
}
