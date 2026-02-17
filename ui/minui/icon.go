package minui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/resource"
	"github.com/mechanical-lich/mlge/text"
)

// Icon defines a sprite icon from a sprite sheet
type Icon struct {
	// SpriteSheet is the resource name of the sprite sheet (e.g., "ui", "items")
	SpriteSheet string
	// SrcX, SrcY are the top-left coordinates in the sprite sheet
	SrcX, SrcY int
	// Width, Height are the dimensions of the icon in the sprite sheet
	Width, Height int
	// Scale is the rendering scale (default 1.0)
	Scale float64
}

// NewIcon creates a new icon definition
func NewIcon(spriteSheet string, srcX, srcY, width, height int) *Icon {
	return &Icon{
		SpriteSheet: spriteSheet,
		SrcX:        srcX,
		SrcY:        srcY,
		Width:       width,
		Height:      height,
		Scale:       1.0,
	}
}

// NewIconWithScale creates an icon with a specific scale
func NewIconWithScale(spriteSheet string, srcX, srcY, width, height int, scale float64) *Icon {
	return &Icon{
		SpriteSheet: spriteSheet,
		SrcX:        srcX,
		SrcY:        srcY,
		Width:       width,
		Height:      height,
		Scale:       scale,
	}
}

// ScaledWidth returns the width after applying scale
func (i *Icon) ScaledWidth() int {
	if i == nil {
		return 0
	}
	return int(float64(i.Width) * i.Scale)
}

// ScaledHeight returns the height after applying scale
func (i *Icon) ScaledHeight() int {
	if i == nil {
		return 0
	}
	return int(float64(i.Height) * i.Scale)
}

// Draw draws the icon at the specified position
func (i *Icon) Draw(screen *ebiten.Image, x, y int) {
	if i == nil {
		return
	}

	img := resource.GetSubImage(i.SpriteSheet, i.SrcX, i.SrcY, i.Width, i.Height)
	if img == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}
	if i.Scale != 1.0 {
		op.GeoM.Scale(i.Scale, i.Scale)
	}
	op.GeoM.Translate(float64(x), float64(y))

	screen.DrawImage(img, op)
}

// DrawCentered draws the icon centered at the specified position
func (i *Icon) DrawCentered(screen *ebiten.Image, centerX, centerY int) {
	if i == nil {
		return
	}

	x := centerX - i.ScaledWidth()/2
	y := centerY - i.ScaledHeight()/2
	i.Draw(screen, x, y)
}

// DrawWithOpacity draws the icon with opacity
func (i *Icon) DrawWithOpacity(screen *ebiten.Image, x, y int, opacity float32) {
	if i == nil {
		return
	}

	img := resource.GetSubImage(i.SpriteSheet, i.SrcX, i.SrcY, i.Width, i.Height)
	if img == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}
	if i.Scale != 1.0 {
		op.GeoM.Scale(i.Scale, i.Scale)
	}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleAlpha(opacity)

	screen.DrawImage(img, op)
}

// IconPosition defines where the icon should be placed relative to text
type IconPosition int

const (
	IconLeft   IconPosition = iota // Icon to the left of text
	IconRight                      // Icon to the right of text
	IconTop                        // Icon above text
	IconBottom                     // Icon below text
	IconOnly                       // Icon only, no text
)

// IconLabel is a label with an optional icon
type IconLabel struct {
	*ElementBase
	Text         string
	Icon         *Icon
	IconPosition IconPosition
	IconSpacing  int // Space between icon and text
}

// NewIconLabel creates a new label with an icon
func NewIconLabel(id string, icon *Icon, text string) *IconLabel {
	label := &IconLabel{
		ElementBase:  NewElementBase(id),
		Text:         text,
		Icon:         icon,
		IconPosition: IconLeft,
		IconSpacing:  4,
	}

	// Set default size
	label.updateSize()

	return label
}

// NewIconOnlyLabel creates a label that shows only an icon
func NewIconOnlyLabel(id string, icon *Icon) *IconLabel {
	label := &IconLabel{
		ElementBase:  NewElementBase(id),
		Text:         "",
		Icon:         icon,
		IconPosition: IconOnly,
		IconSpacing:  0,
	}

	label.updateSize()
	return label
}

func (l *IconLabel) updateSize() {
	iconW := 0
	iconH := 0
	if l.Icon != nil {
		iconW = l.Icon.ScaledWidth()
		iconH = l.Icon.ScaledHeight()
	}

	textW := len(l.Text) * 8
	textH := 16

	switch l.IconPosition {
	case IconLeft, IconRight:
		if l.Icon != nil && l.Text != "" {
			l.SetSize(iconW+l.IconSpacing+textW, max(iconH, textH))
		} else if l.Icon != nil {
			l.SetSize(iconW, iconH)
		} else {
			l.SetSize(textW, textH)
		}
	case IconTop, IconBottom:
		if l.Icon != nil && l.Text != "" {
			l.SetSize(max(iconW, textW), iconH+l.IconSpacing+textH)
		} else if l.Icon != nil {
			l.SetSize(iconW, iconH)
		} else {
			l.SetSize(textW, textH)
		}
	case IconOnly:
		l.SetSize(iconW, iconH)
	}
}

// SetIcon sets the icon
func (l *IconLabel) SetIcon(icon *Icon) {
	l.Icon = icon
	l.updateSize()
}

// SetText sets the text
func (l *IconLabel) SetText(text string) {
	l.Text = text
	l.updateSize()
}

