package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/resource"
	"github.com/mechanical-lich/mlge/text"
	theming "github.com/mechanical-lich/mlge/ui/themed/theming"
)

// Represents a single radio button element.
type RadioButton struct {
	ElementBase
	Label           string
	Tooltip         string
	tooltipBg       *ebiten.Image // cache for tooltip background
	tooltipBgWidth  int
	tooltipBgHeight int
	OnClick         func() // Optional onclick handler function
}

// Represents a group of radio buttons where only one can be selected at a time.
type RadioGroup struct {
	ElementBase
	Buttons   []*RadioButton
	Selected  int
	OnChange  func(selected int) // Optional onchange handler function
	OnClicked func(selected int) // Optional onclick handler function
}

// Creates a new Radio group with the given buttons.
func NewRadioGroup(name string, buttons []*RadioButton) *RadioGroup {
	if buttons == nil {
		buttons = make([]*RadioButton, 0)
	}
	return &RadioGroup{
		ElementBase: ElementBase{
			Name: name,
		},
		Buttons:  buttons,
		Selected: -1, // No selection by default
	}
}

// Creates a new Radio button with the given parameters.
func NewRadioButton(name string, x, y int, label string, tooltip string) *RadioButton {
	w, h := text.Measure(label, 16)

	return &RadioButton{
		ElementBase: ElementBase{
			Name:   name,
			X:      x,
			Y:      y,
			Width:  int(w) + 10,
			Height: int(h) + 10, // Add some padding
			Op:     &ebiten.DrawImageOptions{},
		},
		Label:   label,
		Tooltip: tooltip,
	}
}

// Creates a new Radio button with an icon.
func NewIconRadioButton(name string, x int, y, iconX, iconY int, iconResource string, label string, tooltip string) *RadioButton {
	w, h := text.Measure(label, 16)

	b := &RadioButton{
		ElementBase: ElementBase{
			Name:         name,
			X:            x,
			Y:            y,
			IconX:        iconX,
			IconY:        iconY,
			IconResource: iconResource,
			Width:        int(w + 10 + 16),
			Height:       int(h + 16),
			Op:           &ebiten.DrawImageOptions{},
		},
		Label:   label,
		Tooltip: tooltip,
	}

	return b
}

func (rg *RadioGroup) Update() {
	for i, button := range rg.Buttons {
		button.Update()
		if button.IsClicked() {
			if rg.OnChange != nil && rg.Selected != i {
				rg.OnChange(i)
			}

			if rg.OnClicked != nil {
				rg.OnClicked(i)
			}
			rg.Selected = i
			break
		}
	}
}

func (rg *RadioGroup) Draw(screen *ebiten.Image, theme *theming.Theme) {
	for i, button := range rg.Buttons {
		selected := (i == rg.Selected)
		button.Draw(screen, selected, theme)
	}
}

// Get the currently selected radio button.
func (rg *RadioGroup) GetSelected() *RadioButton {
	if rg.Selected < 0 || rg.Selected >= len(rg.Buttons) {
		return nil
	}
	return rg.Buttons[rg.Selected]
}

func (rb *RadioButton) Update() {
	if rb.OnClick != nil && rb.IsJustClicked() {
		rb.OnClick()
	}
}

func (rb *RadioButton) Draw(screen *ebiten.Image, selected bool, theme *theming.Theme) {
	absX, absY := rb.GetAbsolutePosition()
	rb.Op.GeoM.Reset()
	rb.Op.GeoM.Scale(float64(rb.Width)/float64(theme.RadioButton.Width), float64(rb.Height)/float64(theme.RadioButton.Height))
	rb.Op.GeoM.Translate(float64(absX), float64(absY))
	sX := theme.RadioButton.SrcX
	sY := theme.RadioButton.SrcY

	if selected {
		sX += 32
	}

	screen.DrawImage(resource.GetSubImage("ui", sX, sY, theme.RadioButton.Width, theme.RadioButton.Height), rb.Op)
	if rb.IconResource != "" {
		text.Draw(screen, rb.Label, 15, absX+5+16, absY+5, color.White)

		if rb.IconResource != "" {
			rb.Op.GeoM.Reset()
			rb.Op.GeoM.Scale(1.0, 1.0)
			rb.Op.GeoM.Translate(float64(absX+5), float64(absY+5))
			screen.DrawImage(resource.GetSubImage(rb.IconResource, rb.IconX, rb.IconY, spriteSizeW, spriteSizeH), rb.Op)
		}
	} else {
		text.Draw(screen, rb.Label, 15, absX+5, absY+5, color.White)
	}

	// Tooltip rendering
	cX, cY := ebiten.CursorPosition()
	if rb.IsWithin(cX, cY) && rb.Tooltip != "" {
		tw, th := text.Measure(rb.Tooltip, 14)
		tooltipW := int(tw + 10)
		tooltipH := int(th + 8)
		absX, absY := rb.GetAbsolutePosition()
		tooltipX := absX + rb.Width + 8
		tooltipY := absY

		// Only recreate if size changes
		if rb.tooltipBg == nil || rb.tooltipBgWidth != tooltipW || rb.tooltipBgHeight != tooltipH {
			rb.tooltipBg = ebiten.NewImage(tooltipW, tooltipH)
			rb.tooltipBg.Fill(color.RGBA{30, 30, 30, 220})
			rb.tooltipBgWidth = tooltipW
			rb.tooltipBgHeight = tooltipH
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(tooltipX), float64(tooltipY))
		screen.DrawImage(rb.tooltipBg, op)
		text.Draw(screen, rb.Tooltip, 14, tooltipX+5, tooltipY+5, color.White)
	}
}
