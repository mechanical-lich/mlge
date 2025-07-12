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

func (b *Toggle) Update() {
	if b.IsJustClicked() {
		b.On = !b.On
	}
}

func (b Toggle) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(b.Width)/32.0, float64(b.Height)/16.0)
	op.GeoM.Translate(float64(b.X), float64(b.Y))
	sX := 16
	sY := 64

	if (!b.IsClicked() && b.On) || (b.IsClicked() && !b.On) {
		sX += 32
	}

	//s.drawSpriteEx(int32(x), int32(y), sX, sY, 32, 32, 255, 255, 255, 255, s.uiTexture)
	screen.DrawImage(resource.GetSubImage(resource.Textures["ui"], sX, sY, config.SpriteSizeW*2, config.SpriteSizeH), op)
	if b.IconResource != "" {
		text.Draw(screen, b.Text, 15, b.X+5+16, b.Y+5, color.White)

		iconImage := resource.Textures[b.IconResource]
		if iconImage != nil {
			iconOp := &ebiten.DrawImageOptions{}
			iconOp.GeoM.Scale(1.0, 1.0)
			iconOp.GeoM.Translate(float64(b.X+5), float64(b.Y+5))
			screen.DrawImage(resource.GetSubImage(iconImage, b.IconX, b.IconY, config.SpriteSizeW, config.SpriteSizeH), iconOp)
		}
	} else {
		text.Draw(screen, b.Text, 15, b.X+5, b.Y+5, color.White)
	}
}

func (b *Toggle) IsWithin(cX int, cY int) bool {
	if cX >= b.X && cX <= b.X+b.Width && cY >= b.Y && cY <= b.Height+b.Y {
		return true
	}
	return false
}

func (b *Toggle) IsClicked() bool {
	cX, cY := ebiten.CursorPosition()

	if cX >= b.X && cX <= b.X+b.Width && cY >= b.Y && cY <= b.Height+b.Y && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return true
	}
	return false
}

func (b *Toggle) IsJustClicked() bool {
	cX, cY := ebiten.CursorPosition()

	if cX >= b.X && cX <= b.X+b.Width && cY >= b.Y && cY <= b.Height+b.Y && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		return true
	}
	return false
}
