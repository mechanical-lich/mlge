package minui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// GUI is the main UI manager
type GUI struct {
	RootStyle *Style
	Theme     *Theme
	elements  []Element
	modals    []Element
}

// NewGUI creates a new GUI manager
func NewGUI() *GUI {
	return &GUI{
		RootStyle: DefaultStyle(),
		Theme:     nil, // No theme by default (uses vector rendering)
		elements:  make([]Element, 0),
		modals:    make([]Element, 0),
	}
}

// NewGUIWithTheme creates a new GUI manager with a theme
func NewGUIWithTheme(theme *Theme) *GUI {
	return &GUI{
		RootStyle: DefaultStyle(),
		Theme:     theme,
		elements:  make([]Element, 0),
		modals:    make([]Element, 0),
	}
}

// SetTheme sets the theme for the GUI
func (g *GUI) SetTheme(theme *Theme) {
	g.Theme = theme
}

// GetTheme returns the current theme
func (g *GUI) GetTheme() *Theme {
	return g.Theme
}

// AddElement adds an element to the GUI
func (g *GUI) AddElement(element Element) {
	g.elements = append(g.elements, element)
	element.SetParent(nil) // Root elements have no parent

	// Set theme on the element if one is configured
	if g.Theme != nil {
		if setter, ok := element.(interface{ SetTheme(*Theme) }); ok {
			setter.SetTheme(g.Theme)
		}
	}
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

	// Set theme on the modal if one is configured
	if g.Theme != nil {
		if setter, ok := modal.(interface{ SetTheme(*Theme) }); ok {
			setter.SetTheme(g.Theme)
		}
	}
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

// WithinModalBounds returns true if the given coordinates are inside any
// visible modal. Used by game code to decide whether mouse input should be
// intercepted by the UI.
func (g *GUI) WithinModalBounds(mouseX, mouseY int) bool {
	for _, m := range g.modals {
		if !m.IsVisible() {
			continue
		}
		if m.IsWithin(mouseX, mouseY) {
			return true
		}
	}
	return false
}

// GetMouseFocused returns true if any element (including modal contents)
// currently reports hovered or focused state. This indicates the GUI is
// capturing the mouse and game code should avoid handling mouse input.
func (g *GUI) GetMouseFocused() bool {
	// Check top-level elements
	for _, e := range g.elements {
		if elementHasFocusRecursive(e) {
			return true
		}
	}

	// Check visible modals
	for _, m := range g.modals {
		if !m.IsVisible() {
			continue
		}
		if elementHasFocusRecursive(m) {
			return true
		}
	}

	return false
}

// GetKeyboardFocused returns true if any element (including modals) currently
// has keyboard focus (IsFocused). Game code can use this to decide whether
// it should route keyboard input to the GUI or handle it itself.
func (g *GUI) GetKeyboardFocused() bool {
	for _, e := range g.elements {
		if elementHasKeyboardFocusRecursive(e) {
			return true
		}
	}
	for _, m := range g.modals {
		if !m.IsVisible() {
			continue
		}
		if elementHasKeyboardFocusRecursive(m) {
			return true
		}
	}
	return false
}

// elementHasFocusRecursive checks an element and its children for hovered
// or focused state.
func elementHasFocusRecursive(e Element) bool {
	if e == nil || !e.IsVisible() {
		return false
	}
	if e.IsHovered() || e.IsFocused() {
		return true
	}

	// Special case: Check if element is a SelectBox with an expanded dropdown
	if e.GetType() == "SelectBox" {
		if sb, ok := e.(*SelectBox); ok && sb.IsExpanded() {
			mx, my := ebiten.CursorPosition()
			if sb.IsMouseOverDropdown(mx, my) {
				return true
			}
		}
	}

	for _, c := range e.GetChildren() {
		if elementHasFocusRecursive(c) {
			return true
		}
	}
	return false
}

// elementHasKeyboardFocusRecursive checks element and its children for keyboard focus only.
func elementHasKeyboardFocusRecursive(e Element) bool {
	if e == nil || !e.IsVisible() {
		return false
	}
	if e.IsFocused() {
		return true
	}
	for _, c := range e.GetChildren() {
		if elementHasKeyboardFocusRecursive(c) {
			return true
		}
	}
	return false
}
