package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	theming "github.com/mechanical-lich/mlge/ui/themed/theming"
	views "github.com/mechanical-lich/mlge/ui/themed/views"
)

// GUI Main struct that manages the gui for the game. Includes the cursor
type GUI struct {
	State       views.GUIViewInterface
	Theme       *theming.Theme
	CursorImage *ebiten.Image
	op          *ebiten.DrawImageOptions
}

func NewGUI(startingView views.GUIViewInterface, theme *theming.Theme) *GUI {
	if theme == nil {
		theme = &theming.DefaultTheme
	}
	return &GUI{State: startingView, Theme: theme, op: &ebiten.DrawImageOptions{}}
}

func (g *GUI) Update() {
	if g.State != nil {
		g.State.UpdateElements()
		g.State.Update()
	}
}

func (g *GUI) Draw(screen *ebiten.Image) {
	g.State.Draw(screen, g.Theme)
	g.State.DrawElements(screen, g.Theme)
	g.DrawCursor(screen)
}

func (g *GUI) DrawCursor(screen *ebiten.Image) {
	//Cursor logic
	if g.CursorImage != nil {
		ebiten.SetCursorMode(ebiten.CursorModeHidden)
		cX, cY := ebiten.CursorPosition()
		g.op.GeoM.Reset()
		g.op.GeoM.Translate(float64(cX), float64(cY))
		g.op.GeoM.Scale(1.0, 1.0)
		screen.DrawImage(g.CursorImage, g.op)
	} else {
		ebiten.SetCursorMode(ebiten.CursorModeVisible)

	}
}
