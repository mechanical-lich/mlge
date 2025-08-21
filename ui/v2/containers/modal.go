package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	elements "github.com/mechanical-lich/mlge/ui/v2/elements"
	theming "github.com/mechanical-lich/mlge/ui/v2/theming"
	"github.com/mechanical-lich/mlge/utility"
)

type ModalInterface interface {
	Update()
	Draw(screen *ebiten.Image, theme *theming.Theme)
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
	elements.ElementBase

	CloseButton *elements.Button
	OnClose     func()

	dragging    bool
	dragOffsetX int
	dragOffsetY int
	bg          *ebiten.Image // Background image for the modal
	offscreen   *ebiten.Image // Offscreen buffer for double buffering

	Children []elements.ElementInterface // Children elements
}

// NewModal creates a new modal with initial view.
func NewModal(name string, x, y, width, height int) *Modal {
	closeBtn := elements.NewButton("close", 0, 0, "X", "close")
	m := &Modal{
		ElementBase: elements.ElementBase{
			Name:    name,
			X:       x,
			Y:       y,
			Width:   width,
			Height:  height,
			Visible: false,
			Op:      &ebiten.DrawImageOptions{},
		},
		CloseButton: closeBtn,
	}

	closeBtn.Parent = m
	return m
}

// Update handles modal logic and delegates to the current view.
func (m *Modal) Update() {
	if !m.Visible {
		return
	}

	m.handleDragging()

	if m.CloseButton.IsJustClicked() && !m.dragging {
		m.Visible = false
		if m.OnClose != nil {
			m.OnClose()
		}
		return
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
		!m.CloseButton.IsClicked() {
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

// AddChild adds a child element to the modal.
func (m *Modal) AddChild(child elements.ElementInterface) {
	child.SetParent(m)
	m.Children = append(m.Children, child)
}

// RemoveChild removes a child element from the modal.
func (m *Modal) RemoveChild(child elements.ElementInterface) {
	for i, c := range m.Children {
		if c == child {
			m.Children = append(m.Children[:i], m.Children[i+1:]...)
			child.SetParent(nil)
			break
		}
	}
}

// Draw renders the modal background, close button, and children.
func (m *Modal) Draw(screen *ebiten.Image, theme *theming.Theme) {
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
	m.CloseButton.Draw(m.offscreen, theme)

	// Draw children
	for _, child := range m.Children {
		child.Draw(m.offscreen, theme)
	}

	// Draw offscreen buffer to screen at modal position
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(m.X), float64(m.Y))
	screen.DrawImage(m.offscreen, op)
}

// GetInputFocused delegates to the current view.
func (m *Modal) GetInputFocused() bool {
	for _, child := range m.Children {
		if child.GetFocused() {
			return true
		}
	}
	return false
}

func (m *Modal) GetMouseFocused() bool {

	cX, cY := ebiten.CursorPosition()
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

func (m *Modal) GetAbsolutePosition() (int, int) {
	return 0, 0
}
