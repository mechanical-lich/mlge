package minui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mechanical-lich/mlge/event"
	"github.com/mechanical-lich/mlge/text"
)

// TabStripStyle holds the colours and metrics for the shared "notebook" tab
// renderer used by both TabBar and TabPanel. The active tab is filled with
// ActiveFill and drawn without a border on its content-facing edge so it merges
// into the content area — set ActiveFill to the content background for a clean
// join.
type TabStripStyle struct {
	ActiveFill   color.RGBA
	InactiveFill color.RGBA
	HoverFill    color.RGBA
	DisabledText color.RGBA
	Border       color.RGBA
	ActiveText   color.RGBA
	InactiveText color.RGBA
	FontSize     int
	MinTabW      int // minimum tab width (horizontal strips)
	VTabH        int // tab height (vertical strips)
	Pad          int // inner padding along the label axis
}

// DefaultTabStripStyle returns a dark-theme notebook style.
func DefaultTabStripStyle() TabStripStyle {
	return TabStripStyle{
		ActiveFill:   color.RGBA{12, 15, 24, 255},
		InactiveFill: color.RGBA{26, 30, 42, 255},
		HoverFill:    color.RGBA{40, 46, 62, 255},
		DisabledText: color.RGBA{90, 100, 120, 255},
		Border:       color.RGBA{70, 92, 122, 255},
		ActiveText:   color.RGBA{226, 236, 252, 255},
		InactiveText: color.RGBA{150, 172, 200, 255},
		FontSize:     14,
		MinTabW:      64,
		VTabH:        32,
		Pad:          14,
	}
}

// TabItem is one tab's metadata for the shared strip renderer.
type TabItem struct {
	ID      string
	Label   string
	Icon    *Icon
	Enabled bool // false renders dimmed and ignores clicks
}

func tabStripHorizontal(pos TabPosition) bool { return pos == TabsTop || pos == TabsBottom }

func tabItemWidth(t TabItem, st TabStripStyle) int {
	w := st.Pad * 2
	if t.Icon != nil {
		w += t.Icon.ScaledWidth() + 4
	}
	w += len(t.Label) * (st.FontSize*6/10 + 1)
	if w < st.MinTabW {
		w = st.MinTabW
	}
	return w
}

// tabStripLayout returns the absolute rect of each tab within region.
func tabStripLayout(region Rect, pos TabPosition, tabs []TabItem, st TabStripStyle) []Rect {
	rects := make([]Rect, len(tabs))
	if tabStripHorizontal(pos) {
		x := region.X
		for i, t := range tabs {
			w := tabItemWidth(t, st)
			rects[i] = Rect{X: x, Y: region.Y, Width: w, Height: region.Height}
			x += w
		}
	} else {
		y := region.Y
		for i := range tabs {
			rects[i] = Rect{X: region.X, Y: y, Width: region.Width, Height: st.VTabH}
			y += st.VTabH
		}
	}
	return rects
}

// tabStripHitTest returns the id of the tab under the absolute point, or "".
func tabStripHitTest(region Rect, pos TabPosition, tabs []TabItem, style TabStripStyle, mx, my int) string {
	for i, r := range tabStripLayout(region, pos, tabs, style) {
		if mx >= r.X && mx < r.X+r.Width && my >= r.Y && my < r.Y+r.Height {
			return tabs[i].ID
		}
	}
	return ""
}

// drawTabStrip renders the tabs on their edge of region. The active tab merges
// into the content area (no border on the content-facing edge).
func drawTabStrip(dst *ebiten.Image, region Rect, pos TabPosition, tabs []TabItem, activeID, hoveredID string, style TabStripStyle) {
	rects := tabStripLayout(region, pos, tabs, style)
	activeIdx := -1

	// 1) Inactive tabs first (their content-facing edge is later covered by the
	//    baseline, giving them a border there).
	for i, t := range tabs {
		if t.ID == activeID {
			activeIdx = i
			continue
		}
		fill := style.InactiveFill
		if t.Enabled && t.ID == hoveredID {
			fill = style.HoverFill
		}
		drawTabItem(dst, rects[i], pos, t, fill, false, style)
	}

	// 2) Baseline along the strip's content-facing edge.
	drawTabBaseline(dst, region, pos, style)

	// 3) Active tab last, over the baseline, so it opens into the content.
	if activeIdx >= 0 {
		drawTabItem(dst, rects[activeIdx], pos, tabs[activeIdx], style.ActiveFill, true, style)
	}
}

func drawTabBaseline(dst *ebiten.Image, region Rect, pos TabPosition, style TabStripStyle) {
	switch pos {
	case TabsTop:
		DrawRect(dst, Rect{X: region.X, Y: region.Y + region.Height - 1, Width: region.Width, Height: 1}, style.Border)
	case TabsBottom:
		DrawRect(dst, Rect{X: region.X, Y: region.Y, Width: region.Width, Height: 1}, style.Border)
	case TabsLeft:
		DrawRect(dst, Rect{X: region.X + region.Width - 1, Y: region.Y, Width: 1, Height: region.Height}, style.Border)
	case TabsRight:
		DrawRect(dst, Rect{X: region.X, Y: region.Y, Width: 1, Height: region.Height}, style.Border)
	}
}

