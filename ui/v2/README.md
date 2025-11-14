# MLGE UI v2 - Improvements & Features

This document outlines the major improvements made to the MLGE UI v2 package.

## New Features

### 1. UI Event Types (Integrates with MLGE Event System)

UI v2 now integrates with the existing MLGE event system (`github.com/mechanical-lich/mlge/event`):

```go
// Create event manager
eventMgr := &event.EventManager{}

// Create a listener
type MyUIListener struct{}

func (l *MyUIListener) HandleEvent(data event.EventData) error {
    switch e := data.(type) {
    case ui.ClickEventData:
        log.Printf("Button clicked: %s", e.SourceName)
    case ui.ValueChangeEventData:
        log.Printf("Value changed: %v", e.NewValue)
    }
    return nil
}

// Register listener
listener := &MyUIListener{}
eventMgr.RegisterListener(listener, ui.EventTypeUIClick)
eventMgr.RegisterListener(listener, ui.EventTypeUIValueChange)

// Send events from UI components
eventMgr.SendEvent(ui.ClickEventData{
    SourceName: "myButton",
    Data: map[string]interface{}{"action": "submit"},
})
```

**Supported Event Types:**
- `EventTypeUIClick` - Button/element clicks
- `EventTypeUIValueChange` - Input value changes
- `EventTypeUIFocus` - Element gained focus
- `EventTypeUIBlur` - Element lost focus
- `EventTypeUIModalOpen` - Modal opened
- `EventTypeUIModalClose` - Modal closed
- `EventTypeUIHover` - Mouse hover
- `EventTypeUISubmit` - Form submission
- `EventTypeUITabChange` - Tab changed

**Event Data Types:**
- `ClickEventData` - Contains SourceName and Data map
- `ValueChangeEventData` - Contains SourceName, OldValue, NewValue
- `FocusEventData` - Contains SourceName
- `BlurEventData` - Contains SourceName
- `ModalOpenEventData` - Contains ModalName
- `ModalCloseEventData` - Contains ModalName
- `SubmitEventData` - Contains SourceName and Data map

### 2. Validation Framework (`ui/v2/validation`)

Powerful input validation with built-in validators:

```go
// Single validator
nameInput.Validator = validation.Required("Name")

// Combined validators
emailInput.Validator = validation.Combine(
    validation.Required("Email"),
    validation.Email("Email"),
)

// Custom validator
ageInput.Validator = validation.IntegerRange("Age", 18, 120)
```

**Built-in Validators:**
- `Required(fieldName)` - Non-empty value
- `MinLength(fieldName, length)` - Minimum character count
- `MaxLength(fieldName, length)` - Maximum character count
- `Email(fieldName)` - Valid email format
- `Numeric(fieldName)` - Numeric value
- `IntegerRange(fieldName, min, max)` - Integer within range
- `Pattern(fieldName, regex, message)` - Custom regex pattern
- `Custom(fieldName, func, message)` - Custom validation function

### 3. Enhanced Theming System (`ui/v2/theming`)

Color-based themes with predefined color schemes:

```go
// Use predefined themes
gui := ui.NewGUI(view, &theming.DefaultTheme)

// Access theme colors
theme.Colors.Primary      // Primary color
theme.Colors.Error        // Error state color
theme.Colors.Focus        // Focus state color
theme.Colors.Background   // Background color
```

**Color Schemes:**
- `DefaultColors` - Balanced dark theme
- `DarkColors` - Deep dark theme
- `LightColors` - Light theme

**Theme Properties:**
- Sprite coordinates for UI elements
- Color scheme for states and feedback
- Support for custom themes

### 4. New UI Elements (`ui/v2/elements`)

#### Checkbox
```go
checkbox := elements.NewCheckbox("agree", 20, 50, "I agree", false)
checkbox.OnChanged = func(checked bool) {
    log.Printf("Checkbox changed: %v", checked)
}
```

#### Slider
```go
slider := elements.NewSlider("volume", 20, 100, 200, "Volume", 0, 100, 50, 5)
slider.OnChanged = func(value float64) {
    log.Printf("Slider value: %.2f", value)
}
```

Features:
- Draggable thumb
- Min/max values
- Step increments
- Tick marks
- Real-time value display

#### Dropdown
```go
options := []string{"Option 1", "Option 2", "Option 3"}
dropdown := elements.NewDropdown("theme", 20, 150, 200, options, 0)
dropdown.OnChanged = func(index int, value string) {
    log.Printf("Selected: %s", value)
}
```

Features:
- Click to open/close
- Hover highlighting
- Scrollable options list

### 5. Enhanced InputField

