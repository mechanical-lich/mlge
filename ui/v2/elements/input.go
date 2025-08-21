package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mechanical-lich/mlge/resource"
	"github.com/mechanical-lich/mlge/text/v2"
	theming "github.com/mechanical-lich/mlge/ui/v2/theming"
)

type InputField struct {
	ElementBase
	MaxLength int
	Value     []rune
	Cursor    int // index in Value
}

func NewInputField(name string, x, y, width, maxLength int) *InputField {
	// Generate string l
	_, h := text.Measure("A", 16)
	return &InputField{
		ElementBase: ElementBase{
			Name:   name,
			X:      x,
			Y:      y,
			Width:  width,
			Height: int(h + 10), // Add some padding
			Op:     &ebiten.DrawImageOptions{},
		},
		MaxLength: maxLength,
		Value:     []rune{},
		Cursor:    0,
	}
}

func (f *InputField) Update() {
	cX, cY := ebiten.CursorPosition()
	// Focus/unfocus logic
	absX, absY := f.GetScreenPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if cX >= absX && cX <= absX+f.Width && cY >= absY && cY <= absY+16 {
			f.Focused = true
		} else {
			f.Focused = false
		}
	}

	if !f.Focused {
		return
	}

	// Handle text input
	for _, r := range ebiten.AppendInputChars(nil) {
		if r == '\n' || r == '\r' {
			continue
		}
		if len(f.Value) < f.MaxLength {
			// Insert at cursor
			before := f.Value[:f.Cursor]
			after := f.Value[f.Cursor:]
			f.Value = append(append(before, r), after...)
			f.Cursor++
		}
	}

	// Handle backspace/delete
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && f.Cursor > 0 {
		before := f.Value[:f.Cursor-1]
		after := f.Value[f.Cursor:]
		f.Value = append(before, after...)
		f.Cursor--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDelete) && f.Cursor < len(f.Value) {
		before := f.Value[:f.Cursor]
		after := f.Value[f.Cursor+1:]
		f.Value = append(before, after...)
	}

	// Handle arrow keys
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && f.Cursor > 0 {
		f.Cursor--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && f.Cursor < len(f.Value) {
		f.Cursor++
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		f.Focused = false
	}
}

func (f *InputField) Draw(screen *ebiten.Image, theme *theming.Theme) {
	absX, absY := f.GetAbsolutePosition()
	// Stretch the input field sprite horizontally
	f.Op.GeoM.Reset()
	scaleX := float64(f.Width) / 48.0
	f.Op.GeoM.Scale(scaleX, float64(f.Height)/16.0)
	f.Op.GeoM.Translate(float64(absX), float64(absY))
	screen.DrawImage(resource.GetSubImage("ui", theme.InputField.SrcX, theme.InputField.SrcY, theme.InputField.Width, theme.InputField.Height), f.Op)

	// Draw text
	txt := string(f.Value)
	text.Draw(screen, txt, 15, absX+5, absY+5, color.White)

	// Draw cursor if focused
	if f.Focused {
		cursorX := absX + 5
		if f.Cursor > 0 {
			sub := string(f.Value[:f.Cursor])
			w, _ := text.Measure(sub, 16)
			cursorX += int(w)
		}
		// Draw a simple vertical line as cursor
		vector.DrawFilledRect(screen, float32(cursorX), float32(absY+4), 2, 12, color.White, false)
	}
}

func (f *InputField) SetValue(val string) {
	runes := []rune(val)
	if len(runes) > f.MaxLength {
		runes = runes[:f.MaxLength]
	}
	f.Value = runes
	if f.Cursor > len(f.Value) {
		f.Cursor = len(f.Value)
	}
}