func drawTabItem(dst *ebiten.Image, r Rect, pos TabPosition, t TabItem, fill color.RGBA, active bool, style TabStripStyle) {
	DrawRect(dst, r, fill)

	// Borders on every edge except the one facing the content area.
	b := style.Border
	top := Rect{X: r.X, Y: r.Y, Width: r.Width, Height: 1}
	bottom := Rect{X: r.X, Y: r.Y + r.Height - 1, Width: r.Width, Height: 1}
	left := Rect{X: r.X, Y: r.Y, Width: 1, Height: r.Height}
	right := Rect{X: r.X + r.Width - 1, Y: r.Y, Width: 1, Height: r.Height}
	switch pos {
	case TabsTop:
		DrawRect(dst, top, b)
		DrawRect(dst, left, b)
		DrawRect(dst, right, b)
	case TabsBottom:
		DrawRect(dst, bottom, b)
		DrawRect(dst, left, b)
		DrawRect(dst, right, b)
	case TabsLeft:
		DrawRect(dst, top, b)
		DrawRect(dst, bottom, b)
		DrawRect(dst, left, b)
	case TabsRight:
		DrawRect(dst, top, b)
		DrawRect(dst, bottom, b)
		DrawRect(dst, right, b)
	}

	textColor := style.InactiveText
	if active {
		textColor = style.ActiveText
	}
	if !t.Enabled {
		textColor = style.DisabledText
	}

	// Centre the icon+label block within the tab.
	textW := len(t.Label) * (style.FontSize*6/10 + 1)
	iconW := 0
	if t.Icon != nil {
		iconW = t.Icon.ScaledWidth() + 4
	}
	x := r.X + (r.Width-iconW-textW)/2
	if x < r.X+4 {
		x = r.X + 4
	}
	y := r.Y + (r.Height-style.FontSize)/2
	if t.Icon != nil {
		iconY := r.Y + (r.Height-t.Icon.ScaledHeight())/2
		t.Icon.Draw(dst, x, iconY)
		x += t.Icon.ScaledWidth() + 4
	}
	if t.Label != "" {
		text.DrawClipped(dst, t.Label, float64(style.FontSize), x, y, r.X+r.Width-x-4, textColor)
	}
}

// TabBar is a standalone tab strip: it draws the tabs on one edge of its bounds
// and reports selection, leaving the caller to fill the content area. (For a
// widget that also manages per-tab content, see TabPanel.)
type TabBar struct {
	*ElementBase
	Position  TabPosition
	Style     TabStripStyle
	Tabs      []TabItem
	ActiveID  string
	OnChange  func(id string)
	hoveredID string
}

// NewTabBar creates a tab strip. bounds is the strip's own rect (thin along its
// edge); the caller draws content in the remaining area.
func NewTabBar(id string, bounds Rect, pos TabPosition) *TabBar {
	tb := &TabBar{
		ElementBase: NewElementBase(id),
		Position:    pos,
		Style:       DefaultTabStripStyle(),
	}
	tb.SetBounds(bounds)
	return tb
}

// AddTab appends an enabled tab; the first added becomes active.
func (tb *TabBar) AddTab(id, label string) *TabBar {
	tb.Tabs = append(tb.Tabs, TabItem{ID: id, Label: label, Enabled: true})
	if tb.ActiveID == "" {
		tb.ActiveID = id
	}
	return tb
}

// SetTabEnabled toggles a tab's clickability/dimming.
func (tb *TabBar) SetTabEnabled(id string, enabled bool) {
	for i := range tb.Tabs {
		if tb.Tabs[i].ID == id {
			tb.Tabs[i].Enabled = enabled
			return
		}
	}
}

// SetActive switches the active tab (no OnChange fired).
func (tb *TabBar) SetActive(id string) { tb.ActiveID = id }

func (tb *TabBar) region() Rect {
	x, y := tb.GetAbsolutePosition()
	return Rect{X: x, Y: y, Width: tb.bounds.Width, Height: tb.bounds.Height}
}

func (tb *TabBar) Update() {
	if !tb.visible || !tb.enabled {
		return
	}
	mx, my := ebiten.CursorPosition()
	tb.hoveredID = tabStripHitTest(tb.region(), tb.Position, tb.Tabs, tb.Style, mx, my)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && tb.hoveredID != "" {
		for _, t := range tb.Tabs {
			if t.ID == tb.hoveredID && t.Enabled && t.ID != tb.ActiveID {
				oldID := tb.ActiveID
				tb.ActiveID = t.ID
				playInteraction(EventTypeTabBarChange, tb.GetID()) // immediate feedback, before the handler
				if tb.OnChange != nil {
					tb.OnChange(t.ID)
				}
				event.GetQueuedInstance().QueueEvent(TabBarChangeEvent{
					TabBarID: tb.GetID(),
					TabBar:   tb,
					OldTabID: oldID,
					NewTabID: t.ID,
				})
				break
			}
		}
	}
}

func (tb *TabBar) Draw(screen *ebiten.Image) {
	if !tb.visible {
		return
	}
	drawTabStrip(screen, tb.region(), tb.Position, tb.Tabs, tb.ActiveID, tb.hoveredID, tb.Style)
}

func (tb *TabBar) GetType() string { return "TabBar" }
