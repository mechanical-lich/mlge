package minui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mechanical-lich/mlge/event"
	"github.com/mechanical-lich/mlge/text"
)

// IconButton is a clickable button with an icon
type IconButton struct {
	*ElementBase
	Text         string
	Icon         *Icon
	IconPosition IconPosition
	IconSpacing  int
	OnClick      func()
	pressed      bool
}

// NewIconButton creates a new button with an icon
func NewIconButton(id string, icon *Icon, text string) *IconButton {
	button := &IconButton{
		ElementBase:  NewElementBase(id),
		Text:         text,
		Icon:         icon,
		IconPosition: IconLeft,
		IconSpacing:  4,
	}

	// Calculate initial size
	button.updateSize()

	// Set default button style - only structural properties, colors come from theme
	borderWidth := 2
	borderRadius := 4
	padding := NewEdgeInsets(8)

	button.style.BorderWidth = &borderWidth
	button.style.BorderRadius = &borderRadius
	button.style.Padding = padding

	return button
}

// NewIconOnlyButton creates a button with just an icon
func NewIconOnlyButton(id string, icon *Icon) *IconButton {
	button := &IconButton{
		ElementBase:  NewElementBase(id),
		Text:         "",
		Icon:         icon,
		IconPosition: IconOnly,
		IconSpacing:  0,
	}

	button.updateSize()

	// Set default style - only structural properties, colors come from theme
	borderWidth := 2
	borderRadius := 4
	padding := NewEdgeInsets(4)

	button.style.BorderWidth = &borderWidth
	button.style.BorderRadius = &borderRadius
	button.style.Padding = padding

	return button
}

func (b *IconButton) updateSize() {
	iconW := 0
	iconH := 0
	if b.Icon != nil {
		iconW = b.Icon.ScaledWidth()
		iconH = b.Icon.ScaledHeight()
	}

	textW := len(b.Text) * 8
	textH := 16

	var width, height int

	switch b.IconPosition {
	case IconLeft, IconRight:
		if b.Icon != nil && b.Text != "" {
			width = iconW + b.IconSpacing + textW
			height = max(iconH, textH)
		} else if b.Icon != nil {
			width = iconW
			height = iconH
		} else {
			width = textW
			height = textH
		}
	case IconTop, IconBottom:
		if b.Icon != nil && b.Text != "" {
			width = max(iconW, textW)
			height = iconH + b.IconSpacing + textH
		} else if b.Icon != nil {
			width = iconW
			height = iconH
		} else {
			width = textW
			height = textH
		}
	case IconOnly:
		width = iconW
		height = iconH
	}

	// Add default padding
	b.SetSize(width+16, height+8)
}

// SetIcon sets the button's icon
func (b *IconButton) SetIcon(icon *Icon) {
	b.Icon = icon
	b.updateSize()
}

// SetText sets the button's text
func (b *IconButton) SetText(text string) {
	b.Text = text
	b.updateSize()
}

// GetType returns the element type
func (b *IconButton) GetType() string {
	return "IconButton"
}

// Update handles button interaction
func (b *IconButton) Update() {
	if !b.visible || !b.enabled {
		return
	}

	b.UpdateHoverState()

	// Check for click
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if b.hovered {
			b.pressed = true
		}
	} else {
		if b.pressed && b.hovered {
			// Button was clicked
			if b.OnClick != nil {
				b.OnClick()
			}
			// Fire event
			event.GetQueuedInstance().QueueEvent(IconButtonClickEvent{
				ButtonID: b.GetID(),
				Button:   b,
			})
		}
		b.pressed = false
	}
}

// Layout calculates button dimensions
func (b *IconButton) Layout() {
	style := b.GetComputedStyle()

	fontSize := 14
	if style.FontSize != nil {
		fontSize = *style.FontSize
	}

	iconW := 0
	iconH := 0
	if b.Icon != nil {
		iconW = b.Icon.ScaledWidth()
		iconH = b.Icon.ScaledHeight()
	}

	textW := len(b.Text) * fontSize * 6 / 10
	textH := fontSize + 4

	var contentWidth, contentHeight int

	switch b.IconPosition {
	case IconLeft, IconRight:
		if b.Icon != nil && b.Text != "" {
			contentWidth = iconW + b.IconSpacing + textW
			contentHeight = max(iconH, textH)
		} else if b.Icon != nil {
			contentWidth = iconW
			contentHeight = iconH
		} else {
			contentWidth = textW
			contentHeight = textH
		}
	case IconTop, IconBottom:
		if b.Icon != nil && b.Text != "" {
			contentWidth = max(iconW, textW)
			contentHeight = iconH + b.IconSpacing + textH
		} else if b.Icon != nil {
			contentWidth = iconW
			contentHeight = iconH
		} else {
			contentWidth = textW
			contentHeight = textH
		}
	case IconOnly:
		contentWidth = iconW
		contentHeight = iconH
	}

	width := contentWidth
	height := contentHeight

	// Add padding
	if style.Padding != nil {
		width += style.Padding.Left + style.Padding.Right
		height += style.Padding.Top + style.Padding.Bottom
	}

	// Add border
	if style.BorderWidth != nil {
		width += *style.BorderWidth * 2
		height += *style.BorderWidth * 2
	}

	// Apply explicit dimensions from style
	if style.Width != nil {
		width = *style.Width
	}
	if style.Height != nil {
		height = *style.Height
	}

	// Apply constraints
	width, height = ApplySizeConstraints(width, height, style)

	b.bounds.Width = width
	b.bounds.Height = height
}

