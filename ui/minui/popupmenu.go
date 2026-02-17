package minui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mechanical-lich/mlge/event"
	"github.com/mechanical-lich/mlge/text"
)

// PopupMenuItem represents an item in a popup menu
type PopupMenuItem struct {
	ID       string
	Text     string
	Icon     *Icon
	Disabled bool
	Shortcut string      // Keyboard shortcut display (e.g., "Ctrl+S")
	SubMenu  *PopupMenu  // Nested submenu
	Data     interface{} // Custom data
}

// PopupMenu is a contextual menu that can appear anywhere on screen
type PopupMenu struct {
	*ElementBase
	Items         []*PopupMenuItem
	SelectedIndex int
	HoveredIndex  int
	ItemHeight    int
	SubMenu       *PopupMenu // Currently open submenu
	ParentMenu    *PopupMenu // Parent menu if this is a submenu
	OnSelect      func(item *PopupMenuItem)
}

// PopupMenuSelectEvent is emitted when a menu item is selected
type PopupMenuSelectEvent struct {
	Menu *PopupMenu
	Item *PopupMenuItem
}

// GetType returns the event type
func (e PopupMenuSelectEvent) GetType() event.EventType {
	return EventTypePopupMenuSelect
}

// NewPopupMenu creates a new popup menu
func NewPopupMenu(id string) *PopupMenu {
	pm := &PopupMenu{
		ElementBase:   NewElementBase(id),
		Items:         make([]*PopupMenuItem, 0),
		SelectedIndex: -1,
		HoveredIndex:  -1,
		ItemHeight:    28,
	}

	pm.SetSize(200, 100)
	pm.visible = false // Start hidden

	// Set default style - only structural properties, colors come from theme
	borderWidth := 1
	pm.style.BorderWidth = &borderWidth

	return pm
}

// AddItem adds an item to the menu
func (pm *PopupMenu) AddItem(id, text string, icon *Icon) *PopupMenuItem {
	item := &PopupMenuItem{
		ID:   id,
		Text: text,
		Icon: icon,
	}
	pm.Items = append(pm.Items, item)
	pm.updateSize()
	return item
}

// AddItemWithShortcut adds an item with a keyboard shortcut display
func (pm *PopupMenu) AddItemWithShortcut(id, text string, icon *Icon, shortcut string) *PopupMenuItem {
	item := &PopupMenuItem{
		ID:       id,
		Text:     text,
		Icon:     icon,
		Shortcut: shortcut,
	}
	pm.Items = append(pm.Items, item)
	pm.updateSize()
	return item
}

// AddSubmenu adds an item that opens a submenu
func (pm *PopupMenu) AddSubmenu(id, text string, icon *Icon, submenu *PopupMenu) *PopupMenuItem {
	item := &PopupMenuItem{
		ID:      id,
		Text:    text,
		Icon:    icon,
		SubMenu: submenu,
	}
	submenu.ParentMenu = pm
	pm.Items = append(pm.Items, item)
	pm.updateSize()
	return item
}

// AddSeparator adds a visual separator
func (pm *PopupMenu) AddSeparator() {
	item := &PopupMenuItem{
		ID:       "_separator",
		Disabled: true,
	}
	pm.Items = append(pm.Items, item)
	pm.updateSize()
}

// updateSize recalculates the menu size
func (pm *PopupMenu) updateSize() {
	width := 180

	// Calculate width based on items
	for _, item := range pm.Items {
		itemWidth := 40 // Base padding
		if item.Icon != nil {
			itemWidth += item.Icon.ScaledWidth() + 8
		}
		itemWidth += len(item.Text) * 8
		if item.Shortcut != "" {
			itemWidth += len(item.Shortcut)*7 + 16
		}
		if item.SubMenu != nil {
			itemWidth += 20 // Arrow indicator
		}
		if itemWidth > width {
			width = itemWidth
		}
	}

	height := len(pm.Items) * pm.ItemHeight

	pm.SetSize(width, height)
}

