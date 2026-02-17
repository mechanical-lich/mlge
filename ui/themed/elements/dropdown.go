package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mechanical-lich/mlge/resource"
	"github.com/mechanical-lich/mlge/text"
	theming "github.com/mechanical-lich/mlge/ui/themed/theming"
)

// Dropdown represents a dropdown/select UI element
type Dropdown struct {
	ElementBase
	Options       []string
	SelectedIndex int
	IsOpen        bool
	hoverIndex    int
	OnChanged     func(index int, value string)
}

// NewDropdown creates a new dropdown with the given parameters
func NewDropdown(name string, x, y, width int, options []string, selectedIndex int) *Dropdown {
	if selectedIndex < 0 || selectedIndex >= len(options) {
		selectedIndex = 0
	}
	return &Dropdown{
		ElementBase: ElementBase{
			Name:   name,
			X:      x,
			Y:      y,
			Width:  width,
			Height: 20,
			Op:     &ebiten.DrawImageOptions{},
		},
		Options:       options,
		SelectedIndex: selectedIndex,
		hoverIndex:    -1,
	}
}

func (d *Dropdown) Update() {
	absX, absY := d.GetScreenPosition()
	mouseX, mouseY := ebiten.CursorPosition()

	// Check if clicking the dropdown header
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if mouseX >= absX && mouseX <= absX+d.Width &&
			mouseY >= absY && mouseY <= absY+d.Height {
			d.IsOpen = !d.IsOpen
		} else if d.IsOpen {
			// Check if clicking an option
			optionsY := absY + d.Height
			for i := range d.Options {
				optionY := optionsY + i*20
				if mouseX >= absX && mouseX <= absX+d.Width &&
					mouseY >= optionY && mouseY <= optionY+20 {
					d.SelectedIndex = i
					d.IsOpen = false
					if d.OnChanged != nil {
						d.OnChanged(i, d.Options[i])
					}
					break
				}
			}
			// Close if clicked outside
			if mouseY < absY || mouseY > optionsY+len(d.Options)*20 ||
				mouseX < absX || mouseX > absX+d.Width {
				d.IsOpen = false
			}
		}
	}

	// Track hover
	if d.IsOpen {
		d.hoverIndex = -1
		optionsY := absY + d.Height
		for i := range d.Options {
			optionY := optionsY + i*20
			if mouseX >= absX && mouseX <= absX+d.Width &&
				mouseY >= optionY && mouseY <= optionY+20 {
				d.hoverIndex = i
				break
			}
		}
	}
}

func (d *Dropdown) Draw(screen *ebiten.Image, theme *theming.Theme) {
	absX, absY := d.GetAbsolutePosition()

	// Draw dropdown header
	d.Op.GeoM.Reset()
	scaleX := float64(d.Width) / float64(theme.Dropdown.Width)
	d.Op.GeoM.Scale(scaleX, 1.0)
	d.Op.GeoM.Translate(float64(absX), float64(absY))

	screen.DrawImage(
		resource.GetSubImage("ui", theme.Dropdown.SrcX, theme.Dropdown.SrcY,
			theme.Dropdown.Width, theme.Dropdown.Height),
		d.Op,
	)

	// Draw selected option text
	if d.SelectedIndex >= 0 && d.SelectedIndex < len(d.Options) {
		text.Draw(screen, d.Options[d.SelectedIndex], 14, absX+4, absY+3, theme.Colors.Text)
	}

	// Draw arrow indicator
	arrowX := absX + d.Width - 14
	arrowY := absY + 7
	if d.IsOpen {
		// Up arrow
		vector.DrawFilledRect(screen, float32(arrowX+2), float32(arrowY), 6, 1, theme.Colors.Text, false)
		vector.DrawFilledRect(screen, float32(arrowX+3), float32(arrowY+1), 4, 1, theme.Colors.Text, false)
		vector.DrawFilledRect(screen, float32(arrowX+4), float32(arrowY+2), 2, 1, theme.Colors.Text, false)
	} else {
		// Down arrow
		vector.DrawFilledRect(screen, float32(arrowX+2), float32(arrowY+5), 6, 1, theme.Colors.Text, false)
		vector.DrawFilledRect(screen, float32(arrowX+3), float32(arrowY+4), 4, 1, theme.Colors.Text, false)
		vector.DrawFilledRect(screen, float32(arrowX+4), float32(arrowY+3), 2, 1, theme.Colors.Text, false)
	}

	// Draw options if open
	if d.IsOpen {
		optionsY := absY + d.Height
		for i, option := range d.Options {
			optionY := optionsY + i*20

			// Draw option background
			bgColor := theme.Colors.Surface
			if i == d.hoverIndex {
				bgColor = theme.Colors.Primary
			} else if i == d.SelectedIndex {
				bgColor = theme.Colors.Secondary
			}

			optionImg := ebiten.NewImage(d.Width, 20)
			optionImg.Fill(bgColor)
			d.Op.GeoM.Reset()
			d.Op.GeoM.Translate(float64(absX), float64(optionY))
			screen.DrawImage(optionImg, d.Op)

			// Draw border
			vector.StrokeRect(screen, float32(absX), float32(optionY),
				float32(d.Width), 20, 1, theme.Colors.Border, false)

			// Draw option text
			textColor := theme.Colors.Text
			if i == d.hoverIndex {
				textColor = color.RGBA{255, 255, 255, 255}
			}
			text.Draw(screen, option, 14, absX+4, optionY+3, textColor)
		}
	}
}

// SetSelectedIndex sets the selected index
func (d *Dropdown) SetSelectedIndex(index int) {
	if index >= 0 && index < len(d.Options) {
		d.SelectedIndex = index
	}
}

// GetSelectedIndex returns the currently selected index
func (d *Dropdown) GetSelectedIndex() int {
	return d.SelectedIndex
}

// GetSelectedValue returns the currently selected value
func (d *Dropdown) GetSelectedValue() string {
	if d.SelectedIndex >= 0 && d.SelectedIndex < len(d.Options) {
		return d.Options[d.SelectedIndex]
	}
	return ""
}

// SetOptions updates the options list
func (d *Dropdown) SetOptions(options []string) {
	d.Options = options
	if d.SelectedIndex >= len(options) {
		d.SelectedIndex = 0
	}
}
