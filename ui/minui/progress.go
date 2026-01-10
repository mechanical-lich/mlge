package minui

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/text/v2"
)

// ProgressBar is a simple horizontal progress bar element.
// Value is clamped between 0.0 and 1.0. It can be updated at any time
// via SetValue and GetValue returns current value.
// The bar respects the element's style (padding, margin, border, width/height constraints).
type ProgressBar struct {
	*ElementBase

	Value float64 // 0.0 .. 1.0

	// Optional callback when value changes
	OnChange func(v float64)

	// If true, draws percent text inside the bar
	ShowPercent bool
}

func NewProgressBar(id string) *ProgressBar {
	p := &ProgressBar{
		ElementBase: NewElementBase(id),
		Value:       0,
	}
	// sensible default size
	p.SetSize(200, 18)

	borderWidth := 1
	p.GetStyle().BorderWidth = &borderWidth
	return p
}

func (p *ProgressBar) GetType() string { return "ProgressBar" }

// SetValue sets the progress value (0..1). Will clamp and call OnChange if changed.
func (p *ProgressBar) SetValue(v float64) {
	v = math.Max(0, math.Min(1, v))
	if p.Value == v {
		return
	}
	p.Value = v
	if p.OnChange != nil {
		p.OnChange(p.Value)
	}
}

func (p *ProgressBar) GetValue() float64 { return p.Value }

func (p *ProgressBar) Update() {
	// passive element; value is set externally
}

func (p *ProgressBar) Layout() {
	style := p.GetComputedStyle()

	// Default sizes
	w := 200
	h := 18

	if style != nil {
		if style.Width != nil {
			w = *style.Width
		}
		if style.Height != nil {
			h = *style.Height
		}
	}

	// Apply min/max constraints
	w, h = ApplySizeConstraints(w, h, style)

	p.bounds.Width = w
	p.bounds.Height = h
}

func (p *ProgressBar) Draw(screen *ebiten.Image) {
	if !p.visible {
		return
	}

	style := p.GetComputedStyle()
	theme := p.GetTheme()
	absX, absY := p.GetAbsolutePosition()

	absBounds := Rect{X: absX, Y: absY, Width: p.bounds.Width, Height: p.bounds.Height}

	// Draw background with theme support
	DrawBackgroundWithTheme(screen, absBounds, style, theme)

	// Content bounds (accounts for padding & border)
	content := GetContentBounds(absBounds, style)

	// Filled width based on Value
	filledW := int(math.Round(float64(content.Width) * p.Value))
	if filledW < 0 {
		filledW = 0
	}
	if filledW > content.Width {
		filledW = content.Width
	}

	if filledW > 0 {
		filledRect := Rect{X: content.X, Y: content.Y, Width: filledW, Height: content.Height}
		// Get fill color from theme or default
		fillColor := color.NRGBA{R: 0x33, G: 0x99, B: 0xff, A: 0xff}
		if theme != nil {
			fillColor = color.NRGBA{
				R: colorToRGBA(theme.Colors.Primary).R,
				G: colorToRGBA(theme.Colors.Primary).G,
				B: colorToRGBA(theme.Colors.Primary).B,
				A: 0xff,
			}
		}
		DrawRect(screen, filledRect, fillColor)
	}

	// Draw border with theme support
	DrawBorderWithTheme(screen, absBounds, style, theme)

	// Optionally draw percent text centered
	if p.ShowPercent {
		percent := int(math.Round(p.Value * 100))
		textStr := fmt.Sprintf("%d%%", percent)

		fontSize := 14
		if style != nil && style.FontSize != nil {
			fontSize = *style.FontSize
		}

		// Center text
		textX := content.X + (content.Width / 2)
		textY := content.Y + (content.Height-fontSize)/2
		// Get text color from style, then theme, then default
		textColor := color.RGBA{230, 230, 230, 255}
		if style != nil && style.ForegroundColor != nil {
			textColor = colorToRGBA(*style.ForegroundColor)
		} else if theme != nil {
			textColor = colorToRGBA(theme.Colors.Text)
		}

		text.Draw(screen, textStr, float64(fontSize), textX, textY, textColor)
	}
}
