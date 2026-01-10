package minui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// DrawerPosition defines which edge the drawer attaches to
type DrawerPosition int

const (
	DrawerLeft DrawerPosition = iota
	DrawerRight
	DrawerTop
	DrawerBottom
)

// Drawer is a sliding panel that can be shown/hidden from screen edges
type Drawer struct {
	*ElementBase
	Position       DrawerPosition
	Width          int // For left/right drawers
	Height         int // For top/bottom drawers
	Collapsed      bool
	AnimationSpeed int // Pixels per frame when animating

	currentOffset  int // Current animation offset
	targetOffset   int // Target animation offset
	animating      bool
	overlayEnabled bool // Draw overlay behind drawer when open
	closeOnOverlay bool // Close when clicking overlay
}

// NewDrawer creates a new drawer panel
func NewDrawer(id string, position DrawerPosition) *Drawer {
	d := &Drawer{
		ElementBase:    NewElementBase(id),
		Position:       position,
		Width:          280,
		Height:         200,
		Collapsed:      true,
		AnimationSpeed: 20,
		overlayEnabled: true,
		closeOnOverlay: true,
	}

	// Set default style - only structural properties, colors come from theme
	borderWidth := 1
	d.style.BorderWidth = &borderWidth

	return d
}

// AddChild adds a child element to the drawer
func (d *Drawer) AddChild(child Element) {
	d.children = append(d.children, child)
	child.SetParent(d)
}

// SetSize sets the drawer size (width for left/right, height for top/bottom)
func (d *Drawer) SetDrawerSize(size int) {
	switch d.Position {
	case DrawerLeft, DrawerRight:
		d.Width = size
	case DrawerTop, DrawerBottom:
		d.Height = size
	}
}

// Open opens the drawer with animation
func (d *Drawer) Open() {
	if !d.Collapsed {
		return
	}
	d.Collapsed = false
	d.targetOffset = 0
	d.animating = true
	d.visible = true
}

// Close closes the drawer with animation
func (d *Drawer) Close() {
	if d.Collapsed {
		return
	}
	d.Collapsed = true
	switch d.Position {
	case DrawerLeft:
		d.targetOffset = -d.Width
	case DrawerRight:
		d.targetOffset = d.Width
	case DrawerTop:
		d.targetOffset = -d.Height
	case DrawerBottom:
		d.targetOffset = d.Height
	}
	d.animating = true
}

// Toggle toggles the drawer open/closed state
func (d *Drawer) Toggle() {
	if d.Collapsed {
		d.Open()
	} else {
		d.Close()
	}
}

// IsOpen returns whether the drawer is open
func (d *Drawer) IsOpen() bool {
	return !d.Collapsed
}

// SetOverlay enables/disables the darkened overlay behind the drawer
func (d *Drawer) SetOverlay(enabled bool) {
	d.overlayEnabled = enabled
}

// SetCloseOnOverlay sets whether clicking the overlay closes the drawer
func (d *Drawer) SetCloseOnOverlay(close bool) {
	d.closeOnOverlay = close
}

// GetType returns the element type
func (d *Drawer) GetType() string {
	return "Drawer"
}

// Update updates the drawer
func (d *Drawer) Update() {
	if !d.visible && d.Collapsed {
		return
	}

	// Animate
	if d.animating {
		if d.currentOffset < d.targetOffset {
			d.currentOffset += d.AnimationSpeed
			if d.currentOffset > d.targetOffset {
				d.currentOffset = d.targetOffset
				d.animating = false
				if d.Collapsed {
					d.visible = false
				}
			}
		} else if d.currentOffset > d.targetOffset {
			d.currentOffset -= d.AnimationSpeed
			if d.currentOffset < d.targetOffset {
				d.currentOffset = d.targetOffset
				d.animating = false
				if d.Collapsed {
					d.visible = false
				}
			}
		}
	}

	// Handle overlay click
	if !d.Collapsed && d.closeOnOverlay && d.overlayEnabled {
		mx, my := ebiten.CursorPosition()
		drawerBounds := d.getDrawerBounds()

		inDrawer := mx >= drawerBounds.X && mx < drawerBounds.X+drawerBounds.Width &&
			my >= drawerBounds.Y && my < drawerBounds.Y+drawerBounds.Height

		if !inDrawer && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			d.Close()
		}
	}

	// Update children
	for _, child := range d.children {
		child.Update()
	}
}

