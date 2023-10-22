package entity

import "math"

type Damage interface {
	Physical() int32
	Magical() int32
	CalculateWith(r Resist) int32
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

func (d *damage) CalculateWith(r Resist) int32 {
	var value float64
	if r.PhysicalResist() != 0 {
		value += float64(d.Physical()) / r.PhysicalResist()
	}
	if r.MagicalResist() != 0 {
		value += float64(d.Magical()) / r.MagicalResist()
	}
	return int32(math.Round(value))
}
