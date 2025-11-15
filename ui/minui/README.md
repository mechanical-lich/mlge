# Min-UI - Minimal Vector-Based UI Library for MLGE

A lightweight, sprite-independent UI library for the MLGE game engine that uses vector graphics and CSS-like styling with inheritance.

## Features

- **Vector-Based Rendering**: No external sprite dependencies - all UI elements drawn with vector graphics
- **CSS-Like Styling**: Familiar style properties (padding, margin, border, background, etc.)
- **Style Inheritance**: Styles cascade from parent to child elements automatically
- **Interactive Elements**: Buttons, inputs, checkboxes, lists with hover/focus/active states
- **Modal Dialogs**: Draggable modal windows with title bars
- **Flexible Layouts**: HBox, VBox containers with automatic child positioning
- **Background Images**: Optional image backgrounds via resource system

## Quick Start

```go
package main

import (
	"image/color"
	"github.com/hajimehoshi/ebiten/v2"
	minui "github.com/mechanical-lich/mlge/ui/min-ui"
)

func main() {
	// Create GUI manager
	gui := minui.NewGUI()
	
	// Create a button
	button := minui.NewButton("myButton", "Click Me")
	button.SetPosition(100, 100)
	button.OnClick = func() {
		println("Button clicked!")
	}
	
	gui.AddElement(button)
	
	// In your game loop:
	// gui.Update()
	// gui.Layout()
	// gui.Draw(screen)
}
```

## Core Concepts

### Elements

All UI components implement the `Element` interface:
- `Label` - Text display
- `Button` - Clickable button
- `TextInput` - Single-line text input
- `Checkbox` - Toggle checkbox
- `ListBox` - Scrollable list with selection
- `Panel` - Generic container
- `VBox` - Vertical layout container
- `HBox` - Horizontal layout container
- `Modal` - Draggable dialog window

### Styling System

Styles use CSS-like properties and support inheritance:

```go
button := minui.NewButton("btn", "Submit")

// Customize style
bgColor := color.Color(color.RGBA{100, 150, 255, 255})
borderWidth := 2
borderRadius := 8
padding := minui.NewEdgeInsets(12)

style := button.GetStyle()
style.BackgroundColor = &bgColor
style.BorderWidth = &borderWidth
style.BorderRadius = &borderRadius
style.Padding = padding

// Hover state
hoverBg := color.Color(color.RGBA{120, 170, 255, 255})
style.HoverStyle = &minui.Style{
	BackgroundColor: &hoverBg,
}
```

### Style Properties

#### Layout
- `Width`, `Height` - Element dimensions
- `MinWidth`, `MinHeight`, `MaxWidth`, `MaxHeight` - Size constraints
- `Padding` - Inner spacing (EdgeInsets)
- `Margin` - Outer spacing (EdgeInsets)

#### Visual
- `BackgroundColor` - Background color
- `BackgroundImage` - Background image resource ID
- `BorderColor` - Border color
- `BorderWidth` - Border thickness
- `BorderRadius` - Rounded corners
- `ForegroundColor` - Text/content color
- `Opacity` - Transparency (0.0 - 1.0)

#### Typography
- `FontSize` - Text size
- `FontBold` - Bold text
- `FontItalic` - Italic text
- `TextAlign` - Horizontal alignment (Left, Center, Right)
- `VertAlign` - Vertical alignment (Top, Middle, Bottom)

#### States
- `HoverStyle` - Applied when mouse hovers
- `ActiveStyle` - Applied when pressed/active
- `DisabledStyle` - Applied when disabled
- `FocusStyle` - Applied when focused

### Style Inheritance

Styles automatically inherit from parent elements:

```go
// Set panel style
panel := minui.NewPanel("panel")
fontSize := 16
panel.GetStyle().FontSize = &fontSize

// Children inherit font size
label := minui.NewLabel("label", "Text")
panel.AddChild(label) // Label will use 16pt font

// Override in child
childFontSize := 12
label.GetStyle().FontSize = &childFontSize // Now uses 12pt
```

### EdgeInsets Helper

Create spacing with convenience functions:

```go
// All sides equal
padding := minui.NewEdgeInsets(10) // 10px all around

// Vertical and horizontal
padding := minui.NewEdgeInsetsLR(8, 16) // 8px top/bottom, 16px left/right

// Individual sides
padding := minui.NewEdgeInsetsTRBL(4, 8, 12, 16) // top, right, bottom, left
```

## Layout Containers

### Panel

Generic container with manual or automatic layout:

