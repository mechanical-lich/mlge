package minui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Panel is a basic container that can hold other elements
type Panel struct {
	*ElementBase
	layoutDirection LayoutDirection
}

// NewPanel creates a new panel
func NewPanel(id string) *Panel {
	return &Panel{
		ElementBase:     NewElementBase(id),
		layoutDirection: LayoutNone,
	}
}

// GetType returns the element type
func (p *Panel) GetType() string {
	return "Panel"
}

// GetLayoutDirection returns the layout direction
func (p *Panel) GetLayoutDirection() LayoutDirection {
	return p.layoutDirection
}

// SetLayoutDirection sets the layout direction
func (p *Panel) SetLayoutDirection(dir LayoutDirection) {
	p.layoutDirection = dir
}

// AddChild adds a child element to the panel
func (p *Panel) AddChild(child Element) {
	p.children = append(p.children, child)
	child.SetParent(p) // p is the Panel, which implements Element
}

// Update updates the panel and its children
func (p *Panel) Update() {
	if !p.visible {
		return
	}

	p.UpdateHoverState()

	for _, child := range p.children {
		child.Update()
	}
}

// Layout calculates the layout for the panel and its children
func (p *Panel) Layout() {
	if !p.visible {
		return
	}

	style := p.GetComputedStyle()

	// Apply width/height from style if specified
	width := p.bounds.Width
	height := p.bounds.Height

	if style.Width != nil {
		width = *style.Width
	}
	if style.Height != nil {
		height = *style.Height
	}

	// Apply min/max size constraints
	width, height = ApplySizeConstraints(width, height, style)

	p.bounds.Width = width
	p.bounds.Height = height

	contentBounds := GetContentBounds(p.bounds, style)

	switch p.layoutDirection {
	case LayoutVertical:
		p.layoutVertical(contentBounds)
	case LayoutHorizontal:
		p.layoutHorizontal(contentBounds)
	case LayoutNone:
		// Children use their own positioning
	}

	// Layout children
	for _, child := range p.children {
		child.Layout()
	}
}

// layoutVertical arranges children vertically
func (p *Panel) layoutVertical(contentBounds Rect) {
	currentY := contentBounds.Y

	for _, child := range p.children {
		if !child.IsVisible() {
			continue
		}

		childStyle := child.GetComputedStyle()

		// Apply margin
		marginedBounds := Rect{
			X:      contentBounds.X,
			Y:      currentY,
			Width:  contentBounds.Width,
			Height: child.GetHeight(),
		}
		marginedBounds = ApplyMargin(marginedBounds, childStyle)

		// Set child bounds
		width := contentBounds.Width
		if childStyle.Width != nil {
			width = *childStyle.Width
		}

		height := child.GetHeight()
		if childStyle.Height != nil {
			height = *childStyle.Height
		}

		child.SetBounds(Rect{
			X:      marginedBounds.X,
			Y:      marginedBounds.Y,
			Width:  width,
			Height: height,
		})

		currentY = marginedBounds.Y + height
		if childStyle.Margin != nil {
			currentY += childStyle.Margin.Bottom
		}
	}
}

// layoutHorizontal arranges children horizontally
func (p *Panel) layoutHorizontal(contentBounds Rect) {
	currentX := contentBounds.X

	for _, child := range p.children {
		if !child.IsVisible() {
			continue
		}

		childStyle := child.GetComputedStyle()

		// Apply margin
		marginedBounds := Rect{
			X:      currentX,
			Y:      contentBounds.Y,
			Width:  child.GetWidth(),
			Height: contentBounds.Height,
		}
		marginedBounds = ApplyMargin(marginedBounds, childStyle)

		// Set child bounds
		width := child.GetWidth()
		if childStyle.Width != nil {
			width = *childStyle.Width
		}

		height := contentBounds.Height
		if childStyle.Height != nil {
			height = *childStyle.Height
		}

		child.SetBounds(Rect{
			X:      marginedBounds.X,
			Y:      marginedBounds.Y,
			Width:  width,
			Height: height,
		})

		currentX = marginedBounds.X + width
		if childStyle.Margin != nil {
			currentX += childStyle.Margin.Right
		}
	}
}

// Draw draws the panel and its children
func (p *Panel) Draw(screen *ebiten.Image) {
	if !p.visible {
		return
	}

	style := p.GetComputedStyle()

	// Get absolute position for drawing
	absX, absY := p.GetAbsolutePosition()
	absBounds := Rect{
		X:      absX,
		Y:      absY,
		Width:  p.bounds.Width,
		Height: p.bounds.Height,
	}

	// Draw background
	DrawBackground(screen, absBounds, style)

	// Draw border
	DrawBorder(screen, absBounds, style)

	// Draw children
	for _, child := range p.children {
		child.Draw(screen)
	}
}

