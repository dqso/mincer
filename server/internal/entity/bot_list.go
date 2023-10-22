package entity

import (
	"sync"
	"sync/atomic"
)

type BotList interface {
	AcquireID() uint32
	Add(b Bot, stopFunc func())
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

func (l *botList) AcquireID() uint32 {
	return atomic.AddUint32(&l.botLastID, 1)
}

func (l *botList) Add(b Bot, stopFunc func()) {
	l.addBot(b)
	l.addStop(b.BotID(), stopFunc)
}

func (l *botList) addBot(b Bot) {
	l.mxByID.Lock()
	defer l.mxByID.Unlock()
	l.byID[b.BotID()] = b
}

func (l *botList) addStop(id uint32, stop func()) {
	l.mxStops.Lock()
	defer l.mxStops.Unlock()
	l.stops[id] = stop
}
