package minui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mechanical-lich/mlge/text/v2"
)

// Modal is a dialog/modal window container
type Modal struct {
	*Panel
	Title     string
	Closeable bool
	OnClose   func()

	dragging    bool
	dragOffsetX int
	dragOffsetY int

	titleBarHeight int
}

// NewModal creates a new modal dialog
func NewModal(id, title string, width, height int) *Modal {
	modal := &Modal{
		Panel:          NewPanel(id),
		Title:          title,
		Closeable:      true,
		titleBarHeight: 30,
	}

	modal.SetSize(width, height)

	// Set default modal style
	bgColor := color.Color(color.RGBA{240, 240, 245, 255})
	borderColor := color.Color(color.RGBA{100, 100, 110, 255})
	borderWidth := 2
	borderRadius := 6

	modal.style.BackgroundColor = &bgColor
	modal.style.BorderColor = &borderColor
	modal.style.BorderWidth = &borderWidth
	modal.style.BorderRadius = &borderRadius

	return modal
}

// GetType returns the element type
func (m *Modal) GetType() string {
	return "Modal"
}

// AddChild adds a child element to the modal
func (m *Modal) AddChild(child Element) {
	m.children = append(m.children, child)
	child.SetParent(m) // m is the Modal, which implements Element through Panel
}

// Update updates the modal
func (m *Modal) Update() {
	if !m.visible {
		return
	}

	m.UpdateHoverState()

	// Get absolute position for hit detection
	absX, absY := m.GetAbsolutePosition()

	// Handle dragging from title bar
	titleBarBounds := Rect{
		X:      absX,
		Y:      absY,
		Width:  m.bounds.Width,
		Height: m.titleBarHeight,
	}

	mx, my := ebiten.CursorPosition()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if titleBarBounds.Contains(mx, my) {
			m.dragging = true
			m.dragOffsetX = mx - m.bounds.X
			m.dragOffsetY = my - m.bounds.Y
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && m.dragging {
		m.bounds.X = mx - m.dragOffsetX
		m.bounds.Y = my - m.dragOffsetY
	} else {
		m.dragging = false
	}

	// Handle close button
	if m.Closeable {
		closeBtnBounds := Rect{
			X:      absX + m.bounds.Width - 28,
			Y:      absY + 4,
			Width:  24,
			Height: 22,
		}

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if closeBtnBounds.Contains(mx, my) {
				if m.OnClose != nil {
					m.OnClose()
				}
				m.visible = false
				return
			}
		}
	}

	// Update children
	for _, child := range m.children {
		child.Update()
	}
}

// Layout calculates modal layout
func (m *Modal) Layout() {
	if !m.visible {
		return
	}

	style := m.GetComputedStyle()

	// Content area is below title bar
	contentBounds := GetContentBounds(m.bounds, style)
	contentBounds.Y += m.titleBarHeight
	contentBounds.Height -= m.titleBarHeight

	// The children should already have their positions set relative to modal's content area
	// We just need to ensure they're offset by the modal's position
	// Children's X,Y are relative to the modal's content area (0,0) at the content start

	// Just layout the children - their positions are already relative to content area
	for _, child := range m.children {
		child.Layout()
	}
}

// Draw draws the modal
func (m *Modal) Draw(screen *ebiten.Image) {
	if !m.visible {
		return
	}

	style := m.GetComputedStyle()

	// Draw semi-transparent overlay behind modal
	overlayColor := color.RGBA{0, 0, 0, 128}
	overlayBounds := Rect{X: 0, Y: 0, Width: screen.Bounds().Dx(), Height: screen.Bounds().Dy()}
	DrawRect(screen, overlayBounds, overlayColor)

	// Get absolute position
	absX, absY := m.GetAbsolutePosition()
	absBounds := Rect{
		X:      absX,
		Y:      absY,
		Width:  m.bounds.Width,
		Height: m.bounds.Height,
	}

	// Draw modal background
	DrawBackground(screen, absBounds, style)

	// Draw title bar
	titleBarBounds := Rect{
		X:      absX,
		Y:      absY,
		Width:  m.bounds.Width,
		Height: m.titleBarHeight,
	}

	titleBarColor := color.RGBA{100, 120, 180, 255}
	borderRadius := 0
	if style.BorderRadius != nil {
		borderRadius = *style.BorderRadius
	}

	if borderRadius > 0 {
		// Draw rounded top
		DrawRoundedRect(screen, titleBarBounds, borderRadius, titleBarColor)
		// Cover bottom with rectangle
		bottomRect := Rect{
			X:      titleBarBounds.X,
			Y:      titleBarBounds.Y + borderRadius,
			Width:  titleBarBounds.Width,
			Height: titleBarBounds.Height - borderRadius,
		}
		DrawRect(screen, bottomRect, titleBarColor)
	} else {
		DrawRect(screen, titleBarBounds, titleBarColor)
	}

	// Draw title text
	titleColor := color.RGBA{255, 255, 255, 255}
	text.Draw(screen, m.Title, 16.0, absX+8, absY+6, titleColor)

	// Draw close button
	if m.Closeable {
		closeBtnBounds := Rect{
			X:      absX + m.bounds.Width - 28,
			Y:      absY + 4,
			Width:  24,
			Height: 22,
		}

		closeBtnColor := color.RGBA{200, 80, 80, 255}
		mx, my := ebiten.CursorPosition()
		if closeBtnBounds.Contains(mx, my) {
			closeBtnColor = color.RGBA{220, 100, 100, 255}
		}

		DrawRoundedRect(screen, closeBtnBounds, 3, closeBtnColor)

		// Draw X
		xColor := color.RGBA{255, 255, 255, 255}
		text.Draw(screen, "X", 14.0, closeBtnBounds.X+7, closeBtnBounds.Y+3, xColor)
	}

	// Draw children
	for _, child := range m.children {
		child.Draw(screen)
	}

	// Draw border
	DrawBorder(screen, absBounds, style)
}
