package minui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mechanical-lich/mlge/text"
)

// Overflow-tooltip behavior: when a widget renders clipped text, hovering the
// widget for a short delay reveals a small floating tooltip showing the full
// text. State is keyed by the owning element's ID so we can persist hover
// frame counters across redraws without each widget tracking its own.

const (
	overflowTooltipDelayFrames = 30 // ~0.5s at 60fps
	overflowTooltipFontSize    = 13.0
	overflowTooltipPadX        = 6
	overflowTooltipPadY        = 4
	overflowTooltipCursorOff   = 14
)

type overflowTooltipState struct {
	text         string
	hoverFrames  int
	lastSeenFrame uint64
}

var (
	overflowTooltipStates = map[string]*overflowTooltipState{}
	overflowFrameCounter  uint64
)

// tickOverflowTooltips is called once per frame (from flushOverlays) to age
// out states whose owner widget didn't draw this frame.
func tickOverflowTooltips() {
	overflowFrameCounter++
	for id, st := range overflowTooltipStates {
		if st.lastSeenFrame != overflowFrameCounter-1 && st.lastSeenFrame != overflowFrameCounter {
			delete(overflowTooltipStates, id)
		}
	}
}

// DrawClippedWithTooltip draws txt clipped to maxW, and registers an overflow
// tooltip if the text was actually truncated and the owner is hovered.
// Pass the owner element so we can key hover state and bound the hover region.
// If owner is nil this behaves like text.DrawClipped.
func DrawClippedWithTooltip(screen *ebiten.Image, owner Element, txt string, size float64, x, y, maxW int, clr color.Color) {
	text.DrawClipped(screen, txt, size, x, y, maxW, clr)
	if owner == nil || maxW <= 0 || txt == "" {
		return
	}

	// Was the text actually truncated?
	tw, _ := text.Measure(txt, size)
	if int(tw) <= maxW {
		// Reset hover state for this owner so the timer doesn't carry over.
		delete(overflowTooltipStates, owner.GetID())
		return
	}

	// Is the cursor inside the owner?
	mx, my := ebiten.CursorPosition()
	absX, absY := owner.GetAbsolutePosition()
	b := owner.GetBounds()
	hovered := mx >= absX && mx < absX+b.Width && my >= absY && my < absY+b.Height
	if !hovered {
		delete(overflowTooltipStates, owner.GetID())
		return
	}

	st := overflowTooltipStates[owner.GetID()]
	if st == nil {
		st = &overflowTooltipState{text: txt}
		overflowTooltipStates[owner.GetID()] = st
	}
	st.text = txt
	st.lastSeenFrame = overflowFrameCounter
	st.hoverFrames++

	if st.hoverFrames >= overflowTooltipDelayFrames {
		// Capture for the overlay closure.
		fullText := txt
		mouseX, mouseY := mx, my
		QueueOverlay(func(s *ebiten.Image) {
			drawOverflowTooltip(s, fullText, mouseX, mouseY)
		})
	}
}

func drawOverflowTooltip(screen *ebiten.Image, fullText string, mx, my int) {
	tw, th := text.Measure(fullText, overflowTooltipFontSize)
	w := int(tw) + overflowTooltipPadX*2
	h := int(th) + overflowTooltipPadY*2

	// Position near cursor, then nudge inside the screen.
	x := mx + overflowTooltipCursorOff
	y := my + overflowTooltipCursorOff
	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	if x+w > sw {
		x = sw - w
	}
	if y+h > sh {
		y = sh - h
	}
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	bg := color.RGBA{20, 22, 28, 240}
	border := color.RGBA{120, 130, 150, 255}
	textColor := color.RGBA{235, 235, 240, 255}

	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), bg, false)
	vector.StrokeRect(screen, float32(x), float32(y), float32(w), float32(h), 1, border, false)
	text.Draw(screen, fullText, overflowTooltipFontSize, x+overflowTooltipPadX, y+overflowTooltipPadY, textColor)
}
