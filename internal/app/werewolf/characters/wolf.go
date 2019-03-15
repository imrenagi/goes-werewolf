package characters

func NewWolf() *Wolf {
	w := &Wolf{
		Alive: true,
	}
	return w
}

type Wolf struct {
	Alive bool
}

func (w Wolf) VisitWolf(wolf *Wolf) {
	panic("Wolf should not visit wolf")
}

func (w Wolf) VisitVillager(villager *Villager) {
	villager.Die()
}

func (w Wolf) VisitDoctor(doctor *Doctor) {
	doctor.Die()
}

func (w *Wolf) VisitSeer(seer *Seer) {
	seer.Die()
}

func (w *Wolf) Accept(visitor Visitor) {
	visitor.VisitWolf(w)
}

func (w Wolf) String() string {
	return "Wolf"
}

// Die will kills its self
func (w *Wolf) Die() {
	w.Alive = false
}

// IsAlive checks whether it is alive
func (w Wolf) IsAlive() bool {
	return w.Alive
}

// Revive revives a character back to live from death
func (w *Wolf) Revive() {
	w.Alive = true
}
