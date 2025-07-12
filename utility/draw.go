package utility

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/resource"
)

// Draw9Slice draws a 9-slice scaled image to the destination image.
// SrcX, srcY are the top-left corner of the 9-slice source image.
func Draw9Slice(dst *ebiten.Image, x, y, w, h, srcX, srcY, tileSize int, tileScale int) {
	scaledTile := tileSize * tileScale

	// Draw corners
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(tileScale), float64(tileScale))
	op.GeoM.Translate(float64(x), float64(y))
	dst.DrawImage(resource.Textures["ui"].SubImage(image.Rect(srcX, srcY, srcX+tileSize, srcY+tileSize)).(*ebiten.Image), &op)

	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(tileScale), float64(tileScale))
	op.GeoM.Translate(float64(x+w-scaledTile), float64(y))
	dst.DrawImage(resource.Textures["ui"].SubImage(image.Rect(srcX+2*tileSize, srcY, srcX+3*tileSize, srcY+tileSize)).(*ebiten.Image), &op)

	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(tileScale), float64(tileScale))
	op.GeoM.Translate(float64(x), float64(y+h-scaledTile))
	dst.DrawImage(resource.Textures["ui"].SubImage(image.Rect(srcX, srcY+2*tileSize, srcX+tileSize, srcY+3*tileSize)).(*ebiten.Image), &op)

	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(tileScale), float64(tileScale))
	op.GeoM.Translate(float64(x+w-scaledTile), float64(y+h-scaledTile))
	dst.DrawImage(resource.Textures["ui"].SubImage(image.Rect(srcX+2*tileSize, srcY+2*tileSize, srcX+3*tileSize, srcY+3*tileSize)).(*ebiten.Image), &op)

	// Draw edges
	// Top and bottom
	for dx := scaledTile; dx < w-scaledTile; dx += scaledTile {
		op = ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(tileScale), float64(tileScale))
		op.GeoM.Translate(float64(x+dx), float64(y))
		dst.DrawImage(resource.Textures["ui"].SubImage(image.Rect(srcX+tileSize, srcY, srcX+2*tileSize, srcY+tileSize)).(*ebiten.Image), &op)

		op = ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(tileScale), float64(tileScale))
		op.GeoM.Translate(float64(x+dx), float64(y+h-scaledTile))
		dst.DrawImage(resource.Textures["ui"].SubImage(image.Rect(srcX+tileSize, srcY+2*tileSize, srcX+2*tileSize, srcY+3*tileSize)).(*ebiten.Image), &op)
	}
	// Left and right
	for dy := scaledTile; dy < h-scaledTile; dy += scaledTile {
		op = ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(tileScale), float64(tileScale))
		op.GeoM.Translate(float64(x), float64(y+dy))
		dst.DrawImage(resource.Textures["ui"].SubImage(image.Rect(srcX, srcY+tileSize, srcX+tileSize, srcY+2*tileSize)).(*ebiten.Image), &op)

		op = ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(tileScale), float64(tileScale))
		op.GeoM.Translate(float64(x+w-scaledTile), float64(y+dy))
		dst.DrawImage(resource.Textures["ui"].SubImage(image.Rect(srcX+2*tileSize, srcY+tileSize, srcX+3*tileSize, srcY+2*tileSize)).(*ebiten.Image), &op)
	}
	// Center
	for dx := scaledTile; dx < w-scaledTile; dx += scaledTile {
		for dy := scaledTile; dy < h-scaledTile; dy += scaledTile {
			op = ebiten.DrawImageOptions{}
			op.GeoM.Scale(float64(tileScale), float64(tileScale))
			op.GeoM.Translate(float64(x+dx), float64(y+dy))
			dst.DrawImage(resource.Textures["ui"].SubImage(image.Rect(srcX+tileSize, srcY+tileSize, srcX+2*tileSize, srcY+2*tileSize)).(*ebiten.Image), &op)
		}
	}
}