// Show displays the menu at the given position
func (pm *PopupMenu) Show(x, y int) {
	pm.SetPosition(x, y)
	pm.visible = true
	pm.HoveredIndex = -1
	pm.SelectedIndex = -1
}

// Hide hides the menu and any submenus
func (pm *PopupMenu) Hide() {
	pm.visible = false
	if pm.SubMenu != nil {
		pm.SubMenu.Hide()
		pm.SubMenu = nil
	}
}

// GetType returns the element type
func (pm *PopupMenu) GetType() string {
	return "PopupMenu"
}

// Update updates the popup menu
func (pm *PopupMenu) Update() {
	if !pm.visible {
		return
	}

	absX, absY := pm.GetAbsolutePosition()
	mx, my := ebiten.CursorPosition()

	// Check if mouse is within menu bounds
	inBounds := mx >= absX && mx < absX+pm.bounds.Width &&
		my >= absY && my < absY+pm.bounds.Height

	if inBounds {
		// Calculate which item is hovered
		relY := my - absY
		pm.HoveredIndex = relY / pm.ItemHeight
		if pm.HoveredIndex >= len(pm.Items) {
			pm.HoveredIndex = -1
		}

		// Check for separator
		if pm.HoveredIndex >= 0 && pm.Items[pm.HoveredIndex].ID == "_separator" {
			pm.HoveredIndex = -1
		}

		// Handle submenu opening
		if pm.HoveredIndex >= 0 && pm.Items[pm.HoveredIndex].SubMenu != nil {
			subItem := pm.Items[pm.HoveredIndex]
			if pm.SubMenu != subItem.SubMenu {
				if pm.SubMenu != nil {
					pm.SubMenu.Hide()
				}
				pm.SubMenu = subItem.SubMenu
				// Position submenu to the right
				subX := absX + pm.bounds.Width
				subY := absY + pm.HoveredIndex*pm.ItemHeight
				pm.SubMenu.Show(subX, subY)
			}
		} else if pm.SubMenu != nil && pm.HoveredIndex >= 0 {
			pm.SubMenu.Hide()
			pm.SubMenu = nil
		}

		// Handle click
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if pm.HoveredIndex >= 0 {
				item := pm.Items[pm.HoveredIndex]
				if !item.Disabled && item.SubMenu == nil {
					pm.SelectedIndex = pm.HoveredIndex
					if pm.OnSelect != nil {
						pm.OnSelect(item)
					}
					event.GetQueuedInstance().QueueEvent(PopupMenuSelectEvent{
						Menu: pm,
						Item: item,
					})
					// Hide all menus up to root
					pm.hideToRoot()
				}
			}
		}
	} else {
		pm.HoveredIndex = -1
		// Check if clicking outside closes menu
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			// Check if click is in any submenu
			inSubmenu := false
			sub := pm.SubMenu
			for sub != nil {
				subX, subY := sub.GetAbsolutePosition()
				if mx >= subX && mx < subX+sub.bounds.Width &&
					my >= subY && my < subY+sub.bounds.Height {
					inSubmenu = true
					break
				}
				sub = sub.SubMenu
			}
			if !inSubmenu && pm.ParentMenu == nil {
				// Root menu clicked outside - hide
				pm.Hide()
			}
		}
	}

	// Update submenu
	if pm.SubMenu != nil {
		pm.SubMenu.Update()
	}
}

// hideToRoot hides this menu and all parent menus
func (pm *PopupMenu) hideToRoot() {
	current := pm
	for current != nil {
		current.visible = false
		if current.SubMenu != nil {
			current.SubMenu.Hide()
			current.SubMenu = nil
		}
		current = current.ParentMenu
	}
}

// Layout calculates dimensions
func (pm *PopupMenu) Layout() {
	// Size is calculated in updateSize
}

