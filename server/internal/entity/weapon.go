package entity

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
			newAxe,
			newHammer,
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

func (w *warriorWeapon) Attack(Player, float64, projectileAcquirer) (proj Projectile, isMelee bool) {
	return nil, true
}

type swordWeapon struct {
	warriorWeapon
}

func newSword() Weapon { return &swordWeapon{} }

func (w *swordWeapon) Name() string          { return "Sword" }
func (w *swordWeapon) PhysicalDamage() int32 { return 13 }
func (w *swordWeapon) AttackRadius() float64 { return 15.0 }
func (w *swordWeapon) CoolDown() float64     { return 0.5 }

type axeWeapon struct {
	warriorWeapon
}

func newAxe() Weapon { return &axeWeapon{} }

func (w *axeWeapon) Name() string          { return "Axe" }
func (w *axeWeapon) PhysicalDamage() int32 { return 28 }
func (w *axeWeapon) AttackRadius() float64 { return 20.0 }
func (w *axeWeapon) CoolDown() float64     { return 1.2 }

type hammerWeapon struct {
	warriorWeapon
}

func newHammer() Weapon { return &hammerWeapon{} }

func (w *hammerWeapon) Name() string          { return "Hammer" }
func (w *hammerWeapon) PhysicalDamage() int32 { return 36 }
func (w *hammerWeapon) AttackRadius() float64 { return 20.0 }
func (w *hammerWeapon) CoolDown() float64     { return 1.5 }

type magicalWeapon struct{}

func (w *magicalWeapon) PhysicalDamage() int32 { return 0 }

type wandWeapon struct {
	magicalWeapon
}

func newWand() Weapon { return &wandWeapon{} }

func (w *wandWeapon) Name() string          { return "Wand \"Fireball\"" }
func (w *wandWeapon) MagicalDamage() int32  { return 21 }
func (w *wandWeapon) AttackRadius() float64 { return 40.0 }
func (w *wandWeapon) CoolDown() float64     { return 3.0 }

func (w *wandWeapon) Attack(owner Player, attackDirection float64, acquirer projectileAcquirer) (proj Projectile, isMelee bool) {
	return newFireball(acquirer.AcquireProjectileID(), owner, w, attackDirection), false
}
