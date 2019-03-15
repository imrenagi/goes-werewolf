package states

import (
	"fmt"
)

type End struct {
	State
}

func (i *End) Execute(c Context) {
	fmt.Println("end state")

}

func (i *End) String() string {
	return "end"
}
