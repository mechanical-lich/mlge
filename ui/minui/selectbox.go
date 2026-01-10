package minui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mechanical-lich/mlge/event"
	"github.com/mechanical-lich/mlge/text/v2"
)

// SelectBox is a drop-down selection widget modeled after HTML select
type SelectBox struct {
	*ElementBase
	Items         []string
	SelectedIndex int
	HoveredIndex  int
	OnSelect      func(index int, item string)
	expanded      bool
	listBox       *ListBox
	itemHeight    int
}

// NewSelectBox creates a new SelectBox
func NewSelectBox(id string, items []string) *SelectBox {
	sb := &SelectBox{
		ElementBase:   NewElementBase(id),
		Items:         items,
		SelectedIndex: -1,
		HoveredIndex:  -1,
		expanded:      false,
		itemHeight:    20,
	}

	sb.SetSize(140, 28)

	// Default style - only structural properties, colors come from theme
	borderWidth := 1
	padding := NewEdgeInsets(6)

	sb.style.BorderWidth = &borderWidth
	sb.style.Padding = padding

	// Create internal ListBox (floating) to show when expanded
	lb := NewListBox(id+"_list", items)
	lb.itemHeight = sb.itemHeight
	lb.SetSize(sb.bounds.Width, sb.itemHeight*5) // default visible items: 5
	lb.OnSelect = func(idx int, item string) {
		sb.setSelectedIndex(idx)
		// Collapse dropdown after selecting
		sb.expanded = false
	}
	sb.listBox = lb

	return sb
}

// GetType returns element type
func (sb *SelectBox) GetType() string {
	return "SelectBox"
}

// Update handles interactions
func (sb *SelectBox) Update() {
	if !sb.visible || !sb.enabled {
		return
	}

	sb.UpdateHoverState()

	mx, my := ebiten.CursorPosition()

	// Calculate dropdown bounds
	absX, absY := sb.GetAbsolutePosition()
	dropdownBounds := Rect{
		X:      absX,
		Y:      absY + sb.bounds.Height,
		Width:  sb.bounds.Width,
		Height: sb.listBox.bounds.Height,
	}

	// If expanded, handle dropdown interactions
	if sb.expanded && sb.listBox != nil {
		// Update hover index for items in dropdown
		sb.listBox.HoveredIndex = -1
		if dropdownBounds.Contains(mx, my) {
			// Get content bounds for accurate click detection
			lbStyle := sb.listBox.GetComputedStyle()
			contentBounds := GetContentBounds(dropdownBounds, lbStyle)

			if contentBounds.Contains(mx, my) {
				relativeY := my - contentBounds.Y + sb.listBox.scrollOffset
				index := relativeY / sb.listBox.itemHeight
				if index >= 0 && index < len(sb.listBox.Items) {
					sb.listBox.HoveredIndex = index
				}
			}
		}

		// Handle click on dropdown item
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if sb.listBox.HoveredIndex != -1 {
				sb.setSelectedIndex(sb.listBox.HoveredIndex)
				sb.expanded = false
				return
			}
		}

		// Handle scrolling in dropdown
		_, dy := ebiten.Wheel()
		if dropdownBounds.Contains(mx, my) && dy != 0 {
			lbStyle := sb.listBox.GetComputedStyle()
			contentBounds := GetContentBounds(dropdownBounds, lbStyle)
			totalHeight := len(sb.listBox.Items) * sb.listBox.itemHeight

			if totalHeight > contentBounds.Height {
				sb.listBox.scrollOffset -= int(dy * 20)
				if sb.listBox.scrollOffset < 0 {
					sb.listBox.scrollOffset = 0
				}
				maxScroll := totalHeight - contentBounds.Height
				if sb.listBox.scrollOffset > maxScroll {
					sb.listBox.scrollOffset = maxScroll
				}
			}
		}
	}

	// Handle click on select to toggle dropdown
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if sb.hovered {
			sb.expanded = !sb.expanded
			return
		}

		// If click outside select and outside listbox, collapse
		if sb.expanded {
			if !(dropdownBounds.Contains(mx, my) || sb.IsWithin(mx, my)) {
				sb.expanded = false
			}
		}
	}
}

