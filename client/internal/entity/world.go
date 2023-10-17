package entity

import (
	"sort"
	"sync"
)

type World struct {
	Me        *Me
	players   map[uint64]*Player
	mxPlayers sync.RWMutex
}

func NewWorld() *World {
	return &World{
		Me: &Me{
			Player: &Player{},
		},
		players: make(map[uint64]*Player),
	}
}

func (w *World) PlayerIDs() []uint64 {
	w.mxPlayers.RLock()
	defer w.mxPlayers.RUnlock()
	out := make([]uint64, 0, len(w.players)+1)
	for id := range w.players {
		out = append(out, id)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i] < out[j]
	})
	out = append(out, w.Me.ID)
	return out
}

func (w *World) Player(id uint64) (*Player, bool) {
	if w.Me.ID == id {
		return w.Me.Player, true
	}
	w.mxPlayers.RLock()
	defer w.mxPlayers.RUnlock()
	player, ok := w.players[id]
	return player, ok
}

func (w *World) SetPlayer(p *Player) {
	if p.ID == w.Me.ID {
		w.Me.Player = p
		return
	}
	w.mxPlayers.Lock()
	defer w.mxPlayers.Unlock()
	w.players[p.ID] = p
}
