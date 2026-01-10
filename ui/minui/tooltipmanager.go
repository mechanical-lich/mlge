package minui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// TooltipManager manages tooltips for multiple elements
type TooltipManager struct {
	tooltips      map[string]*Tooltip
	activeTooltip *Tooltip
	enabled       bool
}

// NewTooltipManager creates a new tooltip manager
func NewTooltipManager() *TooltipManager {
	return &TooltipManager{
		tooltips: make(map[string]*Tooltip),
		enabled:  true,
	}
}

// Register attaches a tooltip to an element
func (tm *TooltipManager) Register(element Element, title, text string, icon *Icon) *Tooltip {
	tooltip := NewTooltip(element.GetID() + "_tooltip")
	tooltip.SetContent(title, text, icon)
	tooltip.SetTarget(element)
	tm.tooltips[element.GetID()] = tooltip
	return tooltip
}

// RegisterSimple attaches a simple text-only tooltip
func (tm *TooltipManager) RegisterSimple(element Element, text string) *Tooltip {
	return tm.Register(element, "", text, nil)
}

// RegisterWithTitle attaches a tooltip with title and text
func (tm *TooltipManager) RegisterWithTitle(element Element, title, text string) *Tooltip {
	return tm.Register(element, title, text, nil)
}

// Unregister removes a tooltip for an element
func (tm *TooltipManager) Unregister(element Element) {
	delete(tm.tooltips, element.GetID())
}

// UnregisterByID removes a tooltip by element ID
func (tm *TooltipManager) UnregisterByID(elementID string) {
	delete(tm.tooltips, elementID)
}

// Get retrieves a tooltip for an element
func (tm *TooltipManager) Get(element Element) *Tooltip {
	return tm.tooltips[element.GetID()]
}

// GetByID retrieves a tooltip by element ID
func (tm *TooltipManager) GetByID(elementID string) *Tooltip {
	return tm.tooltips[elementID]
}

// SetEnabled enables or disables all tooltips
func (tm *TooltipManager) SetEnabled(enabled bool) {
	tm.enabled = enabled
	if !enabled {
		tm.HideAll()
	}
}

// IsEnabled returns whether tooltips are enabled
func (tm *TooltipManager) IsEnabled() bool {
	return tm.enabled
}

// SetGlobalDelay sets the delay for all managed tooltips
func (tm *TooltipManager) SetGlobalDelay(frames int) {
	for _, tooltip := range tm.tooltips {
		tooltip.Delay = frames
	}
}

// SetGlobalPosition sets the position for all managed tooltips
func (tm *TooltipManager) SetGlobalPosition(position TooltipPosition) {
	for _, tooltip := range tm.tooltips {
		tooltip.Position = position
	}
}

// HideAll hides all tooltips
func (tm *TooltipManager) HideAll() {
	for _, tooltip := range tm.tooltips {
		tooltip.Hide()
	}
	tm.activeTooltip = nil
}

// Update updates all managed tooltips
func (tm *TooltipManager) Update() {
	if !tm.enabled {
		return
	}

	var newActive *Tooltip

	for _, tooltip := range tm.tooltips {
		tooltip.Update()
		if tooltip.showing {
			newActive = tooltip
		}
	}

	// Only show one tooltip at a time
	if newActive != nil && tm.activeTooltip != nil && newActive != tm.activeTooltip {
		tm.activeTooltip.Hide()
	}
	tm.activeTooltip = newActive
}

// Draw draws the active tooltip (if any)
func (tm *TooltipManager) Draw(screen *ebiten.Image) {
	if !tm.enabled {
		return
	}

	// Draw active tooltip last (on top)
	if tm.activeTooltip != nil {
		tm.activeTooltip.Draw(screen)
	}
}

// UpdateForElement updates tooltip content for a specific element
func (tm *TooltipManager) UpdateForElement(element Element, title, text string) {
	tooltip := tm.tooltips[element.GetID()]
	if tooltip != nil {
		tooltip.Title = title
		tooltip.Text = text
	}
}

// Clear removes all tooltips
func (tm *TooltipManager) Clear() {
	tm.tooltips = make(map[string]*Tooltip)
	tm.activeTooltip = nil
}

// Count returns the number of registered tooltips
func (tm *TooltipManager) Count() int {
	return len(tm.tooltips)
}
