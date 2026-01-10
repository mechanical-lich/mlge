package minui

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/text/v2"
)

// ResourceItem represents a single resource with icon and value
type ResourceItem struct {
	ID    string
	Icon  *Icon
	Value int
	Color color.Color // Optional color for the value text
}

// ResourceBar is a horizontal bar displaying resource icons and values
type ResourceBar struct {
	*ElementBase
	Items       []*ResourceItem
	ItemSpacing int // Space between items
	IconSpacing int // Space between icon and value
}

// NewResourceBar creates a new resource bar
func NewResourceBar(id string) *ResourceBar {
	rb := &ResourceBar{
		ElementBase: NewElementBase(id),
		Items:       make([]*ResourceItem, 0),
		ItemSpacing: 24,
		IconSpacing: 4,
	}

	rb.SetSize(400, 32)

	// Set default style - only structural properties, colors come from theme

	return rb
}

// AddResource adds a resource to display
func (rb *ResourceBar) AddResource(id string, icon *Icon, value int) *ResourceItem {
	item := &ResourceItem{
		ID:    id,
		Icon:  icon,
		Value: value,
	}
	rb.Items = append(rb.Items, item)
	return item
}

// SetResourceValue updates a resource's value by ID
func (rb *ResourceBar) SetResourceValue(id string, value int) {
	for _, item := range rb.Items {
		if item.ID == id {
			item.Value = value
			return
		}
	}
}

// GetResourceValue gets a resource's value by ID
func (rb *ResourceBar) GetResourceValue(id string) int {
	for _, item := range rb.Items {
		if item.ID == id {
			return item.Value
		}
	}
	return 0
}

// SetResourceColor sets a specific color for a resource's value
func (rb *ResourceBar) SetResourceColor(id string, clr color.Color) {
	for _, item := range rb.Items {
		if item.ID == id {
			item.Color = clr
			return
		}
	}
}

// Clear removes all resources
func (rb *ResourceBar) Clear() {
	rb.Items = make([]*ResourceItem, 0)
}

// GetType returns the element type
func (rb *ResourceBar) GetType() string {
	return "ResourceBar"
}

// Update updates the resource bar
func (rb *ResourceBar) Update() {
	if !rb.visible {
		return
	}
	rb.UpdateHoverState()
}

// Layout calculates dimensions
func (rb *ResourceBar) Layout() {
	style := rb.GetComputedStyle()

	// Calculate required width based on items
	width := 0
	height := 32

	for i, item := range rb.Items {
		if i > 0 {
			width += rb.ItemSpacing
		}
		if item.Icon != nil {
			width += item.Icon.ScaledWidth() + rb.IconSpacing
		}
		// Estimate value text width
		valueStr := formatNumber(item.Value)
		width += len(valueStr) * 8
	}

	// Add padding
	if style.Padding != nil {
		width += style.Padding.Left + style.Padding.Right
		height += style.Padding.Top + style.Padding.Bottom
	}

	// Apply explicit dimensions from style
	if style.Width != nil {
		width = *style.Width
	}
	if style.Height != nil {
		height = *style.Height
	}

	width, height = ApplySizeConstraints(width, height, style)

	rb.bounds.Width = width
	rb.bounds.Height = height
}

// Draw draws the resource bar
func (rb *ResourceBar) Draw(screen *ebiten.Image) {
	if !rb.visible {
		return
	}

	style := rb.GetComputedStyle()
	theme := rb.GetTheme()
	absX, absY := rb.GetAbsolutePosition()
	absBounds := Rect{
		X:      absX,
		Y:      absY,
		Width:  rb.bounds.Width,
		Height: rb.bounds.Height,
	}

	// Draw background with theme support
	DrawBackgroundWithTheme(screen, absBounds, style, theme)

	contentBounds := GetContentBounds(absBounds, style)

	fontSize := 14
	if style.FontSize != nil {
		fontSize = *style.FontSize
	}

	// Get default text color from style, then theme, then default
	defaultTextColor := color.RGBA{255, 255, 255, 255}
	if style.ForegroundColor != nil {
		defaultTextColor = colorToRGBA(*style.ForegroundColor)
	} else if theme != nil {
		defaultTextColor = colorToRGBA(theme.Colors.Text)
	}

	x := contentBounds.X

	for i, item := range rb.Items {
		if i > 0 {
			x += rb.ItemSpacing
		}

		// Draw icon
		if item.Icon != nil {
			iconY := contentBounds.Y + (contentBounds.Height-item.Icon.ScaledHeight())/2
			item.Icon.Draw(screen, x, iconY)
			x += item.Icon.ScaledWidth() + rb.IconSpacing
		}

		// Draw value
		valueStr := formatNumber(item.Value)
		textColor := defaultTextColor
		if item.Color != nil {
			r, g, b, a := item.Color.RGBA()
			textColor = color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
		}

		textY := contentBounds.Y + (contentBounds.Height-fontSize)/2
		text.Draw(screen, valueStr, float64(fontSize), x, textY, textColor)
		x += len(valueStr) * fontSize * 6 / 10
	}

	// Draw border
	DrawBorder(screen, absBounds, style)
}

// formatNumber formats a number with thousands separator
func formatNumber(n int) string {
	if n < 1000 {
		return strconv.Itoa(n)
	}

	// Add commas for readability
	s := strconv.Itoa(n)
	result := ""
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result += ","
		}
		result += string(c)
	}
	return result
}
