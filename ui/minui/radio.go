package minui

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mechanical-lich/mlge/event"
	"github.com/mechanical-lich/mlge/text/v2"
)

// RadioButton is a circular selectable option
type RadioButton struct {
	*ElementBase
	Label    string
	Selected bool
	OnChange func(selected bool)
	group    *RadioGroup // Reference to parent group, if any
}

// NewRadioButton creates a new radio button
func NewRadioButton(id, label string) *RadioButton {
	rb := &RadioButton{
		ElementBase: NewElementBase(id),
		Label:       label,
		Selected:    false,
	}

	rb.SetSize(20, 20)

	// Set default style
	borderColor := color.Color(color.RGBA{100, 100, 120, 255})
	borderWidth := 2
	bgColor := color.Color(color.RGBA{255, 255, 255, 255})

	rb.style.BorderColor = &borderColor
	rb.style.BorderWidth = &borderWidth
	rb.style.BackgroundColor = &bgColor

	return rb
}

// GetType returns the element type
func (rb *RadioButton) GetType() string {
	return "RadioButton"
}

// SetGroup sets the radio button's group
func (rb *RadioButton) SetGroup(group *RadioGroup) {
	rb.group = group
}

// Update handles radio button interaction
func (rb *RadioButton) Update() {
	if !rb.visible {
		return
	}

	rb.UpdateHoverState()

	// Get absolute position for hit detection
	absX, absY := rb.GetAbsolutePosition()
	absBounds := Rect{
		X:      absX,
		Y:      absY,
		Width:  rb.bounds.Width,
		Height: rb.bounds.Height,
	}

	mx, my := ebiten.CursorPosition()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if absBounds.Contains(mx, my) {
			oldSelected := rb.Selected

			if rb.group != nil {
				// In a group - selecting this radio button
				rb.group.Select(rb)
			} else {
				// Not in a group - toggle like a round checkbox
				rb.Selected = !rb.Selected
			}

			// Fire callback if state changed
			if rb.Selected != oldSelected {
				if rb.OnChange != nil {
					rb.OnChange(rb.Selected)
				}

				// Fire event
				event.GetQueuedInstance().QueueEvent(RadioButtonChangeEvent{
					RadioButtonID: rb.GetID(),
					RadioButton:   rb,
					Selected:      rb.Selected,
				})
			}
		}
	}
}

// Layout calculates dimensions
func (rb *RadioButton) Layout() {
	style := rb.GetComputedStyle()

	// Start with default radio button size
	width := rb.bounds.Width
	height := rb.bounds.Height

	// Apply width/height from style if specified
	if style.Width != nil {
		width = *style.Width
	}
	if style.Height != nil {
		height = *style.Height
	}

	// Apply min/max size constraints
	width, height = ApplySizeConstraints(width, height, style)

	rb.bounds.Width = width
	rb.bounds.Height = height
}

