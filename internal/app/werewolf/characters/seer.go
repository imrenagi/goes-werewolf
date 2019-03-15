package characters

func NewSeer() *Seer {
	return &Seer{

		Alive: true,
	}
}

type Seer struct {
	Alive bool
}

func (s *Seer) VisitWolf(wolf *Wolf) {
}
func (s *Seer) VisitVillager(villager *Villager) {
}
func (s *Seer) VisitSeer(seer *Seer) {
	panic("There must be only 1 seer in the game.")
}
func (s *Seer) VisitDoctor(doctor *Doctor) {
}
func (s *Seer) Accept(visitor Visitor) {
	visitor.VisitSeer(s)
}
func (s *Seer) String() string {
	return "Seer"
}

// Die will kills its self
func (s *Seer) Die() {
	s.Alive = false
}

// IsAlive checks whether it is alive
func (s Seer) IsAlive() bool {
	return s.Alive
}

// Revive revives a character back to live from death
func (s *Seer) Revive() {
	s.Alive = true
}
