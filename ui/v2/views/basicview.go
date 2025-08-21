package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	theming "github.com/mechanical-lich/mlge/ui/v2/theming"
)

// BasicView is useful for basic views.
// By default it provides empty Update and Draw calls to match the View interface.
// These are overwritable by setting UpdateFunc and DrawFunc function pointers.
type BasicView struct {
	GUIViewBase
	UpdateFunc func()
	DrawFunc   func(screen *ebiten.Image, theme *theming.Theme)
}

func NewBasicView() *BasicView {
	view := &BasicView{
		UpdateFunc: func() {},
		DrawFunc:   func(screen *ebiten.Image, theme *theming.Theme) {},
	}
	return view
}

func (v *BasicView) Update() {
	if v.UpdateFunc == nil {
		return
	}
	v.UpdateFunc()
}

func (v *BasicView) Draw(screen *ebiten.Image, theme *theming.Theme) {
	if v.DrawFunc == nil {
		return
	}
	v.DrawFunc(screen, theme)
}
