package minui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

	// Set default style
	bgColor := color.Color(color.RGBA{255, 255, 255, 255})
	borderColor := color.Color(color.RGBA{0, 0, 0, 255})
	borderWidth := 1

	lb.style.BackgroundColor = &bgColor
	lb.style.BorderColor = &borderColor
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
		if lb.HoveredIndex != -1 {
			lb.SelectedIndex = lb.HoveredIndex
			if lb.OnSelect != nil {
				lb.OnSelect(lb.SelectedIndex, lb.Items[lb.SelectedIndex])
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

	// Get absolute position for drawing
	absX, absY := lb.GetAbsolutePosition()
	absBounds := Rect{
		X:      absX,
		Y:      absY,
		Width:  lb.bounds.Width,
		Height: lb.bounds.Height,
	}

	// Draw background
	DrawBackground(screen, absBounds, style)

	// Draw items
	contentBounds := GetContentBounds(absBounds, style)

	// Create clipping region
	clipArea := CreateSubImage(screen, contentBounds)

	startIndex := lb.scrollOffset / lb.itemHeight
	endIndex := startIndex + lb.visibleItems + 1
	if endIndex > len(lb.Items) {
		endIndex = len(lb.Items)
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
			highlightColor := color.RGBA{0, 100, 200, 255}
			DrawRect(clipArea, itemBounds, highlightColor)
		} else if i == lb.HoveredIndex {
			hoverColor := color.RGBA{200, 220, 255, 255}
			DrawRect(clipArea, itemBounds, hoverColor)
		}

		// Draw item text
		textColor := color.RGBA{0, 0, 0, 255}
		if i == lb.SelectedIndex {
			textColor = color.RGBA{255, 255, 255, 255}
		}

		text.Draw(clipArea, lb.Items[i], 14.0, itemBounds.X+4, itemBounds.Y+3, textColor)
	}

	// Draw border
	DrawBorder(screen, absBounds, style)

	// Draw scrollbar if needed
	if len(lb.Items)*lb.itemHeight > contentBounds.Height {
		lb.drawScrollbar(screen, contentBounds, absBounds)
	}
}

// drawScrollbar draws the scrollbar
func (lb *ListBox) drawScrollbar(screen *ebiten.Image, contentBounds Rect, absBounds Rect) {
	scrollbarWidth := 16
	scrollbarX := absBounds.X + absBounds.Width - scrollbarWidth
	scrollbarHeight := contentBounds.Height // Use content height, not absolute bounds height

	// Draw scrollbar track
	trackBounds := Rect{
		X:      scrollbarX,
		Y:      contentBounds.Y, // Start at content area
		Width:  scrollbarWidth,
		Height: scrollbarHeight,
	}
	trackColor := color.RGBA{200, 200, 200, 255}
	DrawRect(screen, trackBounds, trackColor)

	// Calculate thumb size and position
	totalHeight := len(lb.Items) * lb.itemHeight
	thumbHeight := (contentBounds.Height * scrollbarHeight) / totalHeight
	if thumbHeight < 20 {
		thumbHeight = 20
	}

	thumbY := contentBounds.Y + (lb.scrollOffset*(scrollbarHeight-thumbHeight))/(totalHeight-contentBounds.Height)

	// Draw scrollbar thumb
	thumbBounds := Rect{
		X:      scrollbarX + 2,
		Y:      thumbY,
		Width:  scrollbarWidth - 4,
		Height: thumbHeight,
	}
	thumbColor := color.RGBA{120, 120, 120, 255}
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
