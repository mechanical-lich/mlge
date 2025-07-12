package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/state"
)

// Base GUIView interface.
// Since we are dealing with interfaces the GUIView is being passed around by value instead of reference
type GUIViewInterface interface {
	Update(state state.StateInterface)
	Draw(screen *ebiten.Image, s state.StateInterface)
	GetInputFocused() bool
}

// GUIViewBase gives views some basic functionality when inherited.
type GUIViewBase struct {
	Buttons     map[string]*Button
	RadioGroups map[string]*RadioGroup
	Toggles     map[string]*Toggle
	Inputs      map[string]*InputField
}

func (g *GUIViewBase) AddButton(button *Button) {
	if g.Buttons == nil {
		g.Buttons = make(map[string]*Button, 0)
	}
	g.Buttons[button.Name] = button
}

func (g *GUIViewBase) AddRadioGroup(group *RadioGroup) {
	if g.RadioGroups == nil {
		g.RadioGroups = make(map[string]*RadioGroup, 0)
	}
	g.RadioGroups[group.Name] = group
}

func (g *GUIViewBase) AddToggle(toggle *Toggle) {
	if g.Toggles == nil {
		g.Toggles = make(map[string]*Toggle, 0)
	}
	g.Toggles[toggle.Name] = toggle
}

func (g *GUIViewBase) AddInputField(input *InputField) {
	if g.Inputs == nil {
		g.Inputs = make(map[string]*InputField, 0)
	}
	g.Inputs[input.Name] = input
}

func (g *GUIViewBase) UpdateElements() {
	for _, group := range g.RadioGroups {
		group.Update()
	}

	for _, toggle := range g.Toggles {
		toggle.Update()
	}

	for _, input := range g.Inputs {
		input.Update()
	}
}

func (g *GUIViewBase) DrawElements(screen *ebiten.Image) {
	// Draw buttons
	for _, b := range g.Buttons {
		b.Draw(screen)
	}

	// Draw radio groups
	for _, rg := range g.RadioGroups {
		rg.Draw(screen)
	}

	// Draw toggles
	for _, t := range g.Toggles {
		t.Draw(screen)
	}

	// Draw input fields
	for _, input := range g.Inputs {
		input.Draw(screen)
	}
}

func (g *GUIViewBase) GetInputFocused() bool {
	for _, input := range g.Inputs {
		if input.Focused {
			return true
		}
	}
	return false
}
