package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/state"
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
}

// NewModal creates a new modal with initial view.
func NewModal(name string, x, y, width, height int, initialView string, views map[string]GUIViewInterface) *Modal {
	closeBtn := NewButton("close", x+width-24, y+8, "X")
	return &Modal{
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
	m.CloseButton.X = m.X + m.Width - 24
	m.CloseButton.Y = m.Y + 8

	if m.CloseButton.IsClicked() {
		m.Visible = false
		if m.OnClose != nil {
			m.OnClose()
		}
		return
	}
	if v, ok := m.Views[m.CurrentView]; ok {
		v.Update(s)
	}
}

// Draw renders the modal background, close button, and current view.
func (m *Modal) Draw(screen *ebiten.Image, s state.StateInterface) {
	if !m.Visible {
		return
	}
	// Draw modal background (simple rectangle, replace with sprite if needed)
	bg := ebiten.NewImage(m.Width, m.Height)
	bg.Fill(color.RGBA{30, 30, 30, 240})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(m.X), float64(m.Y))
	screen.DrawImage(bg, op)

	// Draw close button
	m.CloseButton.Draw(screen)

	// Draw current view inside modal
	if v, ok := m.Views[m.CurrentView]; ok {
		v.Draw(screen, s)
	}
}

// GetInputFocused delegates to the current view.
func (m *Modal) GetInputFocused() bool {
	if v, ok := m.Views[m.CurrentView]; ok {
		return v.GetInputFocused()
	}
	return false
}
