package resource

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Define a generic function to get a sub-image
func GetSubImage(texture *ebiten.Image, x, y, width, height int) *ebiten.Image {
	return texture.SubImage(image.Rect(x, y, x+width, y+height)).(*ebiten.Image)
}
