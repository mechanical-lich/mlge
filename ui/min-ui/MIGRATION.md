# Migrating from ui/v2 to min-ui

## Quick Comparison

### ui/v2 (Sprite-Based)
```go
import ui "github.com/mechanical-lich/mlge/ui/v2"

button := ui.NewButton("btn", 100, 100, 120, 32, "Click", theme)
button.OnClick = func(b *ui.Button) {
    println("clicked")
}
```

### min-ui (Vector-Based)
```go
import minui "github.com/mechanical-lich/mlge/ui/min-ui"

button := minui.NewButton("btn", "Click")
button.SetBounds(minui.Rect{X: 100, Y: 100, Width: 120, Height: 32})
button.OnClick = func() {
    println("clicked")
}
```

## Key Differences

### 1. No Theme Parameter

**ui/v2**: Pass theme to element constructors
```go
button := ui.NewButton("btn", x, y, w, h, "Text", theme)
```

**min-ui**: Style is part of the element
```go
button := minui.NewButton("btn", "Text")
button.SetPosition(x, y)
button.SetSize(w, h)
```

### 2. Styling Approach

**ui/v2**: Global theme with sprite coordinates
```go
theme.Button.SrcX = 16
theme.Button.SrcY = 32
```

**min-ui**: Per-element CSS-like styles
```go
bgColor := color.RGBA{100, 150, 255, 255}
button.GetStyle().BackgroundColor = &bgColor
button.GetStyle().BorderRadius = &8
```

### 3. GUI Management

**ui/v2**: Often manual element tracking
```go
elements := []ui.Element{button, label, input}
for _, el := range elements {
    el.Update()
}
```

**min-ui**: Centralized GUI manager
```go
gui := minui.NewGUI()
gui.AddElement(button)
gui.AddElement(label)
gui.Update() // Updates all
```

### 4. Containers

**ui/v2**: Container with manual child positioning
```go
container := ui.NewContainer(x, y, w, h)
container.AddChild(button) // Button position is absolute
```

**min-ui**: Auto-layout containers
```go
vbox := minui.NewVBox("vbox")
vbox.AddChild(button1) // Automatically stacked
vbox.AddChild(button2) // Below button1
```

## Migration Steps

### Step 1: Replace Imports

```go
// OLD
import ui "github.com/mechanical-lich/mlge/ui/v2"
import theming "github.com/mechanical-lich/mlge/ui/v2/theming"

// NEW
import minui "github.com/mechanical-lich/mlge/ui/min-ui"
```

### Step 2: Remove Theme Loading

```go
// OLD - Delete this
theme := theming.LoadTheme("assets/ui.png")

// NEW - No theme needed!
```

### Step 3: Convert Elements

#### Button
```go
// OLD
btn := ui.NewButton("btn", 100, 100, 120, 32, "Click", theme)

// NEW
btn := minui.NewButton("btn", "Click")
btn.SetBounds(minui.Rect{X: 100, Y: 100, Width: 120, Height: 32})
```

#### Label
```go
// OLD
lbl := ui.NewLabel("lbl", 50, 50, "Text")
lbl.Draw(screen, theme)

// NEW
lbl := minui.NewLabel("lbl", "Text")
lbl.SetPosition(50, 50)
lbl.Draw(screen) // No theme needed
```

#### Input Field
```go
// OLD
input := ui.NewInputField("input", 20, 80, 200, 30, theme)
input.Placeholder = "Enter text"

// NEW
input := minui.NewTextInput("input", "Enter text")
input.SetBounds(minui.Rect{X: 20, Y: 80, Width: 200, Height: 30})
```

### Step 4: Create GUI Manager

```go
// NEW - Add to your game struct
type Game struct {
    gui *minui.GUI
}

func NewGame() *Game {
    g := &Game{
        gui: minui.NewGUI(),
    }
    
    // Add elements
    button := minui.NewButton("btn", "Click")
    g.gui.AddElement(button)
    
    return g
}
```

### Step 5: Update Game Loop

```go
// OLD
func (g *Game) Update() error {
    for _, element := range g.elements {
        element.Update()
    }
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    for _, element := range g.elements {
        element.Draw(screen, g.theme)
    }
}

// NEW
func (g *Game) Update() error {
    g.gui.Update()
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    g.gui.Layout()
    g.gui.Draw(screen)
}
```

## Element Mapping

| ui/v2 | min-ui | Notes |
|-------|--------|-------|
| Button | Button | Callbacks changed |
| Label | Label | No theme param |
| InputField | TextInput | Renamed |
| Checkbox | Checkbox | Vector-rendered |
| RadioButton | *(not yet)* | Use Checkbox for now |
| Toggle | *(not yet)* | Use Checkbox |
| Slider | *(not yet)* | Coming soon |
| Dropdown | ListBox | Similar functionality |
| Container | Panel | Now with layouts |
| Modal | Modal | Draggable, styled |

