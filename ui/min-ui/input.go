package minui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mechanical-lich/mlge/event"
	"github.com/mechanical-lich/mlge/text/v2"
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

	// Set default style
	bgColor := color.Color(color.RGBA{255, 255, 255, 255})
	borderColor := color.Color(color.RGBA{0, 0, 0, 255})
	borderWidth := 1
	padding := NewEdgeInsetsLR(4, 8)

	ti.style.BackgroundColor = &bgColor
	ti.style.BorderColor = &borderColor
	ti.style.BorderWidth = &borderWidth
	ti.style.Padding = padding

	// Focus style
	focusBorderColor := color.Color(color.RGBA{0, 100, 200, 255})
	ti.style.FocusStyle = &Style{
		BorderColor: &focusBorderColor,
	}

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

	// Get absolute position for drawing
	absX, absY := ti.GetAbsolutePosition()
	absBounds := Rect{
		X:      absX,
		Y:      absY,
		Width:  ti.bounds.Width,
		Height: ti.bounds.Height,
	}

	// Draw background
	DrawBackground(screen, absBounds, style)

	// Draw text or placeholder
	contentBounds := GetContentBounds(absBounds, style)

	displayText := ti.Text
	textColor := color.RGBA{0, 0, 0, 255}

	if displayText == "" && ti.Placeholder != "" {
		displayText = ti.Placeholder
		textColor = color.RGBA{150, 150, 150, 255}
	}

	if style.ForegroundColor != nil && ti.Text != "" {
		r, g, b, a := (*style.ForegroundColor).RGBA()
		textColor = color.RGBA{
			R: uint8(r >> 8),
			G: uint8(g >> 8),
			B: uint8(b >> 8),
			A: uint8(a >> 8),
		}
	}

	text.Draw(screen, displayText, 14.0, contentBounds.X, contentBounds.Y+6, textColor)

	// Draw cursor if focused
	if ti.focused {
		cursorX := contentBounds.X + ti.cursorPos*8
		cursorBounds := Rect{
			X:      cursorX,
			Y:      contentBounds.Y + 2,
			Width:  2,
			Height: contentBounds.Height - 4,
		}
		cursorColor := color.RGBA{0, 0, 0, 255}

		// Blink cursor
		if (ebiten.TPS()*2/3)%2 == 0 {
			DrawRect(screen, cursorBounds, cursorColor)
		}
	}

	// Draw border
	DrawBorder(screen, absBounds, style)
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
	bgColor := color.RGBA{255, 255, 255, 255}
	if cb.hovered {
		bgColor = color.RGBA{240, 240, 255, 255}
	}
	DrawRect(screen, boxBounds, bgColor)

	// Border
	borderColor := color.RGBA{0, 0, 0, 255}
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
		labelColor := color.RGBA{0, 0, 0, 255}
		text.Draw(screen, cb.Label, 14.0, absX+24, absY+2, labelColor)
	}
}
