package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameInput struct {
	Left   int
	Up     int
	Right  int
	Down   int
	Attack int
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
