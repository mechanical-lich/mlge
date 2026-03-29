package minui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/text"
)

// TextSpan is a single styled line within a RichText widget.
type TextSpan struct {
	Text          string
	Color         color.Color // nil = use widget DefaultColor
	Size          int         // 0 = use widget DefaultSize
	Strikethrough bool
	Indent        int // extra left indent in pixels
}

// RichText displays a list of independently styled text spans.
// Each span renders as one line. Call Layout() after modifying Spans
// to recalculate the widget height before the next Draw.
type RichText struct {
	*ElementBase
	Spans        []TextSpan
	DefaultSize  int         // font size used when span.Size == 0 (default 14)
	DefaultColor color.Color // color used when span.Color == nil (default white)
	LineHeight   int         // pixels between lines (0 = DefaultSize + 4)
}

// NewRichText creates a new RichText widget with the given width.
// Height is calculated dynamically from content.
func NewRichText(id string, width int) *RichText {
	rt := &RichText{
		ElementBase:  NewElementBase(id),
		DefaultSize:  14,
		DefaultColor: color.RGBA{255, 255, 255, 255},
		LineHeight:   18,
	}
	rt.SetSize(width, 0)
	return rt
}

// GetType returns the element type.
func (rt *RichText) GetType() string { return "RichText" }

// AddSpan appends a span and recalculates height.
func (rt *RichText) AddSpan(span TextSpan) {
	rt.Spans = append(rt.Spans, span)
	rt.recalcHeight()
}

// Clear removes all spans.
func (rt *RichText) Clear() {
	rt.Spans = rt.Spans[:0]
	rt.bounds.Height = 0
}

func (rt *RichText) lineH() int {
	if rt.LineHeight > 0 {
		return rt.LineHeight
	}
	return rt.DefaultSize + 4
}

func (rt *RichText) recalcHeight() {
	rt.bounds.Height = len(rt.Spans) * rt.lineH()
}

// Update is a no-op; RichText is display-only.
func (rt *RichText) Update() {}

// Layout recalculates the widget height from current spans.
func (rt *RichText) Layout() {
	rt.recalcHeight()
}

// Draw renders all spans.
func (rt *RichText) Draw(screen *ebiten.Image) {
	if !rt.visible || len(rt.Spans) == 0 {
		return
	}
	absX, absY := rt.GetAbsolutePosition()
	lh := rt.lineH()

	for i, span := range rt.Spans {
		sz := rt.DefaultSize
		if span.Size > 0 {
			sz = span.Size
		}

		col := rt.DefaultColor
		if span.Color != nil {
			col = span.Color
		}
		rgba := colorToRGBA(col)

		x := absX + span.Indent
		y := absY + i*lh

		text.Draw(screen, span.Text, float64(sz), x, y, rgba)

		if span.Strikethrough {
			w, _ := text.Measure(span.Text, float64(sz))
			DrawRect(screen, Rect{
				X:      x,
				Y:      y + sz/2,
				Width:  int(w),
				Height: 1,
			}, rgba)
		}
	}
}
