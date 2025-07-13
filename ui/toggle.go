package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/config"
	"github.com/mechanical-lich/mlge/resource"
	"github.com/mechanical-lich/mlge/text/v2"
)

type Toggle struct {
	ElementBase // Embed ElementBase for common properties

	Text string

	On bool
}

func NewToggle(name string, x int, y int, txt string) *Toggle {
	w, h := text.Measure(txt, 16)

	b := &Toggle{
		ElementBase: ElementBase{
			Name:   name,
			X:      x,
			Y:      y,
			Width:  int(w + 10),
			Height: int(h + 10),
			op:     &ebiten.DrawImageOptions{},
		},
		Text: txt,
		On:   false,
	}

	return b
}

func NewIconToggle(name string, x int, y, iconX, iconY int, iconResource string, txt string) *Toggle {
	w, h := text.Measure(txt, 16)

	b := &Toggle{
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

func (b *Toggle) Update(parentX, parentY int) {
	if b.IsJustClicked(parentX, parentY) {
		b.On = !b.On
	}
}

func (b *Toggle) Draw(screen *ebiten.Image, parentX, parentY int, theme *Theme) {
	b.op.GeoM.Reset()
	b.op.GeoM.Scale(float64(b.Width)/float64(theme.Toggle.Width), float64(b.Height)/float64(theme.Toggle.Height))
	b.op.GeoM.Translate(float64(b.X+parentX), float64(b.Y+parentY))
	sX := theme.Toggle.SrcX
	sY := theme.Toggle.SrcY

	if (!b.IsClicked(parentX, parentY) && b.On) || (b.IsClicked(parentX, parentY) && !b.On) {
		sX += 32
	}

	//s.drawSpriteEx(int32(x), int32(y), sX, sY, 32, 32, 255, 255, 255, 255, s.uiTexture)
	screen.DrawImage(resource.GetSubImage("ui", sX, sY, theme.Toggle.Width, theme.Toggle.Height), b.op)
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
