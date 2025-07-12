package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/state"
	"github.com/mechanical-lich/mlge/utility"
)

// Modal wraps a view and provides modal behavior with a close button and view states.
type Modal struct {
	Name          string
	Views         map[string]GUIViewInterface
	CurrentView   string
	X, Y          int
	Width, Height int
	Visible       bool
	CloseButton   *Button
	OnClose       func()

	dragging    bool
	dragOffsetX int
	dragOffsetY int
	bg          *ebiten.Image // Background image for the modal
}

// NewModal creates a new modal with initial view.
func NewModal(name string, x, y, width, height int, initialView string, views map[string]GUIViewInterface) *Modal {
	closeBtn := NewButton("close", width-24, 8, "X")
	return &Modal{
		Name:        name,
		Views:       views,
		CurrentView: initialView,
		X:           x,
		Y:           y,
		Width:       width,
		Height:      height,
		Visible:     true,
		CloseButton: closeBtn,
	}
}

// SetView switches the modal to a different view state.
func (m *Modal) SetView(name string) {
	if _, ok := m.Views[name]; ok {
		m.CurrentView = name
	}
}

// Update handles modal logic and delegates to the current view.
func (m *Modal) Update(s state.StateInterface) {
	if !m.Visible {
		return
	}

	m.HandleDragging()

	if m.CloseButton.IsJustClicked(m.X, m.Y) && !m.dragging {
		m.Visible = false
		if m.OnClose != nil {
			m.OnClose()
		}
		return
	}
	if v, ok := m.Views[m.CurrentView]; ok {
		v.SetPosition(m.X, m.Y)
		v.UpdateElements(s)
		v.Update(s)
	}
}

func (m *Modal) HandleDragging() {
	// Get mouse position and button state
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		m.dragging = false
		return
	}

	mouseX, mouseY := ebiten.CursorPosition()

	// Define a draggable area (e.g., top 32px of modal, excluding close button)
	dragAreaX := m.X
	dragAreaY := m.Y
	dragAreaW := m.Width - 32 // leave space for close button
	dragAreaH := 32

	// Start dragging if mouse pressed in drag area and not already dragging
	if !m.dragging &&
		mouseX >= dragAreaX && mouseX < dragAreaX+dragAreaW &&
		mouseY >= dragAreaY && mouseY < dragAreaY+dragAreaH &&
		!m.CloseButton.IsClicked(m.X, m.Y) {
		m.dragging = true
		m.dragOffsetX = mouseX - m.X
		m.dragOffsetY = mouseY - m.Y
	}

	// If dragging, update modal position
	if m.dragging {
		m.X = mouseX - m.dragOffsetX
		m.Y = mouseY - m.dragOffsetY
	}
}

// Draw renders the modal background, close button, and current view.
func (m *Modal) Draw(screen *ebiten.Image, s state.StateInterface, theme *Theme) {
	if !m.Visible {
		return
	}

	if m.bg == nil {
		m.bg = ebiten.NewImage(m.Width, m.Height)
		utility.Draw9Slice(m.bg, 0, 0, m.Width, m.Height, theme.ModalNineSlice.SrcX, theme.ModalNineSlice.SrcY, theme.ModalNineSlice.TileSize, theme.ModalNineSlice.TileScale)
	}
	// Draw the modal background
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(m.X), float64(m.Y))
	screen.DrawImage(m.bg, op)

	m.CloseButton.Draw(screen, m.X, m.Y, theme)

	if v, ok := m.Views[m.CurrentView]; ok {
		// If the view is a GUIViewBase, set its X/Y to modal's X/Y
		v.SetPosition(m.X, m.Y)
		v.Draw(screen, s, theme)
		v.DrawElements(screen, s, theme)
	}
}

// GetInputFocused delegates to the current view.
func (m *Modal) GetInputFocused() bool {
	if v, ok := m.Views[m.CurrentView]; ok {
		return v.GetInputFocused()
	}
	return false
}