// getDrawerBounds calculates the current drawer bounds with animation offset
func (d *Drawer) getDrawerBounds() Rect {
	screenW, screenH := ebiten.WindowSize()

	var bounds Rect
	switch d.Position {
	case DrawerLeft:
		bounds = Rect{
			X:      d.currentOffset,
			Y:      0,
			Width:  d.Width,
			Height: screenH,
		}
	case DrawerRight:
		bounds = Rect{
			X:      screenW - d.Width + d.currentOffset,
			Y:      0,
			Width:  d.Width,
			Height: screenH,
		}
	case DrawerTop:
		bounds = Rect{
			X:      0,
			Y:      d.currentOffset,
			Width:  screenW,
			Height: d.Height,
		}
	case DrawerBottom:
		bounds = Rect{
			X:      0,
			Y:      screenH - d.Height + d.currentOffset,
			Width:  screenW,
			Height: d.Height,
		}
	}

	return bounds
}

// Layout calculates dimensions
func (d *Drawer) Layout() {
	bounds := d.getDrawerBounds()
	d.bounds = bounds

	// Layout children
	for _, child := range d.children {
		child.Layout()
	}
}

// GetAbsolutePosition returns absolute screen position for children
func (d *Drawer) GetAbsolutePosition() (int, int) {
	bounds := d.getDrawerBounds()
	return bounds.X, bounds.Y
}

// Draw draws the drawer
func (d *Drawer) Draw(screen *ebiten.Image) {
	if !d.visible {
		return
	}

	theme := d.GetTheme()
	screenW, screenH := ebiten.WindowSize()

	// Draw overlay
	if d.overlayEnabled && !d.Collapsed {
		// Calculate opacity based on animation progress
		var maxOffset int
		switch d.Position {
		case DrawerLeft:
			maxOffset = d.Width
		case DrawerRight:
			maxOffset = d.Width
		case DrawerTop:
			maxOffset = d.Height
		case DrawerBottom:
			maxOffset = d.Height
		}

		var progress float64
		if maxOffset > 0 {
			progress = 1.0 - (float64(abs(d.currentOffset)) / float64(maxOffset))
		}
		if progress < 0 {
			progress = 0
		}
		if progress > 1 {
			progress = 1
		}

		overlayAlpha := uint8(float64(120) * progress)
		// Get overlay color from theme or default to black
		overlayColor := color.RGBA{0, 0, 0, overlayAlpha}
		if theme != nil {
			bg := colorToRGBA(theme.Colors.Background)
			overlayColor = color.RGBA{bg.R / 2, bg.G / 2, bg.B / 2, overlayAlpha}
		}
		DrawRect(screen, Rect{X: 0, Y: 0, Width: screenW, Height: screenH}, overlayColor)
	}

	// Draw drawer background with theme support
	bounds := d.getDrawerBounds()
	style := d.GetComputedStyle()

	DrawBackgroundWithTheme(screen, bounds, style, theme)

	// Draw children
	for _, child := range d.children {
		child.Draw(screen)
	}

	// Draw border (on the visible edge) with theme support
	DrawBorderWithTheme(screen, bounds, style, theme)
}

// abs returns absolute value of int
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Initialize sets up the drawer with initial closed state
func (d *Drawer) Initialize() {
	switch d.Position {
	case DrawerLeft:
		d.currentOffset = -d.Width
		d.targetOffset = -d.Width
	case DrawerRight:
		d.currentOffset = d.Width
		d.targetOffset = d.Width
	case DrawerTop:
		d.currentOffset = -d.Height
		d.targetOffset = -d.Height
	case DrawerBottom:
		d.currentOffset = d.Height
		d.targetOffset = d.Height
	}
	d.visible = false
}
