package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/config"
	"github.com/mechanical-lich/mlge/state"
)

// GUI Main struct that manages the gui for the game. Includes the cursor
type GUI struct {
	State       GUIViewInterface
	Theme       *Theme
	CursorImage *ebiten.Image
}

func NewGUI(startingView GUIViewInterface, theme *Theme) *GUI {
	if theme == nil {
		theme = &DefaultTheme
	}
	return &GUI{State: startingView, Theme: theme}
}

func (g *GUI) Update(s state.StateInterface) {
	if g.State != nil {
		g.State.UpdateElements(s)
		g.State.Update(s)
	}
}

func (g *GUI) Draw(screen *ebiten.Image, s state.StateInterface) {
	g.State.Draw(screen, s, g.Theme)
	g.State.DrawElements(screen, s, g.Theme)
	g.DrawCursor(screen, s)
}

func (g *GUI) DrawCursor(screen *ebiten.Image, s state.StateInterface) {
	//Cursor logic
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(config.TileSizeW/config.SpriteSizeW), float64(config.TileSizeH/config.SpriteSizeH))
	if g.CursorImage != nil {
		ebiten.SetCursorMode(ebiten.CursorModeHidden)
		cX, cY := ebiten.CursorPosition()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(cX), float64(cY))
		op.GeoM.Scale(1.0, 1.0)
		screen.DrawImage(g.CursorImage, op)
	} else {
		ebiten.SetCursorMode(ebiten.CursorModeVisible)

	}
}
