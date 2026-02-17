package minui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mechanical-lich/mlge/event"
	"github.com/mechanical-lich/mlge/text"
)

// TreeNode represents a node in a tree view
type TreeNode struct {
	ID       string
	Text     string
	Icon     *Icon
	Children []*TreeNode
	Expanded bool
	Selected bool
	Data     interface{} // User data

	// Internal
	parent *TreeNode
	depth  int
}

// NewTreeNode creates a new tree node
func NewTreeNode(id, text string) *TreeNode {
	return &TreeNode{
		ID:       id,
		Text:     text,
		Children: make([]*TreeNode, 0),
		Expanded: false,
		Selected: false,
	}
}

// NewTreeNodeWithIcon creates a new tree node with an icon
func NewTreeNodeWithIcon(id, text string, icon *Icon) *TreeNode {
	return &TreeNode{
		ID:       id,
		Text:     text,
		Icon:     icon,
		Children: make([]*TreeNode, 0),
		Expanded: false,
		Selected: false,
	}
}

// AddChild adds a child node
func (n *TreeNode) AddChild(child *TreeNode) {
	child.parent = n
	child.depth = n.depth + 1
	n.Children = append(n.Children, child)
}

// IsLeaf returns true if this node has no children
func (n *TreeNode) IsLeaf() bool {
	return len(n.Children) == 0
}

// Toggle toggles the expanded state
func (n *TreeNode) Toggle() {
	n.Expanded = !n.Expanded
}

// TreeView is a hierarchical list with collapsible nodes
type TreeView struct {
	*ElementBase
	Roots          []*TreeNode
	RowHeight      int
	IndentWidth    int
	IconSpacing    int
	ExpandIcon     *Icon // Icon for collapsed nodes (e.g., ">")
	CollapseIcon   *Icon // Icon for expanded nodes (e.g., "v")
	OnSelect       func(node *TreeNode)
	OnToggle       func(node *TreeNode)
	selectedNode   *TreeNode
	hoveredNode    *TreeNode
	visibleNodes   []*TreeNode // Flattened list of visible nodes
	scrollOffset   int
	maxVisibleRows int
	draggingScroll bool
	dragOffsetY    int
}

// NewTreeView creates a new tree view
func NewTreeView(id string, width, height int) *TreeView {
	tv := &TreeView{
		ElementBase: NewElementBase(id),
		Roots:       make([]*TreeNode, 0),
		RowHeight:   24,
		IndentWidth: 16,
		IconSpacing: 4,
	}

	tv.SetSize(width, height)

	// Set default style - only structural properties, colors come from theme
	borderWidth := 1

	tv.style.BorderWidth = &borderWidth

	return tv
}

// AddRoot adds a root node to the tree
func (tv *TreeView) AddRoot(node *TreeNode) {
	node.depth = 0
	node.parent = nil
	tv.Roots = append(tv.Roots, node)
	tv.updateVisibleNodes()
}

// Clear removes all nodes from the tree
func (tv *TreeView) Clear() {
	tv.Roots = make([]*TreeNode, 0)
	tv.selectedNode = nil
	tv.hoveredNode = nil
	tv.visibleNodes = nil
	tv.scrollOffset = 0
}

// GetSelectedNode returns the currently selected node
func (tv *TreeView) GetSelectedNode() *TreeNode {
	return tv.selectedNode
}

// SelectNode selects a specific node
func (tv *TreeView) SelectNode(node *TreeNode) {
	if tv.selectedNode != nil {
		tv.selectedNode.Selected = false
	}
	tv.selectedNode = node
	if node != nil {
		node.Selected = true
	}
}

// SelectByID selects a node by its ID
func (tv *TreeView) SelectByID(id string) {
	node := tv.findNodeByID(id)
	tv.SelectNode(node)
}

func (tv *TreeView) findNodeByID(id string) *TreeNode {
	var search func(nodes []*TreeNode) *TreeNode
	search = func(nodes []*TreeNode) *TreeNode {
		for _, n := range nodes {
			if n.ID == id {
				return n
			}
			if found := search(n.Children); found != nil {
				return found
			}
		}
		return nil
	}
	return search(tv.Roots)
}

// ExpandAll expands all nodes
func (tv *TreeView) ExpandAll() {
	var expand func(nodes []*TreeNode)
	expand = func(nodes []*TreeNode) {
		for _, n := range nodes {
			n.Expanded = true
			expand(n.Children)
		}
	}
	expand(tv.Roots)
	tv.updateVisibleNodes()
}

// CollapseAll collapses all nodes
func (tv *TreeView) CollapseAll() {
	var collapse func(nodes []*TreeNode)
	collapse = func(nodes []*TreeNode) {
		for _, n := range nodes {
			n.Expanded = false
			collapse(n.Children)
		}
	}
	collapse(tv.Roots)
	tv.updateVisibleNodes()
}

