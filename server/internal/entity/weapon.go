package entity

type Weapon interface {
	Name() string
	Damage() Damage
	AttackRadius() float64
	CoolDown() float64

	Attack(owner Player, attackDirection float64, acquirer projectileAcquirer) (proj Projectile, isMelee bool)
}

func Weapons(class Class) []func() Weapon {
	switch class {
	case ClassWarrior:
		return []func() Weapon{
			newSword,
			newAxe,
			newHammer,
		}
	case ClassMage:
		return []func() Weapon{
			newWand,
		}
	case ClassRanger:
		return []func() Weapon{
			newBow,
		}
	default:
		return []func() Weapon{
			newNoWeapon,
		}
	}
}

type weapon struct {
	damage Damage
}

func (w *weapon) Damage() Damage {
	return w.damage
}

type noWeapon struct {
	weapon
}

func newNoWeapon() Weapon {
	return &noWeapon{
		weapon: weapon{
			damage: newDamage(0, 0),
		},
	}
}

func (noWeapon) Name() string          { return "No weapon" }
func (noWeapon) AttackRadius() float64 { return 0.0 }
func (noWeapon) CoolDown() float64     { return 100.0 }

func (noWeapon) Attack(Player, float64, projectileAcquirer) (Projectile, bool) {
	return nil, true
}

type swordWeapon struct {
	weapon
}

func newSword() Weapon {
	w := &swordWeapon{}
	w.damage = newDamage(13, 0)
	return w
}

func (w *swordWeapon) Name() string          { return "Sword" }
func (w *swordWeapon) AttackRadius() float64 { return 15.0 }
func (w *swordWeapon) CoolDown() float64     { return 0.5 }

func (w *swordWeapon) Attack(Player, float64, projectileAcquirer) (proj Projectile, isMelee bool) {
	return nil, true
}

type axeWeapon struct {
	weapon
}

func newAxe() Weapon {
	w := &axeWeapon{}
	w.damage = newDamage(28, 0)
	return w
}

func (w *axeWeapon) Name() string          { return "Axe" }
func (w *axeWeapon) AttackRadius() float64 { return 20.0 }
func (w *axeWeapon) CoolDown() float64     { return 1.2 }

func (w *axeWeapon) Attack(Player, float64, projectileAcquirer) (proj Projectile, isMelee bool) {
	return nil, true
}

type hammerWeapon struct {
	weapon
}

func newHammer() Weapon {
	w := &hammerWeapon{}
	w.damage = newDamage(36, 0)
	return w
}

func (w *hammerWeapon) Name() string          { return "Hammer" }
func (w *hammerWeapon) AttackRadius() float64 { return 20.0 }
func (w *hammerWeapon) CoolDown() float64     { return 1.5 }

func (w *hammerWeapon) Attack(Player, float64, projectileAcquirer) (proj Projectile, isMelee bool) {
	return nil, true
}

type wandWeapon struct {
	weapon
}

func newWand() Weapon {
	w := &wandWeapon{}
	w.damage = newDamage(0, 21)
	return w
}

func (w *wandWeapon) Name() string          { return "Wand \"Fireball\"" }
func (w *wandWeapon) AttackRadius() float64 { return 40.0 }
func (w *wandWeapon) CoolDown() float64     { return 3.0 }

func (w *wandWeapon) Attack(owner Player, attackDirection float64, acquirer projectileAcquirer) (proj Projectile, isMelee bool) {
	return newFireball(acquirer.AcquireProjectileID(), owner, w, attackDirection), false
}

type bowWeapon struct {
	weapon
}

func newBow() Weapon {
	w := &bowWeapon{}
	w.damage = newDamage(26, 0)
	return w
}

func (w *bowWeapon) Name() string          { return "Bow" }
func (w *bowWeapon) AttackRadius() float64 { return 0.0 }
func (w *bowWeapon) CoolDown() float64     { return 1.8 }

func (w *bowWeapon) Attack(owner Player, attackDirection float64, acquirer projectileAcquirer) (proj Projectile, isMelee bool) {
	return newArrow(acquirer.AcquireProjectileID(), owner, w, attackDirection), false
}
