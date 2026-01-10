package minui

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/text/v2"
)

// Label is a text display element
type Label struct {
	*ElementBase
	Text string
}

// NewLabel creates a new label
func NewLabel(id, text string) *Label {
	label := &Label{
		ElementBase: NewElementBase(id),
		Text:        text,
	}

	// Set default size based on text
	label.SetSize(len(text)*8, 20)

	return label
}

// GetType returns the element type
func (l *Label) GetType() string {
	return "Label"
}

// Update updates the label
func (l *Label) Update() {
	if !l.visible {
		return
	}
	l.UpdateHoverState()
}

// Layout calculates label dimensions
func (l *Label) Layout() {
	style := l.GetComputedStyle()

	fontSize := 14
	if style.FontSize != nil {
		fontSize = *style.FontSize
	}

	// Split text by newlines to calculate proper dimensions
	lines := strings.Split(l.Text, "\n")
	lineHeight := fontSize + 4 // Line spacing

	// Find the longest line for width calculation
	maxLineLen := 0
	for _, line := range lines {
		if len(line) > maxLineLen {
			maxLineLen = len(line)
		}
	}

	// Estimate text size based on longest line and number of lines
	textWidth := maxLineLen * fontSize * 6 / 10
	textHeight := (len(lines) * lineHeight) - 4 + 6 // Remove spacing from last line, add base padding

	// Add padding
	if style.Padding != nil {
		textWidth += style.Padding.Left + style.Padding.Right
		textHeight += style.Padding.Top + style.Padding.Bottom
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

	l.bounds.Width = textWidth
	l.bounds.Height = textHeight
}

// Draw draws the label
func (l *Label) Draw(screen *ebiten.Image) {
	if !l.visible {
		return
	}

	style := l.GetComputedStyle()
	theme := l.GetTheme()

	// Get absolute position for drawing
	absX, absY := l.GetAbsolutePosition()
	absBounds := Rect{
		X:      absX,
		Y:      absY,
		Width:  l.bounds.Width,
		Height: l.bounds.Height,
	}

	// Draw background with theme support
	DrawBackgroundWithTheme(screen, absBounds, style, theme)

	// Draw border with theme support
	DrawBorderWithTheme(screen, absBounds, style, theme)

	// Draw text
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

	// Calculate text position based on alignment
	textX := contentBounds.X
	textY := contentBounds.Y

	// Split text by newlines
	lines := strings.Split(l.Text, "\n")
	lineHeight := fontSize + 4 // Add some line spacing

	if style.TextAlign != nil {
		switch *style.TextAlign {
		case TextAlignCenter:
			textX = contentBounds.X + (contentBounds.Width / 2) - (len(l.Text)*fontSize*6/10)/2
		case TextAlignRight:
			textX = contentBounds.X + contentBounds.Width - (len(l.Text) * fontSize * 6 / 10)
		}
	}

	if style.VertAlign != nil {
		switch *style.VertAlign {
		case VertAlignMiddle:
			totalHeight := len(lines) * lineHeight
			textY = contentBounds.Y + (contentBounds.Height / 2) - (totalHeight / 2)
		case VertAlignBottom:
			totalHeight := len(lines) * lineHeight
			textY = contentBounds.Y + contentBounds.Height - totalHeight
		}
	}

	// Draw each line
	for i, line := range lines {
		lineX := textX
		lineY := textY + (i * lineHeight)

		// Recalculate X position for each line if centered or right-aligned
		if style.TextAlign != nil {
			switch *style.TextAlign {
			case TextAlignCenter:
				lineX = contentBounds.X + (contentBounds.Width / 2) - (len(line)*fontSize*6/10)/2
			case TextAlignRight:
				lineX = contentBounds.X + contentBounds.Width - (len(line) * fontSize * 6 / 10)
			}
		}

		text.Draw(screen, line, float64(fontSize), lineX, lineY, textColor)
	}
}
