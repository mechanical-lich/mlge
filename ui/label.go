package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/text/v2"
)

type Label struct {
	ElementBase
	Text     string
	FontSize float64
	Color    color.Color
}

func NewLabel(name string, x int, y int, txt string) *Label {
	size := 14.0
	w, h := text.Measure(txt, size)

	l := &Label{
		ElementBase: ElementBase{
			Name:   name,
			X:      x,
			Y:      y,
			Width:  int(w + 2),
			Height: int(h),
			op:     &ebiten.DrawImageOptions{},
		},
		Text:     txt,
		FontSize: size,
		Color:    color.White,
	}

	return l
}

func NewLabelWithSizeAndColor(name string, x int, y int, txt string, size float64, col color.Color) *Label {
	w, h := text.Measure(txt, size)

	l := &Label{
		ElementBase: ElementBase{
			Name:   name,
			X:      x,
			Y:      y,
			Width:  int(w + 2),
			Height: int(h),
			op:     &ebiten.DrawImageOptions{},
		},
		Text:     txt,
		FontSize: size,
		Color:    col,
	}

	return l
}

func (l *Label) Update(parentX, parentY int) {
	// Labels are static by default; keep for API consistency
}

func (l *Label) Draw(screen *ebiten.Image, parentX, parentY int, theme *Theme) {
	// Draw text at label position (no internal padding by default)
	text.Draw(screen, l.Text, l.FontSize, l.X+parentX, l.Y+parentY, l.Color)
}

func (l *Label) SetText(txt string) {
	l.Text = txt
	w, h := text.Measure(txt, l.FontSize)
	l.Width = int(w + 2)
	l.Height = int(h)
}