**New Properties:**
- `Placeholder` - Hint text when empty
- `IsPassword` - Hide characters with bullets
- `Validator` - Validation rules
- `ValidationError` - Current validation state
- `OnSubmit` - Enter key handler

**New Keyboard Shortcuts:**
- `Ctrl+A` - Select all
- `Ctrl+C/V/X` - Copy/paste/cut (interface ready)
- `Home/End` - Jump to start/end
- `Enter` - Submit
- `Escape` - Unfocus

**Visual States:**
- Focus border (blue)
- Error border (red)
- Validation error message
- Password masking

## Usage Examples

### Form with Validation

```go
// Create input with validation
nameInput := elements.NewInputField("name", 20, 50, 200, 50)
nameInput.Placeholder = "Enter your name"
nameInput.Validator = validation.Combine(
    validation.Required("Name"),
    validation.MinLength("Name", 2),
)

// Password input
passwordInput := elements.NewInputField("password", 20, 100, 200, 50)
passwordInput.Placeholder = "Password"
passwordInput.IsPassword = true
passwordInput.Validator = validation.MinLength("Password", 6)

// Submit handler
submitBtn := elements.NewButton("submit", 20, 150, "Submit", "")
submitBtn.OnClicked = func() {
    // Validate all fields
    if nameInput.Validator != nil {
        if err := nameInput.Validator(string(nameInput.Value)); err != nil {
            log.Printf("Validation error: %v", err)
            return
        }
    }
    // Process form...
}
```

### Event-Driven UI

```go
// Create event manager
eventMgr := &event.EventManager{}

// Create listener
type FormListener struct{}

func (l *FormListener) HandleEvent(data event.EventData) error {
    switch e := data.(type) {
    case ui.ClickEventData:
        if e.SourceName == "saveButton" {
            saveData()
        } else if e.SourceName == "loadButton" {
            loadData()
        }
    }
    return nil
}

// Register listener
listener := &FormListener{}
eventMgr.RegisterListener(listener, ui.EventTypeUIClick)

// Create button that emits events
saveBtn := elements.NewButton("saveButton", 20, 20, "Save", "")
saveBtn.OnClicked = func() {
    eventMgr.SendEvent(ui.ClickEventData{
        SourceName: "saveButton",
        Data: map[string]interface{}{
            "timestamp": time.Now(),
        },
    })
}
```

### Theme Switching

```go
// Define themes
themes := []*theming.Theme{
    &theming.DefaultTheme,
    {
        Name: "Dark",
        Colors: theming.DarkColors,
        // ... sprite coordinates
    },
    {
        Name: "Light",
        Colors: theming.LightColors,
        // ... sprite coordinates
    },
}

// Switch theme
currentTheme := 0
themeBtn.OnClicked = func() {
    currentTheme = (currentTheme + 1) % len(themes)
    gui.Theme = themes[currentTheme]
}
```

## Migration from v1

Key differences from UI v1:

1. **No State Dependency** - v2 elements don't require `state.StateInterface`
2. **Parent-Child Relationships** - Elements track their parents via `SetParent()`
3. **Absolute Positioning** - Use `GetAbsolutePosition()` instead of parent offsets
4. **Event System** - Use event bus instead of direct callbacks
5. **Validation** - Built-in validation framework
6. **Theming** - Color-based themes with runtime switching

## File Structure

```
ui/v2/
├── gui.go                    # Main GUI manager
├── uievents.go               # UI event types (integrates with mlge/event)
├── validation/
│   └── validation.go         # Validation framework
├── theming/
│   └── theme.go              # Theme and color definitions
├── elements/
│   ├── button.go
│   ├── checkbox.go           # NEW
│   ├── slider.go             # NEW
│   ├── dropdown.go           # NEW
│   ├── input.go              # ENHANCED
│   ├── label.go
│   ├── radio.go
│   ├── select.go
│   ├── toggle.go
│   ├── scrollingTextArea.go
│   └── elementBase.go
├── containers/
│   ├── modal.go
│   ├── gridContainer.go
│   ├── absolutePositionContainer.go
│   ├── tabbedContainer.go
│   └── drawerModal.go
└── views/
    ├── guiview.go
    └── basicview.go
```

## Future Improvements

Planned enhancements for future versions:

- FlexContainer for flexible layouts
- Keyboard focus management with tab order
- Clipboard support (copy/paste)
- Text selection in inputs
- Drag & drop support
- Context menus
- TreeView component
- ColorPicker component
- DatePicker component
- Animation/transition system
- Accessibility features
- Performance optimizations (dirty rectangles, caching)

## Examples

See `examples/ui-demo/` for a comprehensive demonstration of all features.
