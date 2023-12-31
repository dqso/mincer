package scene

import (
	"github.com/dqso/mincer/client/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Manager struct {
	events  events
	current Scene
	next    Scene
	world   entity.World

	transition     int
	transitionFrom *ebiten.Image
	transitionTo   *ebiten.Image
}

const transitionMaxCount = 20

type Scene interface {
	Update(state State) error
	Draw(screen *ebiten.Image)
}

type events interface {
	Connected() chan struct{}
	ConnectingInformation() chan string
	Disconnected() chan struct{}

	MustDisconnect()
	BeReborn()
}

func NewManager(initial Scene, events events, world entity.World) *Manager {
	m := &Manager{
		current: initial,
		events:  events,
		world:   world,
	}

	w, h := ebiten.WindowSize()
	m.transitionFrom = ebiten.NewImage(w, h)
	m.transitionTo = ebiten.NewImage(w, h)

	return m
}

type State struct {
	manager *Manager
	events  events
	world   entity.World
}

func (m *Manager) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		m.Go(NewCloseScene())
		m.events.MustDisconnect()
		return nil
	}

	if m.transition == 0 {
		return m.current.Update(State{
			manager: m,
			events:  m.events,
			world:   m.world,
		})
	}

	m.transition--
	if m.transition > 0 {
		return nil
	}

	m.current, m.next = m.next, nil
	return nil
}

func (m *Manager) Draw(screen *ebiten.Image) {
	if m.transition <= 0 {
		m.current.Draw(screen)
		return
	}

	m.transitionFrom.Clear()
	m.current.Draw(m.transitionFrom)

	m.transitionTo.Clear()
	m.next.Draw(m.transitionTo)

	screen.DrawImage(m.transitionFrom, nil)

	alpha := 1 - float32(m.transition)/float32(transitionMaxCount)
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(alpha)
	screen.DrawImage(m.transitionTo, op)
}

func (m *Manager) Go(next Scene) {
	if m.current == nil {
		m.current = next
	} else {
		m.next = next
		m.transition = transitionMaxCount
	}
}
