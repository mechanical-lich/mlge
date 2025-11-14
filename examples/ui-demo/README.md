# UI v2 Feature Demo

This example demonstrates all the new features added to the MLGE UI v2 package.

## Features Demonstrated

### 1. **Event System Integration**
- Uses MLGE's built-in event system (`github.com/mechanical-lich/mlge/event`)
- UI event types defined in `ui/v2/uievents.go`
- EventListener implementation for UI events
- Event emission from UI components
- Click, value change, focus, and blur events

### 2. **Validation Framework**
- Input field validation with visual feedback
- Built-in validators (Required, MinLength, Email, etc.)
- Custom validation rules
- Error messages displayed below inputs

### 3. **Enhanced Theming**
- Color-based themes (Default, Dark, Light)
- Customizable color schemes
- Focus and error state colors
- Validation error styling

### 4. **New UI Elements**

#### Checkbox
- Boolean toggle with label
- OnChanged callback
- Visual checked/unchecked states

#### Slider
- Numeric range input with draggable thumb
- Min/max values with optional step increments
- Real-time value display
- Tick marks for stepped values
- OnChanged callback

#### Dropdown
- Select from list of options
- Opens/closes on click
- Hover highlighting
- OnChanged callback with selected value

### 5. **Enhanced Input Field**
- **Placeholder text** - Shows hint when empty
- **Password mode** - Hides characters with bullets (●)
- **Validation** - Real-time validation with error feedback
- **Keyboard shortcuts**:
  - Ctrl+A: Select all
  - Ctrl+C/V/X: Copy/paste/cut (interface ready)
  - Home/End: Jump to start/end
  - Enter: Submit with OnSubmit callback
  - Escape: Unfocus
- **Visual states** - Focus border and error highlighting

### 6. **Modal Dialogs**
- Draggable modal windows
- Child element support
- OnClose callback
- Visual overlay

## Running the Example

```bash
cd examples/ui-demo
go run main.go
```

## Controls

- **Click** inputs to focus and type
- **Click** buttons to trigger actions
- **Drag** sliders to adjust values
- **Click** checkboxes to toggle
- **Click** dropdown to open/close options
- **Click** "Open Modal" to show a modal dialog
- **Drag** modal titlebar to reposition

## Form Elements

1. **Name Input** - Required, minimum 2 characters
2. **Email Input** - Required, valid email format
3. **Password Input** - Minimum 6 characters, hidden text
4. **Age Slider** - 0-100 range with 1-unit steps
5. **Volume Slider** - 0-100 range with 5-unit steps
6. **Agreement Checkbox** - Terms acceptance
7. **Notifications Checkbox** - Email notifications
8. **Theme Dropdown** - Select theme variant
9. **Submit Button** - Submit form (validates all fields)
10. **Clear Button** - Reset all fields to defaults
11. **Open Modal** - Display modal dialog

## Event Logging

All UI events are logged to the console, showing:
- Button clicks
- Value changes
- Input focus/blur
- Modal open/close

Check the terminal for real-time event notifications!

## Assets

The example requires `assets/ux.png` - a sprite sheet containing UI element graphics.
This should be placed in the mlge root directory.