// updateVisibleNodes rebuilds the flattened list of visible nodes
func (tv *TreeView) updateVisibleNodes() {
	tv.visibleNodes = make([]*TreeNode, 0)
	var collect func(nodes []*TreeNode)
	collect = func(nodes []*TreeNode) {
		for _, n := range nodes {
			tv.visibleNodes = append(tv.visibleNodes, n)
			if n.Expanded && len(n.Children) > 0 {
				collect(n.Children)
			}
		}
	}
	collect(tv.Roots)

	// Calculate max visible rows
	tv.maxVisibleRows = tv.bounds.Height / tv.RowHeight
}

// GetType returns the element type
func (tv *TreeView) GetType() string {
	return "TreeView"
}

// Update handles tree view interaction
func (tv *TreeView) Update() {
	if !tv.visible || !tv.enabled {
		return
	}

	tv.UpdateHoverState()

	if len(tv.visibleNodes) == 0 {
		tv.updateVisibleNodes()
	}

	mx, my := ebiten.CursorPosition()
	absX, absY := tv.GetAbsolutePosition()

	// Find hovered node
	tv.hoveredNode = nil
	if tv.hovered {
		relY := my - absY
		rowIndex := (relY / tv.RowHeight) + tv.scrollOffset
		if rowIndex >= 0 && rowIndex < len(tv.visibleNodes) {
			tv.hoveredNode = tv.visibleNodes[rowIndex]
		}
	}

	// Handle click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && tv.hovered {
		if tv.hoveredNode != nil {
			node := tv.hoveredNode

			// Check if clicked on expand/collapse area
			nodeIndent := absX + node.depth*tv.IndentWidth
			expandAreaWidth := 16 // Area for the expand/collapse toggle

			if !node.IsLeaf() && mx >= nodeIndent && mx < nodeIndent+expandAreaWidth {
				// Toggle expand/collapse
				node.Toggle()
				tv.updateVisibleNodes()
				if tv.OnToggle != nil {
					tv.OnToggle(node)
				}
				event.GetQueuedInstance().QueueEvent(TreeViewToggleEvent{
					TreeViewID: tv.GetID(),
					TreeView:   tv,
					Node:       node,
					Expanded:   node.Expanded,
				})
			} else {
				// Select the node
				tv.SelectNode(node)
				if tv.OnSelect != nil {
					tv.OnSelect(node)
				}
				event.GetQueuedInstance().QueueEvent(TreeViewSelectEvent{
					TreeViewID: tv.GetID(),
					TreeView:   tv,
					Node:       node,
				})
			}
		}
	}

	// Mouse wheel scroll
	if tv.hovered {
		_, yoff := ebiten.Wheel()
		if yoff != 0 {
			tv.scrollOffset -= int(yoff * 3)
			tv.clampScrollOffset()
		}
	}
}

func (tv *TreeView) clampScrollOffset() {
	if tv.scrollOffset < 0 {
		tv.scrollOffset = 0
	}
	maxOffset := len(tv.visibleNodes) - tv.maxVisibleRows
	if maxOffset < 0 {
		maxOffset = 0
	}
	if tv.scrollOffset > maxOffset {
		tv.scrollOffset = maxOffset
	}
}

// Layout calculates dimensions
func (tv *TreeView) Layout() {
	style := tv.GetComputedStyle()

	width := tv.bounds.Width
	height := tv.bounds.Height

	if style.Width != nil {
		width = *style.Width
	}
	if style.Height != nil {
		height = *style.Height
	}

	width, height = ApplySizeConstraints(width, height, style)

	tv.bounds.Width = width
	tv.bounds.Height = height
	tv.maxVisibleRows = height / tv.RowHeight
}

