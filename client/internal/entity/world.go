package entity

type World interface {
	IsLoaded() bool
	SetSize(x1, y1, x2, y2 float64)
	Size() (float64, float64, float64, float64)
	Players() Players
	ProjectileList() ProjectileList

	AddNewKill(playerID uint64, playerClass Class, killerID uint64, killerClass Class)
	ActionTable() ActionTable
}

type world struct {
	loaded         bool
	west, north    float64
	east, south    float64
	players        Players
	projectileList ProjectileList
	actionTable    ActionTable
}

const maxActionTableElements = 5

func NewWorld() World {
	return &world{
		players:        NewPlayers(),
		projectileList: NewProjectileList(),
		actionTable:    newActionTable(maxActionTableElements),
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

func (w *world) AddNewKill(playerID uint64, playerClass Class, killerID uint64, killerClass Class) {
	w.actionTable.AddKill(playerID, playerClass, killerID, killerClass)
}

func (w *world) ActionTable() ActionTable { return w.actionTable }
