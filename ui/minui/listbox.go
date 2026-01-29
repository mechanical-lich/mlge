package minui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mechanical-lich/mlge/event"
	"github.com/mechanical-lich/mlge/text/v2"
)

// ListBox displays a scrollable list of items
type ListBox struct {
	*ElementBase
	Items         []string
	SelectedIndex int
	HoveredIndex  int
	OnSelect      func(index int, item string)
	scrollOffset  int
	itemHeight    int
	visibleItems  int
}

// NewListBox creates a new list box
func NewListBox(id string, items []string) *ListBox {
	lb := &ListBox{
		ElementBase:   NewElementBase(id),
		Items:         items,
		SelectedIndex: -1,
		HoveredIndex:  -1,
		itemHeight:    20,
	}

	// Set default size
	lb.SetSize(200, 150)

	// Set default style - only structural properties, colors come from theme
	borderWidth := 1

	lb.style.BorderWidth = &borderWidth

	return lb
}

// GetType returns the element type
func (lb *ListBox) GetType() string {
	return "ListBox"
}

// Update updates the list box
func (lb *ListBox) Update() {
	if !lb.visible || !lb.enabled {
		return
	}

	lb.UpdateHoverState()

	// Calculate which item is being hovered
	mx, my := ebiten.CursorPosition()
	lb.HoveredIndex = -1

	if lb.IsWithin(mx, my) {
		// Get absolute position for mouse comparison
		absX, absY := lb.GetAbsolutePosition()
		absBounds := Rect{
			X:      absX,
			Y:      absY,
			Width:  lb.bounds.Width,
			Height: lb.bounds.Height,
		}
		contentBounds := GetContentBounds(absBounds, lb.GetComputedStyle())
		relativeY := my - contentBounds.Y + lb.scrollOffset
		index := relativeY / lb.itemHeight

		if index >= 0 && index < len(lb.Items) {
			lb.HoveredIndex = index
		}
	}

	// Handle click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if lb.HoveredIndex >= 0 && lb.HoveredIndex < len(lb.Items) {
			lb.SelectedIndex = lb.HoveredIndex
			if lb.OnSelect != nil && lb.SelectedIndex >= 0 && lb.SelectedIndex < len(lb.Items) {
				lb.OnSelect(lb.SelectedIndex, lb.Items[lb.SelectedIndex])
			}
			// Fire event
			if lb.SelectedIndex >= 0 && lb.SelectedIndex < len(lb.Items) {
				event.GetQueuedInstance().QueueEvent(ListBoxSelectEvent{
					ListBoxID:     lb.GetID(),
					ListBox:       lb,
					SelectedIndex: lb.SelectedIndex,
					SelectedItem:  lb.Items[lb.SelectedIndex],
				})
			}
		}
	}

	// Handle scrolling
	_, dy := ebiten.Wheel()
	if lb.hovered && dy != 0 {
		// Get content bounds to calculate proper scroll limits
		absX, absY := lb.GetAbsolutePosition()
		absBounds := Rect{
			X:      absX,
			Y:      absY,
			Width:  lb.bounds.Width,
			Height: lb.bounds.Height,
		}
		contentBounds := GetContentBounds(absBounds, lb.GetComputedStyle())

		totalHeight := len(lb.Items) * lb.itemHeight

		// Only allow scrolling if content is larger than visible area
		if totalHeight > contentBounds.Height {
			lb.scrollOffset -= int(dy * 20)
			if lb.scrollOffset < 0 {
				lb.scrollOffset = 0
			}
			maxScroll := totalHeight - contentBounds.Height
			if lb.scrollOffset > maxScroll {
				lb.scrollOffset = maxScroll
			}
		}
	}
}

// Layout calculates dimensions
func (lb *ListBox) Layout() {
	style := lb.GetComputedStyle()

	// Apply width/height from style if specified
	width := lb.bounds.Width
	height := lb.bounds.Height

	if style.Width != nil {
		width = *style.Width
	}
	if style.Height != nil {
		height = *style.Height
	}

	// Apply min/max size constraints
	width, height = ApplySizeConstraints(width, height, style)

	lb.bounds.Width = width
	lb.bounds.Height = height

	// Calculate visible items
	contentBounds := GetContentBounds(lb.bounds, style)
	lb.visibleItems = contentBounds.Height / lb.itemHeight
}

