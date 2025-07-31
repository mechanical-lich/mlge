package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/state"
)

// Base GUIView interface.
// Since we are dealing with interfaces the GUIView is being passed around by value instead of reference
type GUIViewInterface interface {
	Update(state state.StateInterface)
	UpdateElements(state state.StateInterface)
	Draw(screen *ebiten.Image, s state.StateInterface, theme *Theme)
	DrawElements(screen *ebiten.Image, s state.StateInterface, theme *Theme)
	GetInputFocused() bool
	GetModalFocused() bool
	OpenModal(name string)
	CloseModal(name string)
	ModalOpen(name string) bool
	GetPosition() (int, int)
	SetPosition(x, y int)
	WithinModalBounds(mouseX, mouseY int) bool
	GetMouseFocused() bool
}

// GUIViewBase gives views some basic functionality when inherited.
type GUIViewBase struct {
	Elements map[string]ElementInterface
	Modals   map[string]ModalInterface
	X, Y     int // Add offset for the view
}

func (g *GUIViewBase) initElements() {
	if g.Elements == nil {
		g.Elements = make(map[string]ElementInterface, 0)
	}
}

func (g *GUIViewBase) GetPosition() (int, int) {
	return g.X, g.Y
}

func (g *GUIViewBase) SetPosition(x, y int) {
	g.X = x
	g.Y = y
}

func (g *GUIViewBase) AddElement(element ElementInterface) {
	g.initElements()
	g.Elements[element.GetName()] = element
}

func (g *GUIViewBase) AddModal(modal ModalInterface) {
	if g.Modals == nil {
		g.Modals = make(map[string]ModalInterface)
	}

	g.Modals[modal.GetName()] = modal
}

func (g *GUIViewBase) UpdateElements(s state.StateInterface) {
	for _, element := range g.Elements {
		element.Update(g.X, g.Y)
	}

	for _, modal := range g.Modals {
		modal.Update(s)
	}

}

func (g *GUIViewBase) DrawElements(screen *ebiten.Image, s state.StateInterface, theme *Theme) {
	// Draw buttons
	for _, e := range g.Elements {
		e.Draw(screen, g.X, g.Y, theme)
	}

	// Draw modals
	for _, modal := range g.Modals {
		if modal.IsVisible() {
			modal.Draw(screen, s, theme)
		}
	}

}

func (g *GUIViewBase) GetInputFocused() bool {
	for _, e := range g.Elements {
		if e.GetFocused() {
			return true
		}
	}

	for _, modal := range g.Modals {
		if modal.IsOpen() && modal.GetInputFocused() {
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

	for _, modal := range g.Modals {
		if modal.IsOpen() && modal.GetMouseFocused() {
			return true
		}
	}
	return false
}

func (g *GUIViewBase) GetModalFocused() bool {
	for _, modal := range g.Modals {
		if modal.IsVisible() {
			return true
		}
	}
	return false
}

func (g *GUIViewBase) ModalOpen(name string) bool {
	if modal, exists := g.Modals[name]; exists {
		return modal.IsOpen()
	}
	return false
}

// Opens a modal by name.  Does nothing if the modal does not exist.
func (g *GUIViewBase) OpenModal(name string) {
	if modal, exists := g.Modals[name]; exists {
		modal.OpenModal()
	}
}

// Closes a modal by name.  Does nothing if the modal does not exist.
func (g *GUIViewBase) CloseModal(name string) {
	if modal, exists := g.Modals[name]; exists {
		modal.CloseModal()
	}
}

func (g *GUIViewBase) WithinModalBounds(mouseX, mouseY int) bool {
	for _, modal := range g.Modals {
		if modal.IsVisible() && modal.WithinBounds(mouseX, mouseY) {
			return true
		}
	}
	return false
}
