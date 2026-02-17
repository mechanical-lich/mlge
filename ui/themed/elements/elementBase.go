package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	theming "github.com/mechanical-lich/mlge/ui/themed/theming"
)

type ElementInterface interface {
	Update()
	Draw(screen *ebiten.Image, theme *theming.Theme)

	IsWithin(cX int, cY int) bool
	IsClicked() bool
	IsJustClicked() bool
	GetName() string
	GetPosition() (int, int)
	SetPosition(x, y int)
	GetWidth() int
	GetHeight() int
	GetFocused() bool
	GetAbsolutePosition() (int, int)
	GetScreenPosition() (int, int)
	SetParent(p ElementInterface)
}

type ElementBase struct {
	Name          string
	X, Y          int
	Width, Height int
	IconX, IconY  int
	IconResource  string
	Visible       bool
	Focused       bool
	Op            *ebiten.DrawImageOptions
	Parent        ElementInterface
}

func (b *ElementBase) GetName() string {
	return b.Name
}

func (b *ElementBase) IsWithin(cX int, cY int) bool {
	absX, absY := b.GetScreenPosition()
	if cX >= absX && cX <= absX+b.Width && cY >= absY && cY <= absY+b.Height {
		return true
	}
	return false
}

func (b *ElementBase) IsClicked() bool {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		cX, cY := ebiten.CursorPosition()
		absX, absY := b.GetScreenPosition()
		if cX >= absX && cX <= absX+b.Width && cY >= absY && cY <= absY+b.Height {
			return true
		}
	}
	return false
}

// Modals override this to do rendering on smaller sub images
func (e *ElementBase) GetAbsolutePosition() (int, int) {
	if e.Parent != nil {
		px, py := e.Parent.GetAbsolutePosition()
		return px + e.X, py + e.Y
	}
	return e.X, e.Y
}

// Intended for input interactions that don't care about rendering.
func (e *ElementBase) GetScreenPosition() (int, int) {
	if e.Parent != nil {
		px, py := e.Parent.GetScreenPosition()
		return px + e.X, py + e.Y
	}
	return e.X, e.Y
}

// Returns if the element was just clicked.  Requires the element to become unclicked to register a new click.
func (b *ElementBase) IsJustClicked() bool {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		cX, cY := ebiten.CursorPosition()
		absX, absY := b.GetScreenPosition()
		if cX >= absX && cX <= absX+b.Width && cY >= absY && cY <= absY+b.Height {
			return true
		}
	}
	return false
}

// Get the position of the element
func (b *ElementBase) GetPosition() (int, int) {
	return b.X, b.Y
}

// Set the position of the element.
func (b *ElementBase) SetPosition(x, y int) {
	b.X = x
	b.Y = y
}

// Get the element's width
func (b *ElementBase) GetWidth() int {
	return b.Width
}

// Get the element's height
func (b *ElementBase) GetHeight() int {
	return b.Height
}

// Get if the element is currently focused on.   Ex.  Input is being typed into.
func (b *ElementBase) GetFocused() bool {
	return b.Focused
}

// Set the element's parent.  Normally a container.
func (b *ElementBase) SetParent(p ElementInterface) {
	b.Parent = p
}

// Get's the parent's current position.
func (b *ElementBase) GetParentPosition() (int, int) {
	if b.Parent != nil {
		return b.Parent.GetPosition()
	}
	return 0, 0
}
