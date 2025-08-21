package ui

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mechanical-lich/mlge/resource"
	"github.com/mechanical-lich/mlge/text/v2"
	theming "github.com/mechanical-lich/mlge/ui/v2/theming"
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

	// Cached images to avoid per-frame allocations
	dropdownBg     *ebiten.Image
	dropdownBgW    int
	dropdownBgH    int
	highlightImg   *ebiten.Image
	highlightImgW  int
	highlightImgH  int
	scrollBarBg    *ebiten.Image
	scrollBarBgW   int
	scrollBarBgH   int
	scrollThumbImg *ebiten.Image
	scrollThumbW   int
	scrollThumbH   int
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
			Op:     &ebiten.DrawImageOptions{},
		},
		Options:      options,
		Selected:     0,
		Open:         false,
		ScrollOffset: 0,
		VisibleItems: visibleItems,
		ItemHeight:   int(h + 6),
	}
}

func (s *SelectBox) Update() {
	cX, cY := ebiten.CursorPosition()
	absX, absY := s.GetScreenPosition()
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

func (s *SelectBox) Draw(screen *ebiten.Image, theme *theming.Theme) {
	absX, absY := s.GetAbsolutePosition()

	// Draw main box (closed or open)
	s.Op.GeoM.Reset()
	s.Op.GeoM.Scale(float64(s.Width)/float64(theme.InputField.Width), float64(s.Height)/float64(theme.InputField.Height))
	s.Op.GeoM.Translate(float64(absX), float64(absY))
	screen.DrawImage(resource.GetSubImage("ui", theme.InputField.SrcX, theme.InputField.SrcY, theme.InputField.Width, theme.InputField.Height), s.Op)

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
		// Cache/recreate dropdown background
		if s.dropdownBg == nil || s.dropdownBgW != s.Width || s.dropdownBgH != listH {
			s.dropdownBg = ebiten.NewImage(s.Width, listH)
			s.dropdownBg.Fill(color.RGBA{40, 40, 40, 240})
			s.dropdownBgW = s.Width
			s.dropdownBgH = listH
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(absX), float64(absY+s.Height))
		screen.DrawImage(s.dropdownBg, op)

		start := s.ScrollOffset
		end := int(math.Min(float64(start+s.VisibleItems), float64(len(s.Options))))
		for i := start; i < end; i++ {
			optY := absY + s.Height + (i-start)*s.ItemHeight
			var optColor color.Color
			optColor = color.White
			if i == s.Selected {
				// Cache/recreate highlight image
				if s.highlightImg == nil || s.highlightImgW != s.Width || s.highlightImgH != s.ItemHeight {
					s.highlightImg = ebiten.NewImage(s.Width, s.ItemHeight)
					s.highlightImg.Fill(color.RGBA{80, 80, 120, 180})
					s.highlightImgW = s.Width
					s.highlightImgH = s.ItemHeight
				}
				op2 := &ebiten.DrawImageOptions{}
				op2.GeoM.Translate(float64(absX), float64(optY))
				screen.DrawImage(s.highlightImg, op2)
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

func (s *SelectBox) drawScrollbar(screen *ebiten.Image, x, y, listH int, theme *theming.Theme) {
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
	// Cache/recreate scrollbar background
	if s.scrollBarBg == nil || s.scrollBarBgW != barW || s.scrollBarBgH != barH {
		s.scrollBarBg = ebiten.NewImage(barW, barH)
		s.scrollBarBg.Fill(color.RGBA{60, 60, 60, 180})
		s.scrollBarBgW = barW
		s.scrollBarBgH = barH
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(barX), float64(barY))
	screen.DrawImage(s.scrollBarBg, op)
	// Cache/recreate thumb
	if s.scrollThumbImg == nil || s.scrollThumbW != barW || s.scrollThumbH != thumbH {
		s.scrollThumbImg = ebiten.NewImage(barW, thumbH)
		s.scrollThumbImg.Fill(color.RGBA{160, 160, 160, 220})
		s.scrollThumbW = barW
		s.scrollThumbH = thumbH
	}
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(float64(barX), float64(thumbY))
	screen.DrawImage(s.scrollThumbImg, op2)
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
