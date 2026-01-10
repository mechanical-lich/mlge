package minui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Element is the base interface for all UI elements
type Element interface {
	// Lifecycle
	Update()
	Draw(screen *ebiten.Image)
	Layout() // Calculate layout/dimensions

	// Hierarchy
	GetParent() Element
	SetParent(parent Element)
	GetChildren() []Element
	AddChild(child Element)
	RemoveChild(child Element)

	// Style
	GetStyle() *Style
	SetStyle(style *Style)
	GetComputedStyle() *Style // Inherited + merged style

	// Position & Size
	GetBounds() Rect
	SetBounds(bounds Rect)
	GetX() int
	GetY() int
	GetWidth() int
	GetHeight() int
	SetPosition(x, y int)
	SetSize(width, height int)
	GetAbsolutePosition() (int, int) // Get position in screen coordinates

	// Interaction
	IsWithin(x, y int) bool
	IsHovered() bool
	IsFocused() bool
	IsEnabled() bool
	SetEnabled(enabled bool)

	// Identity
	GetID() string
	SetID(id string)
	GetType() string

	// Visibility
	IsVisible() bool
	SetVisible(visible bool)
}

// Rect represents a rectangle with position and size
type Rect struct {
	X, Y, Width, Height int
}

// Contains checks if a point is within the rectangle
func (r Rect) Contains(x, y int) bool {
	return x >= r.X && x < r.X+r.Width && y >= r.Y && y < r.Y+r.Height
}

// Container is an element that can contain other elements
type Container interface {
	Element

	// Layout direction
	GetLayoutDirection() LayoutDirection
	SetLayoutDirection(dir LayoutDirection)
}

// LayoutDirection defines how children are arranged
type LayoutDirection int

const (
	LayoutVertical LayoutDirection = iota
	LayoutHorizontal
	LayoutNone // Manual/absolute positioning
)

// ElementBase provides common functionality for all elements
type ElementBase struct {
	id       string
	parent   Element
	children []Element
	style    *Style
	bounds   Rect
	enabled  bool
	visible  bool
	hovered  bool
	focused  bool

	// Theme for sprite-based rendering (set by GUI when element is added)
	theme *Theme

	// Cached computed style
	computedStyle *Style
	styleDirty    bool
}

// NewElementBase creates a new element base
func NewElementBase(id string) *ElementBase {
	return &ElementBase{
		id:         id,
		children:   make([]Element, 0),
		style:      &Style{},
		enabled:    true,
		visible:    true,
		styleDirty: true,
		theme:      nil,
	}
}

// GetID returns the element's ID
func (e *ElementBase) GetID() string {
	return e.id
}

// SetID sets the element's ID
func (e *ElementBase) SetID(id string) {
	e.id = id
}

// GetParent returns the parent element
func (e *ElementBase) GetParent() Element {
	return e.parent
}

// SetParent sets the parent element
func (e *ElementBase) SetParent(parent Element) {
	e.parent = parent
	e.styleDirty = true
}

// GetChildren returns child elements
func (e *ElementBase) GetChildren() []Element {
	return e.children
}

func (e *ElementBase) FindChildByID(id string) Element {
	for _, child := range e.children {
		if child.GetID() == id {
			return child
		}
	}
	return nil
}

// AddChild adds a child element
func (e *ElementBase) AddChild(child Element) {
	e.children = append(e.children, child)
	// Parent should be set by the concrete container implementation that knows the
	// correct Element type (see Panel.AddChild, Modal.AddChild, etc.).
}

// RemoveChild removes a child element
func (e *ElementBase) RemoveChild(child Element) {
	for i, c := range e.children {
		if c == child {
			e.children = append(e.children[:i], e.children[i+1:]...)
			child.SetParent(nil)
			break
		}
	}
}

// GetStyle returns the element's style
func (e *ElementBase) GetStyle() *Style {
	return e.style
}

// SetStyle sets the element's style
func (e *ElementBase) SetStyle(style *Style) {
	e.style = style
	e.styleDirty = true
}

// GetComputedStyle returns the computed style (with inheritance)
func (e *ElementBase) GetComputedStyle() *Style {
	if e.styleDirty {
		e.computedStyle = e.computeStyle()
		e.styleDirty = false
	}
	return e.computedStyle
}

// computeStyle computes the final style with inheritance
func (e *ElementBase) computeStyle() *Style {
	var parentStyle *Style
	if e.parent != nil {
		parentStyle = e.parent.GetComputedStyle()
	} else {
		parentStyle = DefaultStyle()
	}

	// Merge with parent
	merged := e.style.Merge(parentStyle)

	// Apply state-based styles
	return merged.GetComputedStyle(e.hovered, false, !e.enabled, e.focused)
}

// GetBounds returns the element's bounds
func (e *ElementBase) GetBounds() Rect {
	return e.bounds
}

