package ui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/mechanical-lich/mlge/config"
	"github.com/mechanical-lich/mlge/resource"
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

func NewButton(x int, y int, text string) *Button {
	b := &Button{
		X:      x,
		Y:      y,
		Width:  64,
		Height: 32,
		Text:   text,
	}

	return b
}

func (b Button) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(float64(b.X), float64(b.Y))
	sX := 16
	sY := 64

	if b.IsClicked() {
		sX += 32
	}
	//s.drawSpriteEx(int32(x), int32(y), sX, sY, 32, 32, 255, 255, 255, 255, s.uiTexture)
	screen.DrawImage(resource.Textures["ui"].SubImage(image.Rect(sX, sY, sX+config.SpriteSizeW*2, sY+config.SpriteSizeH)).(*ebiten.Image), op)
	text.Draw(screen, b.Text, resource.Fonts["main"], b.X, b.Y+20, color.White)
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
