package input

import (
	"github.com/dqso/mincer/client/internal/entity"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math"
)

type GameInput struct {
	Left          int
	Up            int
	Right         int
	Down          int
	Attack        int
	beReborn      bool
	mousePosition entity.Point
}

func NewGameInput() *GameInput {
	return &GameInput{}
}

func (i *GameInput) Update() {
	i.Left = max(inpututil.KeyPressDuration(ebiten.KeyLeft), inpututil.KeyPressDuration(ebiten.KeyA))
	i.Up = max(inpututil.KeyPressDuration(ebiten.KeyUp), inpututil.KeyPressDuration(ebiten.KeyW))
	i.Right = max(inpututil.KeyPressDuration(ebiten.KeyRight), inpututil.KeyPressDuration(ebiten.KeyD))
	i.Down = max(inpututil.KeyPressDuration(ebiten.KeyDown), inpututil.KeyPressDuration(ebiten.KeyS))

	i.Attack = inpututil.KeyPressDuration(ebiten.KeySpace)

	i.beReborn = inpututil.IsKeyJustPressed(ebiten.KeyR)

	mx, my := ebiten.CursorPosition()
	i.mousePosition = entity.Point{X: float64(mx), Y: float64(my)}
}

func (i *GameInput) BeReborn() bool {
	if i.beReborn {
		i.beReborn = false
		return true
	}
	return false
}

func (i *GameInput) Direction() (direction float64, isMoving bool) {
	if i.Up > i.Down {
		if i.Left > i.Right {
			return 315, true // ↖
		} else if i.Left < i.Right {
			return 45, true // ↗
		}
		return 0, true // ⬆
	} else if i.Up < i.Down {
		if i.Left > i.Right {
			return 225, true // ↙
		} else if i.Left < i.Right {
			return 135, true // ↘
		}
		return 180, true // ⬇
	}
	if i.Left > i.Right {
		return 270, true // ⬅
	} else if i.Left < i.Right {
		return 90, true // ➡
	}
	return 0, false
}

func (i *GameInput) MouseDirection(mePos entity.Point) float64 {
	return math.Mod(360.0-math.Atan2(mePos.X-i.mousePosition.X, mePos.Y-i.mousePosition.Y)*180.0/math.Pi, 360.0)
}
