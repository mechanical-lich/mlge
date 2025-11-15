package minui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mechanical-lich/mlge/event"
	"github.com/mechanical-lich/mlge/text/v2"
)

// Button is a clickable button element
type Button struct {
	*ElementBase
	Text    string
	OnClick func()
	pressed bool
}

// NewButton creates a new button
func NewButton(id, text string) *Button {
	button := &Button{
		ElementBase: NewElementBase(id),
		Text:        text,
	}

	// Set default size
	button.SetSize(len(text)*10+20, 32)

	// Set default button style
	bgColor := color.Color(color.RGBA{90, 90, 100, 255})
	borderColor := color.Color(color.RGBA{120, 120, 130, 255})
	borderWidth := 2
	borderRadius := 4
	padding := NewEdgeInsets(8)

	button.style.BackgroundColor = &bgColor
	button.style.BorderColor = &borderColor
	button.style.BorderWidth = &borderWidth
	button.style.BorderRadius = &borderRadius
	button.style.Padding = padding

	// Hover style
	hoverBg := color.Color(color.RGBA{110, 110, 120, 255})
	button.style.HoverStyle = &Style{
		BackgroundColor: &hoverBg,
	}

	// Active/pressed style
	activeBg := color.Color(color.RGBA{70, 70, 80, 255})
	button.style.ActiveStyle = &Style{
		BackgroundColor: &activeBg,
	}

	return button
}

// GetType returns the element type
func (b *Button) GetType() string {
	return "Button"
}

// Update updates the button state
func (b *Button) Update() {
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
			event.GetQueuedInstance().QueueEvent(ButtonClickEvent{
				ButtonID: b.GetID(),
				Button:   b,
			})
		}
		b.pressed = false
	}
}

// Layout calculates button dimensions
func (b *Button) Layout() {
	style := b.GetComputedStyle()

	fontSize := 14
	if style.FontSize != nil {
		fontSize = *style.FontSize
	}

	// Calculate text size
	textWidth := len(b.Text) * fontSize * 6 / 10
	textHeight := fontSize + 6

	// Add padding
	if style.Padding != nil {
		textWidth += style.Padding.Left + style.Padding.Right
		textHeight += style.Padding.Top + style.Padding.Bottom
	}

	// Add border
	if style.BorderWidth != nil {
		textWidth += *style.BorderWidth * 2
		textHeight += *style.BorderWidth * 2
	}

	// Apply width/height constraints
	if style.Width != nil {
		textWidth = *style.Width
	}
	if style.Height != nil {
		textHeight = *style.Height
	}

	// Apply min/max size constraints
	textWidth, textHeight = ApplySizeConstraints(textWidth, textHeight, style)

	b.bounds.Width = textWidth
	b.bounds.Height = textHeight
}

// Draw draws the button
func (b *Button) Draw(screen *ebiten.Image) {
	if !b.visible {
		return
	}

	// Get style based on state
	style := b.GetComputedStyle()
	if b.pressed {
		if style.ActiveStyle != nil {
			style = style.ActiveStyle.Merge(style)
		}
	}

	// Get absolute position for drawing
	absX, absY := b.GetAbsolutePosition()
	absBounds := Rect{
		X:      absX,
		Y:      absY,
		Width:  b.bounds.Width,
		Height: b.bounds.Height,
	}

	// Draw background
	DrawBackground(screen, absBounds, style)

	// Draw border
	DrawBorder(screen, absBounds, style)

	// Draw text
	contentBounds := GetContentBounds(absBounds, style)

	fontSize := 14
	if style.FontSize != nil {
		fontSize = *style.FontSize
	}

	textColor := color.RGBA{255, 255, 255, 255}
	if style.ForegroundColor != nil {
		r, g, b, a := (*style.ForegroundColor).RGBA()
		textColor = color.RGBA{
			R: uint8(r >> 8),
			G: uint8(g >> 8),
			B: uint8(b >> 8),
			A: uint8(a >> 8),
		}
	}

	// Center text in button
	textWidth := len(b.Text) * fontSize * 6 / 10
	textX := contentBounds.X + (contentBounds.Width-textWidth)/2
	textY := contentBounds.Y + (contentBounds.Height-fontSize)/2

	text.Draw(screen, b.Text, float64(fontSize), textX, textY, textColor)
}

// IsPressed returns if the button is currently pressed
func (b *Button) IsPressed() bool {
	return b.pressed
}

// IsJustClicked returns true if the button was just clicked this frame
func (b *Button) IsJustClicked() bool {
	if !b.visible || !b.enabled {
		return false
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		return b.IsWithin(mx, my)
	}
	return false
}
