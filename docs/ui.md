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

Multi-line scrollable text display:

```go
textArea := minui.NewScrollingTextArea("log-area", 400, 200)
textArea.AddLine("Player entered the dungeon")
textArea.AddLine("A wild goblin appears!")
```

### TabPanel

Tabbed container:

```go
tabPanel := minui.NewTabPanel("settings-tabs", 600, 400)
tabPanel.AddTab("General", generalPanel)
tabPanel.AddTab("Audio", audioPanel)
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
