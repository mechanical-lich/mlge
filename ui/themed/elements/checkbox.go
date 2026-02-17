package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/resource"
	"github.com/mechanical-lich/mlge/text"
	theming "github.com/mechanical-lich/mlge/ui/themed/theming"
)

// Checkbox represents a checkbox UI element
type Checkbox struct {
	ElementBase
	Label     string
	Checked   bool
	OnChanged func(checked bool) // Callback when checkbox state changes
}

// NewCheckbox creates a new checkbox with the given parameters
func NewCheckbox(name string, x, y int, label string, checked bool) *Checkbox {
	return &Checkbox{
		ElementBase: ElementBase{
			Name:   name,
			X:      x,
			Y:      y,
			Width:  16,
			Height: 16,
			Op:     &ebiten.DrawImageOptions{},
		},
		Label:   label,
		Checked: checked,
	}
}

func (c *Checkbox) Update() {
	if c.IsJustClicked() {
		c.Checked = !c.Checked
		if c.OnChanged != nil {
			c.OnChanged(c.Checked)
		}
	}
}

func (c *Checkbox) Draw(screen *ebiten.Image, theme *theming.Theme) {
	absX, absY := c.GetAbsolutePosition()

	// Draw checkbox background
	c.Op.GeoM.Reset()
	c.Op.GeoM.Translate(float64(absX), float64(absY))

	srcY := theme.Checkbox.SrcY
	if c.Checked {
		srcY += theme.Checkbox.Height // Use second row for checked state
	}

	screen.DrawImage(
		resource.GetSubImage("ui", theme.Checkbox.SrcX, srcY, theme.Checkbox.Width, theme.Checkbox.Height),
		c.Op,
	)

	// Draw label
	if c.Label != "" {
		text.Draw(screen, c.Label, 14, absX+20, absY+2, theme.Colors.Text)
	}
}

// SetChecked sets the checked state
func (c *Checkbox) SetChecked(checked bool) {
	c.Checked = checked
}

// IsChecked returns the checked state
func (c *Checkbox) IsChecked() bool {
	return c.Checked
}
