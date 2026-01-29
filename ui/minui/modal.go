package minui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mechanical-lich/mlge/event"
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
	initialWidth   int // Store initial width as minimum
	initialHeight  int // Store initial height as minimum

	// Cached background for sprite-based rendering
	bgCache       *ebiten.Image
	bgCacheWidth  int
	bgCacheHeight int
}

// NewModal creates a new modal dialog
func NewModal(id, title string, width, height int) *Modal {
	modal := &Modal{
		Panel:          NewPanel(id),
		Title:          title,
		Closeable:      true,
		titleBarHeight: 30,
		initialWidth:   width,
		initialHeight:  height,
	}

	modal.SetSize(width, height)

	// Set default modal style - only border properties, colors come from theme
	borderWidth := 2
	borderRadius := 6

	modal.style.BorderWidth = &borderWidth
	modal.style.BorderRadius = &borderRadius

	// Default opaque background for modal (prevent transparent interior)
	bgColor := color.Color(color.RGBA{30, 30, 35, 255})
	modal.style.BackgroundColor = &bgColor

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

	// Propagate theme to child if we have one
	if m.theme != nil {
		if setter, ok := child.(interface{ SetTheme(*Theme) }); ok {
			setter.SetTheme(m.theme)
		}
	}
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
				// Fire event
				event.GetQueuedInstance().QueueEvent(ModalCloseEvent{
					ModalID: m.GetID(),
					Modal:   m,
				})
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

	// First, layout all children to get their sizes
	for _, child := range m.children {
		child.Layout()
	}

	// Calculate required size based on children
	contentWidth, contentHeight := m.calculateContentSize()

	// Add padding to content size
	if style.Padding != nil {
		contentWidth += style.Padding.Left + style.Padding.Right
		contentHeight += style.Padding.Top + style.Padding.Bottom
	}

	// Add border width
	if style.BorderWidth != nil {
		contentWidth += *style.BorderWidth * 2
		contentHeight += *style.BorderWidth * 2
	}

	// Add title bar height
	contentHeight += m.titleBarHeight

	// Start with the larger of content size or initial size
	width := m.initialWidth
	height := m.initialHeight

	if contentWidth > width {
		width = contentWidth
	}
	if contentHeight > height {
		height = contentHeight
	}

	// Apply explicit width/height from style if specified (overrides auto-sizing)
	if style.Width != nil {
		width = *style.Width
	}
	if style.Height != nil {
		height = *style.Height
	}

	// Apply min size constraints (use initial size as minimum if no explicit min set)
	if style.MinWidth == nil {
		minWidth := m.initialWidth
		style.MinWidth = &minWidth
	}
	if style.MinHeight == nil {
		minHeight := m.initialHeight
		style.MinHeight = &minHeight
	}

	// Apply min/max size constraints
	width, height = ApplySizeConstraints(width, height, style)

	m.bounds.Width = width
	m.bounds.Height = height
}

// calculateContentSize calculates the required size to fit all children
func (m *Modal) calculateContentSize() (int, int) {
	if len(m.children) == 0 {
		return 0, 0
	}

	maxRight := 0
	maxBottom := 0

	for _, child := range m.children {
		childBounds := child.GetBounds()
		childStyle := child.GetComputedStyle()

		// Start with the child's position and size
		right := childBounds.X + childBounds.Width
		bottom := childBounds.Y + childBounds.Height

		// Add child's padding (already included in child bounds via GetContentBounds, but borders/margins are not)
		// Note: padding is already accounted for in the child's width/height from Layout()

		// Add child's border width
		if childStyle.BorderWidth != nil {
			// Borders are typically already included in the child's Layout calculation
			// but we should verify the child's right edge includes it
		}

		// Add child's margin
		if childStyle.Margin != nil {
			right += childStyle.Margin.Right
			bottom += childStyle.Margin.Bottom
		}

		if right > maxRight {
			maxRight = right
		}
		if bottom > maxBottom {
			maxBottom = bottom
		}
	}

	// Add some default padding to ensure content isn't flush against edges
	const defaultContentPadding = 0 // Don't add extra - rely on modal's padding

	return maxRight + defaultContentPadding, maxBottom + defaultContentPadding
}

// Draw draws the modal
func (m *Modal) Draw(screen *ebiten.Image) {
	if !m.visible {
		return
	}

	style := m.GetComputedStyle()
	theme := m.GetTheme()

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

	// Check if we should use sprite-based rendering
	if theme != nil && theme.HasModalSprites() {
		// Use 9-slice sprite rendering
		// Cache the background if needed
		if m.bgCache == nil || m.bgCacheWidth != m.bounds.Width || m.bgCacheHeight != m.bounds.Height {
			m.bgCache = ebiten.NewImage(m.bounds.Width, m.bounds.Height)
			Draw9SliceToImage(m.bgCache, theme.SpriteSheet, theme.ModalNineSlice)
			m.bgCacheWidth = m.bounds.Width
			m.bgCacheHeight = m.bounds.Height
		}

		// Draw cached background
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(absX), float64(absY))
		screen.DrawImage(m.bgCache, op)

		// Draw title text
		titleColor := color.RGBA{255, 255, 255, 255}
		text.Draw(screen, m.Title, 16.0, absX+8, absY+6, titleColor)

		// Draw close button (use sprite if available)
		if m.Closeable {
			m.drawCloseButton(screen, absX, absY, theme)
		}
	} else {
		// Use vector-based rendering with theme support
		DrawBackgroundWithTheme(screen, absBounds, style, theme)

		// Draw title bar
		titleBarBounds := Rect{
			X:      absX,
			Y:      absY,
			Width:  m.bounds.Width,
			Height: m.titleBarHeight,
		}

		// Use theme primary color for title bar, fallback to default
		titleBarColor := color.RGBA{100, 120, 180, 255}
		if theme != nil {
			titleBarColor = colorToRGBA(theme.Colors.Primary)
		}

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

		// Draw title text - use contrasting color based on theme
		titleColor := color.RGBA{255, 255, 255, 255}
		text.Draw(screen, m.Title, 16.0, absX+8, absY+6, titleColor)

		// Draw close button
		if m.Closeable {
			m.drawCloseButton(screen, absX, absY, theme)
		}

		// Draw border with theme support
		DrawBorderWithTheme(screen, absBounds, style, theme)
	}

	// Draw children
	for _, child := range m.children {
		child.Draw(screen)
	}
}

// drawCloseButton draws the close button, using sprites if available
func (m *Modal) drawCloseButton(screen *ebiten.Image, absX, absY int, theme *Theme) {
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
