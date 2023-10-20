package entity

type World interface {
	Players() Players

	AddNewKill(playerID, killerID uint64)
	KillTable() KillTable
}

type world struct {
	players   Players
	killTable KillTable
}

const maxKillTableElements = 5

func NewWorld() World {
	return &world{
		players:   NewPlayers(),
		killTable: newKillTable(maxKillTableElements),
	}
}

func (w *world) Players() Players { return w.players }

func (w *world) AddNewKill(playerID, killerID uint64) {
	player, _ := w.players.Get(playerID)
	killer, _ := w.players.Get(killerID)
	w.killTable.Add(playerID, player, killerID, killer)
}

func (w *world) KillTable() KillTable { return w.killTable }
