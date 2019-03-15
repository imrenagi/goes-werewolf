package services

import (
	"sort"

	"fmt"

	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/models"
)

type PollDAO interface {
	GetPolls(gameID, pollID string) ([]models.Poll, error)
}

var (
	ErrNoAgreement = fmt.Errorf("No agreement reached")
	ErrNoVotes     = fmt.Errorf("No votes given")
)

type PollingService struct {
	PollDAO PollDAO
}

func (ps PollingService) NewPollID(gameID, state string, day int, character string) string {
	return fmt.Sprintf("%s#%s#%d#%s", gameID, state, day, character)
}

func (ps PollingService) GetMostVotedPlayer(gameID, pollID string) (*models.Player, error) {

	polls, err := ps.PollDAO.GetPolls(gameID, pollID)
	if err != nil {
		return nil, err
	}

	if len(polls) == 0 {
		return nil, ErrNoVotes
	}

	playerMap := make(map[string]PlayerPollCount)

	for _, poll := range polls {
		chosenPlayerID := poll.Choice.ID
		if val, ok := playerMap[chosenPlayerID]; !ok {
			playerMap[chosenPlayerID] = PlayerPollCount{
				Player: poll.Choice,
				Count:  1,
			}
		} else {
			val.Count++
			playerMap[chosenPlayerID] = val
		}
	}

	var arr PollList = make([]PlayerPollCount, 0)
	for _, v := range playerMap {
		arr = append(arr, v)
	}
	sort.Sort(arr)

	if len(arr) == 1 {
		return &arr[0].Player, nil
	} else if len(arr) > 1 {
		if arr[0].Count != arr[1].Count {
			return &arr[0].Player, nil
		} else {
			return nil, ErrNoAgreement
		}
	}

	return nil, nil
}

type PlayerPollCount struct {
	Player models.Player
	Count  int
}

type PollList []PlayerPollCount

func (p PollList) Len() int           { return len(p) }
func (p PollList) Less(i, j int) bool { return p[i].Count > p[j].Count }
func (p PollList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
