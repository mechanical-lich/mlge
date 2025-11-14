package minui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// GUI is the main UI manager
type GUI struct {
	RootStyle *Style
	elements  []Element
	modals    []Element
}

// NewGUI creates a new GUI manager
func NewGUI() *GUI {
	return &GUI{
		RootStyle: DefaultStyle(),
		elements:  make([]Element, 0),
		modals:    make([]Element, 0),
	}
}

// AddElement adds an element to the GUI
func (g *GUI) AddElement(element Element) {
	g.elements = append(g.elements, element)
	element.SetParent(nil) // Root elements have no parent
}

// RemoveElement removes an element from the GUI
func (g *GUI) RemoveElement(element Element) {
	for i, e := range g.elements {
		if e == element {
			g.elements = append(g.elements[:i], g.elements[i+1:]...)
			break
		}
	}
}

// AddModal adds a modal dialog
func (g *GUI) AddModal(modal Element) {
	g.modals = append(g.modals, modal)
}

// RemoveModal removes a modal dialog
func (g *GUI) RemoveModal(modal Element) {
	for i, m := range g.modals {
		if m == modal {
			g.modals = append(g.modals[:i], g.modals[i+1:]...)
			break
		}
	}
}

// Update updates all elements
func (g *GUI) Update() {
	// Check if there are any visible modals
	hasVisibleModal := false
	for _, modal := range g.modals {
		if modal.IsVisible() {
			hasVisibleModal = true
			break
		}
	}

	if hasVisibleModal {
		// Only update the topmost visible modal
		for i := len(g.modals) - 1; i >= 0; i-- {
			if g.modals[i].IsVisible() {
				g.modals[i].Update()
				return
			}
		}
	}

	// Update regular elements
	for _, element := range g.elements {
		element.Update()
	}
}

// Layout performs layout calculations for all elements
func (g *GUI) Layout() {
	// Layout regular elements
	for _, element := range g.elements {
		element.Layout()
	}

	// Layout modals
	for _, modal := range g.modals {
		modal.Layout()
	}
}

// Draw draws all elements
func (g *GUI) Draw(screen *ebiten.Image) {
	// Draw regular elements
	for _, element := range g.elements {
		element.Draw(screen)
	}

	// Draw modals on top
	for _, modal := range g.modals {
		modal.Draw(screen)
	}
}

// FindElementByID finds an element by its ID
func (g *GUI) FindElementByID(id string) Element {
	return g.findInElements(id, g.elements)
}

// findInElements recursively searches for an element by ID
func (g *GUI) findInElements(id string, elements []Element) Element {
	for _, element := range elements {
		if element.GetID() == id {
			return element
		}

		// Search children
		children := element.GetChildren()
		if found := g.findInElements(id, children); found != nil {
			return found
		}
	}
	return nil
}
