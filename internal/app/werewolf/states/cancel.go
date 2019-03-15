package states

import "fmt"

type Cancel struct {
	State
}

func (s *Cancel) Execute(c Context) {
	fmt.Println("cancel state")

}

func (s *Cancel) String() string {
	return "cancel"
}
