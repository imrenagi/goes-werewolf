package states

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/characters"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/models"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/services"
)

func NewInitialState(waitDuration time.Duration) *Initial {
	return &Initial{
		WaitDuration: waitDuration,
	}
}

type Initial struct {
	State
	PlayerDAO    services.PlayerDAO
	WaitDuration time.Duration
}

// Execute assign all registered players to a character. If the number of players is not enough, game state changes
// to Cancel state. If internal error happened, state goes to Error state
func (i *Initial) Execute(c Context) {

	log := logrus.WithFields(map[string]interface{}{
		"gameID": c.GameID(),
		"day":    c.GetDay(),
		"state":  i.String(),
	})

	log.Info("Starting initial state")

	//TODO Notify game has been started. Waiting for all other players

	<-time.After(i.WaitDuration)

	players, err := i.PlayerDAO.GetPlayers(c.GameID())
	if err != nil {
		log.WithError(err).Error("Can't retrieve all game players.")
		c.SetState(&Error{Error: err})
		return
	}

	allocator := CharacterAllocator{}
	players, err = allocator.Assign(players)
	if err == ErrNotEnoughPlayer {
		c.SetState(&Cancel{})
		return
	}

	if err != nil {
		log.WithError(err).Error("Unidentified error from character allocation")
		c.SetState(&Error{Error: err})
		return
	}

	for _, p := range players {
		err := i.PlayerDAO.SavePlayer(c.GameID(), p)
		if err != nil {
			log.WithError(err).Error("Can't save game player data.")
			c.SetState(&Error{Error: err})
			return
		}

		//TODO Publish the role to users
	}

	c.SetState(&NightTime{})

	log.Info("Initial state complete.")
}

func (i *Initial) String() string {
	return "initial"
}

const (
	MIN_PLAYER                 = 7
	MIN_PLAYER_ADDITIONAL_WOLF = 16
	MIN_WOLVES                 = 2
	MAX_WOLVES                 = 3
	MAX_SEER                   = 1
	MAX_DOCTOR                 = 1
)

var (
	ErrNotEnoughPlayer = fmt.Errorf("The number of players are not enough")
)

type CharacterAllocator struct {
	numWolves, numVillagers, numDoctor, numSeer int
	totalAssigned                               int
}

//Assign assigns character to every given players and returns the same players, but with it's new character
func (ca *CharacterAllocator) Assign(players []models.Player) ([]models.Player, error) {

	l := len(players)

	if len(players) < MIN_PLAYER {
		return nil, ErrNotEnoughPlayer
	}

	result := make([]models.Player, 0)

	nWolf := MIN_WOLVES
	if l > MIN_PLAYER_ADDITIONAL_WOLF {
		nWolf = MAX_WOLVES
	}

	for i := 0; i < nWolf; i++ {
		idx, wolf := ca.assignRandomPlayer(players, characters.NewWolf())
		result = append(result, wolf)
		players = append(players[:idx], players[idx+1:]...)
		ca.numWolves++
	}

	idx, doctor := ca.assignRandomPlayer(players, characters.NewDoctor())
	result = append(result, doctor)
	players = append(players[:idx], players[idx+1:]...)
	ca.numDoctor++

	idx, seer := ca.assignRandomPlayer(players, characters.NewSeer())
	result = append(result, seer)
	players = append(players[:idx], players[idx+1:]...)
	ca.numSeer++

	for i := 0; i < l-(ca.numWolves+ca.numSeer+ca.numDoctor); i++ {
		idx, villager := ca.assignRandomPlayer(players, characters.NewVillager())
		result = append(result, villager)
		players = append(players[:idx], players[idx+1:]...)
		ca.numVillagers++
	}

	return result, nil
}

func (ca CharacterAllocator) assignRandomPlayer(p []models.Player, character characters.Interface) (int, models.Player) {
	rand.Seed(time.Now().UTC().UnixNano())
	idx := rand.Intn(len(p))
	player := p[idx]
	player.Character = character

	return idx, player
}
