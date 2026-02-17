package minui

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mechanical-lich/mlge/event"
	"github.com/mechanical-lich/mlge/text"
)

// TextInput is a single-line text input field
type TextInput struct {
	*ElementBase
	Text        string
	Placeholder string
	OnChange    func(text string)
	OnSubmit    func(text string)
	cursorPos   int
	focused     bool
}

// NewTextInput creates a new text input
func NewTextInput(id, placeholder string) *TextInput {
	ti := &TextInput{
		ElementBase: NewElementBase(id),
		Placeholder: placeholder,
		Text:        "",
		cursorPos:   0,
	}

	// Set default size
	ti.SetSize(200, 28)

	// Set default style - only structural properties, colors come from theme
	borderWidth := 1
	padding := NewEdgeInsetsLR(4, 8)

	ti.style.BorderWidth = &borderWidth
	ti.style.Padding = padding

	return ti
}

// GetType returns the element type
func (ti *TextInput) GetType() string {
	return "TextInput"
}

// Update updates the text input
func (ti *TextInput) Update() {
	if !ti.visible || !ti.enabled {
		return
	}

	ti.UpdateHoverState()

	// Handle focus
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		wasFocused := ti.focused
		ti.focused = ti.IsWithin(mx, my)

		if ti.focused != wasFocused {
			ti.SetFocused(ti.focused)
		}
	}

	if !ti.focused {
		return
	}

	oldText := ti.Text
	textChanged := false

	// Handle text input
	runes := ebiten.AppendInputChars(nil)
	for _, r := range runes {
		if r >= 32 && r != 127 { // Printable characters
			ti.Text = ti.Text[:ti.cursorPos] + string(r) + ti.Text[ti.cursorPos:]
			ti.cursorPos++
			textChanged = true
		}
	}

	// Handle backspace
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && ti.cursorPos > 0 {
		ti.Text = ti.Text[:ti.cursorPos-1] + ti.Text[ti.cursorPos:]
		ti.cursorPos--
		textChanged = true
	}

	// Handle delete
	if inpututil.IsKeyJustPressed(ebiten.KeyDelete) && ti.cursorPos < len(ti.Text) {
		ti.Text = ti.Text[:ti.cursorPos] + ti.Text[ti.cursorPos+1:]
		textChanged = true
	}

	// Fire change event if text changed
	if textChanged {
		if ti.OnChange != nil {
			ti.OnChange(ti.Text)
		}
		event.GetQueuedInstance().QueueEvent(TextInputChangeEvent{
			InputID: ti.GetID(),
			Input:   ti,
			Text:    ti.Text,
			OldText: oldText,
		})
	}

	// Handle left/right arrows
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && ti.cursorPos > 0 {
		ti.cursorPos--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && ti.cursorPos < len(ti.Text) {
		ti.cursorPos++
	}

	// Handle Home/End
	if inpututil.IsKeyJustPressed(ebiten.KeyHome) {
		ti.cursorPos = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnd) {
		ti.cursorPos = len(ti.Text)
	}

	// Handle Enter
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if ti.OnSubmit != nil {
			ti.OnSubmit(ti.Text)
		}
		// Fire submit event
		event.GetQueuedInstance().QueueEvent(TextInputSubmitEvent{
			InputID: ti.GetID(),
			Input:   ti,
			Text:    ti.Text,
		})
		ti.focused = false
		ti.SetFocused(false)
	}

	// Handle Escape
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		ti.focused = false
		ti.SetFocused(false)
	}
}

// Layout calculates dimensions
func (ti *TextInput) Layout() {
	style := ti.GetComputedStyle()

	// Start with current bounds
	width := ti.bounds.Width
	height := ti.bounds.Height

	// Apply width/height from style
	if style.Width != nil {
		width = *style.Width
	}
	if style.Height != nil {
		height = *style.Height
	}

	// Apply min/max size constraints
	width, height = ApplySizeConstraints(width, height, style)

	ti.bounds.Width = width
	ti.bounds.Height = height
}

