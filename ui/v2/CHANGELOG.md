# UI v2 Changelog

## Version 2.1.0 - November 2025

### 🎉 New Features

#### Event Integration
- Added `ui/v2/uievents.go` with UI event types for MLGE's existing event system
- Event types implement `event.EventData` interface
- Support for Click, ValueChange, Focus, Blur, ModalOpen, ModalClose, and Submit events
- Uses MLGE's `event.EventManager` for consistent event handling
- Compatible with `event.EventListener` interface

#### Validation Framework  
- Added `ui/v2/validation` package with reusable validators
- Built-in validators: Required, MinLength, MaxLength, Email, Numeric, IntegerRange, Pattern, Custom
- Validator composition with `Combine()`
- ValidationError type with field and message

#### Enhanced Theming
- Added `Colors` struct to theme with semantic color properties
- Three predefined color schemes: Default, Dark, Light
- Theme colors for: Primary, Secondary, Background, Surface, Text, Border, Focus, Error, Success, Warning, Disabled
- Runtime theme switching capability
- Backward compatible with existing sprite-based themes

#### New UI Elements

**Checkbox** (`ui/v2/elements/checkbox.go`)
- Boolean toggle with label
- Visual checked/unchecked states
- `OnChanged` callback
- Theme-aware rendering

**Slider** (`ui/v2/elements/slider.go`)
- Horizontal numeric range input
- Min/max values with step increments
- Draggable thumb with visual feedback
- Real-time value display
- Optional tick marks
- `OnChanged` callback

**Dropdown** (`ui/v2/elements/dropdown.go`)
- Select from list of options
- Click to open/close
- Hover highlighting
- Visual arrow indicator
- `OnChanged` callback with index and value
- Supports many options with scrolling

### 🔧 Enhancements

#### InputField (`ui/v2/elements/input.go`)
- **Placeholder text** - Shows hint when empty
- **Password mode** - `IsPassword` property hides characters with bullets
- **Validation integration** - `Validator` property with visual error feedback
- **Visual states** - Focus border (blue), error border (red)
- **Error messages** - Displayed below field
- **Enhanced keyboard support**:
  - Ctrl+A - Select all (interface ready)
  - Ctrl+C/V/X - Copy/paste/cut (interface ready)  
  - Home/End - Jump to start/end
  - Enter - OnSubmit callback
  - Escape - Unfocus
- **Validation on blur** - Automatically validates when focus is lost
- **Error clearing** - Clears validation error when value changes

### 📚 Documentation

- Added `ui/v2/README.md` - Comprehensive feature documentation
- Added `ui/v2/IMPLEMENTATION_SUMMARY.md` - Implementation details
- Added `UI_V2_QUICK_START.md` - Quick start guide
- Added `examples/ui-demo/README.md` - Example usage guide

### 🎮 Examples

- Added `examples/ui-demo/main.go` - Full-featured demo application showing:
  - Form with validation
  - All new UI elements
  - Event integration with `event.EventManager`
  - `UIEventLogger` implementing `event.EventListener`
  - Theme application
  - Modal dialogs
  - Real-time event logging to console

### 📝 Theme Updates

#### New Theme Properties
- `Name` - Theme name
- `Colors` - Color scheme struct
- `Checkbox` - Checkbox sprite coordinates
- `Slider` - Slider track and thumb coordinates
- `Dropdown` - Dropdown sprite coordinates

#### Color Properties
- Primary, Secondary, Background, Surface
- Text, TextSecondary
- Border, Focus
- Error, Success, Warning
- Disabled

### 🐛 Bug Fixes
- None (new features only)

### ⚠️ Breaking Changes
- None - All changes are additive and backward compatible

### 🔮 Deprecated
- None

### 🚧 Known Limitations

- Clipboard operations (Ctrl+C/V/X) have interface but no platform-specific implementation
- Text selection in InputField not yet implemented
- No keyboard tab navigation between elements yet

### 📋 Migration Guide

No migration needed. All v2 code continues to work. To use new features:

1. Import packages as needed:
   ```go
   import "github.com/mechanical-lich/mlge/event"
   import "github.com/mechanical-lich/mlge/ui/v2"
   import "github.com/mechanical-lich/mlge/ui/v2/validation"
   ```

2. Create event manager and register listeners:
   ```go
   eventMgr := event.NewEventManager()
   eventMgr.RegisterListener(myListener, ui.EventTypeUIClick)
   ```

3. Use new elements:
   ```go
   checkbox := elements.NewCheckbox(...)
   slider := elements.NewSlider(...)
   dropdown := elements.NewDropdown(...)
   ```

4. Add validation to inputs:
   ```go
   input.Validator = validation.Required("Name")
   ```

### 🎯 What's Next

Planned for future versions:
- FlexContainer for flexible layouts
- Full keyboard navigation with tab order
- Actual clipboard support (platform-specific)
- Text selection with mouse in inputs
- TreeView, ColorPicker, DatePicker components
- Animation/transition system
- Performance optimizations

---

## Version 2.0.0 - Previous

- Initial v2 release
- Removed state dependency
- Added parent-child relationships
- Reorganized into packages
- Enhanced containers and modals

See git history for complete v2.0.0 changes.
