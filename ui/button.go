package ui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/config"
	"github.com/mechanical-lich/mlge/resource"
	"github.com/mechanical-lich/mlge/text/v2"
)

type Button struct {
	X      int
	Y      int
	Width  int
	Height int
	Text   string
	IconX  int
	IconY  int
}

func NewButton(x int, y int, txt string) *Button {
	w, h := text.Measure(txt, 16)

	b := &Button{
		X:      x,
		Y:      y,
		Width:  int(w + 10),
		Height: int(h + 10),
		Text:   txt,
	}

	return b
}

func (b Button) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(b.Width)/32.0, float64(b.Height)/16.0)
	op.GeoM.Translate(float64(b.X), float64(b.Y))
	sX := 16
	sY := 64

	if b.IsClicked() {
		sX += 32
	}
	//s.drawSpriteEx(int32(x), int32(y), sX, sY, 32, 32, 255, 255, 255, 255, s.uiTexture)
	screen.DrawImage(resource.Textures["ui"].SubImage(image.Rect(sX, sY, sX+config.SpriteSizeW*2, sY+config.SpriteSizeH)).(*ebiten.Image), op)
	text.Draw(screen, b.Text, 15, b.X+5, b.Y+5, color.White)
}

func (b *Button) IsWithin(cX int, cY int) bool {
	if cX >= b.X && cX <= b.X+b.Width && cY >= b.Y && cY <= b.Height+b.Y {
		return true
	}
	return false
}

func (b *Button) IsClicked() bool {
	cX, cY := ebiten.CursorPosition()

	if cX >= b.X && cX <= b.X+b.Width && cY >= b.Y && cY <= b.Height+b.Y && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return true
	}
	return false
}
