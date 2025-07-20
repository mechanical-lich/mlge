package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/config"
	"github.com/mechanical-lich/mlge/resource"
	"github.com/mechanical-lich/mlge/text/v2"
)

type Button struct {
	ElementBase
	Text string
}

func NewButton(name string, x int, y int, txt string) *Button {
	w, h := text.Measure(txt, 16)

	b := &Button{
		ElementBase: ElementBase{
			Name:         name,
			X:            x,
			Y:            y,
			Width:        int(w + 10),
			Height:       int(h + 10),
			IconX:        0,
			IconY:        0,
			IconResource: "",
			op:           &ebiten.DrawImageOptions{},
		},
		Text: txt,
	}

	return b
}

func NewIconButton(name string, x int, y, iconX, iconY int, iconResource string, txt string) *Button {
	w, h := text.Measure(txt, 16)

	b := &Button{
		ElementBase: ElementBase{
			Name:         name,
			X:            x,
			Y:            y,
			IconX:        iconX,
			IconY:        iconY,
			IconResource: iconResource,
			Width:        int(w + 10 + 16),
			Height:       int(h + 16),
			op:           &ebiten.DrawImageOptions{},
		},
		Text: txt,
	}

	return b
}

func (b *Button) Update(parentX, parentY int) {
}

func (b *Button) Draw(screen *ebiten.Image, parentX, parentY int, theme *Theme) {
	b.op.GeoM.Reset()
	b.op.GeoM.Scale(float64(b.Width)/float64(theme.Button.Width), float64(b.Height)/float64(theme.Button.Height))
	b.op.GeoM.Translate(float64(b.X+parentX), float64(b.Y+parentY))
	sX := theme.Button.SrcX
	sY := theme.Button.SrcY

	if b.IsClicked(parentX, parentY) {
		sX += 32
	}

	screen.DrawImage(resource.GetSubImage("ui", sX, sY, theme.Button.Width, theme.Button.Height), b.op)
	if b.IconResource != "" {
		text.Draw(screen, b.Text, 15, b.X+5+16+parentX, b.Y+5+parentY, color.White)

		if b.IconResource != "" {
			b.op.GeoM.Reset()
			b.op.GeoM.Scale(1.0, 1.0)
			b.op.GeoM.Translate(float64(b.X+5+parentX), float64(b.Y+5+parentY))
			screen.DrawImage(resource.GetSubImage(b.IconResource, b.IconX, b.IconY, config.SpriteSizeW, config.SpriteSizeH), b.op)
		}
	} else {
		text.Draw(screen, b.Text, 15, b.X+5+parentX, b.Y+5+parentY, color.White)
	}
}
