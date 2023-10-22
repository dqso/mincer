package entity

import (
	"github.com/dqso/mincer/server/internal/api"
	"math"
	"sync"
	"time"
)

type Player interface {
	ID() uint64

	GetStats() PlayerStats
	PlayerStats

	Weapon() Weapon
	SetWeapon(w Weapon)

	HP() int32
	IsDead() bool
	SetHP(hp int32) (newHP int32, wasChanged bool)
	DealDamage(killerID uint64, damage Damage)

	Position() Point
	SetPosition(p Point)

	SetDirection(direction float64, isMoving bool)
	Move(lifeTime time.Duration, mapSize Rect) (wasMoved bool)
	SetAttack(a bool, directionAim float64)
	Attack() bool
	DirectionAim() float64
	Relax(lifeTime time.Duration)

	getHorn() Horn
}

type player struct {
	horn Horn

	id uint64

	PlayerStats

	mxWeapon sync.RWMutex
	weapon   Weapon

	mxHP sync.RWMutex
	hp   int32

	mxPosition sync.RWMutex
	x, y       float64
	// for calculating x and y
	direction float64
	isMoving  bool

	mxAttack     sync.RWMutex
	isAttack     bool
	directionAim float64
	coolDown     float64
}

func newPlayer(id uint64, class Class, weapon Weapon, horn Horn) Player {
	p := &player{
		horn:   horn,
		id:     id,
		weapon: weapon,
		PlayerStats: newPlayerStats(
			class,
			defaultPlayerRadius,
			defaultPlayerSpeed,
			defaultPlayerHP,
		),
		x: 0,
		y: 0,
	}
	p.coolDown = weapon.CoolDown()
	p.hp = p.MaxHP()
	return p
}

func (p *player) ID() uint64 { return p.id }

func (p *player) GetStats() PlayerStats {
	return p.PlayerStats
}

func (p *player) Weapon() Weapon {
	p.mxWeapon.RLock()
	defer p.mxWeapon.RUnlock()
	return p.weapon
}

func (p *player) SetWeapon(w Weapon) {
	p.mxWeapon.Lock()
	defer p.mxWeapon.Unlock()
	p.weapon = w
	p.horn.SetPlayerWeapon(p.ID(), p.weapon)
}

func (p *player) SetClass(v Class) {
	p.PlayerStats.SetClass(v)
	p.horn.SetPlayerStats(p.ID(), p.PlayerStats)
}

func (p *player) SetRadius(v float64) {
	p.PlayerStats.SetRadius(v)
	p.horn.SetPlayerStats(p.ID(), p.PlayerStats)
}

func (p *player) SetSpeed(v float64) {
	p.PlayerStats.SetSpeed(v)
	p.horn.SetPlayerStats(p.ID(), p.PlayerStats)
}

func (p *player) SetMaxHP(v int32) {
	p.PlayerStats.SetMaxHP(v)
	p.horn.SetPlayerStats(p.ID(), p.PlayerStats)
}

func (p *player) HP() int32 {
	p.mxHP.RLock()
	defer p.mxHP.RUnlock()
	return p.hp
}

func (p *player) IsDead() bool {
	return p.HP() <= 0
}

func (p *player) SetHP(hp int32) (_ int32, wasChanged bool) {
	if hp < 0 {
		hp = 0
	}
	if p.MaxHP() < hp {
		hp = 100
	}
	func() {
		p.mxHP.Lock()
		defer p.mxHP.Unlock()
		if p.hp != hp {
			wasChanged = true
			p.hp = hp
		}
	}()
	if wasChanged {
		p.horn.SetPlayerHP(p.ID(), hp)
	}
	return hp, wasChanged
}

func (p *player) DealDamage(killerID uint64, damage Damage) {
	newHP, wasChanged := p.SetHP(p.HP() - damage.CalculateWith(p.PlayerStats))
	if wasChanged && newHP == 0 {
		p.horn.OnPlayerWasted(p.ID(), killerID)
	}
}

func (p *player) Position() Point {
	p.mxPosition.RLock()
	defer p.mxPosition.RUnlock()
	return Point{X: p.x, Y: p.y}
}

func (p *player) SetPosition(point Point) {
	func() {
		p.mxPosition.Lock()
		defer p.mxPosition.Unlock()
		p.x, p.y = point.X, point.Y
	}()
	p.horn.SetPlayerPosition(p.ID(), Point{p.x, p.y})
}

func (p *player) SetDirection(direction float64, isMoving bool) {
	p.mxPosition.Lock()
	defer p.mxPosition.Unlock()
	p.direction, p.isMoving = direction, isMoving
}

func (p *player) Move(lifeTime time.Duration, mapSize Rect) (wasMoved bool) {
	if p.IsDead() {
		wasMoved = false
		return
	}
	func() {
		p.mxPosition.Lock()
		defer p.mxPosition.Unlock()
		if !p.isMoving {
			wasMoved = false
			return
		}
		sin, cos := math.Sincos(p.direction * math.Pi / 180)
		x := p.x + p.Speed()*lifeTime.Seconds()*sin
		y := p.y - p.Speed()*lifeTime.Seconds()*cos
		if x < mapSize.LeftUp.X || mapSize.RightDown.X < x ||
			y < mapSize.LeftUp.Y || mapSize.RightDown.Y < y {
			wasMoved = false
			return
		}
		p.x, p.y = x, y
		wasMoved = true
		return
	}()
	if wasMoved {
		p.horn.SetPlayerPosition(p.ID(), Point{X: p.x, Y: p.y})
	}
	return wasMoved
}

func (p *player) SetAttack(a bool, directionAim float64) {
	p.mxAttack.Lock()
	defer p.mxAttack.Unlock()
	p.isAttack = a
	p.directionAim = directionAim
}

func (p *player) DirectionAim() float64 {
	p.mxAttack.RLock()
	defer p.mxAttack.RUnlock()
	return p.directionAim
}

func (p *player) Attack() (isAllowed bool) {
	if p.HP() <= 0 {
		return false
	}
	p.mxAttack.Lock()
	defer p.mxAttack.Unlock()
	defer func() {
		if isAllowed {
			p.isAttack = false
			p.coolDown = 0
		}
	}()
	if !p.isAttack || p.coolDown < p.Weapon().CoolDown() {
		return false
	}
	return true
}

func (p *player) Relax(lifeTime time.Duration) {
	p.mxAttack.Lock()
	defer p.mxAttack.Unlock()
	p.coolDown += lifeTime.Seconds()
}

func (p *player) getHorn() Horn {
	return p.horn
}

type Class int32

const (
	ClassWarrior = Class(api.Class_WARRIOR)
	ClassMage    = Class(api.Class_MAGE)
	ClassRanger  = Class(api.Class_RANGER)
)

func Classes() []Class {
	return []Class{
		ClassWarrior,
		ClassMage,
		ClassRanger,
	}
}

func (c Class) physicalResist() float64 {
	switch c {
	case ClassWarrior:
		return 1.46
	case ClassMage:
		return 1.0
	case ClassRanger:
		return 1.21
	default:
		return 1.0
	}
}

func (c Class) magicalResist() float64 {
	switch c {
	case ClassWarrior:
		return 1.0
	case ClassMage:
		return 1.34
	case ClassRanger:
		return 1.05
	default:
		return 1.0
	}
}
