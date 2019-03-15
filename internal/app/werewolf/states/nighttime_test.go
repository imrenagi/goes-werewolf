package states_test

import (
	"testing"

	"fmt"

	"time"

	"github.com/golang/mock/gomock"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/mocks"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/models"

	"github.com/stretchr/testify/require"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/characters"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/services"
	. "github.com/imrenagi/goes-werewolf/internal/app/werewolf/states"
)

//func TestNightTime_Execute(t *testing.T) {
//
//	t.Run("Should has selected victim", func(t *testing.T) {
//
//	})
//
//}

func TestNightTime_TriggerCharacterPolls(t *testing.T) {

	t.Run("Get wolves returns error, return error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		context := mocks.NewMockContext(mockCtrl)
		context.EXPECT().GameID().Return("anID").AnyTimes()
		context.EXPECT().GetDay().Return(0).AnyTimes()

		playerDAO := mocks.NewMockPlayerDAO(mockCtrl)
		playerDAO.EXPECT().GetWolves(context.GameID()).Return(nil, fmt.Errorf("Any error")).Times(1)

		night := NightTime{
			PlayerDAO:        playerDAO,
			WolfVoteDuration: 1 * time.Millisecond,
		}
		err := night.TriggerCharacterPolls(context)

		require.Equal(t, fmt.Errorf("Any error"), err)
	})

	t.Run("Get doctor returns error, return error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		context := mocks.NewMockContext(mockCtrl)
		context.EXPECT().GameID().Return("anID").AnyTimes()
		context.EXPECT().GetDay().Return(0).AnyTimes()

		pollingService := mocks.NewMockPollingService(mockCtrl)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "wolf").Return("wolfpollid").Times(1)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "doctor").Return("doctorpollid").Times(0)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "seer").Return("seerpollid").Times(0)

		playerDAO := mocks.NewMockPlayerDAO(mockCtrl)
		playerDAO.EXPECT().GetWolves(context.GameID()).Return([]models.Player{}, nil).Times(1)
		playerDAO.EXPECT().GetDoctor(context.GameID()).Return(nil, fmt.Errorf("Any error")).Times(1)
		playerDAO.EXPECT().GetSeer(context.GameID()).Times(0)

		night := NightTime{
			PlayerDAO:        playerDAO,
			PollingService:   pollingService,
			WolfVoteDuration: 1 * time.Millisecond,
		}
		err := night.TriggerCharacterPolls(context)

		require.Equal(t, fmt.Errorf("Any error"), err)
	})

	t.Run("Get seer returns error, return error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		context := mocks.NewMockContext(mockCtrl)
		context.EXPECT().GameID().Return("anID").AnyTimes()
		context.EXPECT().GetDay().Return(0).AnyTimes()

		pollingService := mocks.NewMockPollingService(mockCtrl)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "wolf").Return("wolfpollid").Times(1)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "doctor").Return("doctorpollid").Times(1)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "seer").Return("seerpollid").Times(0)

		playerDAO := mocks.NewMockPlayerDAO(mockCtrl)
		playerDAO.EXPECT().GetWolves(context.GameID()).Return([]models.Player{}, nil).Times(1)
		playerDAO.EXPECT().GetDoctor(context.GameID()).Return(&models.Player{}, nil).Times(1)
		playerDAO.EXPECT().GetSeer(context.GameID()).Return(nil, fmt.Errorf("Any error")).Times(1)

		night := NightTime{
			PlayerDAO:        playerDAO,
			PollingService:   pollingService,
			WolfVoteDuration: 1 * time.Millisecond,
		}
		err := night.TriggerCharacterPolls(context)

		require.Equal(t, fmt.Errorf("Any error"), err)
	})
}

