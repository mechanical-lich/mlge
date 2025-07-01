package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/state"
)

// Base GUIView interface.
// Since we are dealing with interfaces the GUIView is being passed around by value instead of reference
type GUIViewInterface interface {
	Update(state state.StateInterface)
	Draw(screen *ebiten.Image, s state.StateInterface)
}

// GUIViewBase gives views some basic functionality when inherited.
type GUIViewBase struct {
	Buttons []*Button
}

func (g *GUIViewBase) AddButton(button *Button) {
	if g.Buttons == nil {
		g.Buttons = make([]*Button, 0)
	}
	g.Buttons = append(g.Buttons, button)
}
