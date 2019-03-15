package states

import (
	"time"

	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	e "github.com/imrenagi/goes-werewolf/internal/app/events"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/models"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/models/events"
	"github.com/imrenagi/goes-werewolf/internal/app/werewolf/services"
	"github.com/imrenagi/goes-werewolf/pkg/eventbus"
)

type PollingService interface {
	NewPollID(gameID, state string, day int, character string) string
	GetMostVotedPlayer(gameID, pollID string) (*models.Player, error)
}

type NightTime struct {
	State
	PlayerDAO        services.PlayerDAO
	PollingService   PollingService
	WolfVoteDuration time.Duration
}

func (n *NightTime) TriggerCharacterPolls(c Context) error {

	log := logrus.WithFields(map[string]interface{}{
		"gameID": c.GameID(),
		"day":    c.GetDay(),
		"state":  n.String(),
	})

	wolves, err := n.PlayerDAO.GetWolves(c.GameID())
	if err != nil {
		log.WithError(err).Error("Can't get all wolves.")
		return err
	}

	wolvesPollID := n.PollingService.NewPollID(c.GameID(), n.String(), c.GetDay(), "wolf")
	for _, wolf := range wolves {
		//TODO make sure it is called
		eventbus.Publish(events.WolvesVoteStarted{
			DomainEvent: e.DomainEvent{
				ID:           uuid.New().String(),
				EventVersion: 1,
				Timestamp:    time.Now(),
				EventName:    "WolvesVoteStarted",
			},
			GameID:    c.GameID(),
			PollID:    wolvesPollID,
			GameDay:   c.GetDay(),
			GameState: n.String(),
			SendTo:    wolf,
		})
	}

	doctor, err := n.PlayerDAO.GetDoctor(c.GameID())
	if err != nil {
		log.WithError(err).Error("Can't get doctor.")
		return err
	}
	doctorPollID := n.PollingService.NewPollID(c.GameID(), n.String(), c.GetDay(), "doctor")
	if doctor != nil {
		eventbus.Publish(events.DoctorTargetRequested{
			DomainEvent: e.DomainEvent{
				ID:           uuid.New().String(),
				EventVersion: 1,
				Timestamp:    time.Now(),
				EventName:    "DoctorTargetRequested",
			},
			GameID:    c.GameID(),
			PollID:    doctorPollID,
			GameDay:   c.GetDay(),
			GameState: n.String(),
			SendTo:    *doctor,
		})
	}

	seer, err := n.PlayerDAO.GetSeer(c.GameID())
	if err != nil {
		log.WithError(err).Error("Can't get seer.")
		return err
	}
	seerPollID := n.PollingService.NewPollID(c.GameID(), n.String(), c.GetDay(), "seer")
	if seer != nil {
		eventbus.Publish(events.SeerTargetRequested{
			DomainEvent: e.DomainEvent{
				ID:           uuid.New().String(),
				EventVersion: 1,
				Timestamp:    time.Now(),
				EventName:    "SeerTargetRequested",
			},
			GameID:    c.GameID(),
			PollID:    seerPollID,
			GameDay:   c.GetDay(),
			GameState: n.String(),
			SendTo:    *seer,
		})
	}

	<-time.After(n.WolfVoteDuration)

	return nil
}

func (n *NightTime) Action(c Context) error {

	log := logrus.WithFields(map[string]interface{}{
		"gameID": c.GameID(),
		"day":    c.GetDay(),
		"state":  n.String(),
	})

	wolves, err := n.PlayerDAO.GetWolves(c.GameID())
	if err != nil {
		log.WithError(err).Error("Can't get all wolves.")
		return err
	}

	wolvesPollID := n.PollingService.NewPollID(c.GameID(), n.String(), c.GetDay(), "wolf")
	victim, err := n.PollingService.GetMostVotedPlayer(c.GameID(), wolvesPollID)
	if err != nil {
		if err == services.ErrNoVotes || err == services.ErrNoAgreement {
			//TODO publish wolves didnt kill anybody
		} else {
			log.WithError(err).Error("Can't get most voted player for wolf")
			return err
		}
	}
	if victim != nil {
		fmt.Println("killed")
		victim.Character.Accept(wolves[0].Character)

		err := n.PlayerDAO.SavePlayer(c.GameID(), *victim)
		if err != nil {
			log.WithError(err).Error("Can't save victim state to datastore")
		}
	}

	doctor, err := n.PlayerDAO.GetDoctor(c.GameID())
	if err != nil {
		log.WithError(err).Error("Can't get doctor.")
		return err
	}

	//get the saved player
	doctorPollID := n.PollingService.NewPollID(c.GameID(), n.String(), c.GetDay(), "doctor")
	savedPlayer, err := n.PollingService.GetMostVotedPlayer(c.GameID(), doctorPollID)
	if err != nil {
		if err != services.ErrNoVotes && err != services.ErrNoAgreement {
			log.WithError(err).Error("Can't get selected player to be saved by doctor")
			return err
		}
	}

	if doctor != nil && savedPlayer != nil {
		if victim != nil && savedPlayer.ID == victim.ID {
			//TODO publish someone has been saved
			victim.Character.Accept(doctor.Character)

			err := n.PlayerDAO.SavePlayer(c.GameID(), *victim)
			if err != nil {
				log.WithError(err).Error("Can't save victim state to datastore")
			}
		} else {
			savedPlayer.Character.Accept(doctor.Character)

			if !doctor.Character.IsAlive() {
				err := n.PlayerDAO.SavePlayer(c.GameID(), *doctor)
				if err != nil {
					log.WithError(err).Error("Can't save doctor state to datastore")
				}
			}
		}
	}

	//get the seer's selected player
	seer, err := n.PlayerDAO.GetSeer(c.GameID())
	if err != nil {
		log.WithError(err).Error("Can't get seer.")
		return err
	}
	seerPollID := n.PollingService.NewPollID(c.GameID(), n.String(), c.GetDay(), "seer")
	checkedPlayer, err := n.PollingService.GetMostVotedPlayer(c.GameID(), seerPollID)
	if err != nil {
		if err != services.ErrNoVotes && err != services.ErrNoAgreement {
			log.WithError(err).Error("Can't get selected player to be checked by seer")
			return err
		}
	}

	if seer != nil && checkedPlayer != nil {
		//TODO optional: if selected player by seer is the same as wolves's voted player, make seer accept wolves.
		checkedPlayer.Character.Accept(seer.Character)
		//TODO publish by seer. This dude get info of this user. can be inside the models?
	}

	return nil
}

func (n *NightTime) Execute(c Context) {
	log := logrus.WithFields(map[string]interface{}{
		"gameID": c.GameID(),
		"day":    c.GetDay(),
		"state":  n.String(),
	})
	log.Info("Starting nighttime state")

	err := n.TriggerCharacterPolls(c)
	if err != nil {
		c.SetState(&Error{Error: err})
		return
	}

	err = n.Action(c)
	if err != nil {
		c.SetState(&Error{Error: err})
		return
	}

	c.SetState(&DayTime{})

	log.Info("Nighttime state completed")
}

func (n *NightTime) String() string {
	return "nighttime"
}