// GetType returns the element type
func (l *IconLabel) GetType() string {
	return "IconLabel"
}

// Update updates the label
func (l *IconLabel) Update() {
	if !l.visible {
		return
	}
	l.UpdateHoverState()
}

// Layout calculates dimensions
func (l *IconLabel) Layout() {
	style := l.GetComputedStyle()

	fontSize := 14
	if style.FontSize != nil {
		fontSize = *style.FontSize
	}

	iconW := 0
	iconH := 0
	if l.Icon != nil {
		iconW = l.Icon.ScaledWidth()
		iconH = l.Icon.ScaledHeight()
	}

	textW := len(l.Text) * fontSize * 6 / 10
	textH := fontSize + 4

	var width, height int

	switch l.IconPosition {
	case IconLeft, IconRight:
		if l.Icon != nil && l.Text != "" {
			width = iconW + l.IconSpacing + textW
			height = max(iconH, textH)
		} else if l.Icon != nil {
			width = iconW
			height = iconH
		} else {
			width = textW
			height = textH
		}
	case IconTop, IconBottom:
		if l.Icon != nil && l.Text != "" {
			width = max(iconW, textW)
			height = iconH + l.IconSpacing + textH
		} else if l.Icon != nil {
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

	// Add padding
	if style.Padding != nil {
		width += style.Padding.Left + style.Padding.Right
		height += style.Padding.Top + style.Padding.Bottom
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

	l.bounds.Width = width
	l.bounds.Height = height
}

// Draw draws the icon label
func (l *IconLabel) Draw(screen *ebiten.Image) {
	if !l.visible {
		return
	}

	style := l.GetComputedStyle()
	absX, absY := l.GetAbsolutePosition()
	absBounds := Rect{
		X:      absX,
		Y:      absY,
		Width:  l.bounds.Width,
		Height: l.bounds.Height,
	}

	// Draw background if any
	DrawBackground(screen, absBounds, style)

	contentBounds := GetContentBounds(absBounds, style)

	fontSize := 14
	if style.FontSize != nil {
		fontSize = *style.FontSize
	}

	textColor := getTextColor(style)

	iconW := 0
	iconH := 0
	if l.Icon != nil {
		iconW = l.Icon.ScaledWidth()
		iconH = l.Icon.ScaledHeight()
	}

	textH := fontSize

	switch l.IconPosition {
	case IconLeft:
		// Icon on left, text on right
		if l.Icon != nil {
			iconY := contentBounds.Y + (contentBounds.Height-iconH)/2
			l.Icon.Draw(screen, contentBounds.X, iconY)
		}
		if l.Text != "" {
			textX := contentBounds.X
			if l.Icon != nil {
				textX += iconW + l.IconSpacing
			}
			textY := contentBounds.Y + (contentBounds.Height-textH)/2
			drawText(screen, l.Text, fontSize, textX, textY, textColor)
		}

	case IconRight:
		// Text on left, icon on right
		textW := len(l.Text) * fontSize * 6 / 10
		if l.Text != "" {
			textY := contentBounds.Y + (contentBounds.Height-textH)/2
			drawText(screen, l.Text, fontSize, contentBounds.X, textY, textColor)
		}
		if l.Icon != nil {
			iconX := contentBounds.X + textW
			if l.Text != "" {
				iconX += l.IconSpacing
			}
			iconY := contentBounds.Y + (contentBounds.Height-iconH)/2
			l.Icon.Draw(screen, iconX, iconY)
		}

	case IconTop:
		// Icon above text
		if l.Icon != nil {
			iconX := contentBounds.X + (contentBounds.Width-iconW)/2
			l.Icon.Draw(screen, iconX, contentBounds.Y)
		}
		if l.Text != "" {
			textW := len(l.Text) * fontSize * 6 / 10
			textX := contentBounds.X + (contentBounds.Width-textW)/2
			textY := contentBounds.Y + iconH + l.IconSpacing
			drawText(screen, l.Text, fontSize, textX, textY, textColor)
		}

	case IconBottom:
		// Text above icon
		if l.Text != "" {
			textW := len(l.Text) * fontSize * 6 / 10
			textX := contentBounds.X + (contentBounds.Width-textW)/2
			drawText(screen, l.Text, fontSize, textX, contentBounds.Y, textColor)
		}
		if l.Icon != nil {
			iconX := contentBounds.X + (contentBounds.Width-iconW)/2
			iconY := contentBounds.Y + textH + l.IconSpacing
			l.Icon.Draw(screen, iconX, iconY)
		}

	case IconOnly:
		// Just the icon, centered
		if l.Icon != nil {
			iconX := contentBounds.X + (contentBounds.Width-iconW)/2
			iconY := contentBounds.Y + (contentBounds.Height-iconH)/2
			l.Icon.Draw(screen, iconX, iconY)
		}
	}

	// Draw border
	DrawBorder(screen, absBounds, style)
}

// helper function to get text color from style
func getTextColor(style *Style) [4]uint8 {
	textColor := [4]uint8{255, 255, 255, 255}
	if style.ForegroundColor != nil {
		r, g, b, a := (*style.ForegroundColor).RGBA()
		textColor = [4]uint8{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
	}
	return textColor
}

// helper function to draw text
func drawText(screen *ebiten.Image, str string, fontSize int, x, y int, clr [4]uint8) {
	text.Draw(screen, str, float64(fontSize), x, y, color.RGBA{clr[0], clr[1], clr[2], clr[3]})
}

// helper function to max of two ints
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
