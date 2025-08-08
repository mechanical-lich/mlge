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

	draggingThumb bool
	dragOffsetY   int
	bg            *ebiten.Image
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
	// Mouse position
	mx, my := ebiten.CursorPosition()
	s.Focused = s.IsWithin(mx, my, parentX, parentY)
	mx -= parentX
	my -= parentY

	barX := s.X + s.Width - 12 + parentX
	barY := s.Y + 4 + parentY
	barW := 32
	barH := s.Height - 8

	// Thumb calculations
	totalLines := len(s.Lines)
	if totalLines <= s.VisibleLines {
		s.draggingThumb = false
		return
	}
	thumbH := int(math.Max(float64(barH*s.VisibleLines/totalLines), 16))
	maxThumbY := barY + barH - thumbH
	thumbY := barY
	if totalLines > s.VisibleLines {
		thumbY = barY + (barH-thumbH)*s.ScrollOffset/(totalLines-s.VisibleLines)
	}

	left := barX
	right := barX + barW
	top := thumbY
	bottom := thumbY + thumbH

	// Mouse input
	mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	//	mouseJustPressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !s.draggingThumb
	mouseJustReleased := !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && s.draggingThumb

	if s.draggingThumb {
		if mouseJustReleased {
			s.draggingThumb = false
		} else if mousePressed {
			// Calculate new thumb position
			newThumbY := my - s.dragOffsetY
			if newThumbY < barY {
				newThumbY = barY
			}
			if newThumbY > maxThumbY {
				newThumbY = maxThumbY
			}
			// Map thumb position to scroll offset
			scrollRange := barH - thumbH
			if scrollRange > 0 {
				s.ScrollOffset = int(float64(newThumbY-barY) / float64(scrollRange) * float64(totalLines-s.VisibleLines))
			}
			if s.ScrollOffset < 0 {
				s.ScrollOffset = 0
			}
			if s.ScrollOffset > totalLines-s.VisibleLines {
				s.ScrollOffset = totalLines - s.VisibleLines
			}
		}
	} else if mousePressed && mx >= left && mx < right && my >= top && my < bottom {
		s.draggingThumb = true
		s.dragOffsetY = my - thumbY
	}

	// Mouse wheel scroll (only if not dragging)

	if !s.draggingThumb && s.Focused {
		_, yoff := ebiten.Wheel()
		if yoff != 0 {
			s.ScrollOffset -= int(yoff)
			if s.ScrollOffset < 0 {
				s.ScrollOffset = 0
			}
			maxOffset := totalLines - s.VisibleLines
			if maxOffset < 0 {
				maxOffset = 0
			}
			if s.ScrollOffset > maxOffset {
				s.ScrollOffset = maxOffset
			}
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
	if s.bg == nil {
		bg := ebiten.NewImage(s.Width, s.Height)
		s.bg = bg
	}
	utility.Draw9Slice(
		s.bg,
		0, 0, s.Width, s.Height,
		theme.ScrollingTextArea.SrcX,
		theme.ScrollingTextArea.SrcY,
		theme.ScrollingTextArea.TileSize,
		theme.ScrollingTextArea.TileScale,
	)
	screen.DrawImage(s.bg, s.op)

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
	barW := 32
	barH := s.Height - 8

	totalLines := len(s.Lines)
	if totalLines <= s.VisibleLines {
		return
	}

	// Draw scrollbar background
	barBg := resource.GetSubImage("ui", theme.ScrollingTextArea.ScrollBarX, theme.ScrollingTextArea.ScrollBarY, theme.ScrollingTextArea.ScrollBarWidth, theme.ScrollingTextArea.ScrollBarHeight)
	s.op.GeoM.Reset()
	s.op.GeoM.Scale(float64(barW)/16.0, float64(barH)/48.0)
	s.op.GeoM.Translate(float64(barX), float64(barY))
	screen.DrawImage(barBg, s.op)

	// Draw thumb
	thumbH := int(math.Max(float64(barH*s.VisibleLines/totalLines), 16))
	scrollRange := barH - thumbH
	var thumbY int
	if totalLines > s.VisibleLines && scrollRange > 0 {
		thumbY = barY + (scrollRange*s.ScrollOffset)/(totalLines-s.VisibleLines)
	} else {
		thumbY = barY
	}
	thumb := resource.GetSubImage("ui", theme.ScrollingTextArea.ThumbX, theme.ScrollingTextArea.ThumbY, theme.ScrollingTextArea.ThumbWidth, theme.ScrollingTextArea.ThumbHeight)
	s.op.GeoM.Reset()
	s.op.GeoM.Scale(float64(barW)/16.0, float64(thumbH)/16.0)
	s.op.GeoM.Translate(float64(barX), float64(thumbY))
	screen.DrawImage(thumb, s.op)
}

func (s *ScrollingTextArea) AddText(txt string) {
	s.Text += "\n" + txt
	// Wrap text to width of scrolling area
	s.Lines = text.Wrap(s.Text, (s.Width-32)/5, 15) // Subtract padding for wrap

	if len(s.Lines) > s.VisibleLines {
		s.ScrollOffset = len(s.Lines) - s.VisibleLines // Auto-scroll to bottom
	}
}
