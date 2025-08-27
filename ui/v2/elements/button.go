package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/config"
	"github.com/mechanical-lich/mlge/resource"
	"github.com/mechanical-lich/mlge/text/v2"
	theming "github.com/mechanical-lich/mlge/ui/v2/theming"
)

// Button - represents a clickable button UI element
type Button struct {
	ElementBase
	Text            string
	Tooltip         string
	pressed         bool
	tooltipBg       *ebiten.Image // cache for tooltip background
	tooltipBgWidth  int
	tooltipBgHeight int
	OnClicked       func() // Optional onclick handler function
}

// Creates a new button with given name, position, text, and tooltip.
func NewButton(name string, x int, y int, txt string, tooltip string) *Button {
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
			Op:           &ebiten.DrawImageOptions{},
		},
		Text:    txt,
		Tooltip: tooltip,
	}

	return b
}

// Creates a new button with an icon, given name, position, icon coordinates, icon resource, text, and tooltip.
func NewIconButton(name string, x int, y, iconX, iconY int, iconResource string, txt string, tooltip string) *Button {
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
			Op:           &ebiten.DrawImageOptions{},
		},
		Text:    txt,
		Tooltip: tooltip,
	}

	return b
}

func (b *Button) Update() {
	b.pressed = b.IsClicked()
	if b.pressed && b.OnClicked != nil {
		b.OnClicked()
	}
}

func (b *Button) Draw(screen *ebiten.Image, theme *theming.Theme) {
	absX, absY := b.GetAbsolutePosition()

	b.Op.GeoM.Reset()
	b.Op.GeoM.Scale(float64(b.Width)/float64(theme.Button.Width), float64(b.Height)/float64(theme.Button.Height))
	b.Op.GeoM.Translate(float64(absX), float64(absY))
	sX := theme.Button.SrcX
	sY := theme.Button.SrcY

	if b.pressed {
		sX += 32
	}

	screen.DrawImage(resource.GetSubImage("ui", sX, sY, theme.Button.Width, theme.Button.Height), b.Op)
	if b.IconResource != "" {
		text.Draw(screen, b.Text, 15, absX+5+16, absY+5, color.White)

		if b.IconResource != "" {
			b.Op.GeoM.Reset()
			b.Op.GeoM.Scale(1.0, 1.0)
			b.Op.GeoM.Translate(float64(absX+5), float64(absY+5))
			screen.DrawImage(resource.GetSubImage(b.IconResource, b.IconX, b.IconY, config.SpriteSizeW, config.SpriteSizeH), b.Op)
		}
	} else {
		text.Draw(screen, b.Text, 15, absX+5, absY+5, color.White)
	}

	// Tooltip rendering
	cX, cY := ebiten.CursorPosition()
	if b.IsWithin(cX, cY) && b.Tooltip != "" {
		tw, th := text.Measure(b.Tooltip, 14)
		tooltipW := int(tw + 10)
		tooltipH := int(th + 8)
		tooltipX := absX + b.Width + 8
		tooltipY := absY
		// Only recreate if size changes
		if b.tooltipBg == nil || b.tooltipBgWidth != tooltipW || b.tooltipBgHeight != tooltipH {
			b.tooltipBg = ebiten.NewImage(tooltipW, tooltipH)
			b.tooltipBg.Fill(color.RGBA{30, 30, 30, 220})
			b.tooltipBgWidth = tooltipW
			b.tooltipBgHeight = tooltipH
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(tooltipX), float64(tooltipY))
		screen.DrawImage(b.tooltipBg, op)
		// Draw tooltip text
		text.Draw(screen, b.Tooltip, 14, tooltipX+5, tooltipY+5, color.White)
	}
}