// Draw draws the button
func (b *IconButton) Draw(screen *ebiten.Image) {
	if !b.visible {
		return
	}

	style := b.GetComputedStyle()
	theme := b.GetTheme()

	if b.pressed {
		if style.ActiveStyle != nil {
			style = style.ActiveStyle.Merge(style)
		}
	}

	absX, absY := b.GetAbsolutePosition()
	absBounds := Rect{
		X:      absX,
		Y:      absY,
		Width:  b.bounds.Width,
		Height: b.bounds.Height,
	}

	// Draw background
	if theme != nil && theme.HasButtonSprites() {
		var coords *SpriteCoords
		if b.pressed && theme.ButtonPressed != nil {
			coords = theme.ButtonPressed
		} else {
			coords = theme.Button
		}
		DrawSprite(screen, theme.SpriteSheet, coords, absBounds)
	} else {
		// Use vector-based rendering with theme support
		DrawBackgroundWithTheme(screen, absBounds, style, theme)
		DrawBorderWithTheme(screen, absBounds, style, theme)
	}

	// Draw content
	contentBounds := GetContentBounds(absBounds, style)

	fontSize := 14
	if style.FontSize != nil {
		fontSize = *style.FontSize
	}

	// Get text color from style, then theme, then default
	textColor := color.RGBA{255, 255, 255, 255}
	if style.ForegroundColor != nil {
		textColor = colorToRGBA(*style.ForegroundColor)
	} else if theme != nil {
		textColor = colorToRGBA(theme.Colors.Text)
	}

	iconW := 0
	iconH := 0
	if b.Icon != nil {
		iconW = b.Icon.ScaledWidth()
		iconH = b.Icon.ScaledHeight()
	}

	textH := fontSize

	switch b.IconPosition {
	case IconLeft:
		// Calculate total content width for centering
		textW := len(b.Text) * fontSize * 6 / 10
		totalW := iconW
		if b.Text != "" {
			totalW += b.IconSpacing + textW
		}
		startX := contentBounds.X + (contentBounds.Width-totalW)/2

		if b.Icon != nil {
			iconY := contentBounds.Y + (contentBounds.Height-iconH)/2
			b.Icon.Draw(screen, startX, iconY)
		}
		if b.Text != "" {
			textX := startX
			if b.Icon != nil {
				textX += iconW + b.IconSpacing
			}
			textY := contentBounds.Y + (contentBounds.Height-textH)/2
			text.Draw(screen, b.Text, float64(fontSize), textX, textY, textColor)
		}

	case IconRight:
		textW := len(b.Text) * fontSize * 6 / 10
		totalW := textW
		if b.Icon != nil {
			totalW += b.IconSpacing + iconW
		}
		startX := contentBounds.X + (contentBounds.Width-totalW)/2

		if b.Text != "" {
			textY := contentBounds.Y + (contentBounds.Height-textH)/2
			text.Draw(screen, b.Text, float64(fontSize), startX, textY, textColor)
		}
		if b.Icon != nil {
			iconX := startX + textW
			if b.Text != "" {
				iconX += b.IconSpacing
			}
			iconY := contentBounds.Y + (contentBounds.Height-iconH)/2
			b.Icon.Draw(screen, iconX, iconY)
		}

	case IconTop:
		textW := len(b.Text) * fontSize * 6 / 10
		totalH := iconH
		if b.Text != "" {
			totalH += b.IconSpacing + textH
		}
		startY := contentBounds.Y + (contentBounds.Height-totalH)/2

		if b.Icon != nil {
			iconX := contentBounds.X + (contentBounds.Width-iconW)/2
			b.Icon.Draw(screen, iconX, startY)
		}
		if b.Text != "" {
			textX := contentBounds.X + (contentBounds.Width-textW)/2
			textY := startY + iconH + b.IconSpacing
			text.Draw(screen, b.Text, float64(fontSize), textX, textY, textColor)
		}

	case IconBottom:
		textW := len(b.Text) * fontSize * 6 / 10
		totalH := textH
		if b.Icon != nil {
			totalH += b.IconSpacing + iconH
		}
		startY := contentBounds.Y + (contentBounds.Height-totalH)/2

		if b.Text != "" {
			textX := contentBounds.X + (contentBounds.Width-textW)/2
			text.Draw(screen, b.Text, float64(fontSize), textX, startY, textColor)
		}
		if b.Icon != nil {
			iconX := contentBounds.X + (contentBounds.Width-iconW)/2
			iconY := startY + textH + b.IconSpacing
			b.Icon.Draw(screen, iconX, iconY)
		}

	case IconOnly:
		if b.Icon != nil {
			iconX := contentBounds.X + (contentBounds.Width-iconW)/2
			iconY := contentBounds.Y + (contentBounds.Height-iconH)/2
			b.Icon.Draw(screen, iconX, iconY)
		}
	}
}

// IsPressed returns if the button is currently pressed
func (b *IconButton) IsPressed() bool {
	return b.pressed
}

// IsJustClicked returns true if the button was just clicked this frame
func (b *IconButton) IsJustClicked() bool {
	if !b.visible || !b.enabled {
		return false
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		return b.IsWithin(mx, my)
	}
	return false
}

// IconButtonClickEvent is fired when an icon button is clicked
type IconButtonClickEvent struct {
	ButtonID string
	Button   *IconButton
}

func (e IconButtonClickEvent) GetType() event.EventType {
	return EventTypeIconButtonClick
}
