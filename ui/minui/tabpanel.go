package minui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mechanical-lich/mlge/event"
)

// TabPosition defines where tabs are placed
type TabPosition int

const (
	TabsTop    TabPosition = iota // Tabs at the top
	TabsBottom                    // Tabs at the bottom
	TabsLeft                      // Tabs on the left
	TabsRight                     // Tabs on the right
)

// Tab represents a single tab with optional icon
type Tab struct {
	ID      string
	Text    string
	Icon    *Icon
	Content Element // The content panel for this tab
}

// TabPanel is a container with tabs for switching between content panels
type TabPanel struct {
	*ElementBase
	Tabs         []*Tab
	ActiveTabID  string
	TabPosition  TabPosition
	TabHeight    int // Height of tab buttons (or width for left/right)
	TabSpacing   int // Space between tabs
	IconSpacing  int // Space between icon and text in tab
	OnTabChange  func(tabID string)
	hoveredTabID string
}

// NewTabPanel creates a new tab panel
func NewTabPanel(id string, width, height int) *TabPanel {
	tp := &TabPanel{
		ElementBase: NewElementBase(id),
		Tabs:        make([]*Tab, 0),
		TabPosition: TabsTop,
		TabHeight:   32,
		TabSpacing:  2,
		IconSpacing: 4,
	}

	tp.SetSize(width, height)

	// Set default style - only structural properties, colors come from theme
	borderWidth := 1

	tp.style.BorderWidth = &borderWidth

	return tp
}

// AddChild adds a child element to the tab panel
func (tp *TabPanel) AddChild(child Element) {
	tp.children = append(tp.children, child)
	child.SetParent(tp)
}

// AddTab adds a new tab
func (tp *TabPanel) AddTab(id, text string, content Element) *Tab {
	tab := &Tab{
		ID:      id,
		Text:    text,
		Content: content,
	}
	tp.Tabs = append(tp.Tabs, tab)

	// Set first tab as active
	if tp.ActiveTabID == "" {
		tp.ActiveTabID = id
	}

	// Add content as child
	if content != nil {
		tp.AddChild(content)
		content.SetVisible(id == tp.ActiveTabID)
	}

	return tab
}

// AddIconTab adds a new tab with an icon
func (tp *TabPanel) AddIconTab(id, text string, icon *Icon, content Element) *Tab {
	tab := &Tab{
		ID:      id,
		Text:    text,
		Icon:    icon,
		Content: content,
	}
	tp.Tabs = append(tp.Tabs, tab)

	// Set first tab as active
	if tp.ActiveTabID == "" {
		tp.ActiveTabID = id
	}

	// Add content as child
	if content != nil {
		tp.AddChild(content)
		content.SetVisible(id == tp.ActiveTabID)
	}

	return tab
}

// AddIconOnlyTab adds a tab with just an icon (no text)
func (tp *TabPanel) AddIconOnlyTab(id string, icon *Icon, content Element) *Tab {
	return tp.AddIconTab(id, "", icon, content)
}

// SetActiveTab changes the active tab, emitting a TabChangeEvent as if the user
// clicked it.
func (tp *TabPanel) SetActiveTab(id string) {
	tp.setActiveTab(id, true)
}

// SetActiveTabQuiet changes the active tab programmatically WITHOUT emitting a
// TabChangeEvent. Use it when restoring or defaulting the tab so listeners that
// treat the event as a user action (e.g. audio feedback) aren't triggered by
// setup. The OnTabChange callback still fires so app state stays in sync.
func (tp *TabPanel) SetActiveTabQuiet(id string) {
	tp.setActiveTab(id, false)
}

