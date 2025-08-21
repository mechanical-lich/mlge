package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	elements "github.com/mechanical-lich/mlge/ui/v2/elements"

	theming "github.com/mechanical-lich/mlge/ui/v2/theming"
)

// Base GUIView interface.
// Since we are dealing with interfaces the GUIView is being passed around by value instead of reference
type GUIViewInterface interface {
	Update()
	UpdateElements()
	Draw(screen *ebiten.Image, theme *theming.Theme)
	DrawElements(screen *ebiten.Image, theme *theming.Theme)
	GetInputFocused() bool
	GetPosition() (int, int)
	SetPosition(x, y int)
	GetMouseFocused() bool
}

// GUIViewBase gives views some basic functionality when inherited.
type GUIViewBase struct {
	elements.ElementBase
	Elements map[string]elements.ElementInterface
	X, Y     int // Add offset for the view
}

func (g *GUIViewBase) initElements() {
	if g.Elements == nil {
		g.Elements = make(map[string]elements.ElementInterface, 0)
	}
}

func (g *GUIViewBase) GetPosition() (int, int) {
	return g.X, g.Y
}

func (g *GUIViewBase) SetPosition(x, y int) {
	g.X = x
	g.Y = y
}

func (g *GUIViewBase) AddElement(element elements.ElementInterface) {
	g.initElements()
	element.SetParent(g)
	g.Elements[element.GetName()] = element
}

func (g *GUIViewBase) UpdateElements() {
	for _, element := range g.Elements {
		element.Update()
	}
}

func (g *GUIViewBase) DrawElements(screen *ebiten.Image, theme *theming.Theme) {
	// Draw buttons
	for _, e := range g.Elements {
		e.Draw(screen, theme)
	}
}

func (g *GUIViewBase) Update() {
}

func (g *GUIViewBase) Draw(screen *ebiten.Image, theme *theming.Theme) {

}

func (g *GUIViewBase) GetInputFocused() bool {
	for _, e := range g.Elements {
		if e.GetFocused() {
			return true
		}
	}

	return false
}

func (g *GUIViewBase) GetMouseFocused() bool {
	for _, e := range g.Elements {
		if e.GetFocused() {
			return true
		}
	}

	return false
}
