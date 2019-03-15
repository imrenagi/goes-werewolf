package states_test

import (
	"fmt"
	"testing"

	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/mocks"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/models"
	. "github.com/imrenagi/goes-werewolf/internal/app/werewolf/states"
)

func TestCharacterAssigner_Assign(t *testing.T) {
	t.Run("player < 7 return error", func(t *testing.T) {
		assigner := CharacterAllocator{}
		_, err := assigner.Assign(make([]models.Player, 5))
		require.Equal(t, fmt.Errorf("The number of players are not enough"), err)
	})

	t.Run("Should has required characters", func(t *testing.T) {

		tests := []struct {
			totalPlayer, expNumVillagers           int
			expNumWolves, expNumSeer, expNumDoctor int
		}{
			{totalPlayer: 7, expNumVillagers: 3, expNumWolves: MIN_WOLVES, expNumSeer: MAX_SEER, expNumDoctor: MAX_DOCTOR},
			{totalPlayer: 10, expNumVillagers: 6, expNumWolves: MIN_WOLVES, expNumSeer: MAX_SEER, expNumDoctor: MAX_DOCTOR},
			{totalPlayer: 12, expNumVillagers: 8, expNumWolves: MIN_WOLVES, expNumSeer: MAX_SEER, expNumDoctor: MAX_DOCTOR},
			{totalPlayer: 15, expNumVillagers: 12, expNumWolves: MIN_WOLVES, expNumSeer: MAX_SEER, expNumDoctor: MAX_DOCTOR},
			{totalPlayer: 16, expNumVillagers: 13, expNumWolves: MIN_WOLVES, expNumSeer: MAX_SEER, expNumDoctor: MAX_DOCTOR},

			{totalPlayer: 17, expNumVillagers: 12, expNumWolves: MAX_WOLVES, expNumSeer: MAX_SEER, expNumDoctor: MAX_DOCTOR},
			{totalPlayer: 20, expNumVillagers: 15, expNumWolves: MAX_WOLVES, expNumSeer: MAX_SEER, expNumDoctor: MAX_DOCTOR},
		}

		for _, test := range tests {
			t.Run(fmt.Sprintf("assigning %d characters", test.totalPlayer), func(t *testing.T) {
				assigner := CharacterAllocator{}
				inPlayer := make([]models.Player, test.totalPlayer)

				assignedPlayers, _ := assigner.Assign(inPlayer)
				require.Len(t, assignedPlayers, test.totalPlayer)

				var wolves, seer, doctor, villagers int

				for i := 0; i < test.totalPlayer; i++ {
					switch assignedPlayers[i].Character.String() {
					case "Wolf":
						wolves++
					case "Seer":
						seer++
					case "Doctor":
						doctor++
					case "Villager":
						villagers++
					}
				}

				require.Equal(t, test.expNumWolves, wolves)
				require.Equal(t, test.expNumSeer, seer)
				require.Equal(t, test.expNumDoctor, doctor)
				require.Equal(t, test.totalPlayer-(test.expNumWolves+test.expNumDoctor+test.expNumSeer), villagers)
			})
		}
	})
}

func TestInitial_String(t *testing.T) {
	initialState := Initial{
		WaitDuration: 1 * time.Millisecond,
	}

	require.Equal(t, "initial", initialState.String())
}

func TestInitial_Execute(t *testing.T) {
	initialState := Initial{
		WaitDuration: 1 * time.Millisecond,
	}

	players := []models.Player{
		models.NewPlayer("a", "a"),
		models.NewPlayer("b", "b"),
		models.NewPlayer("c", "c"),
		models.NewPlayer("d", "d"),
		models.NewPlayer("e", "e"),
		models.NewPlayer("f", "f"),
		models.NewPlayer("g", "g"),
	}

	t.Run("Smoothly execute initial state", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		context := mocks.NewMockContext(mockCtrl)
		context.EXPECT().GameID().Return("anID").AnyTimes()
		context.EXPECT().GetDay().Return(0).AnyTimes()
		context.EXPECT().SetState(&NightTime{}).Times(1)

		playerDAO := mocks.NewMockPlayerDAO(mockCtrl)
		playerDAO.EXPECT().GetPlayers("anID").Times(1).Return(players, nil)
		playerDAO.EXPECT().SavePlayer("anID", gomock.Any()).Times(len(players))
		initialState.PlayerDAO = playerDAO

		initialState.Execute(context)
	})

	t.Run("Cancelling game if the number of player is not enough", func(t *testing.T) {
		players := []models.Player{
			models.NewPlayer("a", "a"),
			models.NewPlayer("b", "b"),
			models.NewPlayer("c", "c"),
		}

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		context := mocks.NewMockContext(mockCtrl)
		context.EXPECT().GameID().Return("anID").AnyTimes()
		context.EXPECT().GetDay().Return(0).AnyTimes()
		context.EXPECT().SetState(&Cancel{}).Times(1)

		playerDAO := mocks.NewMockPlayerDAO(mockCtrl)
		playerDAO.EXPECT().GetPlayers("anID").Times(1).Return(players, nil)
		playerDAO.EXPECT().SavePlayer("anID", gomock.Any()).Times(0)
		initialState.PlayerDAO = playerDAO

		initialState.Execute(context)
	})

	t.Run("Set state to error when error in getting all players", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		context := mocks.NewMockContext(mockCtrl)
		context.EXPECT().GameID().Return("anID").AnyTimes()
		context.EXPECT().GetDay().Return(0).AnyTimes()
		context.EXPECT().SetState(gomock.AssignableToTypeOf(&Error{})).Times(1)

		playerDAO := mocks.NewMockPlayerDAO(mockCtrl)
		playerDAO.EXPECT().GetPlayers("anID").Times(1).Return(nil, fmt.Errorf("any error"))
		playerDAO.EXPECT().SavePlayer("anID", gomock.Any()).Times(0)
		initialState.PlayerDAO = playerDAO

		initialState.Execute(context)
	})

	t.Run("Set state to error when error in save the player", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		context := mocks.NewMockContext(mockCtrl)
		context.EXPECT().GameID().Return("anID").AnyTimes()
		context.EXPECT().GetDay().Return(0).AnyTimes()
		context.EXPECT().SetState(gomock.AssignableToTypeOf(&Error{})).Times(1)

		playerDAO := mocks.NewMockPlayerDAO(mockCtrl)
		playerDAO.EXPECT().GetPlayers("anID").Times(1).Return(players, nil)
		playerDAO.EXPECT().SavePlayer("anID", gomock.Any()).Return(fmt.Errorf("database error")).Times(1)
		initialState.PlayerDAO = playerDAO

		initialState.Execute(context)
	})

}
