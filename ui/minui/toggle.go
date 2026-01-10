package minui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mechanical-lich/mlge/event"
	"github.com/mechanical-lich/mlge/text/v2"
)

// ToggleChangeEvent is fired when a toggle state changes
type ToggleChangeEvent struct {
	ToggleID string
	Toggle   *Toggle
	On       bool
}

func (e ToggleChangeEvent) GetType() event.EventType {
	return EventTypeToggleChange
}

// Toggle is a toggle button element (on/off switch)
type Toggle struct {
	*ElementBase
	Text     string
	On       bool
	OnChange func(on bool)
}

// NewToggle creates a new toggle button
func NewToggle(id, text string) *Toggle {
	toggle := &Toggle{
		ElementBase: NewElementBase(id),
		Text:        text,
		On:          false,
	}

	// Set default size
	toggle.SetSize(len(text)*10+40, 28)

	// Set default style - only structural properties, colors come from theme
	borderWidth := 2
	borderRadius := 4
	padding := NewEdgeInsets(6)

	toggle.style.BorderWidth = &borderWidth
	toggle.style.BorderRadius = &borderRadius
	toggle.style.Padding = padding

	return toggle
}

// GetType returns the element type
func (t *Toggle) GetType() string {
	return "Toggle"
}

// Update handles toggle interaction
func (t *Toggle) Update() {
	if !t.visible || !t.enabled {
		return
	}

	t.UpdateHoverState()

	// Handle click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		absX, absY := t.GetAbsolutePosition()
		absBounds := Rect{
			X:      absX,
			Y:      absY,
			Width:  t.bounds.Width,
			Height: t.bounds.Height,
		}

		mx, my := ebiten.CursorPosition()
		if absBounds.Contains(mx, my) {
			t.On = !t.On
			if t.OnChange != nil {
				t.OnChange(t.On)
			}
			// Fire event
			event.GetQueuedInstance().QueueEvent(ToggleChangeEvent{
				ToggleID: t.GetID(),
				Toggle:   t,
				On:       t.On,
			})
		}
	}
}

// Layout calculates toggle dimensions
func (t *Toggle) Layout() {
	style := t.GetComputedStyle()

	fontSize := 14
	if style.FontSize != nil {
		fontSize = *style.FontSize
	}

	// Calculate size: toggle switch (24px) + spacing (8px) + text
	switchWidth := 24
	textWidth := len(t.Text) * fontSize * 6 / 10
	width := switchWidth + 8 + textWidth
	height := 24

	// Add padding
	if style.Padding != nil {
		width += style.Padding.Left + style.Padding.Right
		height += style.Padding.Top + style.Padding.Bottom
	}

	// Apply width/height from style if specified
	if style.Width != nil {
		width = *style.Width
	}
	if style.Height != nil {
		height = *style.Height
	}

	// Apply min/max constraints
	width, height = ApplySizeConstraints(width, height, style)

	t.bounds.Width = width
	t.bounds.Height = height
}

// Draw draws the toggle
func (t *Toggle) Draw(screen *ebiten.Image) {
	if !t.visible {
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

	// Draw background (if any)
	DrawBackground(screen, absBounds, style)

	contentBounds := GetContentBounds(absBounds, style)

	// Draw toggle switch
	switchWidth := 36
	switchHeight := 20
	switchX := contentBounds.X
	switchY := contentBounds.Y + (contentBounds.Height-switchHeight)/2
	switchBounds := Rect{X: switchX, Y: switchY, Width: switchWidth, Height: switchHeight}

	// Check if we should use sprite-based rendering
	if theme != nil && theme.HasToggleSprites() {
		var coords *SpriteCoords
		if t.On && theme.ToggleOn != nil {
			coords = theme.ToggleOn
		} else {
			coords = theme.Toggle
		}
		DrawSprite(screen, theme.SpriteSheet, coords, switchBounds)
	} else {
		// Use vector-based rendering with theme colors
		// Draw switch track
		trackColor := color.RGBA{80, 80, 90, 255}
		if theme != nil {
			trackColor = colorToRGBA(theme.Colors.Surface)
		}
		if t.On {
			if theme != nil {
				trackColor = colorToRGBA(theme.Colors.Primary)
			} else {
				trackColor = color.RGBA{80, 140, 200, 255}
			}
		}
		if t.hovered {
			trackColor.R = min(trackColor.R+20, 255)
			trackColor.G = min(trackColor.G+20, 255)
			trackColor.B = min(trackColor.B+20, 255)
		}

		trackBounds := Rect{X: switchX, Y: switchY, Width: switchWidth, Height: switchHeight}
		DrawRoundedRect(screen, trackBounds, switchHeight/2, trackColor)

		// Draw switch knob
		knobSize := switchHeight - 4
		knobX := switchX + 2
		if t.On {
			knobX = switchX + switchWidth - knobSize - 2
		}
		knobY := switchY + 2

		knobColor := color.RGBA{255, 255, 255, 255}
		if theme != nil {
			knobColor = colorToRGBA(theme.Colors.Text)
		}
		vector.DrawFilledCircle(screen, float32(knobX)+float32(knobSize)/2, float32(knobY)+float32(knobSize)/2, float32(knobSize)/2, knobColor, true)
	}

	// Draw text label
	if t.Text != "" {
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

		textX := switchX + switchWidth + 8
		textY := contentBounds.Y + (contentBounds.Height-fontSize)/2
		text.Draw(screen, t.Text, float64(fontSize), textX, textY, textColor)
	}
}

// SetOn sets the toggle state
func (t *Toggle) SetOn(on bool) {
	if t.On != on {
		t.On = on
		if t.OnChange != nil {
			t.OnChange(t.On)
		}
	}
}

// IsOn returns the current toggle state
func (t *Toggle) IsOn() bool {
	return t.On
}
