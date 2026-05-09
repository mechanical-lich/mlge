# minui — Minimal Vector-Based UI for MLGE

A lightweight, sprite-independent UI library for the MLGE game engine. Vector-rendered widgets, CSS-like styling with cascading inheritance, theme support, and a small set of opinionated conventions for hover/focus/disabled state and overflow handling.

## Quick Start

```go
package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    minui "github.com/mechanical-lich/mlge/ui/minui"
)

func main() {
    gui := minui.NewGUIWithTheme(minui.NewDarkTheme())

    btn := minui.NewButton("save", "Save Game")
    btn.SetPosition(100, 100)
    btn.OnClick = func() { /* ... */ }
    gui.AddElement(btn)

    // Game loop:
    //   gui.Update()
    //   gui.Layout()
    //   gui.Draw(screen)
}
```

## Core Concepts

### Element interface

Every widget implements `Element`: `Update()`, `Layout()`, `Draw(screen)`, plus state queries (`IsVisible`, `IsEnabled`, `GetBounds`, `GetAbsolutePosition`, …). `ElementBase` provides the common implementation; widgets embed it.

### Lifecycle

1. **Update** — handle input, advance hover/focus/scroll state. Skip if `!visible || !enabled`.
2. **Layout** — compute bounds from style + children.
3. **Draw** — render. Disabled widgets render dimmed; truncated text auto-registers an overflow tooltip (see below).

GUI runs all three for you in order. If you draw widgets manually outside `GUI.Draw`, call `minui.FlushOverlays(screen)` at the end of your draw routine.

### Styling

CSS-like properties on each element, with cascading inheritance from the parent. Common pattern:

```go
fontSize := 14
borderRadius := 6
padding := minui.NewEdgeInsets(8)

style := btn.GetStyle()
style.FontSize = &fontSize
style.BorderRadius = &borderRadius
style.Padding = padding

style.HoverStyle = &minui.Style{ /* override fields */ }
```

Property categories:
- **Layout** — `Width`, `Height`, `Min/MaxWidth`, `Min/MaxHeight`, `Padding`, `Margin`
- **Visual** — `BackgroundColor`, `BackgroundImage`, `BorderColor`, `BorderWidth`, `BorderRadius`, `ForegroundColor`, `Opacity`
- **Typography** — `FontSize`, `FontBold`, `FontItalic`, `TextAlign`, `VertAlign`
- **State styles** — `HoverStyle`, `ActiveStyle`, `DisabledStyle`, `FocusStyle`

Use the helpers `NewEdgeInsets(n)`, `NewEdgeInsetsLR(v, h)`, `NewEdgeInsetsTRBL(t, r, b, l)` for padding/margin.

### Themes

A `Theme` carries the palette (Primary, Surface, Text, Border, Focus, …) and is propagated through the element tree. `NewDarkTheme()` returns a sensible default. Widgets prefer theme colors; per-element `Style` overrides win.

```go
gui := minui.NewGUIWithTheme(minui.NewDarkTheme())
```

## Widgets

### Containers
- `Panel` — generic container with optional layout direction
- `VBox`, `HBox` — auto-stacking containers with `Spacing`
- `Modal` — draggable, closeable dialog with title bar; add via `gui.AddModal`
- `ScrollPanel` — vertical-scrolling container; mouse-wheel scrolling, themed scrollbar, content auto-clipped to bounds. Children may be any Element and report correct hit-test coordinates while scrolled.
- `TabPanel` — tabbed container with top/left tab strips
- `Drawer` — slide-out side panel

### Text & input
- `Label` — single- or multi-line text
- `RichText` — labels with mixed colors/sizes/bold spans
- `TextInput` — single-line input with cursor, selection, submit
- `Checkbox` — boolean checkbox with label
- `Toggle` — switch-style toggle
- `RadioButton` / `RadioGroup` — exclusive selection

### Buttons & menus
- `Button` — basic clickable button
- `IconButton` — icon + optional text in 4 layouts (left/right/top/bottom/icon-only)
- `MenuItem`, `MenuHeader` — for sidebars and popup menus
- `PopupMenu` — context/dropdown menu

### Lists & selection
- `ListBox` — scrollable list of strings with hover/select
- `SelectBox` — HTML-style dropdown; expanded list draws on the overlay layer (always on top)
- `TreeView` — hierarchical tree

