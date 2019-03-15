package states

import (
	"time"

	"github.com/sirupsen/logrus"
)

type Dusk struct {
	State

	VillagerVoteDuration time.Duration
}

func (d *Dusk) Execute(c Context) {
	log := logrus.WithFields(map[string]interface{}{
		"gameID": c.GameID(),
		"day":    c.GetDay(),
		"state":  d.String(),
	})
	log.Info("Starting dusk state")

	//notify people to start voting
	//send vote option

	//when people send their vote, group will know

	<-time.After(d.VillagerVoteDuration)

	//get vote result
	//calculate who needs to be lynch
	//if votes are equal, lynch nobody

	c.SetState(&NightTime{})
}

func (d *Dusk) String() string {
	return "dusk"
}
