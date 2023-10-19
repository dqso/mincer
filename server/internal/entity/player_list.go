package entity

import "sync"

type PlayerList interface {
	Get(id uint64) (Player, bool)
	Add(p Player)
	Remove(id uint64)
	Slice() []Player
}

type playerList struct {
	byID   map[uint64]Player
	mxByID sync.RWMutex
}

func NewPlayers() PlayerList {
	return &playerList{
		byID: make(map[uint64]Player),
	}
}

func (pp *playerList) Get(id uint64) (Player, bool) {
	pp.mxByID.RLock()
	defer pp.mxByID.RUnlock()
	p, ok := pp.byID[id]
	return p, ok
}

func (pp *playerList) Add(p Player) {
	pp.mxByID.Lock()
	defer pp.mxByID.Unlock()
	pp.byID[p.ID()] = p
}

func (pp *playerList) Remove(id uint64) {
	pp.mxByID.Lock()
	defer pp.mxByID.Unlock()
	delete(pp.byID, id)
}

func (pp *playerList) Slice() []Player {
	pp.mxByID.RLock()
	defer pp.mxByID.RUnlock()
	out := make([]Player, 0, len(pp.byID))
	for _, p := range pp.byID {
		out = append(out, p)
	}
	return out
}
