package ui

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mechanical-lich/mlge/resource"
	"github.com/mechanical-lich/mlge/text/v2"
)

type SelectBox struct {
	ElementBase
	Options      []string
	Selected     int
	Open         bool
	ScrollOffset int
	VisibleItems int
	ItemHeight   int
	OnChange     func(selected int)
}

func NewSelectBox(name string, x, y, width, visibleItems int, options []string) *SelectBox {
	_, h := text.Measure("A", 16)
	return &SelectBox{
		ElementBase: ElementBase{
			Name:   name,
			X:      x,
			Y:      y,
			Width:  width,
			Height: int(h + 10),
			op:     &ebiten.DrawImageOptions{},
		},
		Options:      options,
		Selected:     0,
		Open:         false,
		ScrollOffset: 0,
		VisibleItems: visibleItems,
		ItemHeight:   int(h + 6),
	}
}

func (s *SelectBox) Update(parentX, parentY int) {
	cX, cY := ebiten.CursorPosition()
	absX := s.X + parentX
	absY := s.Y + parentY

	// Toggle open/close on click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if !s.Open && cX >= absX && cX <= absX+s.Width && cY >= absY && cY <= absY+s.Height {
			s.Open = true
			s.Focused = true
		} else if s.Open {
			// Click inside options list?
			listY := absY + s.Height
			if cX >= absX && cX <= absX+s.Width && cY >= listY && cY <= listY+s.ItemHeight*s.VisibleItems {
				idx := (cY - listY) / s.ItemHeight
				optIdx := s.ScrollOffset + idx
				if optIdx >= 0 && optIdx < len(s.Options) {
					s.Selected = optIdx
					if s.OnChange != nil {
						s.OnChange(s.Selected)
					}
				}
			}
			s.Open = false
			s.Focused = false
		} else {
			s.Open = false
			s.Focused = false
		}
	}

	// Scroll if open and mouse wheel used
	if s.Open && s.Focused {
		_, yoff := ebiten.Wheel()
		if yoff != 0 {
			s.ScrollOffset -= int(yoff)
			maxOffset := len(s.Options) - s.VisibleItems
			if maxOffset < 0 {
				maxOffset = 0
			}
			if s.ScrollOffset < 0 {
				s.ScrollOffset = 0
			}
			if s.ScrollOffset > maxOffset {
				s.ScrollOffset = maxOffset
			}
		}
	}
}

func (s *SelectBox) Draw(screen *ebiten.Image, parentX, parentY int, theme *Theme) {
	absX := s.X + parentX
	absY := s.Y + parentY

	// Draw main box (closed or open)
	s.op.GeoM.Reset()
	s.op.GeoM.Scale(float64(s.Width)/float64(theme.InputField.Width), float64(s.Height)/float64(theme.InputField.Height))
	s.op.GeoM.Translate(float64(absX), float64(absY))
	screen.DrawImage(resource.GetSubImage("ui", theme.InputField.SrcX, theme.InputField.SrcY, theme.InputField.Width, theme.InputField.Height), s.op)

	// Draw selected option
	if s.Selected >= 0 && s.Selected < len(s.Options) {
		text.Draw(screen, s.Options[s.Selected], 15, absX+6, absY+5, color.White)
	}

	// Draw dropdown arrow
	arrowX := absX + s.Width - 18
	arrowY := absY + (s.Height / 2) - 3
	// Simple triangle for arrow
	ebitenutilDrawTriangle(screen, arrowX, arrowY, 12, color.White)

	// Draw dropdown list if open
	if s.Open {
		listH := s.ItemHeight * s.VisibleItems
		bg := ebiten.NewImage(s.Width, listH)
		bg.Fill(color.RGBA{40, 40, 40, 240})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(absX), float64(absY+s.Height))
		screen.DrawImage(bg, op)

		start := s.ScrollOffset
		end := int(math.Min(float64(start+s.VisibleItems), float64(len(s.Options))))
		for i := start; i < end; i++ {
			optY := absY + s.Height + (i-start)*s.ItemHeight
			var optColor color.Color
			optColor = color.White
			if i == s.Selected {
				// Highlight selected
				highlight := ebiten.NewImage(s.Width, s.ItemHeight)
				highlight.Fill(color.RGBA{80, 80, 120, 180})
				op2 := &ebiten.DrawImageOptions{}
				op2.GeoM.Translate(float64(absX), float64(optY))
				screen.DrawImage(highlight, op2)
				optColor = color.RGBA{255, 255, 180, 255}
			}
			text.Draw(screen, s.Options[i], 15, absX+6, optY+4, optColor)
		}

		// Draw scrollbar if needed
		if len(s.Options) > s.VisibleItems {
			s.drawScrollbar(screen, absX, absY+s.Height, listH, theme)
		}
	}
}

func (s *SelectBox) drawScrollbar(screen *ebiten.Image, x, y, listH int, theme *Theme) {
	barX := x + s.Width - 12
	barY := y
	barW := 10
	barH := listH

	total := len(s.Options)
	if total <= s.VisibleItems {
		return
	}
	thumbH := int(math.Max(float64(barH*s.VisibleItems/total), 16))
	scrollRange := barH - thumbH
	var thumbY int
	if total > s.VisibleItems && scrollRange > 0 {
		thumbY = barY + (scrollRange*s.ScrollOffset)/(total-s.VisibleItems)
	} else {
		thumbY = barY
	}
	// Draw scrollbar background
	barBg := ebiten.NewImage(barW, barH)
	barBg.Fill(color.RGBA{60, 60, 60, 180})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(barX), float64(barY))
	screen.DrawImage(barBg, op)
	// Draw thumb
	thumb := ebiten.NewImage(barW, thumbH)
	thumb.Fill(color.RGBA{160, 160, 160, 220})
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(float64(barX), float64(thumbY))
	screen.DrawImage(thumb, op2)
}

// Helper: draw a downward triangle for the dropdown arrow
func ebitenutilDrawTriangle(screen *ebiten.Image, x, y, size int, clr color.Color) {
	vs := []ebiten.Vertex{
		{DstX: float32(x), DstY: float32(y)},
		{DstX: float32(x + size), DstY: float32(y)},
		{DstX: float32(x + size/2), DstY: float32(y + size)},
	}
	is := []uint16{0, 1, 2}
	screen.DrawTriangles(vs, is, ebiten.NewImage(1, 1), &ebiten.DrawTrianglesOptions{FillRule: ebiten.EvenOdd})
}

func (s *SelectBox) SetOptions(options []string) {
	s.Options = options
	if s.Selected >= len(options) {
		s.Selected = 0
	}
	if s.ScrollOffset > len(options)-s.VisibleItems {
		s.ScrollOffset = 0
	}
}