// HBox is a container that automatically arranges children horizontally
type HBox struct {
	*ElementBase
	Spacing int // Space between children in pixels
}

// NewHBox creates a new horizontal box container
func NewHBox(id string) *HBox {
	return &HBox{
		ElementBase: NewElementBase(id),
		Spacing:     5, // Default spacing
	}
}

// GetType returns the element type
func (h *HBox) GetType() string {
	return "HBox"
}

// AddChild adds a child and sets its parent
func (h *HBox) AddChild(child Element) {
	h.children = append(h.children, child)
	child.SetParent(h)
}

// Update updates all children
func (h *HBox) Update() {
	if !h.visible {
		return
	}
	h.UpdateHoverState()
	for _, child := range h.GetChildren() {
		child.Update()
	}
}

// Layout arranges children horizontally with spacing
func (h *HBox) Layout() {
	if !h.visible {
		return
	}

	style := h.GetComputedStyle()

	// First, layout all children to get their sizes
	for _, child := range h.GetChildren() {
		child.Layout()
	}

	// Calculate total width needed and max height
	totalWidth := 0
	maxHeight := 0
	visibleCount := 0

	for _, child := range h.GetChildren() {
		if !child.IsVisible() {
			continue
		}
		visibleCount++

		childStyle := child.GetComputedStyle()
		childWidth := child.GetWidth()
		childHeight := child.GetHeight()

		// Add margin to width calculation
		if childStyle != nil && childStyle.Margin != nil {
			childWidth += childStyle.Margin.Left + childStyle.Margin.Right
			childHeight += childStyle.Margin.Top + childStyle.Margin.Bottom
		}

		totalWidth += childWidth
		if childHeight > maxHeight {
			maxHeight = childHeight
		}
	}

	// Add spacing between visible children (n-1 gaps)
	if visibleCount > 1 {
		totalWidth += h.Spacing * (visibleCount - 1)
	}

	// Add padding to dimensions
	if style != nil && style.Padding != nil {
		totalWidth += style.Padding.Left + style.Padding.Right
		maxHeight += style.Padding.Top + style.Padding.Bottom
	}

	// Set HBox size (or use explicit size from style)
	width := totalWidth
	height := maxHeight

	if style != nil {
		if style.Width != nil {
			width = *style.Width
		}
		if style.Height != nil {
			height = *style.Height
		}
	}

	// Apply min/max constraints
	width, height = ApplySizeConstraints(width, height, style)

	h.bounds.Width = width
	h.bounds.Height = height

	// Get content bounds - this includes the container's position
	contentBounds := GetContentBounds(h.bounds, style)

	// Calculate padding/border offset (relative to container origin)
	offsetX := contentBounds.X - h.bounds.X // Just the padding/border
	offsetY := contentBounds.Y - h.bounds.Y

	// Position children horizontally, relative to container's origin
	currentX := offsetX

	for _, child := range h.GetChildren() {
		if !child.IsVisible() {
			continue
		}

		childStyle := child.GetComputedStyle()
		childWidth := child.GetWidth()
		childHeight := child.GetHeight()

		// Apply margin
		marginedBounds := Rect{
			X:      currentX,
			Y:      offsetY,
			Width:  childWidth,
			Height: childHeight,
		}
		marginedBounds = ApplyMargin(marginedBounds, childStyle)

		// Set child bounds (relative to container's origin)
		child.SetBounds(Rect{
			X:      marginedBounds.X,
			Y:      marginedBounds.Y,
			Width:  childWidth,
			Height: childHeight,
		})

		// Move to next position
		currentX = marginedBounds.X + childWidth
		if childStyle != nil && childStyle.Margin != nil {
			currentX += childStyle.Margin.Right
		}
		currentX += h.Spacing
	}

	// Layout children after positioning (like Panel does)
	for _, child := range h.GetChildren() {
		child.Layout()
	}
}

// Draw draws the HBox background and children
func (h *HBox) Draw(screen *ebiten.Image) {
	if !h.visible {
		return
	}

	style := h.GetComputedStyle()
	absX, absY := h.GetAbsolutePosition()
	absBounds := Rect{X: absX, Y: absY, Width: h.bounds.Width, Height: h.bounds.Height}

	// Draw background if specified
	DrawBackground(screen, absBounds, style)
	DrawBorder(screen, absBounds, style)

	// Draw children
	for _, child := range h.GetChildren() {
		child.Draw(screen)
	}
}

