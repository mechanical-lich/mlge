package minui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mechanical-lich/mlge/event"
	"github.com/mechanical-lich/mlge/text"
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

// SetActiveTab changes the active tab
func (tp *TabPanel) SetActiveTab(id string) {
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

	if tp.OnTabChange != nil {
		tp.OnTabChange(id)
	}

	event.GetQueuedInstance().QueueEvent(TabChangeEvent{
		TabPanelID: tp.GetID(),
		TabPanel:   tp,
		OldTabID:   oldTabID,
		NewTabID:   id,
	})
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
	absX, absY := tp.GetAbsolutePosition()

	// Find hovered tab
	tp.hoveredTabID = ""
	if tp.hovered {
		tp.hoveredTabID = tp.getTabAtPosition(mx-absX, my-absY)
	}

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

func (tp *TabPanel) getTabAtPosition(relX, relY int) string {
	switch tp.TabPosition {
	case TabsTop:
		if relY > tp.TabHeight {
			return ""
		}
		x := 0
		for _, tab := range tp.Tabs {
			tabW := tp.calculateTabWidth(tab)
			if relX >= x && relX < x+tabW {
				return tab.ID
			}
			x += tabW + tp.TabSpacing
		}
	case TabsBottom:
		if relY < tp.bounds.Height-tp.TabHeight {
			return ""
		}
		x := 0
		for _, tab := range tp.Tabs {
			tabW := tp.calculateTabWidth(tab)
			if relX >= x && relX < x+tabW {
				return tab.ID
			}
			x += tabW + tp.TabSpacing
		}
	case TabsLeft:
		if relX > tp.TabHeight {
			return ""
		}
		y := 0
		for _, tab := range tp.Tabs {
			tabH := tp.calculateTabHeight(tab)
			if relY >= y && relY < y+tabH {
				return tab.ID
			}
			y += tabH + tp.TabSpacing
		}
	case TabsRight:
		if relX < tp.bounds.Width-tp.TabHeight {
			return ""
		}
		y := 0
		for _, tab := range tp.Tabs {
			tabH := tp.calculateTabHeight(tab)
			if relY >= y && relY < y+tabH {
				return tab.ID
			}
			y += tabH + tp.TabSpacing
		}
	}
	return ""
}

func (tp *TabPanel) calculateTabWidth(tab *Tab) int {
	width := 16 // Base padding
	if tab.Icon != nil {
		width += tab.Icon.ScaledWidth()
		if tab.Text != "" {
			width += tp.IconSpacing
		}
	}
	if tab.Text != "" {
		width += len(tab.Text) * 8
	}
	return width
}

func (tp *TabPanel) calculateTabHeight(tab *Tab) int {
	height := 8 // Base padding
	if tab.Icon != nil {
		height += tab.Icon.ScaledHeight()
		if tab.Text != "" {
			height += tp.IconSpacing
		}
	}
	if tab.Text != "" {
		height += 14
	}
	return max(height, tp.TabHeight)
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
	absX, absY := tp.GetAbsolutePosition()
	absBounds := Rect{
		X:      absX,
		Y:      absY,
		Width:  tp.bounds.Width,
		Height: tp.bounds.Height,
	}

	// Draw background with theme support
	DrawBackgroundWithTheme(screen, absBounds, style, theme)

	// Draw tabs
	tp.drawTabs(screen, absX, absY, style, theme)

	// Draw content area border with theme support
	contentBounds := tp.getContentBounds()
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

func (tp *TabPanel) drawTabs(screen *ebiten.Image, absX, absY int, style *Style, theme *Theme) {
	fontSize := 14
	if style.FontSize != nil {
		fontSize = *style.FontSize
	}

	// Get text color from style, then theme, then default
	textColor := color.RGBA{255, 255, 255, 255}
	if style.ForegroundColor != nil {
		textColor = colorToRGBA(*style.ForegroundColor)
	} else if theme != nil {
		textColor = colorToRGBA(theme.Colors.Text)
	}

	// Get tab colors from theme
	activeBg := color.RGBA{70, 70, 90, 255}
	hoverBg := color.RGBA{60, 60, 75, 255}
	inactiveBg := color.RGBA{40, 40, 50, 255}
	if theme != nil {
		activeBg = colorToRGBA(theme.Colors.Primary)
		hoverBg = colorToRGBA(theme.Colors.Surface)
		hoverBg.R = min(hoverBg.R+10, 255)
		hoverBg.G = min(hoverBg.G+10, 255)
		hoverBg.B = min(hoverBg.B+10, 255)
		inactiveBg = colorToRGBA(theme.Colors.Surface)
	}

	switch tp.TabPosition {
	case TabsTop:
		x := absX
		for _, tab := range tp.Tabs {
			tabW := tp.calculateTabWidth(tab)
			tabBounds := Rect{X: x, Y: absY, Width: tabW, Height: tp.TabHeight}

			// Draw tab background
			bgColor := inactiveBg
			if tab.ID == tp.ActiveTabID {
				bgColor = activeBg
			} else if tab.ID == tp.hoveredTabID {
				bgColor = hoverBg
			}
			DrawRect(screen, tabBounds, bgColor)

			// Draw tab content
			contentX := x + 8
			contentY := absY + (tp.TabHeight-fontSize)/2

			if tab.Icon != nil {
				iconY := absY + (tp.TabHeight-tab.Icon.ScaledHeight())/2
				tab.Icon.Draw(screen, contentX, iconY)
				contentX += tab.Icon.ScaledWidth() + tp.IconSpacing
			}

			if tab.Text != "" {
				text.Draw(screen, tab.Text, float64(fontSize), contentX, contentY, textColor)
			}

			x += tabW + tp.TabSpacing
		}

	case TabsLeft:
		y := absY
		for _, tab := range tp.Tabs {
			tabH := tp.calculateTabHeight(tab)
			tabBounds := Rect{X: absX, Y: y, Width: tp.TabHeight, Height: tabH}

			// Draw tab background
			bgColor := inactiveBg
			if tab.ID == tp.ActiveTabID {
				bgColor = activeBg
			} else if tab.ID == tp.hoveredTabID {
				bgColor = hoverBg
			}
			DrawRect(screen, tabBounds, bgColor)

			// Draw tab content (centered, icon above text)
			if tab.Icon != nil {
				iconX := absX + (tp.TabHeight-tab.Icon.ScaledWidth())/2
				iconY := y + 4
				if tab.Text == "" {
					iconY = y + (tabH-tab.Icon.ScaledHeight())/2
				}
				tab.Icon.Draw(screen, iconX, iconY)
			}

			if tab.Text != "" {
				textX := absX + (tp.TabHeight-len(tab.Text)*fontSize*6/10)/2
				textY := y + tabH - fontSize - 4
				text.Draw(screen, tab.Text, float64(fontSize-2), textX, textY, textColor)
			}

			y += tabH + tp.TabSpacing
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
