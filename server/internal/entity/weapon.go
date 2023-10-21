package entity

import (
	"image/color"
	"sync"
)

const (
	maxDistanceProjectile   = 500.0
	defaultProjectileRadius = 4.0
	defaultProjectileSpeed  = 10.0
)

type Weapon interface {
	PhysicalDamage() float64
	MagicalDamage() float64
	AttackRadius() float64
	CoolDown() float64

	Attack(owner Player, attackDirection float64, acquirer projectileAcquirer) (proj Projectile, isMelee bool)
}

func Weapons(class Class) []func() Weapon {
	switch class {
	case ClassWarrior:
		return []func() Weapon{
			newSword,
		}
	case ClassMage:
		return []func() Weapon{
			newWand,
		}
	default:
		return []func() Weapon{
			newNoWeapon,
		}
	}
}

type noWeapon struct{}

func newNoWeapon() Weapon { return &noWeapon{} }

func (noWeapon) PhysicalDamage() float64 { return 0.0 }
func (noWeapon) MagicalDamage() float64  { return 0.0 }
func (noWeapon) AttackRadius() float64   { return 0.0 }
func (noWeapon) CoolDown() float64       { return 100.0 }

func (noWeapon) Attack(owner Player, attackDirection float64, acquirer projectileAcquirer) (proj Projectile, isMelee bool) {
	return nil, true
}

type warriorWeapon struct{}

func (w *warriorWeapon) MagicalDamage() float64 { return 0.0 }

type swordWeapon struct {
	warriorWeapon
}

func newSword() Weapon { return &swordWeapon{} }

func (w *swordWeapon) PhysicalDamage() float64 { return 10.0 }
func (w *swordWeapon) AttackRadius() float64   { return 40.0 }
func (w *swordWeapon) CoolDown() float64       { return 0.5 }

func (w *swordWeapon) Attack(owner Player, attackDirection float64, acquirer projectileAcquirer) (proj Projectile, isMelee bool) {
	return nil, true
}

type Projectile interface {
	Color() color.NRGBA
	Position() Point
	PhysicalDamage() float64
	MagicalDamage() float64
	Radius() float64
	Speed() float64
	Direction() float64
	Distance() float64
	Owner() uint64
}

type magicalWeapon struct{}

func (w *magicalWeapon) PhysicalDamage() float64 { return 0.0 }

type wandWeapon struct {
	magicalWeapon
}

func newWand() Weapon { return &wandWeapon{} }

func (w *wandWeapon) MagicalDamage() float64 { return 10.0 }
func (w *wandWeapon) AttackRadius() float64  { return 40.0 }
func (w *wandWeapon) CoolDown() float64      { return 3.0 }

type projectileAcquirer interface {
	AcquireProjectileID() uint64
}

func (w *wandWeapon) Attack(owner Player, attackDirection float64, acquirer projectileAcquirer) (proj Projectile, isMelee bool) {
	projectileID := acquirer.AcquireProjectileID()
	return newFireball(projectileID, owner, w, attackDirection), false
}

type projectile struct {
	id uint64

	position   Point
	distance   float64
	mxPosition sync.RWMutex

	direction      float64
	radius         float64
	speed          float64
	physicalDamage float64
	magicalDamage  float64
	owner          uint64
}

func newProjectile(id uint64, owner Player, weapon Weapon, attackDirection float64) *projectile {
	return &projectile{
		id:             id,
		position:       owner.Position(),
		distance:       maxDistanceProjectile,
		radius:         defaultProjectileRadius,
		speed:          defaultProjectileSpeed,
		direction:      attackDirection,
		physicalDamage: weapon.PhysicalDamage(),
		magicalDamage:  weapon.MagicalDamage(),
		owner:          owner.ID(),
	}
}

func (p *projectile) Position() Point {
	p.mxPosition.RLock()
	defer p.mxPosition.RUnlock()
	return p.position
}

func (p *projectile) Distance() float64 {
	p.mxPosition.RLock()
	defer p.mxPosition.RUnlock()
	return p.distance
}

func (p *projectile) Radius() float64         { return p.radius }
func (p *projectile) Speed() float64          { return p.speed }
func (p *projectile) Direction() float64      { return p.direction }
func (p *projectile) PhysicalDamage() float64 { return p.physicalDamage }
func (p *projectile) MagicalDamage() float64  { return p.magicalDamage }
func (p *projectile) Owner() uint64           { return p.owner }

type fireball struct {
	*projectile
}

func newFireball(id uint64, owner Player, weapon Weapon, attackDirection float64) *fireball {
	return &fireball{
		projectile: newProjectile(id, owner, weapon, attackDirection),
	}
}

func (p *fireball) Color() color.NRGBA { return color.NRGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF} }