// Layout calculates dimensions
func (sb *SelectBox) Layout() {
	style := sb.GetComputedStyle()

	// Apply width/height from style if specified
	width := sb.bounds.Width
	height := sb.bounds.Height
	if style.Width != nil {
		width = *style.Width
	}
	if style.Height != nil {
		height = *style.Height
	}
	width, height = ApplySizeConstraints(width, height, style)
	sb.bounds.Width = width
	sb.bounds.Height = height

	// Update listbox size to match select width if listbox exists
	if sb.listBox != nil {
		// Height is either preconfigured or based on items (max 5 visible)
		maxVisible := 5
		maxHeight := sb.itemHeight * maxVisible
		totalHeight := len(sb.Items) * sb.itemHeight
		lbHeight := totalHeight
		if lbHeight > maxHeight {
			lbHeight = maxHeight
		}
		sb.listBox.SetSize(sb.bounds.Width, lbHeight)
		sb.listBox.SetItems(sb.Items)
		// Keep listbox selected index in sync
		if sb.SelectedIndex >= 0 && sb.SelectedIndex < len(sb.Items) {
			sb.listBox.SelectedIndex = sb.SelectedIndex
		} else {
			sb.listBox.SelectedIndex = -1
		}
	}
}

// Draw draws the select box and optionally the dropdown list
func (sb *SelectBox) Draw(screen *ebiten.Image) {
	if !sb.visible {
		return
	}

	style := sb.GetComputedStyle()
	theme := sb.GetTheme()
	absX, absY := sb.GetAbsolutePosition()
	absBounds := Rect{X: absX, Y: absY, Width: sb.bounds.Width, Height: sb.bounds.Height}

	// Draw background and border with theme support
	DrawBackgroundWithTheme(screen, absBounds, style, theme)
	DrawBorderWithTheme(screen, absBounds, style, theme)

	// Draw selected item text
	contentBounds := GetContentBounds(absBounds, style)
	// Get text color from style, then theme, then default
	textColor := color.RGBA{255, 255, 255, 255}
	if style.ForegroundColor != nil {
		textColor = colorToRGBA(*style.ForegroundColor)
	} else if theme != nil {
		textColor = colorToRGBA(theme.Colors.Text)
	}

	selectedText := ""
	if sb.SelectedIndex >= 0 && sb.SelectedIndex < len(sb.Items) {
		selectedText = sb.Items[sb.SelectedIndex]
	}
	fontSize := 14
	if style.FontSize != nil {
		fontSize = *style.FontSize
	}
	text.Draw(screen, selectedText, float64(fontSize), contentBounds.X+4, contentBounds.Y+3, textColor)

	// Draw arrow on right with theme color
	arrowColor := textColor
	arrowX := absBounds.X + absBounds.Width - 16
	arrowY := absBounds.Y + (absBounds.Height/2 - 6)
	text.Draw(screen, ">", 12.0, arrowX, arrowY, arrowColor)

	// If expanded, draw the list box (floating) at absolute coordinates
	if sb.expanded && sb.listBox != nil {
		// Calculate where listbox should appear (directly below select)
		absX, absY := sb.GetAbsolutePosition()
		dropdownBounds := Rect{
			X:      absX,
			Y:      absY + sb.bounds.Height,
			Width:  sb.bounds.Width,
			Height: sb.listBox.bounds.Height,
		}

		// Draw listbox background and border with theme support
		lbStyle := sb.listBox.GetComputedStyle()
		DrawBackgroundWithTheme(screen, dropdownBounds, lbStyle, theme)

		// Draw listbox items
		contentBounds := GetContentBounds(dropdownBounds, lbStyle)
		clipArea := CreateSubImage(screen, contentBounds)

		startIndex := sb.listBox.scrollOffset / sb.listBox.itemHeight
		visibleItems := contentBounds.Height / sb.listBox.itemHeight
		endIndex := startIndex + visibleItems + 1
		if endIndex > len(sb.listBox.Items) {
			endIndex = len(sb.listBox.Items)
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
			itemY := contentBounds.Y + (i * sb.listBox.itemHeight) - sb.listBox.scrollOffset
			itemBounds := Rect{
				X:      contentBounds.X,
				Y:      itemY,
				Width:  contentBounds.Width,
				Height: sb.listBox.itemHeight,
			}

			// Draw selection highlight
			if i == sb.listBox.SelectedIndex {
				DrawRect(clipArea, itemBounds, highlightColor)
			} else if i == sb.listBox.HoveredIndex {
				DrawRect(clipArea, itemBounds, hoverColor)
			}

			// Draw item text with theme colors
			itemTextColor := color.RGBA{255, 255, 255, 255}
			if theme != nil {
				itemTextColor = colorToRGBA(theme.Colors.Text)
			}
			if i == sb.listBox.SelectedIndex {
				itemTextColor = color.RGBA{255, 255, 255, 255}
			}

			text.Draw(clipArea, sb.listBox.Items[i], 14.0, itemBounds.X+4, itemBounds.Y+3, itemTextColor)
		}

		// Draw border with theme support
		DrawBorderWithTheme(screen, dropdownBounds, lbStyle, theme)

		// Draw scrollbar if needed
		totalHeight := len(sb.listBox.Items) * sb.listBox.itemHeight
		if totalHeight > contentBounds.Height {
			sb.drawScrollbar(screen, contentBounds, dropdownBounds, theme)
		}
	}
}