func TestNightTime_Action(t *testing.T) {

	wolves := []models.Player{
		{ID: "wolf1", Character: characters.NewWolf()},
		{ID: "wolf2", Character: characters.NewWolf()},
	}

	doctor := models.Player{
		ID: "doctor", Character: characters.NewDoctor(),
	}

	seer := models.Player{
		ID: "seer", Character: characters.NewSeer(),
	}

	t.Run("wolves kill other than doctor, not saved by doctor, then target dies", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		context := mocks.NewMockContext(mockCtrl)
		context.EXPECT().GameID().Return("anID").AnyTimes()
		context.EXPECT().GetDay().Return(0).AnyTimes()

		wolvesTarget := models.Player{ID: "villager1", Character: characters.NewVillager()}
		doctorTarget := models.Player{ID: "villager2", Character: characters.NewVillager()}

		pollingService := mocks.NewMockPollingService(mockCtrl)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "wolf").Return("wolfpollid").Times(1)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "doctor").Return("doctorpollid").Times(1)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "seer").Return("seerpollid").Times(1)

		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "wolfpollid").Return(&wolvesTarget, nil).Times(1)
		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "doctorpollid").Return(&doctorTarget, nil).Times(1)
		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "seerpollid").Times(1)

		playerDAO := mocks.NewMockPlayerDAO(mockCtrl)
		playerDAO.EXPECT().GetWolves(context.GameID()).Return(wolves, nil).Times(1)
		playerDAO.EXPECT().GetDoctor(context.GameID()).Return(&doctor, nil).Times(1)
		playerDAO.EXPECT().GetSeer(context.GameID()).Return(&seer, nil).Times(1)

		wolvesTarget.Character.Die()

		playerDAO.EXPECT().SavePlayer(context.GameID(), wolvesTarget).Return(nil).Times(1)

		night := NightTime{
			PlayerDAO:        playerDAO,
			PollingService:   pollingService,
			WolfVoteDuration: 1 * time.Millisecond,
		}

		night.Action(context)
	})

	t.Run("wolves kill other than doctor, saved by doctor, then target doesn't die", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		context := mocks.NewMockContext(mockCtrl)
		context.EXPECT().GameID().Return("anID").AnyTimes()
		context.EXPECT().GetDay().Return(0).AnyTimes()

		wolvesTarget := models.Player{ID: "villager1", Character: characters.NewVillager()}

		pollingService := mocks.NewMockPollingService(mockCtrl)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "wolf").Return("wolfpollid").Times(1)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "doctor").Return("doctorpollid").Times(1)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "seer").Return("seerpollid").Times(1)

		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "wolfpollid").Return(&wolvesTarget, nil).Times(1)
		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "doctorpollid").Return(&wolvesTarget, nil).Times(1)
		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "seerpollid").Times(1)

		playerDAO := mocks.NewMockPlayerDAO(mockCtrl)
		playerDAO.EXPECT().GetWolves(context.GameID()).Return(wolves, nil).Times(1)
		playerDAO.EXPECT().GetDoctor(context.GameID()).Return(&doctor, nil).Times(1)
		playerDAO.EXPECT().GetSeer(context.GameID()).Return(&seer, nil).Times(1)

		wolvesTarget.Character.Die()
		playerDAO.EXPECT().SavePlayer(context.GameID(), wolvesTarget).Return(nil).Times(1)

		wolvesTarget.Character.Revive()
		playerDAO.EXPECT().SavePlayer(context.GameID(), wolvesTarget).Return(nil).Times(1)

		night := NightTime{
			PlayerDAO:        playerDAO,
			PollingService:   pollingService,
			WolfVoteDuration: 1 * time.Millisecond,
		}

		night.Action(context)
	})

	t.Run("wolves kill doctor, doctor saved himself, then doctor doesn't die", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		context := mocks.NewMockContext(mockCtrl)
		context.EXPECT().GameID().Return("anID").AnyTimes()
		context.EXPECT().GetDay().Return(0).AnyTimes()

		wolvesTarget := doctor

		pollingService := mocks.NewMockPollingService(mockCtrl)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "wolf").Return("wolfpollid").Times(1)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "doctor").Return("doctorpollid").Times(1)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "seer").Return("seerpollid").Times(1)

		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "wolfpollid").Return(&wolvesTarget, nil).Times(1)
		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "doctorpollid").Return(&doctor, nil).Times(1)
		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "seerpollid").Times(1) //does

		playerDAO := mocks.NewMockPlayerDAO(mockCtrl)
		playerDAO.EXPECT().GetWolves(context.GameID()).Return(wolves, nil).Times(1)
		playerDAO.EXPECT().GetDoctor(context.GameID()).Return(&doctor, nil).Times(1)
		playerDAO.EXPECT().GetSeer(context.GameID()).Return(&seer, nil).Times(1)

		wolvesTarget.Character.Die()
		playerDAO.EXPECT().SavePlayer(context.GameID(), wolvesTarget).Return(nil).Times(1)

		wolvesTarget.Character.Revive()
		playerDAO.EXPECT().SavePlayer(context.GameID(), wolvesTarget).Return(nil).Times(1)

		night := NightTime{
			PlayerDAO:        playerDAO,
			PollingService:   pollingService,
			WolfVoteDuration: 1 * time.Millisecond,
		}

		night.Action(context)
	})

	t.Run("wolves kill doctor, doctor tried to save other, then doctor dies", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		context := mocks.NewMockContext(mockCtrl)
		context.EXPECT().GameID().Return("anID").AnyTimes()
		context.EXPECT().GetDay().Return(0).AnyTimes()

		wolvesTarget := doctor
		doctorTarget := models.Player{ID: "villager1", Character: characters.NewVillager()}

		pollingService := mocks.NewMockPollingService(mockCtrl)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "wolf").Return("wolfpollid").Times(1)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "doctor").Return("doctorpollid").Times(1)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "seer").Return("seerpollid").Times(1)

		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "wolfpollid").Return(&wolvesTarget, nil).Times(1)
		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "doctorpollid").Return(&doctorTarget, nil).Times(1)
		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "seerpollid").Times(1) //does

		playerDAO := mocks.NewMockPlayerDAO(mockCtrl)
		playerDAO.EXPECT().GetWolves(context.GameID()).Return(wolves, nil).Times(1)
		playerDAO.EXPECT().GetDoctor(context.GameID()).Return(nil, nil).Times(1)
		playerDAO.EXPECT().GetSeer(context.GameID()).Return(&seer, nil).Times(1)

		wolvesTarget.Character.Die()
		playerDAO.EXPECT().SavePlayer(context.GameID(), wolvesTarget).Return(nil).Times(1)

		night := NightTime{
			PlayerDAO:        playerDAO,
			PollingService:   pollingService,
			WolfVoteDuration: 1 * time.Millisecond,
		}

		night.Action(context)
	})

	t.Run("wolves didnt vote at all", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		context := mocks.NewMockContext(mockCtrl)
		context.EXPECT().GameID().Return("anID").AnyTimes()
		context.EXPECT().GetDay().Return(0).AnyTimes()

		pollingService := mocks.NewMockPollingService(mockCtrl)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "wolf").Return("wolfpollid").Times(1)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "doctor").Return("doctorpollid").Times(1)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "seer").Return("seerpollid").Times(1)

		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "wolfpollid").Return(nil, services.ErrNoVotes).Times(1)
		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "doctorpollid").Times(1)
		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "seerpollid").Times(1) //does

		playerDAO := mocks.NewMockPlayerDAO(mockCtrl)
		playerDAO.EXPECT().GetWolves(context.GameID()).Return(wolves, nil).Times(1)
		playerDAO.EXPECT().GetDoctor(context.GameID()).Return(&doctor, nil).Times(1)
		playerDAO.EXPECT().GetSeer(context.GameID()).Return(&seer, nil).Times(1)

		night := NightTime{
			PlayerDAO:        playerDAO,
			PollingService:   pollingService,
			WolfVoteDuration: 1 * time.Millisecond,
		}

		night.Action(context)
	})

	t.Run("wolves votes are draw, no one gets killed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		context := mocks.NewMockContext(mockCtrl)
		context.EXPECT().GameID().Return("anID").AnyTimes()
		context.EXPECT().GetDay().Return(0).AnyTimes()

		pollingService := mocks.NewMockPollingService(mockCtrl)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "wolf").Return("wolfpollid").Times(1)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "doctor").Return("doctorpollid").Times(1)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "seer").Return("seerpollid").Times(1)

		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "wolfpollid").Return(nil, services.ErrNoAgreement).Times(1)
		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "doctorpollid").Times(1)
		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "seerpollid").Times(1) //does

		playerDAO := mocks.NewMockPlayerDAO(mockCtrl)
		playerDAO.EXPECT().GetWolves(context.GameID()).Return(wolves, nil).Times(1)
		playerDAO.EXPECT().GetDoctor(context.GameID()).Return(&doctor, nil).Times(1)
		playerDAO.EXPECT().GetSeer(context.GameID()).Return(&seer, nil).Times(1)

		night := NightTime{
			PlayerDAO:        playerDAO,
			PollingService:   pollingService,
			WolfVoteDuration: 1 * time.Millisecond,
		}

		night.Action(context)
	})

	t.Run("seer checked wolf, seer doesnt dies", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		context := mocks.NewMockContext(mockCtrl)
		context.EXPECT().GameID().Return("anID").AnyTimes()
		context.EXPECT().GetDay().Return(0).AnyTimes()

		seerTarget := wolves[0]

		pollingService := mocks.NewMockPollingService(mockCtrl)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "wolf").Return("wolfpollid").Times(1)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "doctor").Return("doctorpollid").Times(1)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "seer").Return("seerpollid").Times(1)

		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "wolfpollid").Times(1)
		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "doctorpollid").Times(1)
		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "seerpollid").Return(&seerTarget, nil).Times(1)

		playerDAO := mocks.NewMockPlayerDAO(mockCtrl)
		playerDAO.EXPECT().GetWolves(context.GameID()).Return(wolves, nil).Times(1)
		playerDAO.EXPECT().GetDoctor(context.GameID()).Return(&doctor, nil).Times(1)
		playerDAO.EXPECT().GetSeer(context.GameID()).Return(&seer, nil).Times(1)

		night := NightTime{
			PlayerDAO:        playerDAO,
			PollingService:   pollingService,
			WolfVoteDuration: 1 * time.Millisecond,
		}

		night.Action(context)
	})

	t.Run("doctor tried to save wolf, then doctor dies", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		context := mocks.NewMockContext(mockCtrl)
		context.EXPECT().GameID().Return("anID").AnyTimes()
		context.EXPECT().GetDay().Return(0).AnyTimes()

		wolfTarget := models.Player{ID: "villager1", Character: characters.NewVillager()}
		doctorTarget := wolves[0]

		pollingService := mocks.NewMockPollingService(mockCtrl)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "wolf").Return("wolfpollid").Times(1)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "doctor").Return("doctorpollid").Times(1)
		pollingService.EXPECT().NewPollID(context.GameID(), "nighttime", context.GetDay(), "seer").Return("seerpollid").Times(1)

		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "wolfpollid").Return(&wolfTarget, nil).Times(1)
		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "doctorpollid").Return(&doctorTarget, nil).Times(1)
		pollingService.EXPECT().GetMostVotedPlayer(context.GameID(), "seerpollid").Times(1)

		playerDAO := mocks.NewMockPlayerDAO(mockCtrl)
		playerDAO.EXPECT().GetWolves(context.GameID()).Return(wolves, nil).Times(1)
		playerDAO.EXPECT().GetDoctor(context.GameID()).Return(&doctor, nil).Times(1)
		playerDAO.EXPECT().GetSeer(context.GameID()).Return(&seer, nil).Times(1)

		wolfTarget.Character.Die()
		playerDAO.EXPECT().SavePlayer(context.GameID(), wolfTarget).Return(nil).Times(1)

		doctor.Character.Die()
		playerDAO.EXPECT().SavePlayer(context.GameID(), doctor).Return(nil).Times(1)

		night := NightTime{
			PlayerDAO:        playerDAO,
			PollingService:   pollingService,
			WolfVoteDuration: 1 * time.Millisecond,
		}

		night.Action(context)
	})
}