// SetBounds sets the element's bounds
func (e *ElementBase) SetBounds(bounds Rect) {
	e.bounds = bounds
}

// GetX returns the X position
func (e *ElementBase) GetX() int {
	return e.bounds.X
}

// GetY returns the Y position
func (e *ElementBase) GetY() int {
	return e.bounds.Y
}

// GetWidth returns the width
func (e *ElementBase) GetWidth() int {
	return e.bounds.Width
}

// GetHeight returns the height
func (e *ElementBase) GetHeight() int {
	return e.bounds.Height
}

// SetPosition sets the position
func (e *ElementBase) SetPosition(x, y int) {
	e.bounds.X = x
	e.bounds.Y = y
}

// GetPosition returns the position
func (e *ElementBase) GetPosition() (int, int) {
	return e.bounds.X, e.bounds.Y
}

// SetSize sets the size
func (e *ElementBase) SetSize(width, height int) {
	e.bounds.Width = width
	e.bounds.Height = height
}

// GetSize returns the size
func (e *ElementBase) GetSize() (int, int) {
	return e.bounds.Width, e.bounds.Height
}

// GetAbsolutePosition returns the absolute screen position
func (e *ElementBase) GetAbsolutePosition() (int, int) {
	if e.parent == nil {
		return e.bounds.X, e.bounds.Y
	}

	// Get parent's absolute position
	px, py := e.parent.GetAbsolutePosition()

	// Get parent's style to calculate content bounds
	parentStyle := e.parent.GetComputedStyle()
	parentBounds := e.parent.GetBounds()
	contentOffset := GetContentBounds(parentBounds, parentStyle)

	// The content offset gives us the padding/border offset
	offsetX := contentOffset.X - parentBounds.X
	offsetY := contentOffset.Y - parentBounds.Y

	// Check if parent is a modal - add title bar height
	if e.parent.GetType() == "Modal" {
		offsetY += 30 // Title bar height
	}

	return px + offsetX + e.bounds.X, py + offsetY + e.bounds.Y
}

// IsWithin checks if a point is within the element (in screen coordinates)
func (e *ElementBase) IsWithin(x, y int) bool {
	absX, absY := e.GetAbsolutePosition()
	return x >= absX && x < absX+e.bounds.Width && y >= absY && y < absY+e.bounds.Height
}

// IsHovered returns if the element is hovered
func (e *ElementBase) IsHovered() bool {
	return e.hovered
}

// SetHovered sets the hovered state
func (e *ElementBase) SetHovered(hovered bool) {
	if e.hovered != hovered {
		e.hovered = hovered
		e.styleDirty = true
	}
}

// IsFocused returns if the element is focused
func (e *ElementBase) IsFocused() bool {
	return e.focused
}

// SetFocused sets the focused state
func (e *ElementBase) SetFocused(focused bool) {
	if e.focused != focused {
		e.focused = focused
		e.styleDirty = true
	}
}

// IsEnabled returns if the element is enabled
func (e *ElementBase) IsEnabled() bool {
	return e.enabled
}

// SetEnabled sets the enabled state
func (e *ElementBase) SetEnabled(enabled bool) {
	if e.enabled != enabled {
		e.enabled = enabled
		e.styleDirty = true
	}
}

// IsVisible returns if the element is visible
func (e *ElementBase) IsVisible() bool {
	return e.visible
}

// SetVisible sets the visibility
func (e *ElementBase) SetVisible(visible bool) {
	e.visible = visible
}

// UpdateHoverState updates hover state based on mouse position
func (e *ElementBase) UpdateHoverState() {
	mx, my := ebiten.CursorPosition()
	e.SetHovered(e.IsWithin(mx, my))
}

// MarkStyleDirty marks the style as needing recomputation
func (e *ElementBase) MarkStyleDirty() {
	e.styleDirty = true
	for _, child := range e.children {
		// Recursively mark children dirty
		if childBase := getElementBase(child); childBase != nil {
			childBase.MarkStyleDirty()
		}
	}
}

// GetTheme returns the theme for this element
func (e *ElementBase) GetTheme() *Theme {
	return e.theme
}

// SetTheme sets the theme for this element and all children
func (e *ElementBase) SetTheme(theme *Theme) {
	e.theme = theme
	for _, child := range e.children {
		setThemeRecursive(child, theme)
	}
}

// setThemeRecursive sets the theme on an element and its children
func setThemeRecursive(elem Element, theme *Theme) {
	// Try to get ElementBase methods
	if setter, ok := elem.(interface{ SetTheme(*Theme) }); ok {
		setter.SetTheme(theme)
	}
}

// getElementBase is a helper to access ElementBase from an Element
// Since all elements embed ElementBase, we can use type switches
func getElementBase(elem Element) *ElementBase {
	// This is a workaround since we can't directly access embedded struct
	// Each concrete type will need to provide access if needed
	// For now, we'll skip recursive marking
	return nil
}
