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
	}

	return b
}

func (rg *RadioGroup) Update() {
	for i, button := range rg.Buttons {
		if button.IsClicked() {
			rg.Selected = i
			break
		}
	}
}

func (rg *RadioGroup) Draw(screen *ebiten.Image) {
	for i, button := range rg.Buttons {
		selected := (i == rg.Selected)
		button.Draw(screen, selected)
	}
}

func (rb *RadioButton) Draw(screen *ebiten.Image, selected bool) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(rb.Width)/32.0, float64(rb.Height)/16.0)
	op.GeoM.Translate(float64(rb.X), float64(rb.Y))
	sX := 16
	sY := 64

	if selected {
		sX += 32
	}

	screen.DrawImage(resource.GetSubImage(resource.Textures["ui"], sX, sY, config.SpriteSizeW*2, config.SpriteSizeH), op)
	if rb.IconResource != "" {
		text.Draw(screen, rb.Label, 15, rb.X+5+16, rb.Y+5, color.White)

		iconImage := resource.Textures[rb.IconResource]
		if iconImage != nil {
			iconOp := &ebiten.DrawImageOptions{}
			iconOp.GeoM.Scale(1.0, 1.0)
			iconOp.GeoM.Translate(float64(rb.X+5), float64(rb.Y+5))
			screen.DrawImage(resource.GetSubImage(iconImage, rb.IconX, rb.IconY, config.SpriteSizeW, config.SpriteSizeH), iconOp)
		}
	} else {
		text.Draw(screen, rb.Label, 15, rb.X+5, rb.Y+5, color.White)
	}
}

func (b *RadioButton) IsWithin(cX int, cY int) bool {
	if cX >= b.X && cX <= b.X+b.Width && cY >= b.Y && cY <= b.Height+b.Y {
		return true
	}
	return false
}

func (b *RadioButton) IsClicked() bool {
	cX, cY := ebiten.CursorPosition()

	if cX >= b.X && cX <= b.X+b.Width && cY >= b.Y && cY <= b.Height+b.Y && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return true
	}
	return false
}