// Draw draws the popup menu
func (pm *PopupMenu) Draw(screen *ebiten.Image) {
	if !pm.visible {
		return
	}

	style := pm.GetComputedStyle()
	theme := pm.GetTheme()
	absX, absY := pm.GetAbsolutePosition()
	absBounds := Rect{
		X:      absX,
		Y:      absY,
		Width:  pm.bounds.Width,
		Height: pm.bounds.Height,
	}

	// Draw background with theme support
	DrawBackgroundWithTheme(screen, absBounds, style, theme)

	fontSize := 14
	if style.FontSize != nil {
		fontSize = *style.FontSize
	}

	// Get colors from style, then theme, then defaults
	textColor := color.RGBA{230, 230, 230, 255}
	if style.ForegroundColor != nil {
		textColor = colorToRGBA(*style.ForegroundColor)
	} else if theme != nil {
		textColor = colorToRGBA(theme.Colors.Text)
	}

	hoverColor := color.RGBA{70, 100, 150, 255}
	disabledColor := color.RGBA{100, 100, 100, 255}
	separatorColor := color.RGBA{60, 60, 70, 255}
	if theme != nil {
		hoverColor = colorToRGBA(theme.Colors.Primary)
		disabledColor = colorToRGBA(theme.Colors.Disabled)
		separatorColor = colorToRGBA(theme.Colors.Border)
		separatorColor.A = 128
	}

	for i, item := range pm.Items {
		itemY := absY + i*pm.ItemHeight

		// Handle separator
		if item.ID == "_separator" {
			sepY := float32(itemY + pm.ItemHeight/2)
			var path vector.Path
			path.MoveTo(float32(absX+8), sepY)
			path.LineTo(float32(absX+pm.bounds.Width-8), sepY)
			vertices, indices := path.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
				Width: 1,
			})
			for i := range vertices {
				vertices[i].ColorR = 1
				vertices[i].ColorG = 1
				vertices[i].ColorB = 1
				vertices[i].ColorA = 1
			}
			colorImg := ebiten.NewImage(1, 1)
			colorImg.Fill(separatorColor)
			screen.DrawTriangles(vertices, indices, colorImg, &ebiten.DrawTrianglesOptions{})
			continue
		}

		// Draw hover highlight
		if i == pm.HoveredIndex && !item.Disabled {
			DrawRect(screen, Rect{
				X:      absX,
				Y:      itemY,
				Width:  pm.bounds.Width,
				Height: pm.ItemHeight,
			}, hoverColor)
		}

		x := absX + 8
		centerY := itemY + pm.ItemHeight/2

		// Draw icon
		if item.Icon != nil {
			iconY := centerY - item.Icon.ScaledHeight()/2
			if item.Disabled {
				item.Icon.DrawWithOpacity(screen, x, iconY, 0.5)
			} else {
				item.Icon.Draw(screen, x, iconY)
			}
			x += item.Icon.ScaledWidth() + 8
		}

		// Draw text
		clr := textColor
		if item.Disabled {
			clr = disabledColor
		}
		textY := centerY - fontSize/2
		text.Draw(screen, item.Text, float64(fontSize), x, textY, clr)

		// Draw shortcut (right-aligned)
		if item.Shortcut != "" {
			shortcutX := absX + pm.bounds.Width - len(item.Shortcut)*7 - 12
			text.Draw(screen, item.Shortcut, float64(fontSize-2), shortcutX, textY+1, disabledColor)
		}

		// Draw submenu arrow
		if item.SubMenu != nil {
			arrowX := absX + pm.bounds.Width - 16
			text.Draw(screen, "▶", float64(fontSize-4), arrowX, textY, textColor)
		}
	}

	// Draw border with theme support
	DrawBorderWithTheme(screen, absBounds, style, theme)

	// Draw submenu
	if pm.SubMenu != nil {
		pm.SubMenu.Draw(screen)
	}
}