// setActiveTab switches tabs and always fires the OnTabChange callback. It emits
// the TabChangeEvent only when notify is true — the click path and SetActiveTab
// do; the quiet setter does not.
func (tp *TabPanel) setActiveTab(id string, notify bool) {
	if tp.ActiveTabID == id {
		return
	}

	oldTabID := tp.ActiveTabID
	tp.ActiveTabID = id

	// Update content visibility
	for _, tab := range tp.Tabs {
		if tab.Content != nil {
			tab.Content.SetVisible(tab.ID == id)
		}
	}

	if notify {
		playInteraction(EventTypeTabChange, tp.GetID()) // immediate feedback, before the handler
	}

	if tp.OnTabChange != nil {
		tp.OnTabChange(id)
	}

	if notify {
		event.GetQueuedInstance().QueueEvent(TabChangeEvent{
			TabPanelID: tp.GetID(),
			TabPanel:   tp,
			OldTabID:   oldTabID,
			NewTabID:   id,
		})
	}
}

// GetActiveTab returns the currently active tab
func (tp *TabPanel) GetActiveTab() *Tab {
	for _, tab := range tp.Tabs {
		if tab.ID == tp.ActiveTabID {
			return tab
		}
	}
	return nil
}

// GetType returns the element type
func (tp *TabPanel) GetType() string {
	return "TabPanel"
}

// Update handles tab panel interaction
func (tp *TabPanel) Update() {
	if !tp.visible || !tp.enabled {
		return
	}

	tp.UpdateHoverState()

	mx, my := ebiten.CursorPosition()
	items := tp.tabItems()
	tp.hoveredTabID = tabStripHitTest(tp.stripRegion(), tp.TabPosition, items, tp.tabStyle(), mx, my)

	// Handle click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && tp.hoveredTabID != "" {
		tp.SetActiveTab(tp.hoveredTabID)
	}

	// Update active content
	for _, tab := range tp.Tabs {
		if tab.Content != nil && tab.ID == tp.ActiveTabID {
			tab.Content.Update()
		}
	}
}

// tabItems adapts the panel's tabs to the shared strip renderer's TabItem.
func (tp *TabPanel) tabItems() []TabItem {
	items := make([]TabItem, len(tp.Tabs))
	for i, t := range tp.Tabs {
		items[i] = TabItem{ID: t.ID, Label: t.Text, Icon: t.Icon, Enabled: true}
	}
	return items
}

// stripRegion is the absolute rect occupied by the tab strip on its edge.
func (tp *TabPanel) stripRegion() Rect {
	absX, absY := tp.GetAbsolutePosition()
	switch tp.TabPosition {
	case TabsBottom:
		return Rect{X: absX, Y: absY + tp.bounds.Height - tp.TabHeight, Width: tp.bounds.Width, Height: tp.TabHeight}
	case TabsLeft:
		return Rect{X: absX, Y: absY, Width: tp.TabHeight, Height: tp.bounds.Height}
	case TabsRight:
		return Rect{X: absX + tp.bounds.Width - tp.TabHeight, Y: absY, Width: tp.TabHeight, Height: tp.bounds.Height}
	default: // TabsTop
		return Rect{X: absX, Y: absY, Width: tp.bounds.Width, Height: tp.TabHeight}
	}
}

// tabStyle returns the notebook style for this panel's strip, tinting the active
// tab to the content surface so it merges into the panel below.
func (tp *TabPanel) tabStyle() TabStripStyle {
	s := DefaultTabStripStyle()
	if theme := tp.GetTheme(); theme != nil {
		s.ActiveFill = colorToRGBA(theme.Colors.Surface)
		s.InactiveFill = colorToRGBA(theme.Colors.Background)
	}
	return s
}

