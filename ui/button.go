package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mechanical-lich/mlge/config"
	"github.com/mechanical-lich/mlge/resource"
	"github.com/mechanical-lich/mlge/text/v2"
)

type Button struct {
	Name         string
	X            int
	Y            int
	Width        int
	Height       int
	Text         string
	IconX        int
	IconY        int
	IconResource string
}

func NewButton(name string, x int, y int, txt string) *Button {
	w, h := text.Measure(txt, 16)

	b := &Button{
		Name:   name,
		X:      x,
		Y:      y,
		Width:  int(w + 10),
		Height: int(h + 10),
		Text:   txt,
	}

	return b
}

func NewIconButton(name string, x int, y, iconX, iconY int, iconResource string, txt string) *Button {
	w, h := text.Measure(txt, 16)

	b := &Button{
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

func (b Button) Draw(screen *ebiten.Image, parentX, parentY int, theme *Theme) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(b.Width)/float64(theme.Button.Width), float64(b.Height)/float64(theme.Button.Height))
	op.GeoM.Translate(float64(b.X+parentX), float64(b.Y+parentY))
	sX := theme.Button.SrcX
	sY := theme.Button.SrcY

	if b.IsClicked(parentX, parentY) {
		sX += 32
	}

	screen.DrawImage(resource.GetSubImage("ui", sX, sY, theme.Button.Width, theme.Button.Height), op)
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

func (b *Button) IsWithin(cX int, cY int, parentX, parentY int) bool {
	if cX >= b.X+parentX && cX <= b.X+b.Width+parentX && cY >= b.Y+parentY && cY <= b.Height+b.Y+parentY {
		return true
	}
	return false
}

func (b *Button) IsClicked(parentX, parentY int) bool {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		cX, cY := ebiten.CursorPosition()

		if cX >= b.X+parentX && cX <= b.X+b.Width+parentX && cY >= b.Y+parentY && cY <= b.Height+b.Y+parentY {
			return true
		}
	}
	return false
}

func (b *Button) IsJustClicked(parentX, parentY int) bool {
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		cX, cY := ebiten.CursorPosition()

		if cX >= b.X+parentX && cX <= b.X+b.Width+parentX && cY >= b.Y+parentY && cY <= b.Height+b.Y+parentY {
			return true
		}
	}
	return false
}
