package utility

import (
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
	dst.DrawImage(resource.GetSubImage("ui", srcX, srcY, tileSize, tileSize), &op)

	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(tileScale), float64(tileScale))
	op.GeoM.Translate(float64(x+w-scaledTile), float64(y))
	dst.DrawImage(resource.GetSubImage("ui", srcX+2*tileSize, srcY, tileSize, tileSize), &op)

	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(tileScale), float64(tileScale))
	op.GeoM.Translate(float64(x), float64(y+h-scaledTile))
	dst.DrawImage(resource.GetSubImage("ui", srcX, srcY+2*tileSize, tileSize, tileSize), &op)

	op = ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(tileScale), float64(tileScale))
	op.GeoM.Translate(float64(x+w-scaledTile), float64(y+h-scaledTile))
	dst.DrawImage(resource.GetSubImage("ui", srcX+2*tileSize, srcY+2*tileSize, tileSize, tileSize), &op)

	// Draw edges
	// Top and bottom
	for dx := scaledTile; dx < w-scaledTile; dx += scaledTile {
		op = ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(tileScale), float64(tileScale))
		op.GeoM.Translate(float64(x+dx), float64(y))
		dst.DrawImage(resource.GetSubImage("ui", srcX+tileSize, srcY, tileSize, tileSize), &op)

		op = ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(tileScale), float64(tileScale))
		op.GeoM.Translate(float64(x+dx), float64(y+h-scaledTile))
		dst.DrawImage(resource.GetSubImage("ui", srcX+tileSize, srcY+2*tileSize, tileSize, tileSize), &op)
	}
	// Left and right
	for dy := scaledTile; dy < h-scaledTile; dy += scaledTile {
		op = ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(tileScale), float64(tileScale))
		op.GeoM.Translate(float64(x), float64(y+dy))
		dst.DrawImage(resource.GetSubImage("ui", srcX, srcY+tileSize, tileSize, tileSize), &op)

		op = ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(tileScale), float64(tileScale))
		op.GeoM.Translate(float64(x+w-scaledTile), float64(y+dy))
		dst.DrawImage(resource.GetSubImage("ui", srcX+2*tileSize, srcY+tileSize, tileSize, tileSize), &op)
	}
	// Center
	for dx := scaledTile; dx < w-scaledTile; dx += scaledTile {
		for dy := scaledTile; dy < h-scaledTile; dy += scaledTile {
			op = ebiten.DrawImageOptions{}
			op.GeoM.Scale(float64(tileScale), float64(tileScale))
			op.GeoM.Translate(float64(x+dx), float64(y+dy))
			dst.DrawImage(resource.GetSubImage("ui", srcX+tileSize, srcY+tileSize, tileSize, tileSize), &op)
		}
	}
}