// VBox is a container that automatically arranges children vertically
type VBox struct {
	*ElementBase
	Spacing int // Space between children in pixels
}

// NewVBox creates a new vertical box container
func NewVBox(id string) *VBox {
	return &VBox{
		ElementBase: NewElementBase(id),
		Spacing:     5, // Default spacing
	}
}

// GetType returns the element type
func (v *VBox) GetType() string {
	return "VBox"
}

// AddChild adds a child and sets its parent
func (v *VBox) AddChild(child Element) {
	v.children = append(v.children, child)
	child.SetParent(v)
}

// Update updates all children
func (v *VBox) Update() {
	if !v.visible {
		return
	}
	v.UpdateHoverState()
	for _, child := range v.GetChildren() {
		child.Update()
	}
}

// Layout arranges children vertically with spacing
func (v *VBox) Layout() {
	if !v.visible {
		return
	}

	style := v.GetComputedStyle()

	// First, layout all children to get their sizes
	for _, child := range v.GetChildren() {
		child.Layout()
	}

	// Calculate max width needed and total height
	maxWidth := 0
	totalHeight := 0
	visibleCount := 0

	for _, child := range v.GetChildren() {
		if !child.IsVisible() {
			continue
		}
		visibleCount++

		childStyle := child.GetComputedStyle()
		childWidth := child.GetWidth()
		childHeight := child.GetHeight()

		// Add margin to dimensions
		if childStyle != nil && childStyle.Margin != nil {
			childWidth += childStyle.Margin.Left + childStyle.Margin.Right
			childHeight += childStyle.Margin.Top + childStyle.Margin.Bottom
		}

		totalHeight += childHeight
		if childWidth > maxWidth {
			maxWidth = childWidth
		}
	}

	// Add spacing between visible children (n-1 gaps)
	if visibleCount > 1 {
		totalHeight += v.Spacing * (visibleCount - 1)
	}

	// Add padding to dimensions
	if style != nil && style.Padding != nil {
		maxWidth += style.Padding.Left + style.Padding.Right
		totalHeight += style.Padding.Top + style.Padding.Bottom
	}

	// Set VBox size (or use explicit size from style)
	width := maxWidth
	height := totalHeight

	if style != nil {
		if style.Width != nil {
			width = *style.Width
		}
		if style.Height != nil {
			height = *style.Height
		}
	}

	// Apply min/max constraints
	width, height = ApplySizeConstraints(width, height, style)

	v.bounds.Width = width
	v.bounds.Height = height

	// Get content bounds - this includes the container's position
	contentBounds := GetContentBounds(v.bounds, style)

	// Calculate padding/border offset (relative to container origin)
	offsetX := contentBounds.X - v.bounds.X // Just the padding/border
	offsetY := contentBounds.Y - v.bounds.Y

	// Position children vertically, relative to container's origin
	currentY := offsetY

	for _, child := range v.GetChildren() {
		if !child.IsVisible() {
			continue
		}

		childStyle := child.GetComputedStyle()
		childWidth := child.GetWidth()
		childHeight := child.GetHeight()

		// Apply margin
		marginedBounds := Rect{
			X:      offsetX,
			Y:      currentY,
			Width:  childWidth,
			Height: childHeight,
		}
		marginedBounds = ApplyMargin(marginedBounds, childStyle)

		// Set child bounds (relative to container's origin)
		child.SetBounds(Rect{
			X:      marginedBounds.X,
			Y:      marginedBounds.Y,
			Width:  childWidth,
			Height: childHeight,
		}) // Move to next position
		currentY = marginedBounds.Y + childHeight
		if childStyle != nil && childStyle.Margin != nil {
			currentY += childStyle.Margin.Bottom
		}
		currentY += v.Spacing
	}

	// Layout children after positioning (like Panel does)
	for _, child := range v.GetChildren() {
		child.Layout()
	}
}

// Draw draws the VBox background and children
func (v *VBox) Draw(screen *ebiten.Image) {
	if !v.visible {
		return
	}

	style := v.GetComputedStyle()
	absX, absY := v.GetAbsolutePosition()
	absBounds := Rect{X: absX, Y: absY, Width: v.bounds.Width, Height: v.bounds.Height}

	// Draw background if specified
	DrawBackground(screen, absBounds, style)
	DrawBorder(screen, absBounds, style)

	// Draw children
	for _, child := range v.GetChildren() {
		child.Draw(screen)
	}
}