// drawScrollbar draws the scrollbar for the dropdown list
func (sb *SelectBox) drawScrollbar(screen *ebiten.Image, contentBounds Rect, dropdownBounds Rect, theme *Theme) {
	scrollbarWidth := 16
	scrollbarX := dropdownBounds.X + dropdownBounds.Width - scrollbarWidth
	scrollbarHeight := contentBounds.Height

	// Draw scrollbar track with theme color
	trackBounds := Rect{
		X:      scrollbarX,
		Y:      contentBounds.Y,
		Width:  scrollbarWidth,
		Height: scrollbarHeight,
	}
	trackColor := color.RGBA{40, 40, 50, 255}
	if theme != nil {
		trackColor = colorToRGBA(theme.Colors.Surface)
	}
	DrawRect(screen, trackBounds, trackColor)

	// Calculate thumb size and position
	totalHeight := len(sb.listBox.Items) * sb.listBox.itemHeight
	thumbHeight := (contentBounds.Height * scrollbarHeight) / totalHeight
	if thumbHeight < 20 {
		thumbHeight = 20
	}

	thumbY := contentBounds.Y + (sb.listBox.scrollOffset*(scrollbarHeight-thumbHeight))/(totalHeight-contentBounds.Height)

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

// SelectByIndex selects a specified index programmatically
func (sb *SelectBox) SelectByIndex(index int) {
	sb.setSelectedIndex(index)
}

// setSelectedIndex sets index and notifies listeners
func (sb *SelectBox) setSelectedIndex(index int) {
	if index < 0 || index >= len(sb.Items) {
		sb.SelectedIndex = -1
		return
	}
	sb.SelectedIndex = index
	if sb.listBox != nil {
		sb.listBox.SelectedIndex = index
	}
	if sb.OnSelect != nil {
		sb.OnSelect(index, sb.Items[index])
	}
	// Queue event
	event.GetQueuedInstance().QueueEvent(SelectBoxChangeEvent{
		SelectBoxID:   sb.GetID(),
		SelectBox:     sb,
		SelectedIndex: sb.SelectedIndex,
		SelectedItem:  sb.Items[sb.SelectedIndex],
	})
}

// SetItems sets available options
func (sb *SelectBox) SetItems(items []string) {
	sb.Items = items
	if sb.listBox != nil {
		sb.listBox.SetItems(items)
	}
}

// GetSelectedItem returns selected index and value
func (sb *SelectBox) GetSelectedItem() (int, string) {
	if sb.SelectedIndex >= 0 && sb.SelectedIndex < len(sb.Items) {
		return sb.SelectedIndex, sb.Items[sb.SelectedIndex]
	}
	return -1, ""
}

// IsExpanded returns whether the dropdown is open
func (sb *SelectBox) IsExpanded() bool {
	return sb.expanded
}

// Open opens the dropdown
func (sb *SelectBox) Open() {
	sb.expanded = true
}

// Close closes the dropdown
func (sb *SelectBox) Close() {
	sb.expanded = false
}

// Toggle toggles the dropdown
func (sb *SelectBox) Toggle() {
	sb.expanded = !sb.expanded
}

// IsMouseOverDropdown returns true if the mouse is over the dropdown area
func (sb *SelectBox) IsMouseOverDropdown(mx, my int) bool {
	if !sb.expanded {
		return false
	}
	absX, absY := sb.GetAbsolutePosition()
	dropdownBounds := Rect{
		X:      absX,
		Y:      absY + sb.bounds.Height,
		Width:  sb.bounds.Width,
		Height: sb.listBox.bounds.Height,
	}
	return dropdownBounds.Contains(mx, my)
}
