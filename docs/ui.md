---
layout: default
title: UI Framework
nav_order: 12
---

# UI Framework (minui)

`github.com/mechanical-lich/mlge/ui/minui`

A comprehensive retained-mode UI library with theming, style inheritance, dual rendering paths (vector and sprite-based), and a rich widget set.

## Setup

```go
import "github.com/mechanical-lich/mlge/ui/minui"

gui := minui.NewGUI()
```

## Core Concepts

### Elements and Containers

All UI widgets implement the `Element` interface. Containers (like `Panel`) hold child elements and manage layout.

### Layout

Containers support three layout directions:

| Constant | Description |
|----------|-------------|
| `LayoutVertical` | Children arranged top to bottom |
| `LayoutHorizontal` | Children arranged left to right |
| `LayoutNone` | No automatic layout (manual positioning) |

### Styling

Every element can have a `Style` applied with CSS-like properties:

```go
style := minui.Style{
    Width:           200,
    Height:          40,
    Padding:         minui.EdgeInsets{Top: 8, Right: 12, Bottom: 8, Left: 12},
    BackgroundColor: color.RGBA{40, 40, 40, 255},
    BorderColor:     color.RGBA{100, 100, 100, 255},
    BorderWidth:     1,
    FontColor:       color.White,
    FontSize:        14,
    HoverStyle: &minui.Style{
        BackgroundColor: color.RGBA{60, 60, 60, 255},
    },
}
```

### Theming

Apply a consistent look across all widgets with a `Theme`:

```go
theme := &minui.Theme{
    Name: "dark",
    Colors: minui.Colors{
        Primary:    color.RGBA{70, 130, 180, 255},
        Secondary:  color.RGBA{60, 60, 60, 255},
        Background: color.RGBA{30, 30, 30, 255},
        Surface:    color.RGBA{45, 45, 45, 255},
        Text:       color.White,
        Border:     color.RGBA{80, 80, 80, 255},
        Focus:      color.RGBA{70, 130, 180, 255},
    },
}
gui.SetTheme(theme)
```

When a `Theme.SpriteSheet` is set, widgets render using 9-slice sprites instead of vector drawing.

## Widgets

### Panel

Container that arranges child elements:

```go
panel := minui.NewPanel("main-panel")
panel.SetLayoutDirection(minui.LayoutVertical)
panel.AddChild(button)
panel.AddChild(label)
gui.AddElement(panel)
```

### Button

```go
button := minui.NewButton("submit-btn", "Submit")
// Listen for clicks via the event system
// Event type: "ui.button.click"
```

### IconButton

Button with an icon and text:

```go
icon := minui.NewIcon(spriteSheet, 0, 0, 16, 16)
iconBtn := minui.NewIconButton("save-btn", icon, "Save")
```

### Label

Static text display with multiline support:

```go
label := minui.NewLabel("status-label", "Health: 100")
label.SetText("Health: 80")

// Create with a specific color, or change color at runtime
label := minui.NewLabelWithColor("status-label", "Health: 100", color.RGBA{0, 255, 0, 255})
label.SetColor(color.RGBA{255, 0, 0, 255})
```

### TextInput

Single-line text input with cursor:

```go
input := minui.NewTextInput("name-input", "Enter name...")
// Event type: "ui.textinput.change"
```

### ListBox

Scrollable selectable list:

```go
listBox := minui.NewListBox("inventory-list", []string{"Sword", "Shield", "Potion"})
// Event type: "ui.listbox.select"
```

### SelectBox

Dropdown selection:

```go
selectBox := minui.NewSelectBox("difficulty-select", []string{"Easy", "Normal", "Hard"})
```

### RadioButton / RadioGroup

```go
group := minui.NewRadioGroup("mode-group")
radio1 := minui.NewRadioButton("mode-easy", "Easy")
radio2 := minui.NewRadioButton("mode-hard", "Hard")
group.AddChild(radio1)
group.AddChild(radio2)
```

### Toggle

On/off switch:

```go
toggle := minui.NewToggle("music-toggle", "Music")
```

### ProgressBar

```go
bar := minui.NewProgressBar("health-bar")
bar.SetValue(0.75) // 75%
```

