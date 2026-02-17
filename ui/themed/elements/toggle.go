package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/resource"
	"github.com/mechanical-lich/mlge/text"
	theming "github.com/mechanical-lich/mlge/ui/themed/theming"
)

// Represents a toggle switch element.  Clicking it toggles the "On" value true and false.
type Toggle struct {
	ElementBase // Embed ElementBase for common properties

	Text string

	On       bool
	OnChange func(value bool) // Optional onchange handler function
}

// Creates a new Toggle with the given parameters.
func NewToggle(name string, x int, y int, txt string) *Toggle {
	w, h := text.Measure(txt, 16)

	b := &Toggle{
		ElementBase: ElementBase{
			Name:   name,
			X:      x,
			Y:      y,
			Width:  int(w + 10),
			Height: int(h + 10),
			Op:     &ebiten.DrawImageOptions{},
		},
		Text: txt,
		On:   false,
	}

	return b
}

// Creates a new toggle with an icon.
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
			Op:           &ebiten.DrawImageOptions{},
		},
		Text: txt,
	}

	return b
}

func (b *Toggle) Update() {
	if b.IsJustClicked() {
		b.On = !b.On
		if b.OnChange != nil {
			b.OnChange(b.On)
		}
	}
}

func (b *Toggle) Draw(screen *ebiten.Image, theme *theming.Theme) {
	absX, absY := b.GetAbsolutePosition()
	b.Op.GeoM.Reset()
	b.Op.GeoM.Scale(float64(b.Width)/float64(theme.Toggle.Width), float64(b.Height)/float64(theme.Toggle.Height))
	b.Op.GeoM.Translate(float64(absX), float64(absY))
	sX := theme.Toggle.SrcX
	sY := theme.Toggle.SrcY

	if (!b.IsClicked() && b.On) || (b.IsClicked() && !b.On) {
		sX += 32
	}

	//s.drawSpriteEx(int32(x), int32(y), sX, sY, 32, 32, 255, 255, 255, 255, s.uiTexture)
	screen.DrawImage(resource.GetSubImage("ui", sX, sY, theme.Toggle.Width, theme.Toggle.Height), b.Op)
	if b.IconResource != "" {
		text.Draw(screen, b.Text, 15, absX+5+16, absY+5, color.White)

		if b.IconResource != "" {
			b.Op.GeoM.Reset()
			b.Op.GeoM.Scale(1.0, 1.0)
			b.Op.GeoM.Translate(float64(absX+5), float64(absY+5))
			screen.DrawImage(resource.GetSubImage(b.IconResource, b.IconX, b.IconY, spriteSizeW, spriteSizeH), b.Op)
		}
	} else {
		text.Draw(screen, b.Text, 15, absX+5, absY+5, color.White)
	}
}
