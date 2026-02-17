package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/resource"
	"github.com/mechanical-lich/mlge/text"
	theming "github.com/mechanical-lich/mlge/ui/themed/theming"
)

// Slider represents a horizontal slider UI element
type Slider struct {
	ElementBase
	Label     string
	Min       float64
	Max       float64
	Value     float64
	Step      float64
	dragging  bool
	OnChanged func(value float64)
}

// NewSlider creates a new slider with the given parameters
func NewSlider(name string, x, y, width int, label string, min, max, value, step float64) *Slider {
	return &Slider{
		ElementBase: ElementBase{
			Name:   name,
			X:      x,
			Y:      y,
			Width:  width,
			Height: 20,
			Op:     &ebiten.DrawImageOptions{},
		},
		Label: label,
		Min:   min,
		Max:   max,
		Value: value,
		Step:  step,
	}
}

func (s *Slider) Update() {
	absX, absY := s.GetScreenPosition()
	mouseX, mouseY := ebiten.CursorPosition()

	// Check if clicking on slider
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// Check if mouse is within bounds
		if mouseX >= absX && mouseX <= absX+s.Width && mouseY >= absY && mouseY <= absY+s.Height {
			s.dragging = true
		}
	} else {
		if s.dragging && s.OnChanged != nil {
			s.OnChanged(s.Value)
		}
		s.dragging = false
	}

	// Update value while dragging
	if s.dragging {
		// Calculate value based on mouse position
		relativeX := mouseX - absX
		if relativeX < 0 {
			relativeX = 0
		}
		if relativeX > s.Width {
			relativeX = s.Width
		}

		// Map to value range
		ratio := float64(relativeX) / float64(s.Width)
		newValue := s.Min + ratio*(s.Max-s.Min)

		// Apply step
		if s.Step > 0 {
			newValue = float64(int(newValue/s.Step+0.5)) * s.Step
		}

		s.Value = newValue
		if s.Value < s.Min {
			s.Value = s.Min
		}
		if s.Value > s.Max {
			s.Value = s.Max
		}
	}
}

func (s *Slider) Draw(screen *ebiten.Image, theme *theming.Theme) {
	absX, absY := s.GetAbsolutePosition()

	// Draw label
	if s.Label != "" {
		labelText := fmt.Sprintf("%s: %.2f", s.Label, s.Value)
		text.Draw(screen, labelText, 12, absX, absY-14, theme.Colors.Text)
	}

	// Draw track
	s.Op.GeoM.Reset()
	trackY := absY + 6
	s.Op.GeoM.Scale(float64(s.Width)/float64(theme.Slider.TrackWidth), 1.0)
	s.Op.GeoM.Translate(float64(absX), float64(trackY))
	screen.DrawImage(
		resource.GetSubImage("ui", theme.Slider.TrackSrcX, theme.Slider.TrackSrcY,
			theme.Slider.TrackWidth, theme.Slider.TrackHeight),
		s.Op,
	)

	// Calculate thumb position
	ratio := (s.Value - s.Min) / (s.Max - s.Min)
	if s.Max == s.Min {
		ratio = 0
	}
	thumbX := absX + int(ratio*float64(s.Width)) - theme.Slider.ThumbWidth/2

	// Draw thumb
	s.Op.GeoM.Reset()
	s.Op.GeoM.Translate(float64(thumbX), float64(absY))

	thumbSrcX := theme.Slider.ThumbSrcX
	if s.dragging {
		thumbSrcX += theme.Slider.ThumbWidth // Use different sprite when dragging
	}

	screen.DrawImage(
		resource.GetSubImage("ui", thumbSrcX, theme.Slider.ThumbSrcY,
			theme.Slider.ThumbWidth, theme.Slider.ThumbHeight),
		s.Op,
	)

	// Draw tick marks if applicable
	if s.Step > 0 && s.Max-s.Min > 0 {
		steps := int((s.Max - s.Min) / s.Step)
		if steps <= 20 { // Only draw if reasonable number of steps
			for i := 0; i <= steps; i++ {
				tickRatio := float64(i) / float64(steps)
				tickX := absX + int(tickRatio*float64(s.Width))
				// Draw small tick mark
				tickImg := ebiten.NewImage(1, 4)
				tickImg.Fill(color.RGBA{150, 150, 160, 255})
				s.Op.GeoM.Reset()
				s.Op.GeoM.Translate(float64(tickX), float64(trackY+theme.Slider.TrackHeight+1))
				screen.DrawImage(tickImg, s.Op)
			}
		}
	}
}

// SetValue sets the slider value
func (s *Slider) SetValue(value float64) {
	s.Value = value
	if s.Value < s.Min {
		s.Value = s.Min
	}
	if s.Value > s.Max {
		s.Value = s.Max
	}
}

// GetValue returns the current slider value
func (s *Slider) GetValue() float64 {
	return s.Value
}
