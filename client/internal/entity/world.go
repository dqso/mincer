package entity

type World interface {
	IsLoaded() bool
	SetSize(x1, y1, x2, y2 float64)
	Size() (float64, float64, float64, float64)
	Players() Players
	ProjectileList() ProjectileList

	AddNewKill(playerID, killerID uint64)
	KillTable() KillTable
}

type world struct {
	loaded         bool
	west, north    float64
	east, south    float64
	players        Players
	projectileList ProjectileList
	killTable      KillTable
}

const maxKillTableElements = 5

func NewWorld() World {
	return &world{
		players:        NewPlayers(),
		projectileList: NewProjectileList(),
		killTable:      newKillTable(maxKillTableElements),
	}
}

func (w *world) IsLoaded() bool {
	return w.loaded
}

func (w *world) SetSize(x1, y1, x2, y2 float64) {
	w.west, w.north = x1, y1
	w.east, w.south = x2, y2
	w.loaded = true
}

func (w *world) Size() (float64, float64, float64, float64) {
	return w.west, w.north, w.east, w.south
}

func (w *world) Players() Players               { return w.players }
func (w *world) ProjectileList() ProjectileList { return w.projectileList }

func (w *world) AddNewKill(playerID, killerID uint64) {
	player, _ := w.players.Get(playerID)
	killer, _ := w.players.Get(killerID)
	w.killTable.Add(playerID, player, killerID, killer)
}

func (w *world) KillTable() KillTable { return w.killTable }
