package ui

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/resource"
	"github.com/mechanical-lich/mlge/text/v2"
	"github.com/mechanical-lich/mlge/utility"
)

type ScrollingTextArea struct {
	ElementBase
	Text         string
	ScrollOffset int
	LineHeight   int
	VisibleLines int
	Lines        []string
}

func NewScrollingTextArea(name string, x, y, width, height int, txt string) *ScrollingTextArea {
	lines := text.Wrap(txt, width-16, 15) // Subtract padding for wrap
	lineHeight := 18                      // Adjust as needed for your font
	visibleLines := int(math.Floor(float64(height-10) / float64(lineHeight)))
	return &ScrollingTextArea{
		ElementBase: ElementBase{
			Name:   name,
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
			op:     &ebiten.DrawImageOptions{},
		},
		Text:         txt,
		ScrollOffset: 0,
		LineHeight:   lineHeight,
		VisibleLines: visibleLines,
		Lines:        lines,
	}
}

func (s *ScrollingTextArea) Update(parentX, parentY int) {
	// Simple mouse wheel scroll
	xoff, _ := ebiten.Wheel()
	if xoff != 0 {
		s.ScrollOffset -= int(xoff)
		if s.ScrollOffset < 0 {
			s.ScrollOffset = 0
		}
		maxOffset := len(s.Lines) - s.VisibleLines
		if maxOffset < 0 {
			maxOffset = 0
		}
		if s.ScrollOffset > maxOffset {
			s.ScrollOffset = maxOffset
		}
	}
}

func (s *ScrollingTextArea) Draw(screen *ebiten.Image, parentX, parentY int, theme *Theme) {
	// Draw 9-slice background
	if s.op == nil {
		s.op = &ebiten.DrawImageOptions{}
	}
	s.op.GeoM.Reset()
	s.op.GeoM.Translate(float64(s.X+parentX), float64(s.Y+parentY))
	bg := ebiten.NewImage(s.Width, s.Height)
	utility.Draw9Slice(
		bg,
		0, 0, s.Width, s.Height,
		theme.ScrollingTextArea.SrcX,
		theme.ScrollingTextArea.SrcY,
		theme.ScrollingTextArea.TileSize,
		theme.ScrollingTextArea.TileScale,
	)
	screen.DrawImage(bg, s.op)

	// Draw text lines
	start := s.ScrollOffset
	end := start + s.VisibleLines
	if end > len(s.Lines) {
		end = len(s.Lines)
	}
	for i := start; i < end; i++ {
		text.Draw(
			screen,
			s.Lines[i],
			15,
			s.X+8+parentX,
			s.Y+8+parentY+(i-start)*s.LineHeight,
			color.White,
		)
	}

	// Draw scrollbar if needed
	if len(s.Lines) > s.VisibleLines {
		s.drawScrollbar(screen, parentX, parentY, theme)
	}
}

func (s *ScrollingTextArea) drawScrollbar(screen *ebiten.Image, parentX, parentY int, theme *Theme) {
	barX := s.X + s.Width - 12 + parentX
	barY := s.Y + 4 + parentY
	barW := 8
	barH := s.Height - 8

	// Draw scrollbar background (9-slice not needed for thin bar)
	barBg := resource.GetSubImage("ui", 96, 48, 16, 48)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(barW)/16.0, float64(barH)/48.0)
	op.GeoM.Translate(float64(barX), float64(barY))
	screen.DrawImage(barBg, op)

	// Draw thumb
	thumbH := int(math.Max(float64(barH*s.VisibleLines/len(s.Lines)), 16))
	thumbY := barY + (barH-thumbH)*s.ScrollOffset/(len(s.Lines)-s.VisibleLines)
	thumb := resource.GetSubImage("ui", 112, 48, 8, 16)
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Scale(float64(barW)/8.0, float64(thumbH)/16.0)
	op2.GeoM.Translate(float64(barX), float64(thumbY))
	screen.DrawImage(thumb, op2)
}

func (s *ScrollingTextArea) AddText(txt string) {
	s.Text += "\n" + txt
	s.Lines = text.Wrap(s.Text, 20, 15) // Re-wrap with new text

	if len(s.Lines) > s.VisibleLines {
		s.ScrollOffset = len(s.Lines) - s.VisibleLines // Auto-scroll to bottom
	}
}
