package states

import "fmt"

type State interface {
	fmt.Stringer
	Execute(context Context)
}

type Context interface {
	GameID() string
	GetDay() int
	SetState(s State)
}
