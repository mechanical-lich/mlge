package minui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mechanical-lich/mlge/text"
)

// MenuItem is a simple clickable text item for menus
// It has a clean appearance with just text and a subtle hover effect
type MenuItem struct {
	*ElementBase
	Text     string
	OnClick  func()
	hovered  bool
	selected bool
	clicked  bool // track click state to avoid repeat firing

	// Colors
	TextColor       color.Color
	HoverColor      color.Color
	SelectedColor   color.Color
	HoverBgColor    color.Color
	SelectedBgColor color.Color
}

// NewMenuItem creates a new menu item
func NewMenuItem(id string, text string) *MenuItem {
	mi := &MenuItem{
		ElementBase:     NewElementBase(id),
		Text:            text,
		TextColor:       color.RGBA{220, 220, 220, 255},
		HoverColor:      color.RGBA{255, 255, 255, 255},
		SelectedColor:   color.RGBA{255, 220, 100, 255},
		HoverBgColor:    color.RGBA{60, 60, 70, 255},
		SelectedBgColor: color.RGBA{50, 50, 60, 255},
	}

	mi.SetSize(120, 24)

	// Transparent background by default
	bgColor := color.Color(color.RGBA{0, 0, 0, 0})
	mi.style.BackgroundColor = &bgColor

	return mi
}

// SetSelected marks this item as selected
func (mi *MenuItem) SetSelected(selected bool) {
	mi.selected = selected
}

// IsSelected returns whether this item is selected
func (mi *MenuItem) IsSelected() bool {
	return mi.selected
}

// GetType returns the element type
func (mi *MenuItem) GetType() string {
	return "MenuItem"
}

// Update handles hover detection
func (mi *MenuItem) Update() {
	if !mi.IsVisible() {
		return
	}

	mx, my := ebiten.CursorPosition()
	absX, absY := mi.GetAbsolutePosition()
	bounds := mi.GetBounds()

	mi.hovered = mx >= absX && mx < absX+bounds.Width &&
		my >= absY && my < absY+bounds.Height

	// Only fire click on mouse button release (just pressed then released)
	if mi.hovered && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !mi.clicked {
			mi.clicked = true
		}
	} else if mi.clicked && mi.hovered {
		// Mouse was pressed and now released while still hovering
		mi.clicked = false
		if mi.OnClick != nil {
			mi.OnClick()
		}
	} else {
		mi.clicked = false
	}
}

// Draw renders the menu item
func (mi *MenuItem) Draw(screen *ebiten.Image) {
	if !mi.IsVisible() {
		return
	}

	absX, absY := mi.GetAbsolutePosition()
	bounds := mi.GetBounds()

	// Get colors - use theme if available, otherwise fall back to instance colors
	textColor := mi.TextColor
	hoverColor := mi.HoverColor
	selectedColor := mi.SelectedColor
	hoverBgColor := mi.HoverBgColor
	selectedBgColor := mi.SelectedBgColor

	theme := mi.GetTheme()
	if theme != nil {
		textColor = theme.Colors.Text
		hoverColor = theme.Colors.Primary
		selectedColor = theme.Colors.Focus
		// Derive background colors from theme surface
		if r, g, b, a := theme.Colors.Surface.RGBA(); a > 0 {
			// Slightly lighter for hover
			hoverBgColor = color.RGBA{
				R: uint8(min(int(r>>8)+20, 255)),
				G: uint8(min(int(g>>8)+20, 255)),
				B: uint8(min(int(b>>8)+20, 255)),
				A: uint8(a >> 8),
			}
			// Slightly different for selected
			selectedBgColor = color.RGBA{
				R: uint8(min(int(r>>8)+10, 255)),
				G: uint8(min(int(g>>8)+10, 255)),
				B: uint8(min(int(b>>8)+15, 255)),
				A: uint8(a >> 8),
			}
		}
	}

	// Draw background based on state (selected takes precedence, hover overlays)
	if mi.selected {
		vector.DrawFilledRect(
			screen,
			float32(absX),
			float32(absY),
			float32(bounds.Width),
			float32(bounds.Height),
			selectedBgColor,
			false,
		)
	}
	if mi.hovered {
		vector.DrawFilledRect(
			screen,
			float32(absX),
			float32(absY),
			float32(bounds.Width),
			float32(bounds.Height),
			hoverBgColor,
			false,
		)
	}

	// Draw left accent bar for selected items
	if mi.selected {
		vector.DrawFilledRect(
			screen,
			float32(absX),
			float32(absY+2),
			2,
			float32(bounds.Height-4),
			selectedColor,
			false,
		)
	}

	// Choose text color based on state
	finalTextColor := textColor
	if mi.selected {
		finalTextColor = selectedColor
	} else if mi.hovered {
		finalTextColor = hoverColor
	}

	// Draw text
	text.Draw(screen, mi.Text, 14.0, absX+4, absY+4, finalTextColor)
}

// Layout does nothing for menu items
func (mi *MenuItem) Layout() {}

// MenuHeader is a non-clickable header label for menu sections
type MenuHeader struct {
	*ElementBase
	Text      string
	TextColor color.Color
}

// NewMenuHeader creates a new menu section header
func NewMenuHeader(id string, txt string) *MenuHeader {
	mh := &MenuHeader{
		ElementBase: NewElementBase(id),
		Text:        txt,
		TextColor:   color.RGBA{180, 180, 100, 255},
	}

	mh.SetSize(120, 20)

	// Transparent background
	bgColor := color.Color(color.RGBA{0, 0, 0, 0})
	mh.style.BackgroundColor = &bgColor

	return mh
}

// Update does nothing for headers
func (mh *MenuHeader) Update() {}

// GetType returns the element type
func (mh *MenuHeader) GetType() string {
	return "MenuHeader"
}

// Draw renders the header
func (mh *MenuHeader) Draw(screen *ebiten.Image) {
	if !mh.IsVisible() {
		return
	}

	absX, absY := mh.GetAbsolutePosition()

	// Get text color - use theme if available
	textColor := mh.TextColor
	theme := mh.GetTheme()
	if theme != nil {
		textColor = theme.Colors.TextSecondary
	}

	text.Draw(screen, mh.Text, 14.0, absX+4, absY+2, textColor)
}

// Layout does nothing for headers
func (mh *MenuHeader) Layout() {}