### ScrollingTextArea

Multi-line scrollable text display. Text is automatically word-wrapped. Scrolls to the bottom when new text is added, and supports mouse-wheel scrolling and a draggable scrollbar thumb.

```go
textArea := minui.NewScrollingTextArea("log-area", 400, 200)
textArea.AddText("Player entered the dungeon")
textArea.AddText("A wild goblin appears!")

// Per-line color override (nil = widget default)
textArea.AddColoredText("You are badly wounded!", color.RGBA{255, 60, 60, 255})

textArea.Clear()
```

### RichText

A list of independently styled text spans, each rendered as one line. Height is calculated automatically from content. Useful for HUD panels where different lines need different colors, sizes, or strikethrough.

```go
rt := minui.NewRichText("hover-panel", 200)
rt.LineHeight = 14

rt.AddSpan(minui.TextSpan{
    Text:  "Goblin",
    Color: color.RGBA{255, 200, 100, 255},
    Size:  13,
})
rt.AddSpan(minui.TextSpan{
    Text:   "HP: 4/10",
    Color:  color.RGBA{200, 80, 80, 255},
    Size:   11,
    Indent: 8,
})
rt.AddSpan(minui.TextSpan{
    Text:          "left arm",
    Color:         color.RGBA{120, 80, 160, 255},
    Size:          11,
    Indent:        8,
    Strikethrough: true,
})

rt.Clear()
rt.SetPosition(x, y)
rt.Draw(screen)
```

`RichText` does not scroll — wrap lines yourself before adding spans (see `text.Wrap`). It has no background or border; position it manually and call `Draw` directly (no parent required).

### ImageWidget

Displays an `*ebiten.Image` scaled to the widget's bounds:

```go
widget := minui.NewImageWidget("minimap", 150, 150)
widget.SetPosition(x, y)

// Update the image each frame
widget.Image = generateMinimapImage()
widget.Draw(screen)
```

`ImageWidget` has no parent requirement — call `Draw` directly. If `Image` is nil the widget draws nothing.

### TabPanel

Tabbed container:

```go
tabPanel := minui.NewTabPanel("settings-tabs", 600, 400)
tabPanel.AddTab("general", "General", generalPanel)
tabPanel.AddTab("audio", "Audio", audioPanel)

// Programmatically switch to a tab by ID
tabPanel.SetActiveTab("audio")

// Event type: "ui.tabpanel.change"
```

Tab positions: `TabsTop`, `TabsBottom`, `TabsLeft`, `TabsRight`

### TreeView / TreeNode

Hierarchical tree with expand/collapse:

```go
root := minui.NewTreeNode("root", "Game Objects")
child1 := minui.NewTreeNode("player", "Player")
child2 := minui.NewTreeNode("enemies", "Enemies")
root.AddChild(child1)
root.AddChild(child2)
// Event type: "ui.treeview.select"
```

### Modal

Draggable dialog window:

```go
modal := minui.NewModal("settings-modal", "Settings", 400, 300)
modal.AddChild(settingsPanel)
gui.ShowModal(modal)
// Event type: "ui.modal.close"
```

### PopupMenu

Context/popup menu with nested submenus:

```go
menu := minui.NewPopupMenu("context-menu")
menu.AddItem(minui.NewMenuItem("edit-item", "Edit"))
menu.AddItem(minui.NewMenuItem("delete-item", "Delete"))
// Event type: "ui.popupmenu.select"
```

### Drawer

Sliding panel from screen edges:

```go
drawer := minui.NewDrawer("side-drawer", minui.DrawerLeft)
drawer.AddChild(menuPanel)
```

### ScrollPanel

Vertical-scrolling generic container. Children may be any `Element`; the panel mouse-wheel scrolls when hovered, clips children to its bounds via `SubImage`, and draws a themed scrollbar on the right when content overflows.

```go
scroll := minui.NewScrollPanel("recipes")
scroll.SetBounds(minui.Rect{X: 20, Y: 60, Width: 300, Height: 420})
for i, name := range recipeNames {
    mi := minui.NewMenuItem("r_"+name, name)
    mi.SetBounds(minui.Rect{X: 0, Y: i * 22, Width: 300, Height: 22})
    scroll.AddChild(mi)
}
modal.AddChild(scroll)
```

