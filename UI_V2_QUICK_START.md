# UI v2 Improvements - Quick Start Guide

## What Was Added

### 🎯 High-Priority Features Implemented

1. **Event System** - Decoupled component communication
2. **Validation Framework** - Reusable input validators
3. **Color Theming** - Three color schemes (Default, Dark, Light)
4. **New Elements** - Checkbox, Slider, Dropdown
5. **Enhanced Input** - Placeholder, password mode, validation
6. **Example Application** - Full demo of all features

## Quick Usage

### Event System (Uses MLGE's Event System)

```go
import "github.com/mechanical-lich/mlge/event"
import ui "github.com/mechanical-lich/mlge/ui/v2"

// Create event manager
eventMgr := &event.EventManager{}

// Create listener
type MyListener struct{}

func (l *MyListener) HandleEvent(data event.EventData) error {
    switch e := data.(type) {
    case ui.ClickEventData:
        log.Printf("Clicked: %s", e.SourceName)
    case ui.ValueChangeEventData:
        log.Printf("Value changed to: %v", e.NewValue)
    }
    return nil
}

// Register listener
listener := &MyListener{}
eventMgr.RegisterListener(listener, ui.EventTypeUIClick)
eventMgr.RegisterListener(listener, ui.EventTypeUIValueChange)

// Send events from UI components
button.OnClicked = func() {
    eventMgr.SendEvent(ui.ClickEventData{
        SourceName: "myButton",
        Data: nil,
    })
}
```

### Validation

```go
// Add validation to inputs
input.Validator = validation.Combine(
    validation.Required("Email"),
    validation.Email("Email"),
)
// Visual error shown automatically on blur
```

### Checkbox

```go
checkbox := elements.NewCheckbox("agree", 20, 50, "I agree", false)
checkbox.OnChanged = func(checked bool) {
    // Handle change
}
```

### Slider

```go
slider := elements.NewSlider("volume", 20, 100, 200, "Volume", 0, 100, 50, 5)
slider.OnChanged = func(value float64) {
    // Handle change
}
```

### Dropdown

```go
dropdown := elements.NewDropdown("theme", 20, 150, 200, 
    []string{"Light", "Dark", "Default"}, 0)
dropdown.OnChanged = func(index int, value string) {
    // Handle selection
}
```

### Enhanced Input

```go
input := elements.NewInputField("password", 20, 50, 200, 50)
input.Placeholder = "Enter password"
input.IsPassword = true
input.Validator = validation.MinLength("Password", 6)
input.OnSubmit = func(value string) {
    // Handle Enter key
}
```

## Running the Example

```bash
cd examples/ui-demo
go run main.go
```

The example demonstrates:
- All new UI elements
- Validation with visual feedback
- Event system integration
- Theme colors in action
- Modal dialogs
- Form submission

## File Structure

```
ui/v2/
├── uievents.go               # UI event types (uses mlge/event)
├── validation/validation.go  # Validation framework
├── theming/theme.go          # Enhanced with colors
├── elements/
│   ├── checkbox.go           # NEW
│   ├── slider.go             # NEW
│   ├── dropdown.go           # NEW
│   └── input.go              # ENHANCED
├── README.md                 # Full documentation
└── IMPLEMENTATION_SUMMARY.md # What was built

examples/ui-demo/
├── main.go                   # Working example
└── README.md                 # Usage guide
```

## Next Steps

1. **Test the example** - Run `examples/ui-demo/main.go`
2. **Read the docs** - Check `ui/v2/README.md` for details
3. **Use in your project** - Import and use the new components
4. **Customize themes** - Create your own color schemes

## Assets Required

Make sure `assets/ux.png` exists in the mlge root directory. This sprite sheet contains the UI element graphics referenced in the default theme.

## What's Not Included (Future Work)

- FlexContainer (flexible layouts)
- Full keyboard navigation (tab order)
- Actual clipboard support (platform-specific)
- Text selection with mouse
- Additional elements (TreeView, ColorPicker, etc.)
- Animation system
- Performance optimizations

These can be added incrementally as needed.

## Documentation

- **Full Feature Docs**: `ui/v2/README.md`
- **Example Guide**: `examples/ui-demo/README.md`
- **Implementation Details**: `ui/v2/IMPLEMENTATION_SUMMARY.md`

## Support

All new components follow the existing v2 patterns:
- Use `GetAbsolutePosition()` for rendering
- Use `SetParent()` for hierarchy
- No `state.StateInterface` dependency
- Theme-aware coloring
- Event emission where applicable

Enjoy the new UI features! 🎉
