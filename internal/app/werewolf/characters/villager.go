package characters

func NewVillager() *Villager {
	return &Villager{
		Alive: true,
	}
}

type Villager struct {
	Alive bool
}

func (v Villager) VisitWolf(wolf *Wolf) {
	panic("Villager should not visit wolf")
}
func (v Villager) VisitVillager(villager *Villager) {
	panic("Villager should not visit other villager")
}
func (v Villager) VisitDoctor(doctor *Doctor) {
	panic("Villager should not visit doctor")
}

func (v Villager) VisitSeer(seer *Seer) {
	panic("Villager should not visit seer")
}

func (v *Villager) Accept(visitor Visitor) {
	visitor.VisitVillager(v)
}

func (v Villager) String() string {
	return "Villager"
}

// Die will kills its self
func (v *Villager) Die() {
	v.Alive = false
}

// IsAlive checks whether it is alive
func (v Villager) IsAlive() bool {
	return v.Alive
}

// Revive revives a character back to live from death
func (v *Villager) Revive() {
	v.Alive = true
}
