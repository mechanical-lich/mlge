package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/state"
)

// Base GUIView interface.
// Since we are dealing with interfaces the GUIView is being passed around by value instead of reference
type GUIViewInterface interface {
	Update(state state.StateInterface)
	UpdateElements(state state.StateInterface)
	Draw(screen *ebiten.Image, s state.StateInterface, theme *Theme)
	DrawElements(screen *ebiten.Image, s state.StateInterface, theme *Theme)
	GetInputFocused() bool
	GetModalFocused() bool
	GetPosition() (int, int)
	SetPosition(x, y int)
	WithinModalBounds(mouseX, mouseY int) bool
}

// GUIViewBase gives views some basic functionality when inherited.
type GUIViewBase struct {
	Buttons     map[string]*Button
	RadioGroups map[string]*RadioGroup
	Toggles     map[string]*Toggle
	Inputs      map[string]*InputField
	Modals      map[string]*Modal
	X, Y        int // Add offset for the view
}

func (g *GUIViewBase) GetPosition() (int, int) {
	return g.X, g.Y
}

func (g *GUIViewBase) SetPosition(x, y int) {
	g.X = x
	g.Y = y
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

func (g *GUIViewBase) AddModal(modal *Modal) {
	if g.Modals == nil {
		g.Modals = make(map[string]*Modal)
	}

	g.Modals[modal.Name] = modal
}

func (g *GUIViewBase) UpdateElements(s state.StateInterface) {
	for _, group := range g.RadioGroups {
		group.Update(g.X, g.Y)
	}

	for _, toggle := range g.Toggles {
		toggle.Update(g.X, g.Y)
	}

	for _, input := range g.Inputs {
		input.Update(g.X, g.Y)
	}

	for _, modal := range g.Modals {
		modal.Update(s)
	}
}

func (g *GUIViewBase) DrawElements(screen *ebiten.Image, s state.StateInterface, theme *Theme) {
	// Draw buttons
	for _, b := range g.Buttons {
		b.Draw(screen, g.X, g.Y, theme)
	}

	// Draw radio groups
	for _, rg := range g.RadioGroups {
		rg.Draw(screen, g.X, g.Y, theme)
	}

	// Draw toggles
	for _, t := range g.Toggles {
		t.Draw(screen, g.X, g.Y, theme)
	}

	// Draw input fields
	for _, input := range g.Inputs {
		input.Draw(screen, g.X, g.Y, theme)
	}

	// Draw modals
	for _, modal := range g.Modals {
		if modal.Visible {
			modal.Draw(screen, s, theme)
		}
	}
}

func (g *GUIViewBase) GetInputFocused() bool {
	for _, input := range g.Inputs {
		if input.Focused {
			return true
		}
	}

	for _, modal := range g.Modals {
		if modal.Visible && modal.GetInputFocused() {
			return true
		}
	}
	return false
}

func (g *GUIViewBase) GetModalFocused() bool {
	for _, modal := range g.Modals {
		if modal.Visible {
			return true
		}
	}
	return false
}

func (g *GUIViewBase) WithinModalBounds(mouseX, mouseY int) bool {
	for _, modal := range g.Modals {
		if modal.Visible && modal.WithinBounds(mouseX, mouseY) {
			return true
		}
	}
	return false
}
