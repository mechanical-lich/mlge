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

// VBox is a vertical box container
type VBox struct {
	*Panel
}

// NewVBox creates a new vertical box
func NewVBox(id string) *VBox {
	panel := NewPanel(id)
	panel.layoutDirection = LayoutVertical
	return &VBox{Panel: panel}
}

// GetType returns the element type
func (v *VBox) GetType() string {
	return "VBox"
}

// HBox is a horizontal box container
type HBox struct {
	*Panel
}

// NewHBox creates a new horizontal box
func NewHBox(id string) *HBox {
	panel := NewPanel(id)
	panel.layoutDirection = LayoutHorizontal
	return &HBox{Panel: panel}
}

// GetType returns the element type
func (h *HBox) GetType() string {
	return "HBox"
}