// Draw draws the tree view
func (tv *TreeView) Draw(screen *ebiten.Image) {
	if !tv.visible {
		return
	}

	style := tv.GetComputedStyle()
	theme := tv.GetTheme()
	absX, absY := tv.GetAbsolutePosition()
	absBounds := Rect{
		X:      absX,
		Y:      absY,
		Width:  tv.bounds.Width,
		Height: tv.bounds.Height,
	}

	// Draw background with theme support
	DrawBackgroundWithTheme(screen, absBounds, style, theme)
	DrawBorderWithTheme(screen, absBounds, style, theme)

	// Draw visible nodes
	startIndex := tv.scrollOffset
	endIndex := startIndex + tv.maxVisibleRows + 1
	if endIndex > len(tv.visibleNodes) {
		endIndex = len(tv.visibleNodes)
	}

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

	// Get selection and hover colors from theme
	selectedBg := color.RGBA{80, 100, 140, 255}
	hoverBg := color.RGBA{70, 70, 80, 255}
	if theme != nil {
		selectedBg = colorToRGBA(theme.Colors.Primary)
		hoverBg = colorToRGBA(theme.Colors.Surface)
		hoverBg.R = min(hoverBg.R+20, 255)
		hoverBg.G = min(hoverBg.G+20, 255)
		hoverBg.B = min(hoverBg.B+20, 255)
	}

	for i := startIndex; i < endIndex; i++ {
		node := tv.visibleNodes[i]
		rowY := absY + (i-startIndex)*tv.RowHeight

		// Draw selection/hover background
		rowBounds := Rect{X: absX, Y: rowY, Width: tv.bounds.Width, Height: tv.RowHeight}
		if node.Selected {
			DrawRect(screen, rowBounds, selectedBg)
		} else if node == tv.hoveredNode {
			DrawRect(screen, rowBounds, hoverBg)
		}

		// Calculate content position with indent
		contentX := absX + 4 + node.depth*tv.IndentWidth
		contentY := rowY + (tv.RowHeight-fontSize)/2

		// Draw expand/collapse indicator for non-leaf nodes
		if !node.IsLeaf() {
			indicator := "▶"
			if node.Expanded {
				indicator = "▼"
			}
			if tv.ExpandIcon != nil && !node.Expanded {
				tv.ExpandIcon.Draw(screen, contentX, rowY+(tv.RowHeight-tv.ExpandIcon.ScaledHeight())/2)
			} else if tv.CollapseIcon != nil && node.Expanded {
				tv.CollapseIcon.Draw(screen, contentX, rowY+(tv.RowHeight-tv.CollapseIcon.ScaledHeight())/2)
			} else {
				text.Draw(screen, indicator, float64(fontSize-2), contentX, contentY, textColor)
			}
			contentX += 16
		}

		// Draw node icon
		if node.Icon != nil {
			iconY := rowY + (tv.RowHeight-node.Icon.ScaledHeight())/2
			node.Icon.Draw(screen, contentX, iconY)
			contentX += node.Icon.ScaledWidth() + tv.IconSpacing
		}

		// Draw node text
		text.Draw(screen, node.Text, float64(fontSize), contentX, contentY, textColor)
	}

	// Draw scrollbar if needed
	if len(tv.visibleNodes) > tv.maxVisibleRows {
		tv.drawScrollbar(screen, absBounds)
	}
}

func (tv *TreeView) drawScrollbar(screen *ebiten.Image, absBounds Rect) {
	theme := tv.GetTheme()
	barX := absBounds.X + absBounds.Width - 12
	barY := absBounds.Y + 2
	barW := 8
	barH := absBounds.Height - 4

	totalRows := len(tv.visibleNodes)
	if totalRows <= tv.maxVisibleRows {
		return
	}

	// Draw track with theme color
	trackColor := color.RGBA{60, 60, 70, 200}
	if theme != nil {
		trackColor = colorToRGBA(theme.Colors.Surface)
		trackColor.A = 200
	}
	trackBounds := Rect{X: barX, Y: barY, Width: barW, Height: barH}
	DrawRoundedRect(screen, trackBounds, 3, trackColor)

	// Draw thumb with theme color
	thumbH := max(barH*tv.maxVisibleRows/totalRows, 16)
	scrollRange := barH - thumbH
	var thumbY int
	if totalRows > tv.maxVisibleRows && scrollRange > 0 {
		thumbY = barY + (scrollRange*tv.scrollOffset)/(totalRows-tv.maxVisibleRows)
	} else {
		thumbY = barY
	}

	thumbColor := color.RGBA{120, 120, 140, 255}
	if theme != nil {
		thumbColor = colorToRGBA(theme.Colors.Border)
	}
	if tv.hovered {
		thumbColor.R = min(thumbColor.R+20, 255)
		thumbColor.G = min(thumbColor.G+20, 255)
		thumbColor.B = min(thumbColor.B+20, 255)
	}

	thumbBounds := Rect{X: barX, Y: thumbY, Width: barW, Height: thumbH}
	DrawRoundedRect(screen, thumbBounds, 3, thumbColor)
}

// TreeViewSelectEvent is fired when a tree node is selected
type TreeViewSelectEvent struct {
	TreeViewID string
	TreeView   *TreeView
	Node       *TreeNode
}

func (e TreeViewSelectEvent) GetType() event.EventType {
	return EventTypeTreeViewSelect
}

// TreeViewToggleEvent is fired when a tree node is expanded/collapsed
type TreeViewToggleEvent struct {
	TreeViewID string
	TreeView   *TreeView
	Node       *TreeNode
	Expanded   bool
}

func (e TreeViewToggleEvent) GetType() event.EventType {
	return EventTypeTreeViewToggle
}