## Styling Migration

### ui/v2 Theme Colors
```go
theme.Colors.Primary = color.RGBA{100, 150, 255, 255}
```

### min-ui Element Styles
```go
primaryColor := color.Color(color.RGBA{100, 150, 255, 255})
button.GetStyle().BackgroundColor = &primaryColor
```

### Creating Reusable Styles

```go
// Define common style
func ButtonStyle() *minui.Style {
    bgColor := color.Color(color.RGBA{100, 150, 255, 255})
    borderWidth := 2
    borderRadius := 6
    padding := minui.NewEdgeInsets(10)
    
    return &minui.Style{
        BackgroundColor: &bgColor,
        BorderWidth:     &borderWidth,
        BorderRadius:    &borderRadius,
        Padding:         padding,
    }
}

// Apply to buttons
button1.SetStyle(ButtonStyle())
button2.SetStyle(ButtonStyle())
```

## Common Patterns

### Creating a Form

**ui/v2**
```go
container := ui.NewContainer(x, y, w, h)
label := ui.NewLabel("name", 10, 10, "Name:")
input := ui.NewInputField("nameInput", 10, 35, 180, 28, theme)
container.AddChild(label)
container.AddChild(input)
```

**min-ui**
```go
form := minui.NewVBox("form")
label := minui.NewLabel("name", "Name:")
input := minui.NewTextInput("nameInput", "Enter name")
form.AddChild(label)
form.AddChild(input)
// Labels and inputs automatically stacked vertically
```

### Modal Dialog

**ui/v2**
```go
modal := ui.NewModal("dialog", x, y, w, h, "Title", theme)
modal.AddContent(message)
modal.AddButton(okButton)
```

**min-ui**
```go
modal := minui.NewModal("dialog", "Title", 400, 200)
modal.SetPosition(100, 100)
modal.AddChild(message)
modal.AddChild(okButton)
gui.AddModal(modal)
```

## Benefits of min-ui

1. **No Asset Files**: Remove all UI sprite sheets
2. **Runtime Styling**: Change colors without rebuilding assets
3. **CSS Familiarity**: If you know CSS, you know min-ui styling
4. **Auto Layout**: VBox/HBox save manual positioning work
5. **Style Inheritance**: Set colors once, children inherit
6. **State Styling**: Hover/focus/active styles built-in

## Potential Issues

### Issue 1: Missing Elements

Some ui/v2 elements aren't in min-ui yet:
- Use Button for RadioButton/Toggle for now
- ListBox can replace Dropdown temporarily

### Issue 2: Performance

Vector rendering is slightly slower than sprites:
- Usually not noticeable (<10%)
- Avoid 100+ elements on screen simultaneously
- Consider caching complex UIs

### Issue 3: Font Rendering

min-ui uses MLGE's text system:
- Different text measurement
- May need position adjustments
- Use Layout() to auto-calculate sizes

## Example: Complete Migration

### Before (ui/v2)
```go
type Game struct {
    theme   *theming.Theme
    button  *ui.Button
    label   *ui.Label
}

func NewGame() *Game {
    theme := theming.DefaultTheme
    return &Game{
        theme:  theme,
        button: ui.NewButton("btn", 100, 100, 120, 32, "Click", theme),
        label:  ui.NewLabel("lbl", 100, 50, "Hello"),
    }
}

func (g *Game) Update() error {
    g.button.Update()
    if g.button.OnClick != nil && g.button.IsJustClicked() {
        g.label.Text = "Clicked!"
    }
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    g.button.Draw(screen, g.theme)
    g.label.Draw(screen, g.theme)
}
```

### After (min-ui)
```go
type Game struct {
    gui *minui.GUI
}

func NewGame() *Game {
    gui := minui.NewGUI()
    
    button := minui.NewButton("btn", "Click")
    button.SetBounds(minui.Rect{X: 100, Y: 100, Width: 120, Height: 32})
    button.OnClick = func() {
        label := gui.FindElementByID("lbl").(*minui.Label)
        label.Text = "Clicked!"
    }
    
    label := minui.NewLabel("lbl", "Hello")
    label.SetPosition(100, 50)
    
    gui.AddElement(button)
    gui.AddElement(label)
    
    return &Game{gui: gui}
}

func (g *Game) Update() error {
    g.gui.Update()
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    g.gui.Layout()
    g.gui.Draw(screen)
}
```

## Need Help?

- Check `ui/min-ui/README.md` for full documentation
- See `examples/min-ui-demo/main.go` for complete example
- ui/v2 and min-ui can coexist - migrate incrementally!

## When to Use Each

**Use ui/v2 when:**
- You have existing sprite-based UI assets
- You need exact pixel-perfect reproduction
- Performance is absolutely critical
- You're familiar with the existing system

**Use min-ui when:**
- Starting a new project
- Want easy customization
- Don't have UI sprites
- Prefer CSS-like styling
- Need responsive layouts
