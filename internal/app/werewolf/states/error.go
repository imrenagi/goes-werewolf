package states

import "fmt"

type Error struct {
	State

	Error error
}

func (s *Error) Execute(c Context) {
	fmt.Println("error state")

}

func (s *Error) String() string {
	return "internal error"
}
