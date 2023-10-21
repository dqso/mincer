package entity

type Weapon interface {
	Name() string
	PhysicalDamage() int32
	MagicalDamage() int32
	CoolDown() float64
}

type weapon struct {
	name string

	physicalDamage int32
	magicalDamage  int32

	coolDown float64
}

func NewWeapon(name string, physicalDamage, magicalDamage int32, coolDown float64) Weapon {
	return &weapon{
		name:           name,
		physicalDamage: physicalDamage,
		magicalDamage:  magicalDamage,
		coolDown:       coolDown,
	}
}

func (w *weapon) Name() string          { return w.name }
func (w *weapon) PhysicalDamage() int32 { return w.physicalDamage }
func (w *weapon) MagicalDamage() int32  { return w.magicalDamage }
func (w *weapon) CoolDown() float64     { return w.coolDown }
