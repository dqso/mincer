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
	Name() string
	PhysicalDamage() int32
	MagicalDamage() int32
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

func (noWeapon) Name() string          { return "No weapon" }
func (noWeapon) PhysicalDamage() int32 { return 0 }
func (noWeapon) MagicalDamage() int32  { return 0 }
func (noWeapon) AttackRadius() float64 { return 0.0 }
func (noWeapon) CoolDown() float64     { return 100.0 }

func (noWeapon) Attack(Player, float64, projectileAcquirer) (Projectile, bool) {
	return nil, true
}

type warriorWeapon struct{}

func (w *warriorWeapon) MagicalDamage() int32 { return 0 }

type swordWeapon struct {
	warriorWeapon
}

func newSword() Weapon { return &swordWeapon{} }

func (w *swordWeapon) Name() string          { return "Sword" }
func (w *swordWeapon) PhysicalDamage() int32 { return 10 }
func (w *swordWeapon) AttackRadius() float64 { return 40.0 }
func (w *swordWeapon) CoolDown() float64     { return 0.5 }

func (w *swordWeapon) Attack(owner Player, attackDirection float64, acquirer projectileAcquirer) (proj Projectile, isMelee bool) {
	return nil, true
}

type Projectile interface {
	Color() color.NRGBA
	Position() Point
	PhysicalDamage() int32
	MagicalDamage() int32
	Radius() float64
	Speed() float64
	Direction() float64
	Distance() float64
	Owner() uint64
}

type magicalWeapon struct{}

func (w *magicalWeapon) PhysicalDamage() int32 { return 0 }

type wandWeapon struct {
	magicalWeapon
}

func newWand() Weapon { return &wandWeapon{} }

func (w *wandWeapon) Name() string          { return "Wand \"Fireball\"" }
func (w *wandWeapon) MagicalDamage() int32  { return 10 }
func (w *wandWeapon) AttackRadius() float64 { return 40.0 }
func (w *wandWeapon) CoolDown() float64     { return 3.0 }

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
	physicalDamage int32
	magicalDamage  int32
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

func (p *projectile) Radius() float64       { return p.radius }
func (p *projectile) Speed() float64        { return p.speed }
func (p *projectile) Direction() float64    { return p.direction }
func (p *projectile) PhysicalDamage() int32 { return p.physicalDamage }
func (p *projectile) MagicalDamage() int32  { return p.magicalDamage }
func (p *projectile) Owner() uint64         { return p.owner }

type fireball struct {
	*projectile
}

func newFireball(id uint64, owner Player, weapon Weapon, attackDirection float64) *fireball {
	return &fireball{
		projectile: newProjectile(id, owner, weapon, attackDirection),
	}
}

func (p *fireball) Color() color.NRGBA { return color.NRGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF} }