// Draw draws the radio button
func (rb *RadioButton) Draw(screen *ebiten.Image) {
	if !rb.visible {
		return
	}

	// Get absolute position for drawing
	absX, absY := rb.GetAbsolutePosition()

	style := rb.GetComputedStyle()

	// Get colors
	borderColor := color.RGBA{100, 100, 120, 255}
	if style.BorderColor != nil {
		if rgba, ok := (*style.BorderColor).(color.RGBA); ok {
			borderColor = rgba
		}
	}

	bgColor := color.RGBA{255, 255, 255, 255}
	if style.BackgroundColor != nil {
		if rgba, ok := (*style.BackgroundColor).(color.RGBA); ok {
			bgColor = rgba
		}
	}

	borderWidth := float32(2.0)
	if style.BorderWidth != nil {
		borderWidth = float32(*style.BorderWidth)
	}

	// Calculate center and radius
	centerX := float32(absX) + float32(rb.bounds.Width)/2
	centerY := float32(absY) + float32(rb.bounds.Height)/2
	radius := float32(rb.bounds.Width) / 2

	// Draw background circle
	vector.DrawFilledCircle(screen, centerX, centerY, radius, bgColor, true)

	// Draw border circle
	vector.StrokeCircle(screen, centerX, centerY, radius, borderWidth, borderColor, true)

	// Draw inner filled circle if selected
	if rb.Selected {
		innerRadius := radius * 0.5
		selectedColor := color.RGBA{100, 120, 180, 255}
		if rb.hovered {
			selectedColor = color.RGBA{120, 140, 200, 255}
		}
		vector.DrawFilledCircle(screen, centerX, centerY, innerRadius, selectedColor, true)
	}

	// Highlight on hover
	if rb.hovered && !rb.Selected {
		hoverColor := color.RGBA{220, 220, 230, 255}
		vector.DrawFilledCircle(screen, centerX, centerY, radius-borderWidth, hoverColor, true)
	}

	// Draw label text to the right of the radio button
	if rb.Label != "" {
		fontSize := 14
		if style.FontSize != nil {
			fontSize = *style.FontSize
		}

		textColor := color.RGBA{0, 0, 0, 255}
		if style.ForegroundColor != nil {
			r, g, b, a := (*style.ForegroundColor).RGBA()
			textColor = color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
		}

		// Position text to the right of the circle with some spacing
		textX := absX + rb.bounds.Width + 5
		textY := absY + (rb.bounds.Height-fontSize)/2
		text.Draw(screen, rb.Label, float64(fontSize), textX, textY, textColor)
	}
}

// RadioGroup manages a group of radio buttons where only one can be selected.
// It is a full Element so it can be added to modals/panels directly.
type RadioGroup struct {
	*ElementBase
	selectedID        string
	OnSelectionChange func(selectedID string, selectedButton *RadioButton)
	Vertical          bool // If true, lay out buttons vertically instead of horizontally
	Spacing           int  // Space between buttons
}

// NewRadioGroup creates a new radio button group
func NewRadioGroup(id string) *RadioGroup {
	return &RadioGroup{
		ElementBase: NewElementBase(id),
		Spacing:     8,
	}
}

// GetType returns the element type
func (rg *RadioGroup) GetType() string {
	return "RadioGroup"
}

// AddButton adds a radio button to the group as a child element
func (rg *RadioGroup) AddButton(button *RadioButton) {
	button.SetGroup(rg)
	rg.AddChild(button)
}

// AddChild overrides AddChild to ensure radio buttons added as children
// are also linked to the group and have proper parent reference
func (rg *RadioGroup) AddChild(el Element) {
	if rb, ok := el.(*RadioButton); ok {
		rb.SetGroup(rg)
	}
	rg.ElementBase.AddChild(el)
	el.SetParent(rg)
}

// Select selects a specific radio button and deselects all others
func (rg *RadioGroup) Select(button *RadioButton) {
	// Deselect all buttons (iterate children)
	for _, child := range rg.children {
		if btn, ok := child.(*RadioButton); ok {
			btn.Selected = false
		}
	}

	// Select the specified button
	button.Selected = true
	rg.selectedID = button.GetID()

	// Fire callback
	if rg.OnSelectionChange != nil {
		rg.OnSelectionChange(button.GetID(), button)
	}

	// Fire event
	event.GetQueuedInstance().QueueEvent(RadioGroupChangeEvent{
		RadioGroupID:   rg.GetID(),
		RadioGroup:     rg,
		SelectedID:     button.GetID(),
		SelectedButton: button,
	})
}

// SelectByID selects a radio button by its ID
func (rg *RadioGroup) SelectByID(id string) {
	for _, child := range rg.children {
		if btn, ok := child.(*RadioButton); ok {
			if btn.GetID() == id {
				rg.Select(btn)
				return
			}
		}
	}
}

// GetSelected returns the currently selected radio button
func (rg *RadioGroup) GetSelected() *RadioButton {
	for _, child := range rg.children {
		if btn, ok := child.(*RadioButton); ok {
			if btn.Selected {
				return btn
			}
		}
	}
	return nil
}

// GetSelectedID returns the ID of the currently selected radio button
func (rg *RadioGroup) GetSelectedID() string {
	return rg.selectedID
}

