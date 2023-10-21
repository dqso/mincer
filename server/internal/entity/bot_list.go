package entity

import (
	"sync"
	"sync/atomic"
)

type BotList interface {
	NewBot(w World, class Class, weapon Weapon) Bot
}

type botList struct {
	botLastID uint32
	byID      map[uint32]Bot
	mxByID    sync.RWMutex
	stops     map[uint32]func()
	mxStops   sync.Mutex
}

func NewBotList() BotList {
	return &botList{
		botLastID: 0,
		byID:      make(map[uint32]Bot),
		stops:     make(map[uint32]func()),
	}
}

func (l *botList) NewBot(w World, class Class, weapon Weapon) Bot {
	id := atomic.AddUint32(&l.botLastID, 1)
	b, stop := newBot(w, id, class, weapon)
	l.addBot(b)
	l.addStop(id, stop)
	return b
}

func (l *botList) addStop(id uint32, stop func()) {
	l.mxStops.Lock()
	defer l.mxStops.Unlock()
	l.stops[id] = stop
}

func (l *botList) addBot(b Bot) {
	l.mxByID.Lock()
	defer l.mxByID.Unlock()
	l.byID[b.BotID()] = b
}
