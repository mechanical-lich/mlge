package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/state"
	"github.com/mechanical-lich/mlge/utility"
)

type ModalInterface interface {
	SetView(name string)
	Update(s state.StateInterface)
	Draw(screen *ebiten.Image, s state.StateInterface, theme *Theme)
	GetInputFocused() bool
	GetMouseFocused() bool
	WithinBounds(mouseX, mouseY int) bool
	GetName() string
	GetPosition() (int, int)
	SetPosition(x, y int)
	OpenModal()
	CloseModal()
	IsOpen() bool
	IsVisible() bool
}

// Modal wraps a view and provides modal behavior with a close button and view states.
type Modal struct {
	ElementBase

	Views       map[string]GUIViewInterface
	CurrentView string

	CloseButton *Button
	OnClose     func()

	dragging    bool
	dragOffsetX int
	dragOffsetY int
	bg          *ebiten.Image // Background image for the modal
	offscreen   *ebiten.Image // Offscreen buffer for double buffering
}

// NewModal creates a new modal with initial view.
func NewModal(name string, x, y, width, height int, initialView string, views map[string]GUIViewInterface) *Modal {
	closeBtn := NewButton("close", width-24, 8, "X", "close")
	return &Modal{
		ElementBase: ElementBase{
			Name:    name,
			X:       x,
			Y:       y,
			Width:   width,
			Height:  height,
			Visible: false,
			op:      &ebiten.DrawImageOptions{},
		},
		Views:       views,
		CurrentView: initialView,
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

	m.handleDragging()

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

func (m *Modal) handleDragging() {
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

	// Ensure offscreen buffer is correct size
	if m.bg == nil || m.bg.Bounds().Dx() != m.Width || m.bg.Bounds().Dy() != m.Height {
		m.bg = ebiten.NewImage(m.Width, m.Height)
		utility.Draw9Slice(m.bg, 0, 0, m.Width, m.Height, theme.ModalNineSlice.SrcX, theme.ModalNineSlice.SrcY, theme.ModalNineSlice.TileSize, theme.ModalNineSlice.TileScale)
	}
	if m.offscreen == nil || m.offscreen.Bounds().Dx() != m.Width || m.offscreen.Bounds().Dy() != m.Height {
		m.offscreen = ebiten.NewImage(m.Width, m.Height)
	}

	// Clear offscreen buffer
	m.offscreen.Clear()

	// Draw modal background to offscreen
	opBg := &ebiten.DrawImageOptions{}
	opBg.GeoM.Reset()
	opBg.GeoM.Translate(0, 0)
	m.offscreen.DrawImage(m.bg, opBg)

	// Draw close button to offscreen
	m.CloseButton.Draw(m.offscreen, 0, 0, theme)

	// Draw current view to offscreen
	if v, ok := m.Views[m.CurrentView]; ok {
		v.SetPosition(0, 0)
		v.Draw(m.offscreen, s, theme)
		v.DrawElements(m.offscreen, s, theme)
	}

	// Draw offscreen buffer to screen at modal position
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(m.X), float64(m.Y))
	screen.DrawImage(m.offscreen, op)
}

// GetInputFocused delegates to the current view.
func (m *Modal) GetInputFocused() bool {
	if v, ok := m.Views[m.CurrentView]; ok {
		return v.GetInputFocused()
	}
	return false
}

func (m *Modal) GetMouseFocused() bool {

	cX, cY := ebiten.CursorPosition()
	fmt.Println("Modal GetInputFocused:", m.Name, "Cursor:", cX, cY, "Position:", m.X, m.Y, "Size:", m.Width, m.Height)
	return m.WithinBounds(cX, cY)
}

func (m *Modal) WithinBounds(mouseX, mouseY int) bool {
	return mouseX >= m.X && mouseX <= m.X+m.Width && mouseY >= m.Y && mouseY <= m.Y+m.Height
}

func (m *Modal) GetName() string {
	return m.Name
}

func (m *Modal) GetPosition() (int, int) {
	return m.X, m.Y
}

func (m *Modal) SetPosition(x, y int) {
	m.X = x
	m.Y = y
}

func (m *Modal) OpenModal() {
	m.Visible = true
}

func (m *Modal) CloseModal() {
	m.Visible = false
	if m.OnClose != nil {
		m.OnClose()
	}
}

func (m *Modal) IsOpen() bool {
	return m.Visible
}

func (m *Modal) IsVisible() bool {
	return m.Visible
}