// GetButtons returns all radio buttons in the group
func (rg *RadioGroup) GetButtons() []*RadioButton {
	buttons := make([]*RadioButton, 0)
	for _, child := range rg.children {
		if btn, ok := child.(*RadioButton); ok {
			buttons = append(buttons, btn)
		}
	}
	return buttons
}

// Clear deselects all radio buttons in the group
func (rg *RadioGroup) Clear() {
	for _, child := range rg.children {
		if btn, ok := child.(*RadioButton); ok {
			btn.Selected = false
		}
	}
	rg.selectedID = ""
}

// Update updates all child radio buttons
func (rg *RadioGroup) Update() {
	if !rg.visible {
		return
	}
	for _, child := range rg.children {
		child.Update()
	}
}

// Layout lays out the group (children are positioned manually or by container)
func (rg *RadioGroup) Layout() {
	// Layout children and position them within the group.
	// Supports both horizontal and vertical layout via rg.Vertical flag.
	style := rg.GetComputedStyle()

	paddingLeft := 0
	paddingTop := 0
	spacing := rg.Spacing
	if style != nil {
		if style.Padding != nil {
			paddingLeft = style.Padding.Left
			paddingTop = style.Padding.Top
		}
	}

	x := paddingLeft
	y := paddingTop
	maxRight := 0
	maxBottom := 0

	for _, child := range rg.children {
		if !child.IsVisible() {
			continue
		}

		// Let the child compute its preferred size first
		child.Layout()
		cb := child.GetBounds()

		// Position the child relative to the group's content origin
		if rg.Vertical {
			child.SetBounds(Rect{X: paddingLeft, Y: y, Width: cb.Width, Height: cb.Height})
		} else {
			child.SetBounds(Rect{X: x, Y: paddingTop, Width: cb.Width, Height: cb.Height})
		}

		// Account for margins if present
		cs := child.GetComputedStyle()
		right := x + cb.Width
		bottom := y + cb.Height

		// If the child is a RadioButton, include its label width in
		// spacing/measurement but do NOT change the radio's own size.
		labelExtra := 0
		labelHeight := 0
		if rb, ok := child.(*RadioButton); ok {
			if rb.Label != "" {
				// Determine font size (fallback to 14 as used in Draw)
				fontSize := 14
				if cs != nil && cs.FontSize != nil {
					fontSize = *cs.FontSize
				}
				tw, th := text.Measure(rb.Label, float64(fontSize))
				// Add the small label gap used in Draw (5px)
				labelExtra = int(math.Ceil(tw)) + 5
				labelHeight = int(math.Ceil(th))
				right += labelExtra
				if y+labelHeight > bottom {
					bottom = y + labelHeight
				}
			}
		}

		if cs != nil && cs.Margin != nil {
			right += cs.Margin.Right
			bottom += cs.Margin.Bottom
			if rg.Vertical {
				y += cs.Margin.Top
			} else {
				x += cs.Margin.Left
			}
		}

		if right > maxRight {
			maxRight = right
		}
		if bottom > maxBottom {
			maxBottom = bottom
		}

		// Advance position for next child
		if rg.Vertical {
			y += cb.Height + spacing
			if labelHeight > cb.Height {
				y += labelHeight - cb.Height
			}
		} else {
			x += cb.Width + labelExtra + spacing
		}
	}

	// Add group's own padding/border to measured size
	if style != nil {
		if style.Padding != nil {
			maxRight += style.Padding.Left + style.Padding.Right
			maxBottom += style.Padding.Top + style.Padding.Bottom
		}
		if style.BorderWidth != nil {
			bw := *style.BorderWidth
			maxRight += bw * 2
			maxBottom += bw * 2
		}
	}

	// Apply min/max constraints and assign bounds
	w, h := ApplySizeConstraints(maxRight, maxBottom, style)
	rg.bounds.Width = w
	rg.bounds.Height = h
}

// Draw draws all child radio buttons
func (rg *RadioGroup) Draw(screen *ebiten.Image) {
	if !rg.visible {
		return
	}
	for _, child := range rg.children {
		child.Draw(screen)
	}
}
