package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ElementBase struct {
	Name          string
	X, Y          int
	Width, Height int
	IconX, IconY  int
	IconResource  string
	Visible       bool
	op            *ebiten.DrawImageOptions
}

func (b *ElementBase) IsWithin(cX int, cY int, parentX, parentY int) bool {
	if cX >= b.X+parentX && cX <= b.X+b.Width+parentX && cY >= b.Y+parentY && cY <= b.Height+b.Y+parentY {
		return true
	}
	return false
}

func (b *ElementBase) IsClicked(parentX, parentY int) bool {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		cX, cY := ebiten.CursorPosition()

		if cX >= b.X+parentX && cX <= b.X+b.Width+parentX && cY >= b.Y+parentY && cY <= b.Height+b.Y+parentY {
			return true
		}
	}
	return false
}

func (b *ElementBase) IsJustClicked(parentX, parentY int) bool {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		cX, cY := ebiten.CursorPosition()

		if cX >= b.X+parentX && cX <= b.X+b.Width+parentX && cY >= b.Y+parentY && cY <= b.Height+b.Y+parentY {
			return true
		}
	}
	return false
}