### Display & feedback
- `ProgressBar` — horizontal progress with optional label
- `ResourceBar` — multi-resource HUD bar (icon + numeric value rows)
- `ScrollingTextArea` — auto-scrolling text log (e.g. message history)
- `ImageWidget` — draws a static image or sprite
- `Icon` — themed icon resource
- `Tooltip` / `TooltipManager` — explicit hover tooltips for any element

### Modals
- `FileModal` — in-engine file browser with directory navigation

## Conventions

### Disabled state

Every interactive widget short-circuits in `Update` when `IsEnabled()` is false (no hover, no click). In `Draw`, text and icons render dimmed via the shared `dimColor` helper. Use `SetEnabled(false)` to grey-out without removing the element.

### Text overflow

Most widgets render labels through `DrawClippedWithTooltip(screen, owner, txt, size, x, y, maxW, color)` instead of raw text rendering. If `txt`'s measured width exceeds `maxW`, it's drawn with a trailing ellipsis. If the user then hovers the widget for ~½ second, an overflow tooltip near the cursor reveals the full string. State is automatically per-element and ages out when the widget stops drawing.

If you build a custom widget and want the same behavior, call `DrawClippedWithTooltip` from its `Draw` method, passing the widget itself as `owner`.

### Overlay layer (always-on-top draws)

Some widgets need to paint above everything else (open dropdowns, tooltips, context menus). Rather than juggling z-order, call `minui.QueueOverlay(func(screen) { ... })` from inside your Draw — the closure runs after the GUI has finished its main pass. `GUI.Draw` flushes the queue automatically. If you draw widgets manually (without using `GUI.Draw`), call `minui.FlushOverlays(screen)` at the end of your draw routine instead.

### Coordinate hit-testing

`Element.GetAbsolutePosition()` walks up the parent chain accounting for content padding, modal title bars, and scroll panel offsets, so child widgets work correctly inside any container without per-widget plumbing. Custom containers that want to introduce their own offset (e.g. a viewport) can implement `GetScrollOffsetY() int`.

## Modal example

```go
modal := minui.NewModal("confirm", "Confirm Action", 400, 200)
modal.SetPosition(120, 100)
modal.Closeable = true

okBtn := minui.NewButton("ok", "OK")
okBtn.SetBounds(minui.Rect{X: 150, Y: 140, Width: 80, Height: 32})
okBtn.OnClick = func() { modal.SetVisible(false) }
modal.AddChild(okBtn)

gui.AddModal(modal)
modal.SetVisible(true)
```

## ScrollPanel example

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

The panel reports its own scroll offset to children via the `GetScrollOffsetY` interface, so menu items inside it click correctly even after scrolling.

## Architecture

```
ui/minui/
├── element.go             # Element interface + ElementBase
├── style.go               # Style struct, EdgeInsets, computed style
├── theme.go               # Theme palette + dark default
├── rendering.go           # Vector draw helpers, colorToRGBA, dimColor
├── overlay.go             # Top-level overlay queue (QueueOverlay/FlushOverlays)
├── overflowtooltip.go     # DrawClippedWithTooltip + state tracking
├── gui.go                 # GUI manager (root container)
│
├── containers.go          # Panel, VBox, HBox
├── modal.go               # Modal
├── scrollpanel.go         # ScrollPanel
├── tabpanel.go            # TabPanel
├── drawer.go              # Drawer
│
├── label.go               # Label
├── richtext.go            # RichText
├── input.go               # TextInput + Checkbox
├── toggle.go              # Toggle
├── radio.go               # RadioButton + RadioGroup
│
├── button.go              # Button
├── iconbutton.go          # IconButton
├── menuitem.go            # MenuItem + MenuHeader
├── popupmenu.go           # PopupMenu
│
├── listbox.go             # ListBox
├── selectbox.go           # SelectBox
├── treeview.go            # TreeView
│
├── progress.go            # ProgressBar
├── resourcebar.go         # ResourceBar
├── scrollingtextarea.go   # ScrollingTextArea
├── image.go               # ImageWidget
├── icon.go                # Icon
│
├── tooltip.go             # Tooltip
├── tooltipmanager.go      # TooltipManager
└── filemodal.go           # FileModal
```

## Design principles

1. **No external sprite dependencies** — vector rendering only (themes can opt into sprite-backed buttons).
2. **Cascading style** — parent styles inherit; per-element overrides win.
3. **Predictable state** — visible/enabled/hovered/focused/active are first-class on every element.
4. **Overflow is silent** — never let text spill outside its widget; reveal the full content on hover.
5. **Z-order via overlays, not layering** — defer always-on-top draws to the overlay queue rather than restructuring the tree.

## License

Part of the MLGE game engine project.