// Draw draws the list box
func (lb *ListBox) Draw(screen *ebiten.Image) {
	if !lb.visible {
		return
	}

	style := lb.GetComputedStyle()
	theme := lb.GetTheme()

	// Get absolute position for drawing
	absX, absY := lb.GetAbsolutePosition()
	absBounds := Rect{
		X:      absX,
		Y:      absY,
		Width:  lb.bounds.Width,
		Height: lb.bounds.Height,
	}

	// Draw background with theme support
	DrawBackgroundWithTheme(screen, absBounds, style, theme)

	// Draw items
	contentBounds := GetContentBounds(absBounds, style)

	// Create clipping region
	clipArea := CreateSubImage(screen, contentBounds)

	startIndex := lb.scrollOffset / lb.itemHeight
	endIndex := startIndex + lb.visibleItems + 1
	if endIndex > len(lb.Items) {
		endIndex = len(lb.Items)
	}

	// Get colors from theme
	highlightColor := color.RGBA{0, 100, 200, 255}
	hoverColor := color.RGBA{70, 100, 150, 255}
	if theme != nil {
		highlightColor = colorToRGBA(theme.Colors.Primary)
		hoverColor = colorToRGBA(theme.Colors.Surface)
		hoverColor.R = min(hoverColor.R+20, 255)
		hoverColor.G = min(hoverColor.G+20, 255)
		hoverColor.B = min(hoverColor.B+20, 255)
	}

	for i := startIndex; i < endIndex; i++ {
		itemY := contentBounds.Y + (i * lb.itemHeight) - lb.scrollOffset
		itemBounds := Rect{
			X:      contentBounds.X,
			Y:      itemY,
			Width:  contentBounds.Width,
			Height: lb.itemHeight,
		}

		// Draw selection highlight
		if i == lb.SelectedIndex {
			DrawRect(clipArea, itemBounds, highlightColor)
		} else if i == lb.HoveredIndex {
			DrawRect(clipArea, itemBounds, hoverColor)
		}

		// Draw item text with preference order: explicit style ForegroundColor -> theme Text -> default
		textColor := color.RGBA{255, 255, 255, 255}
		if style != nil && style.ForegroundColor != nil {
			textColor = colorToRGBA(*style.ForegroundColor)
		} else if theme != nil {
			textColor = colorToRGBA(theme.Colors.Text)
		}

		if i == lb.SelectedIndex {
			// Selected items should have contrasting text
			textColor = color.RGBA{255, 255, 255, 255}
		}

		text.Draw(clipArea, lb.Items[i], 14.0, itemBounds.X+4, itemBounds.Y+3, textColor)
	}

	// Draw border with theme support
	DrawBorderWithTheme(screen, absBounds, style, theme)

	// Draw scrollbar if needed
	if len(lb.Items)*lb.itemHeight > contentBounds.Height {
		lb.drawScrollbar(screen, contentBounds, absBounds, theme)
	}
}

// drawScrollbar draws the scrollbar
func (lb *ListBox) drawScrollbar(screen *ebiten.Image, contentBounds Rect, absBounds Rect, theme *Theme) {
	scrollbarWidth := 16
	scrollbarX := absBounds.X + absBounds.Width - scrollbarWidth
	scrollbarHeight := contentBounds.Height // Use content height, not absolute bounds height

	// Draw scrollbar track with theme color
	trackBounds := Rect{
		X:      scrollbarX,
		Y:      contentBounds.Y, // Start at content area
		Width:  scrollbarWidth,
		Height: scrollbarHeight,
	}
	trackColor := color.RGBA{40, 40, 50, 255}
	if theme != nil {
		trackColor = colorToRGBA(theme.Colors.Surface)
	}
	DrawRect(screen, trackBounds, trackColor)

	// Calculate thumb size and position
	totalHeight := len(lb.Items) * lb.itemHeight
	thumbHeight := (contentBounds.Height * scrollbarHeight) / totalHeight
	if thumbHeight < 20 {
		thumbHeight = 20
	}

	thumbY := contentBounds.Y + (lb.scrollOffset*(scrollbarHeight-thumbHeight))/(totalHeight-contentBounds.Height)

	// Draw scrollbar thumb with theme color
	thumbBounds := Rect{
		X:      scrollbarX + 2,
		Y:      thumbY,
		Width:  scrollbarWidth - 4,
		Height: thumbHeight,
	}
	thumbColor := color.RGBA{120, 120, 120, 255}
	if theme != nil {
		thumbColor = colorToRGBA(theme.Colors.Border)
	}
	DrawRoundedRect(screen, thumbBounds, 4, thumbColor)
}

// SetItems sets the list items
func (lb *ListBox) SetItems(items []string) {
	lb.Items = items
	lb.SelectedIndex = -1
	lb.HoveredIndex = -1
	lb.scrollOffset = 0
}

// GetSelectedItem returns the currently selected item
func (lb *ListBox) GetSelectedItem() (int, string) {
	if lb.SelectedIndex >= 0 && lb.SelectedIndex < len(lb.Items) {
		return lb.SelectedIndex, lb.Items[lb.SelectedIndex]
	}
	return -1, ""
}
