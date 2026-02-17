package minui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/text"
)

// TooltipPosition defines where the tooltip appears relative to target
type TooltipPosition int

const (
	TooltipAbove TooltipPosition = iota
	TooltipBelow
	TooltipLeft
	TooltipRight
	TooltipMouse // Follow mouse cursor
)

// Tooltip displays contextual information when hovering over elements
type Tooltip struct {
	*ElementBase
	Title    string
	Text     string
	Icon     *Icon
	Position TooltipPosition
	Delay    int // Frames before showing
	Offset   int // Pixels offset from target

	targetElement  Element
	hoverFrames    int
	showing        bool
	mouseX, mouseY int
}

// NewTooltip creates a new tooltip
func NewTooltip(id string) *Tooltip {
	t := &Tooltip{
		ElementBase: NewElementBase(id),
		Position:    TooltipBelow,
		Delay:       30, // ~0.5 seconds at 60fps
		Offset:      8,
	}

	t.visible = false

	// Set default style - only structural properties, colors come from theme
	borderWidth := 1
	t.style.BorderWidth = &borderWidth
	padding := NewEdgeInsetsTRBL(6, 8, 6, 8)
	t.style.Padding = padding

	return t
}

// SetContent sets the tooltip content
func (t *Tooltip) SetContent(title, text string, icon *Icon) {
	t.Title = title
	t.Text = text
	t.Icon = icon
}

// SetTarget sets the element this tooltip is attached to
func (t *Tooltip) SetTarget(element Element) {
	t.targetElement = element
}

// Show immediately shows the tooltip
func (t *Tooltip) Show() {
	t.showing = true
	t.visible = true
	t.updatePosition()
}

// Hide immediately hides the tooltip
func (t *Tooltip) Hide() {
	t.showing = false
	t.visible = false
	t.hoverFrames = 0
}

// GetType returns the element type
func (t *Tooltip) GetType() string {
	return "Tooltip"
}

// updatePosition calculates tooltip position based on target
func (t *Tooltip) updatePosition() {
	if t.targetElement == nil && t.Position != TooltipMouse {
		return
	}

	// Calculate tooltip size first
	t.Layout()

	var targetX, targetY, targetW, targetH int

	if t.Position == TooltipMouse {
		targetX = t.mouseX
		targetY = t.mouseY
		targetW = 1
		targetH = 1
	} else {
		targetBounds := t.targetElement.GetBounds()
		tpx, tpy := t.targetElement.GetAbsolutePosition()
		targetX = tpx
		targetY = tpy
		targetW = targetBounds.Width
		targetH = targetBounds.Height
	}

	var x, y int
	switch t.Position {
	case TooltipAbove:
		x = targetX + (targetW-t.bounds.Width)/2
		y = targetY - t.bounds.Height - t.Offset
	case TooltipBelow:
		x = targetX + (targetW-t.bounds.Width)/2
		y = targetY + targetH + t.Offset
	case TooltipLeft:
		x = targetX - t.bounds.Width - t.Offset
		y = targetY + (targetH-t.bounds.Height)/2
	case TooltipRight:
		x = targetX + targetW + t.Offset
		y = targetY + (targetH-t.bounds.Height)/2
	case TooltipMouse:
		x = targetX + 16
		y = targetY + 16
	}

	// Keep tooltip on screen
	screenW, screenH := ebiten.WindowSize()
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x+t.bounds.Width > screenW {
		x = screenW - t.bounds.Width
	}
	if y+t.bounds.Height > screenH {
		y = screenH - t.bounds.Height
	}

	t.SetPosition(x, y)
}

// Update updates the tooltip state
func (t *Tooltip) Update() {
	mx, my := ebiten.CursorPosition()
	t.mouseX = mx
	t.mouseY = my

	// Check if target is hovered
	if t.targetElement != nil {
		targetBounds := t.targetElement.GetBounds()
		tpx, tpy := t.targetElement.GetAbsolutePosition()

		isHovering := mx >= tpx && mx < tpx+targetBounds.Width &&
			my >= tpy && my < tpy+targetBounds.Height

		if isHovering {
			t.hoverFrames++
			if t.hoverFrames >= t.Delay && !t.showing {
				t.Show()
			}
			if t.showing {
				t.updatePosition()
			}
		} else {
			t.Hide()
		}
	} else if t.Position == TooltipMouse && t.showing {
		// Follow mouse
		t.updatePosition()
	}
}

// Layout calculates dimensions
func (t *Tooltip) Layout() {
	style := t.GetComputedStyle()

	fontSize := 14
	if style.FontSize != nil {
		fontSize = *style.FontSize
	}
	titleFontSize := fontSize + 2

	paddingH := 0
	paddingV := 0
	if style.Padding != nil {
		paddingH = style.Padding.Left + style.Padding.Right
		paddingV = style.Padding.Top + style.Padding.Bottom
	}

	// Calculate content size
	contentWidth := 0
	contentHeight := 0

	// Icon
	if t.Icon != nil {
		contentWidth += t.Icon.ScaledWidth() + 8
	}

	// Title
	if t.Title != "" {
		titleWidth := len(t.Title) * (titleFontSize * 6 / 10)
		if titleWidth > contentWidth {
			contentWidth = titleWidth
		}
		contentHeight += titleFontSize + 4
	}

	// Text
	if t.Text != "" {
		textWidth := len(t.Text) * (fontSize * 6 / 10)
		if textWidth > contentWidth {
			contentWidth = textWidth
		}
		contentHeight += fontSize
	}

	t.bounds.Width = contentWidth + paddingH
	t.bounds.Height = contentHeight + paddingV
}

// Draw draws the tooltip
func (t *Tooltip) Draw(screen *ebiten.Image) {
	if !t.visible || !t.showing {
		return
	}

	style := t.GetComputedStyle()
	theme := t.GetTheme()
	absX, absY := t.GetAbsolutePosition()
	absBounds := Rect{
		X:      absX,
		Y:      absY,
		Width:  t.bounds.Width,
		Height: t.bounds.Height,
	}

	// Draw background with theme support
	DrawBackgroundWithTheme(screen, absBounds, style, theme)

	contentBounds := GetContentBounds(absBounds, style)

	fontSize := 14
	if style.FontSize != nil {
		fontSize = *style.FontSize
	}
	titleFontSize := fontSize + 2

	// Get text color from style, then theme, then default
	textColor := color.RGBA{255, 255, 255, 255}
	if style.ForegroundColor != nil {
		textColor = colorToRGBA(*style.ForegroundColor)
	} else if theme != nil {
		textColor = colorToRGBA(theme.Colors.Text)
	}

	// Get title color from theme (secondary text) or use warning color for visibility
	titleColor := color.RGBA{255, 220, 120, 255} // Gold for title
	if theme != nil {
		titleColor = colorToRGBA(theme.Colors.Warning)
	}

	x := contentBounds.X
	y := contentBounds.Y

	// Draw icon
	if t.Icon != nil {
		iconY := y + (contentBounds.Height-t.Icon.ScaledHeight())/2
		t.Icon.Draw(screen, x, iconY)
		x += t.Icon.ScaledWidth() + 8
	}

	// Draw title
	if t.Title != "" {
		text.Draw(screen, t.Title, float64(titleFontSize), x, y, titleColor)
		y += titleFontSize + 4
	}

	// Draw text
	if t.Text != "" {
		text.Draw(screen, t.Text, float64(fontSize), x, y, textColor)
	}

	// Draw border with theme support
	DrawBorderWithTheme(screen, absBounds, style, theme)
}
