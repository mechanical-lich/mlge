package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	containers "github.com/mechanical-lich/mlge/ui/v2/containers"
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
	Elements map[string]elements.ElementInterface
	Modals   map[string]containers.ModalInterface
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
	g.Elements[element.GetName()] = element
}

func (g *GUIViewBase) AddModal(modal containers.ModalInterface) {
	if g.Modals == nil {
		g.Modals = make(map[string]containers.ModalInterface)
	}

	g.Modals[modal.GetName()] = modal
}

func (g *GUIViewBase) UpdateElements() {
	for _, element := range g.Elements {
		element.Update()
	}

	for _, modal := range g.Modals {
		modal.Update()
	}

}

func (g *GUIViewBase) DrawElements(screen *ebiten.Image, theme *theming.Theme) {
	// Draw buttons
	for _, e := range g.Elements {
		e.Draw(screen, theme)
	}

	// Draw modals
	for _, modal := range g.Modals {
		if modal.IsVisible() {
			modal.Draw(screen, theme)
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