// Draw draws the text input
func (ti *TextInput) Draw(screen *ebiten.Image) {
	if !ti.visible {
		return
	}

	style := ti.GetComputedStyle()
	theme := ti.GetTheme()

	// Get absolute position for drawing
	absX, absY := ti.GetAbsolutePosition()
	absBounds := Rect{
		X:      absX,
		Y:      absY,
		Width:  ti.bounds.Width,
		Height: ti.bounds.Height,
	}

	// Draw background with theme support
	DrawBackgroundWithTheme(screen, absBounds, style, theme)

	// Draw text or placeholder
	contentBounds := GetContentBounds(absBounds, style)

	displayText := ti.Text
	// Get text color from style, then theme, then default
	textColor := color.RGBA{255, 255, 255, 255}
	if theme != nil {
		textColor = colorToRGBA(theme.Colors.Text)
	}
	if style.ForegroundColor != nil {
		textColor = colorToRGBA(*style.ForegroundColor)
	}

	if displayText == "" && ti.Placeholder != "" {
		displayText = ti.Placeholder
		// Use disabled/secondary color for placeholder
		textColor = color.RGBA{150, 150, 150, 255}
		if theme != nil {
			textColor = colorToRGBA(theme.Colors.TextSecondary)
		}
	}

	// Determine font size (allow style override)
	fontSize := 14.0
	if style != nil && style.FontSize != nil {
		fontSize = float64(*style.FontSize)
	}

	// Measure text height for vertical centering
	_, textH := text.Measure("M", fontSize)
	textY := contentBounds.Y + int(math.Floor((float64(contentBounds.Height)-textH)/2.0))

	// Draw the text starting at the content X and vertically centered
	text.Draw(screen, displayText, fontSize, contentBounds.X, textY, textColor)

	// Draw cursor if focused
	if ti.focused {
		// Compute cursor X by measuring text up to cursorPos
		before := ""
		if ti.cursorPos > 0 && ti.cursorPos <= len(ti.Text) {
			before = ti.Text[:ti.cursorPos]
		}
		tw, _ := text.Measure(before, fontSize)
		cursorX := contentBounds.X + int(math.Ceil(tw))

		// Cursor vertically matches the text height
		cursorBounds := Rect{
			X:      cursorX,
			Y:      textY,
			Width:  2,
			Height: int(math.Ceil(textH)),
		}
		// Get cursor color from theme or default
		cursorColor := color.RGBA{255, 255, 255, 255}
		if theme != nil {
			cursorColor = colorToRGBA(theme.Colors.Text)
		}

		// Blink cursor
		if (ebiten.TPS()*2/3)%2 == 0 {
			DrawRect(screen, cursorBounds, cursorColor)
		}
	}

	// Draw border with theme support (use focus color when focused)
	if ti.focused && theme != nil {
		focusBorderColor := color.Color(colorToRGBA(theme.Colors.Focus))
		focusStyle := &Style{
			BorderColor: &focusBorderColor,
			BorderWidth: style.BorderWidth,
		}
		DrawBorderWithTheme(screen, absBounds, focusStyle, theme)
	} else {
		DrawBorderWithTheme(screen, absBounds, style, theme)
	}
}

// GetText returns the current text
func (ti *TextInput) GetText() string {
	return ti.Text
}

// SetText sets the text
func (ti *TextInput) SetText(text string) {
	ti.Text = text
	if ti.cursorPos > len(ti.Text) {
		ti.cursorPos = len(ti.Text)
	}
}

// Checkbox is a checkbox element
type Checkbox struct {
	*ElementBase
	Label    string
	Checked  bool
	OnChange func(checked bool)
}

// NewCheckbox creates a new checkbox
func NewCheckbox(id, label string) *Checkbox {
	cb := &Checkbox{
		ElementBase: NewElementBase(id),
		Label:       label,
		Checked:     false,
	}

	// Set default size
	cb.SetSize(18, 18)

	return cb
}

// GetType returns the element type
func (cb *Checkbox) GetType() string {
	return "Checkbox"
}

// Update updates the checkbox
func (cb *Checkbox) Update() {
	if !cb.visible || !cb.enabled {
		return
	}

	cb.UpdateHoverState()

	// Handle click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		if cb.IsWithin(mx, my) {
			oldChecked := cb.Checked
			cb.Checked = !cb.Checked
			if cb.OnChange != nil {
				cb.OnChange(cb.Checked)
			}
			// Fire event
			if oldChecked != cb.Checked {
				event.GetQueuedInstance().QueueEvent(CheckboxChangeEvent{
					CheckboxID: cb.GetID(),
					Checkbox:   cb,
					Checked:    cb.Checked,
				})
			}
		}
	}
}

// Layout calculates dimensions
func (cb *Checkbox) Layout() {
	style := cb.GetComputedStyle()

	// Start with default checkbox size
	width := cb.bounds.Width
	height := cb.bounds.Height

	// Apply width/height from style if specified
	if style.Width != nil {
		width = *style.Width
	}
	if style.Height != nil {
		height = *style.Height
	}

	// Apply min/max size constraints
	width, height = ApplySizeConstraints(width, height, style)

	cb.bounds.Width = width
	cb.bounds.Height = height
}

// Draw draws the checkbox
func (cb *Checkbox) Draw(screen *ebiten.Image) {
	if !cb.visible {
		return
	}

	// Get absolute position for drawing
	absX, absY := cb.GetAbsolutePosition()

	// Draw checkbox box
	boxBounds := Rect{
		X:      absX,
		Y:      absY,
		Width:  18,
		Height: 18,
	}

	// Background
	bgColor := color.RGBA{50, 50, 60, 255}
	if cb.hovered {
		bgColor = color.RGBA{60, 60, 75, 255}
	}
	DrawRect(screen, boxBounds, bgColor)

	// Border
	borderColor := color.RGBA{100, 100, 120, 255}
	DrawRectStroke(screen, boxBounds, 1, borderColor)

	// Check mark
	if cb.Checked {
		checkColor := color.RGBA{0, 100, 200, 255}
		// Draw checkmark using simple lines
		checkBounds := Rect{
			X:      absX + 4,
			Y:      absY + 4,
			Width:  10,
			Height: 10,
		}
		DrawRect(screen, checkBounds, checkColor)
	}

	// Draw label if present
	if cb.Label != "" {
		labelColor := color.RGBA{230, 230, 230, 255}
		text.Draw(screen, cb.Label, 14.0, absX+24, absY+2, labelColor)
	}
}
