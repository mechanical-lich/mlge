package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/config"
	"github.com/mechanical-lich/mlge/resource"
	"github.com/mechanical-lich/mlge/text/v2"
)

type RadioButton struct {
	X, Y          int
	Width, Height int
	Label         string
	IconX         int
	IconY         int
	IconResource  string
	op            *ebiten.DrawImageOptions // Options for drawing the radio button background
}

type RadioGroup struct {
	Name     string
	Buttons  []RadioButton
	Selected int
}

func NewRadioGroup(name string, buttons []RadioButton) *RadioGroup {
	if buttons == nil {
		buttons = make([]RadioButton, 0)
	}
	return &RadioGroup{
		Name:     name,
		Buttons:  buttons,
		Selected: -1, // No selection by default
	}
}

func NewRadioButton(x, y int, label string) RadioButton {
	w, h := text.Measure(label, 16)

	return RadioButton{
		X:      x,
		Y:      y,
		Width:  int(w) + 10,
		Height: int(h) + 10, // Add some padding
		Label:  label,
		op:     &ebiten.DrawImageOptions{},
	}
}

func NewIconRadioButton(x int, y, iconX, iconY int, iconResource string, label string) *RadioButton {
	w, h := text.Measure(label, 16)

	b := &RadioButton{
		X:            x,
		Y:            y,
		IconX:        iconX,
		IconY:        iconY,
		IconResource: iconResource,
		Width:        int(w + 10 + 16),
		Height:       int(h + 16),
		Label:        label,
		op:           &ebiten.DrawImageOptions{},
	}

	return b
}

func (rg *RadioGroup) Update(parentX, parentY int) {
	for i, button := range rg.Buttons {
		if button.IsClicked(parentX, parentY) {
			rg.Selected = i
			break
		}
	}
}

func (rg *RadioGroup) Draw(screen *ebiten.Image, parentX, parentY int, theme *Theme) {
	for i, button := range rg.Buttons {
		selected := (i == rg.Selected)
		button.Draw(screen, selected, parentX, parentY, theme)
	}
}

func (rb *RadioButton) Draw(screen *ebiten.Image, selected bool, parentX, parentY int, theme *Theme) {
	rb.op.GeoM.Reset()
	rb.op.GeoM.Scale(float64(rb.Width)/float64(theme.RadioButton.Width), float64(rb.Height)/float64(theme.RadioButton.Height))
	rb.op.GeoM.Translate(float64(rb.X+parentX), float64(rb.Y+parentY))
	sX := theme.RadioButton.SrcX
	sY := theme.RadioButton.SrcY

	if selected {
		sX += 32
	}

	screen.DrawImage(resource.GetSubImage("ui", sX, sY, theme.RadioButton.Width, theme.RadioButton.Height), rb.op)
	if rb.IconResource != "" {
		text.Draw(screen, rb.Label, 15, rb.X+5+16+parentX, rb.Y+5+parentY, color.White)

		if rb.IconResource != "" {
			rb.op.GeoM.Reset()
			rb.op.GeoM.Scale(1.0, 1.0)
			rb.op.GeoM.Translate(float64(rb.X+5+parentX), float64(rb.Y+5+parentY))
			screen.DrawImage(resource.GetSubImage(rb.IconResource, rb.IconX, rb.IconY, config.SpriteSizeW, config.SpriteSizeH), rb.op)
		}
	} else {
		text.Draw(screen, rb.Label, 15, rb.X+5+parentX, rb.Y+5+parentY, color.White)
	}
}

func (b *RadioButton) IsWithin(cX int, cY int, parentX, parentY int) bool {
	if cX >= b.X+parentX && cX <= b.X+b.Width+parentX && cY >= b.Y+parentY && cY <= b.Height+b.Y+parentY {
		return true
	}
	return false
}

func (b *RadioButton) IsClicked(parentX, parentY int) bool {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		cX, cY := ebiten.CursorPosition()

		if cX >= b.X+parentX && cX <= b.X+b.Width+parentX && cY >= b.Y+parentY && cY <= b.Height+b.Y+parentY {
			return true
		}
	}
	return false
}
