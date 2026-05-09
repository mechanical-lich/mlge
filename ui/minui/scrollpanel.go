package minui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ScrollPanel is a vertical-scrolling container. It clips children to its
// bounds and offers a mouse-wheel-driven scrollbar. Children may be any
// Element; their relative Y is interpreted against the un-scrolled origin.
type ScrollPanel struct {
	*ElementBase
	scrollOffset int
	contentH     int
	scrollStep   int
}

// NewScrollPanel creates a new scroll panel.
func NewScrollPanel(id string) *ScrollPanel {
	return &ScrollPanel{
		ElementBase: NewElementBase(id),
		scrollStep:  20,
	}
}

func (s *ScrollPanel) GetType() string { return "ScrollPanel" }

// GetScrollOffsetY exposes the scroll offset for child absolute-position math.
func (s *ScrollPanel) GetScrollOffsetY() int { return s.scrollOffset }

// AddChild appends a child element.
func (s *ScrollPanel) AddChild(child Element) {
	s.children = append(s.children, child)
	child.SetParent(s)
	if s.theme != nil {
		if setter, ok := child.(interface{ SetTheme(*Theme) }); ok {
			setter.SetTheme(s.theme)
		}
	}
}

// RemoveChild removes a child element.
func (s *ScrollPanel) RemoveChild(child Element) {
	for i, c := range s.children {
		if c == child {
			s.children = append(s.children[:i], s.children[i+1:]...)
			child.SetParent(nil)
			return
		}
	}
}

// Update updates children and handles wheel scrolling when hovered.
func (s *ScrollPanel) Update() {
	if !s.visible || !s.enabled {
		return
	}
	s.UpdateHoverState()

	// Recompute content height each frame from children's bottom edges.
	maxBottom := 0
	for _, child := range s.children {
		if !child.IsVisible() {
			continue
		}
		b := child.GetBounds()
		bottom := b.Y + b.Height
		if bottom > maxBottom {
			maxBottom = bottom
		}
	}
	s.contentH = maxBottom

	mx, my := ebiten.CursorPosition()
	hovered := s.IsWithin(mx, my)
	_, dy := ebiten.Wheel()
	if hovered && dy != 0 {
		visibleH := s.bounds.Height
		if s.contentH > visibleH {
			s.scrollOffset -= int(dy * float64(s.scrollStep))
			s.clampScroll()
		}
	}
	s.clampScroll()

	for _, child := range s.children {
		child.Update()
	}
}

func (s *ScrollPanel) clampScroll() {
	if s.scrollOffset < 0 {
		s.scrollOffset = 0
	}
	maxScroll := s.contentH - s.bounds.Height
	if maxScroll < 0 {
		maxScroll = 0
	}
	if s.scrollOffset > maxScroll {
		s.scrollOffset = maxScroll
	}
}

// Layout lays children at their assigned positions (no automatic stacking).
func (s *ScrollPanel) Layout() {
	for _, child := range s.children {
		if l, ok := child.(interface{ Layout() }); ok {
			l.Layout()
		}
	}
}

// Draw renders the panel and its visible children, clipped to bounds.
func (s *ScrollPanel) Draw(screen *ebiten.Image) {
	if !s.visible {
		return
	}
	absX, absY := s.GetAbsolutePosition()
	absBounds := Rect{X: absX, Y: absY, Width: s.bounds.Width, Height: s.bounds.Height}

	// Background (if styled)
	style := s.GetComputedStyle()
	theme := s.GetTheme()
	DrawBackgroundWithTheme(screen, absBounds, style, theme)
	DrawBorderWithTheme(screen, absBounds, style, theme)

	// Clip children to our bounds.
	clipRect := image.Rect(absX, absY, absX+s.bounds.Width, absY+s.bounds.Height)
	clip := screen.SubImage(clipRect).(*ebiten.Image)
	for _, child := range s.children {
		child.Draw(clip)
	}

	// Scrollbar on the right edge if content overflows.
	if s.contentH > s.bounds.Height {
		s.drawScrollbar(screen, absBounds, theme)
	}
}

func (s *ScrollPanel) drawScrollbar(screen *ebiten.Image, absBounds Rect, theme *Theme) {
	const trackW = 6
	trackX := float32(absBounds.X + absBounds.Width - trackW - 2)
	trackY := float32(absBounds.Y + 2)
	trackH := float32(absBounds.Height - 4)

	trackColor := color.RGBA{40, 40, 50, 200}
	thumbColor := color.RGBA{120, 120, 140, 230}
	if theme != nil {
		trackColor = colorToRGBA(theme.Colors.Surface)
		thumbColor = colorToRGBA(theme.Colors.Border)
	}

	vector.DrawFilledRect(screen, trackX, trackY, trackW, trackH, trackColor, false)

	visibleH := absBounds.Height
	thumbH := float32(visibleH) * trackH / float32(s.contentH)
	if thumbH < 12 {
		thumbH = 12
	}
	maxScroll := s.contentH - visibleH
	thumbY := trackY
	if maxScroll > 0 {
		thumbY += (trackH - thumbH) * float32(s.scrollOffset) / float32(maxScroll)
	}
	vector.DrawFilledRect(screen, trackX, thumbY, trackW, thumbH, thumbColor, false)
}

// ScrollTo sets the scroll offset directly.
func (s *ScrollPanel) ScrollTo(y int) {
	s.scrollOffset = y
	s.clampScroll()
}

// ScrollOffsetY returns the current scroll offset.
func (s *ScrollPanel) ScrollOffsetY() int { return s.scrollOffset }
