package characters

func NewDoctor() *Doctor {
	return &Doctor{
		Alive: true,
	}
}

type Doctor struct {
	Alive bool
}

func (d *Doctor) VisitWolf(wolf *Wolf) {
	d.Die()
}

func (d *Doctor) VisitVillager(villager *Villager) {
	villager.Revive()
}

func (d *Doctor) VisitDoctor(doctor *Doctor) {
	doctor.Revive()
}

func (d *Doctor) VisitSeer(seer *Seer) {
	seer.Revive()
}

func (d *Doctor) Accept(visitor Visitor) {
	visitor.VisitDoctor(d)
}

func (d Doctor) String() string {
	return "Doctor"
}

// Die will kills its self
func (d *Doctor) Die() {
	d.Alive = false
}

// IsAlive checks whether it is alive
func (d Doctor) IsAlive() bool {
	return d.Alive
}

// Revive revives a character back to live from death
func (d *Doctor) Revive() {
	d.Alive = true
}
