package entity

import "sync"

type ProjectileList interface {
	Add(p Projectile)
	Remove(id uint64)
	Slice() []Projectile
}

type projectileList struct {
	byID   map[uint64]Projectile
	mxByID sync.RWMutex
}

func newProjectileList() ProjectileList {
	return &projectileList{
		byID: make(map[uint64]Projectile),
	}
}

func (pp *projectileList) Add(p Projectile) {
	pp.mxByID.Lock()
	defer pp.mxByID.Unlock()
	pp.byID[p.ID()] = p
}

func (pp *projectileList) Remove(id uint64) {
	pp.mxByID.Lock()
	defer pp.mxByID.Unlock()
	delete(pp.byID, id)
}

func (pp *projectileList) Slice() []Projectile {
	pp.mxByID.RLock()
	defer pp.mxByID.RUnlock()
	out := make([]Projectile, 0, len(pp.byID))
	for _, p := range pp.byID {
		out = append(out, p)
	}
	return out
}
