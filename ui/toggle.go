package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mechanical-lich/mlge/config"
	"github.com/mechanical-lich/mlge/resource"
	"github.com/mechanical-lich/mlge/text/v2"
)

type Toggle struct {
	Name         string
	X            int
	Y            int
	Width        int
	Height       int
	Text         string
	IconX        int
	IconY        int
	IconResource string
	On           bool
}

func NewToggle(name string, x int, y int, txt string) *Toggle {
	w, h := text.Measure(txt, 16)

	b := &Toggle{
		Name:   name,
		X:      x,
		Y:      y,
		Width:  int(w + 10),
		Height: int(h + 10),
		Text:   txt,
		On:     false,
	}

	return b
}

func NewIconToggle(name string, x int, y, iconX, iconY int, iconResource string, txt string) *Toggle {
	w, h := text.Measure(txt, 16)

	b := &Toggle{
		Name:         name,
		X:            x,
		Y:            y,
		IconX:        iconX,
		IconY:        iconY,
		IconResource: iconResource,
		Width:        int(w + 10 + 16),
		Height:       int(h + 16),
		Text:         txt,
	}

	return b
}

func (b *Toggle) Update(parentX, parentY int) {
	if b.IsJustClicked(parentX, parentY) {
		b.On = !b.On
	}
}

func (b Toggle) Draw(screen *ebiten.Image, parentX, parentY int, theme *Theme) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(b.Width)/float64(theme.Toggle.Width), float64(b.Height)/float64(theme.Toggle.Height))
	op.GeoM.Translate(float64(b.X+parentX), float64(b.Y+parentY))
	sX := theme.Toggle.SrcX
	sY := theme.Toggle.SrcY

	if (!b.IsClicked(parentX, parentY) && b.On) || (b.IsClicked(parentX, parentY) && !b.On) {
		sX += 32
	}

	//s.drawSpriteEx(int32(x), int32(y), sX, sY, 32, 32, 255, 255, 255, 255, s.uiTexture)
	screen.DrawImage(resource.GetSubImage("ui", sX, sY, theme.Toggle.Width, theme.Toggle.Height), op)
	if b.IconResource != "" {
		text.Draw(screen, b.Text, 15, b.X+5+16+parentX, b.Y+5+parentY, color.White)

		if b.IconResource != "" {
			iconOp := &ebiten.DrawImageOptions{}
			iconOp.GeoM.Scale(1.0, 1.0)
			iconOp.GeoM.Translate(float64(b.X+5+parentX), float64(b.Y+5+parentY))
			screen.DrawImage(resource.GetSubImage(b.IconResource, b.IconX, b.IconY, config.SpriteSizeW, config.SpriteSizeH), iconOp)
		}
	} else {
		text.Draw(screen, b.Text, 15, b.X+5+parentX, b.Y+5+parentY, color.White)
	}
}

func (b *Toggle) IsWithin(cX int, cY int, parentX, parentY int) bool {
	if cX >= b.X+parentX && cX <= b.X+b.Width+parentX && cY >= b.Y+parentY && cY <= b.Height+b.Y+parentY {
		return true
	}
	return false
}

func (b *Toggle) IsClicked(parentX, parentY int) bool {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		cX, cY := ebiten.CursorPosition()

		if cX >= b.X+parentX && cX <= b.X+b.Width+parentX && cY >= b.Y+parentY && cY <= b.Height+b.Y+parentY {
			return true
		}
	}
	return false
}

func (b *Toggle) IsJustClicked(parentX, parentY int) bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cX, cY := ebiten.CursorPosition()

		if cX >= b.X+parentX && cX <= b.X+b.Width+parentX && cY >= b.Y+parentY && cY <= b.Height+b.Y+parentY {
			return true
		}
	}
	return false
}
