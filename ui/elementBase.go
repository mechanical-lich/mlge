package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ElementInterface interface {
	Update(parentX, parentY int)
	Draw(screen *ebiten.Image, parentX, parentY int, theme *Theme)

	IsWithin(cX int, cY int, parentX, parentY int) bool
	IsClicked(parentX, parentY int) bool
	IsJustClicked(parentX, parentY int) bool
	GetName() string
	GetPosition() (int, int)
	SetPosition(x, y int)
	GetWidth() int
	GetHeight() int
	GetFocused() bool
}

type ElementBase struct {
	Name          string
	X, Y          int
	Width, Height int
	IconX, IconY  int
	IconResource  string
	Visible       bool
	Focused       bool
	op            *ebiten.DrawImageOptions
}

func (b *ElementBase) GetName() string {
	return b.Name
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

func (b *ElementBase) GetPosition() (int, int) {
	return b.X, b.Y
}
func (b *ElementBase) SetPosition(x, y int) {
	b.X = x
	b.Y = y
}
func (b *ElementBase) GetWidth() int {
	return b.Width
}
func (b *ElementBase) GetHeight() int {
	return b.Height
}

func (b *ElementBase) GetFocused() bool {
	return b.Focused
}
