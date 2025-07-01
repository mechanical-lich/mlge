package ui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/config"
	"github.com/mechanical-lich/mlge/resource"
	"github.com/mechanical-lich/mlge/state"
)

// GUI Main struct that manages the gui for the game. Includes the cursor
type GUI struct {
	State GUIViewInterface
}

func NewGUI(startingView GUIViewInterface) *GUI {
	return &GUI{State: startingView}
}

func (g *GUI) Update(s state.StateInterface) {
	if g.State != nil {
		g.State.Update(s)
	}
}

func (g *GUI) Draw(screen *ebiten.Image, s state.StateInterface) {
	g.State.Draw(screen, s)

	//g.DrawCursor(screen, s)
}

func (g *GUI) DrawCursor(screen *ebiten.Image, s state.StateInterface) {
	//Cursor logic

	cX, cY := ebiten.CursorPosition()

	var cursorY = 128
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		cursorY = 144
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(config.TileSizeW/config.SpriteSizeW), float64(config.TileSizeH/config.SpriteSizeH))

	if cX > config.World_W {
		op.GeoM.Translate(float64(cX), float64(cY))
		screen.DrawImage(resource.Textures["ui"].SubImage(image.Rect(64, cursorY, 64+config.SpriteSizeW, cursorY+config.SpriteSizeH)).(*ebiten.Image), op)
		//s.drawSprite(int32(g.Cursor.X), int32(g.Cursor.Y), 64, cursorY, 255, 255, 255, g.uiTexture) //Cursor?
	} else {
		//This works because the math is being done on ints then turned into a float giving us a nice even number.
		op.GeoM.Translate(float64((cX/config.TileSizeW)*config.TileSizeW), float64((cY/config.TileSizeH)*config.TileSizeH))
		screen.DrawImage(resource.Textures["ui"].SubImage(image.Rect(128, cursorY, 128+config.SpriteSizeW, cursorY+config.SpriteSizeH)).(*ebiten.Image), op)
	}
}
