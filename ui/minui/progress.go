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
	absX, absY := p.GetAbsolutePosition()

	absBounds := Rect{X: absX, Y: absY, Width: p.bounds.Width, Height: p.bounds.Height}

	// Draw background
	DrawBackground(screen, absBounds, style)

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
		DrawRect(screen, filledRect, color.NRGBA{R: 0x33, G: 0x99, B: 0xff, A: 0xff})
	}

	// Draw border
	DrawBorder(screen, absBounds, style)

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
		textColor := color.RGBA{0, 0, 0, 255}
		if style != nil && style.ForegroundColor != nil {
			r, g, b, a := (*style.ForegroundColor).RGBA()
			textColor = color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
		}

		text.Draw(screen, textStr, float64(fontSize), textX, textY, textColor)
	}
}
