package states

import (
	"time"

	"github.com/sirupsen/logrus"
)

type DayTime struct {
	State

	DiscussionDuration time.Duration
}

func (n *DayTime) Execute(c Context) {
	log := logrus.WithFields(map[string]interface{}{
		"gameID": c.GameID(),
		"day":    c.GetDay(),
		"state":  n.String(),
	})
	log.Info("Starting daytime state")

	//let channel knows morning is come and who died or saved, last night

	<-time.After(n.DiscussionDuration)

	c.SetState(&Dusk{})
}

func (n *DayTime) String() string {
	return "daytime"
}
