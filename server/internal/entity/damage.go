package entity

type Damage interface {
	Physical() int32
	Magical() int32
}

type damage struct {
	physical int32
	magical  int32
}

func newDamage(physical, magical int32) *damage {
	return &damage{
		physical: physical,
		magical:  magical,
	}
}

func (d *damage) Physical() int32 { return d.physical }
func (d *damage) Magical() int32  { return d.magical }
