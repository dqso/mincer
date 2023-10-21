package entity

import (
	"image/color"
	"math"
	"sync"
	"time"
)

const (
	maxDistanceProjectile   = 500.0
	defaultProjectileRadius = 4.0
	defaultProjectileSpeed  = 200.0
)

type Projectile interface {
	ID() uint64
	Color() color.NRGBA
	Position() Point
	PhysicalDamage() int32
	MagicalDamage() int32
	AttackRadius() float64
	Radius() float64
	Speed() float64
	Direction() float64
	Distance() float64
	Owner() uint64
	Move(lifeTime time.Duration, masSize Rect) (newPosition Point, outOfRange bool)
	CollisionAnalysis(position Point, players []Player) Player
}

type projectileAcquirer interface {
	AcquireProjectileID() uint64
}

type projectile struct {
	id   uint64
	horn Horn

	position   Point
	distance   float64
	mxPosition sync.RWMutex

	direction      float64
	radius         float64
	speed          float64
	physicalDamage int32
	magicalDamage  int32
	attackRadius   float64
	owner          uint64
}

func newProjectile(id uint64, owner Player, weapon Weapon, attackDirection float64) *projectile {
	return &projectile{
		id:             id,
		horn:           owner.getHorn(),
		position:       owner.Position(),
		distance:       maxDistanceProjectile,
		radius:         defaultProjectileRadius,
		speed:          defaultProjectileSpeed,
		direction:      attackDirection,
		physicalDamage: weapon.PhysicalDamage(),
		magicalDamage:  weapon.MagicalDamage(),
		attackRadius:   weapon.AttackRadius(),
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

func (p *projectile) ID() uint64            { return p.id }
func (p *projectile) Radius() float64       { return p.radius }
func (p *projectile) Speed() float64        { return p.speed }
func (p *projectile) Direction() float64    { return p.direction }
func (p *projectile) PhysicalDamage() int32 { return p.physicalDamage }
func (p *projectile) MagicalDamage() int32  { return p.magicalDamage }
func (p *projectile) AttackRadius() float64 { return p.attackRadius }
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

func (p *projectile) Move(lifeTime time.Duration, mapSize Rect) (newPosition Point, outOfRange bool) {
	wasMoved := false
	func() {
		p.mxPosition.Lock()
		defer p.mxPosition.Unlock()
		sin, cos := math.Sincos(p.direction * math.Pi / 180)
		x := p.position.X + p.Speed()*lifeTime.Seconds()*sin
		y := p.position.Y - p.Speed()*lifeTime.Seconds()*cos
		if x < mapSize.LeftUp.X || mapSize.RightDown.X < x ||
			y < mapSize.LeftUp.Y || mapSize.RightDown.Y < y {
			outOfRange = true
			return
		}
		newPosition = Point{X: x, Y: y}
		p.position = newPosition
		wasMoved = true
		return
	}()
	if wasMoved {
		p.horn.SetProjectilePosition(p.ID(), newPosition)
	}
	return

}

func (p *projectile) CollisionAnalysis(position Point, players []Player) Player {
	for _, player := range players {
		if player.ID() == p.owner {
			continue
		}
		playerPos := player.Position()
		distance := math.Hypot(playerPos.X-position.X, playerPos.Y-position.Y)
		if distance-p.Radius() <= player.Radius() {
			return player
		}
	}
	return nil
}
