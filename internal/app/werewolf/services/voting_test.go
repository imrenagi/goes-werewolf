package services_test

import (
	"testing"

	"sort"

	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/mocks"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/models"
	. "github.com/imrenagi/goes-werewolf/internal/app/werewolf/services"
)

func TestPollList_Sort(t *testing.T) {

	polls := []PlayerPollCount{
		{Player: models.Player{ID: "1"}, Count: 5},
		{Player: models.Player{ID: "2"}, Count: 2},
		{Player: models.Player{ID: "3"}, Count: 9},
	}

	var countList PollList = polls

	sort.Sort(countList)

	require.Len(t, countList, len(polls))
	require.Equal(t, "3", countList[0].Player.ID)
	require.Equal(t, "1", countList[1].Player.ID)
	require.Equal(t, "2", countList[2].Player.ID)
}

func TestPollingService_GetWolvesVictim(t *testing.T) {

	t.Run("Should get a selected victim", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		polls := []models.Poll{
			{PollID: "1", Choice: models.Player{ID: "a"}},
			{PollID: "1", Choice: models.Player{ID: "a"}},
			{PollID: "1", Choice: models.Player{ID: "a"}},
			{PollID: "1", Choice: models.Player{ID: "b"}},
			{PollID: "1", Choice: models.Player{ID: "b"}},
			{PollID: "1", Choice: models.Player{ID: "c"}},
		}

		pollDAO := mocks.NewMockPollDAO(mockCtrl)
		pollDAO.EXPECT().GetPolls("gameID", "1").Return(polls, nil).Times(1)

		pollingService := PollingService{
			PollDAO: pollDAO,
		}

		player, _ := pollingService.GetMostVotedPlayer("gameID", "1")

		require.Equal(t, "a", player.ID)
	})

	t.Run("No votes given, return err wolf no votes", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		polls := []models.Poll{}

		pollDAO := mocks.NewMockPollDAO(mockCtrl)
		pollDAO.EXPECT().GetPolls("gameID", "1").Return(polls, nil).Times(1)

		pollingService := PollingService{
			PollDAO: pollDAO,
		}

		player, err := pollingService.GetMostVotedPlayer("gameID", "1")

		require.Nil(t, player)
		require.Equal(t, ErrNoVotes, err)
	})

	t.Run("Tie votes, return err no aggreement", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		polls := []models.Poll{
			{PollID: "1", Choice: models.Player{ID: "a"}},
			{PollID: "1", Choice: models.Player{ID: "a"}},
			{PollID: "1", Choice: models.Player{ID: "b"}},
			{PollID: "1", Choice: models.Player{ID: "b"}},
		}

		pollDAO := mocks.NewMockPollDAO(mockCtrl)
		pollDAO.EXPECT().GetPolls("gameID", "1").Return(polls, nil).Times(1)

		pollingService := PollingService{
			PollDAO: pollDAO,
		}

		player, err := pollingService.GetMostVotedPlayer("gameID", "1")

		require.Nil(t, player)
		require.Equal(t, ErrNoAgreement, err)
	})

	t.Run("Wolves only vote for 1 player", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		polls := []models.Poll{
			{PollID: "1", Choice: models.Player{ID: "a"}},
			{PollID: "1", Choice: models.Player{ID: "a"}},
			{PollID: "1", Choice: models.Player{ID: "a"}},
			{PollID: "1", Choice: models.Player{ID: "a"}},
		}

		pollDAO := mocks.NewMockPollDAO(mockCtrl)
		pollDAO.EXPECT().GetPolls("gameID", "1").Return(polls, nil).Times(1)

		pollingService := PollingService{
			PollDAO: pollDAO,
		}

		player, err := pollingService.GetMostVotedPlayer("gameID", "1")

		require.Nil(t, err)
		require.Equal(t, "a", player.ID)
	})

	t.Run("Get wolves votes error, return its error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		pollDAO := mocks.NewMockPollDAO(mockCtrl)
		pollDAO.EXPECT().GetPolls("gameID", "1").Return(nil, fmt.Errorf("any error")).Times(1)

		pollingService := PollingService{
			PollDAO: pollDAO,
		}

		player, err := pollingService.GetMostVotedPlayer("gameID", "1")

		require.Nil(t, player)
		require.Equal(t, fmt.Errorf("any error"), err)
	})

}