```go
panel := minui.NewPanel("panel")
panel.SetBounds(minui.Rect{X: 0, Y: 0, Width: 400, Height: 300})

// Add children
panel.AddChild(label)
panel.AddChild(button)
```

### VBox (Vertical Box)

Arranges children vertically:

```go
vbox := minui.NewVBox("vbox")
vbox.AddChild(label1)
vbox.AddChild(label2)
vbox.AddChild(button)
// Children stacked top to bottom
```

### HBox (Horizontal Box)

Arranges children horizontally:

```go
hbox := minui.NewHBox("hbox")
hbox.AddChild(button1)
hbox.AddChild(button2)
// Children placed left to right
```

## Interactive Elements

### Button

```go
button := minui.NewButton("btn", "Click Me")
button.OnClick = func() {
	fmt.Println("Clicked!")
}
```

### TextInput

```go
input := minui.NewTextInput("nameInput", "Enter name...")
input.OnChange = func(text string) {
	fmt.Println("Text:", text)
}
input.OnSubmit = func(text string) {
	fmt.Println("Submitted:", text)
}
```

### Checkbox

```go
checkbox := minui.NewCheckbox("agree", "I agree")
checkbox.OnChange = func(checked bool) {
	fmt.Println("Checked:", checked)
}
```

### ListBox

```go
list := minui.NewListBox("items", []string{
	"Item 1",
	"Item 2",
	"Item 3",
})
list.OnSelect = func(index int, item string) {
	fmt.Println("Selected:", index, item)
}
```

## Modal Dialogs

```go
modal := minui.NewModal("dialog", "Confirm Action", 400, 200)
modal.SetPosition(120, 100)

// Add content
message := minui.NewLabel("msg", "Are you sure?")
message.SetBounds(minui.Rect{X: 20, Y: 20, Width: 360, Height: 30})

okBtn := minui.NewButton("ok", "OK")
okBtn.SetBounds(minui.Rect{X: 150, Y: 140, Width: 80, Height: 32})
okBtn.OnClick = func() {
	modal.SetVisible(false)
}

modal.AddChild(message)
modal.AddChild(okBtn)

// Add to GUI (modals appear above regular elements)
gui.AddModal(modal)
```

## Background Images

Use images from the resource system:

```go
// Load image first
resource.LoadImageAsTexture("panel_bg", "assets/panel_background.png")

// Apply to element
bgImage := "panel_bg"
panel.GetStyle().BackgroundImage = &bgImage
```

## Complete Example

See `examples/min-ui-demo/main.go` for a complete file browser demo application.

## Architecture

### File Structure

```
ui/min-ui/
├── style.go          # Style system and inheritance
├── element.go        # Element interface and base
├── rendering.go      # Vector drawing utilities
├── containers.go     # Panel, VBox, HBox
├── label.go          # Label element
├── button.go         # Button element
├── input.go          # TextInput and Checkbox
├── listbox.go        # ListBox element
├── modal.go          # Modal dialog
└── gui.go            # GUI manager
```

### Rendering Pipeline

1. **Update**: Handle input and state changes
2. **Layout**: Calculate element positions and sizes
3. **Draw**: Render elements with vector graphics

```go
func (g *Game) Update() error {
	gui.Update()  // Process input
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	gui.Layout()  // Calculate layout
	gui.Draw(screen)  // Render
}
```

## Design Principles

1. **No External Dependencies**: All rendering uses vector graphics
2. **Style Inheritance**: Parent styles cascade to children
3. **State Management**: Automatic hover/focus/active state handling
4. **Flexible Sizing**: Support for fixed, min/max, and auto sizing
5. **Event-Driven**: Callback-based interaction model

## Comparison with ui/v2

| Feature | ui/v2 | min-ui |
|---------|-------|--------|
| Rendering | Sprite-based | Vector-based |
| Dependencies | Requires sprite sheets | No external assets |
| Styling | Theme struct | CSS-like Style |
| Inheritance | No | Yes |
| Background Images | No | Yes (optional) |
| Border Radius | Fixed | Configurable |
| State Styles | Limited | Full (hover/focus/active/disabled) |

## Performance Notes

- Vector rendering is fast enough for typical UI needs
- For very complex UIs, consider caching rendered elements
- Modal overlays use semi-transparent overlays (may impact performance with many modals)

## Future Enhancements

- Dropdown menus
- Tree views
- Tab containers
- Tooltip system
- Drag and drop
- Animation/transitions
- Grid layout container
- Scroll containers
- Color picker
- Slider control

## License

Part of the MLGE game engine project.
