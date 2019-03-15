package characters

import "fmt"

type Visitor interface {
	VisitWolf(wolf *Wolf)
	VisitVillager(villager *Villager)
	VisitSeer(seer *Seer)
	VisitDoctor(doctor *Doctor)
}

type Visitable interface {
	Accept(visitor Visitor)
}

type Killable interface {
	Die()
}

type Liveness interface {
	IsAlive() bool
}

type Revivable interface {
	Revive()
}

type Interface interface {
	fmt.Stringer
	Visitor
	Visitable
	Killable
	Revivable
	Liveness
}