Children placed inside a `ScrollPanel` continue to report correct hit-test coordinates while scrolled. The mechanism: `ScrollPanel` exposes `GetScrollOffsetY() int`, and `ElementBase.GetAbsolutePosition` subtracts that offset when its parent satisfies the interface. Custom containers that introduce their own viewport offset can use the same hook.

### Tooltip

```go
tooltipManager := minui.NewTooltipManager()
tooltip := minui.NewTooltip("btn-tooltip")
tooltip.SetText("Click to submit")
tooltipManager.Register(button, tooltip)
```

### ResourceBar

Horizontal bar showing resource icons and values:

```go
bar := minui.NewResourceBar("resources")
```

### FileModal

File browser dialog:

```go
fileModal := minui.NewFileModal("file-browser", "Open File", 500, 400, minui.FileModalLoad)
```

## UI Events

All UI events are dispatched through the queued event system. Common event types:

| Event Type | Triggered By |
|------------|-------------|
| `ui.button.click` | Button press |
| `ui.textinput.change` | Text input value change |
| `ui.listbox.select` | List item selection |
| `ui.modal.close` | Modal closed |
| `ui.tabpanel.change` | Tab switch |
| `ui.treeview.select` | Tree node selection |
| `ui.popupmenu.select` | Menu item selected |

## Rendering

The UI supports two rendering paths:

- **Vector rendering** (default) — Draws widgets using filled rectangles, borders, and text
- **Sprite rendering** — When `Theme.SpriteSheet` is set, draws widgets using 9-slice sprites from a sprite sheet for a pixel-art or custom visual style

## Overlay layer

Some widgets need to paint above everything else — open dropdowns, tooltips, context menus. Rather than juggling z-order or restructuring the element tree, defer the always-on-top portion of a widget's `Draw` to the overlay queue:

```go
func (w *MyWidget) Draw(screen *ebiten.Image) {
    // ...draw the base widget inline...
    if w.expanded {
        minui.QueueOverlay(func(s *ebiten.Image) {
            w.drawPopup(s)
        })
    }
}
```

`GUI.Draw` flushes the overlay queue at the end of its main pass, so overlays land on top of everything (including modals). If you draw widgets manually outside `GUI.Draw`, call `minui.FlushOverlays(screen)` at the end of your draw routine instead — otherwise the overlay queue will never be drained and dropdowns will not appear.

`SelectBox` uses this for its expanded list. Custom popups, autocompletes, and floating tooltips should follow the same pattern.

## Overflow tooltips

Most text-bearing widgets (`Button`, `IconButton`, `Label`, `MenuItem`, `Toggle`, …) render their labels through `DrawClippedWithTooltip(screen, owner, txt, size, x, y, maxW, color)` instead of raw text rendering. If `txt` exceeds `maxW`, it's drawn with a trailing ellipsis. If the user hovers the widget for ~½ second, a small floating tooltip near the cursor reveals the full string. State is keyed per-element and ages out automatically.

Custom widgets that want the same behavior call:

```go
minui.DrawClippedWithTooltip(screen, w /* owner */, w.Text, fontSize, x, y, maxW, textColor)
```

Pure text helpers also exist in the `text` package:

- `text.Measure(txt, size)` — measure a string
- `text.Truncate(txt, size, maxW)` — return the longest prefix + ellipsis that fits
- `text.DrawClipped(screen, txt, size, x, y, maxW, color)` — draw, truncating if needed (no tooltip behavior)

## Disabled state

Every interactive widget short-circuits in `Update` when `IsEnabled()` is false (no hover, no click). In `Draw`, text and icons render dimmed via the shared `dimColor` helper. Use `SetEnabled(false)` to grey out a widget without removing it.

```go
btn.SetEnabled(false)              // grey, non-interactive
btn.SetEnabled(settlement.HasTech("metallurgy"))
```

Widgets with this behavior: `Button`, `IconButton`, `MenuItem`, `Toggle`, `RadioButton`. Containers don't disable their children — call `SetEnabled` on each child you want gated.