// Layout calculates dimensions
func (tp *TabPanel) Layout() {
	style := tp.GetComputedStyle()

	width := tp.bounds.Width
	height := tp.bounds.Height

	if style.Width != nil {
		width = *style.Width
	}
	if style.Height != nil {
		height = *style.Height
	}

	width, height = ApplySizeConstraints(width, height, style)

	tp.bounds.Width = width
	tp.bounds.Height = height

	// Layout content panels - use relative coordinates for children
	contentOffset := tp.getContentOffset()
	for _, tab := range tp.Tabs {
		if tab.Content != nil {
			tab.Content.SetBounds(Rect{
				X:      contentOffset.X,
				Y:      contentOffset.Y,
				Width:  tp.bounds.Width - contentOffset.X,
				Height: tp.bounds.Height - contentOffset.Y,
			})
			tab.Content.Layout()
		}
	}
}

func (tp *TabPanel) getContentBounds() Rect {
	absX, absY := tp.GetAbsolutePosition()
	switch tp.TabPosition {
	case TabsTop:
		return Rect{
			X:      absX,
			Y:      absY + tp.TabHeight,
			Width:  tp.bounds.Width,
			Height: tp.bounds.Height - tp.TabHeight,
		}
	case TabsBottom:
		return Rect{
			X:      absX,
			Y:      absY,
			Width:  tp.bounds.Width,
			Height: tp.bounds.Height - tp.TabHeight,
		}
	case TabsLeft:
		return Rect{
			X:      absX + tp.TabHeight,
			Y:      absY,
			Width:  tp.bounds.Width - tp.TabHeight,
			Height: tp.bounds.Height,
		}
	case TabsRight:
		return Rect{
			X:      absX,
			Y:      absY,
			Width:  tp.bounds.Width - tp.TabHeight,
			Height: tp.bounds.Height,
		}
	}
	return Rect{X: absX, Y: absY, Width: tp.bounds.Width, Height: tp.bounds.Height}
}

// getContentOffset returns the relative offset for content (not absolute)
func (tp *TabPanel) getContentOffset() Rect {
	switch tp.TabPosition {
	case TabsTop:
		return Rect{X: 0, Y: tp.TabHeight, Width: tp.bounds.Width, Height: tp.bounds.Height - tp.TabHeight}
	case TabsBottom:
		return Rect{X: 0, Y: 0, Width: tp.bounds.Width, Height: tp.bounds.Height - tp.TabHeight}
	case TabsLeft:
		return Rect{X: tp.TabHeight, Y: 0, Width: tp.bounds.Width - tp.TabHeight, Height: tp.bounds.Height}
	case TabsRight:
		return Rect{X: 0, Y: 0, Width: tp.bounds.Width - tp.TabHeight, Height: tp.bounds.Height}
	}
	return Rect{X: 0, Y: 0, Width: tp.bounds.Width, Height: tp.bounds.Height}
}

// Draw draws the tab panel
func (tp *TabPanel) Draw(screen *ebiten.Image) {
	if !tp.visible {
		return
	}

	style := tp.GetComputedStyle()
	theme := tp.GetTheme()
	contentBounds := tp.getContentBounds()

	// Fill only the content area — the strip row stays clear so there's no
	// leftover panel background beside the last tab (the panel starts under the
	// tabs).
	DrawBackgroundWithTheme(screen, contentBounds, style, theme)

	// Draw the tab strip via the shared notebook renderer.
	drawTabStrip(screen, tp.stripRegion(), tp.TabPosition, tp.tabItems(), tp.ActiveTabID, tp.hoveredTabID, tp.tabStyle())

	// Draw content area border with theme support
	contentStyle := &Style{}
	borderWidth := 1
	contentStyle.BorderWidth = &borderWidth
	DrawBorderWithTheme(screen, contentBounds, contentStyle, theme)

	// Draw active content
	for _, tab := range tp.Tabs {
		if tab.Content != nil && tab.ID == tp.ActiveTabID {
			tab.Content.Draw(screen)
		}
	}
}

// TabChangeEvent is fired when the active tab changes
type TabChangeEvent struct {
	TabPanelID string
	TabPanel   *TabPanel
	OldTabID   string
	NewTabID   string
}

func (e TabChangeEvent) GetType() event.EventType {
	return EventTypeTabChange
}
