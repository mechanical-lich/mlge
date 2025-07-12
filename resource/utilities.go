package resource

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Subimage cache so we can reuse sub-images without creating new ones
type subCacheKey struct {
	name                string
	x, y, width, height int
}

var subImageCache = make(map[subCacheKey]*ebiten.Image)

// Define a generic function to get a sub-image
func GetSubImage(name string, x, y, width, height int) *ebiten.Image {
	cacheKey := subCacheKey{name: name, x: x, y: y, width: width, height: height}
	if img, found := subImageCache[cacheKey]; found {
		return img
	}
	img := Textures[name].SubImage(image.Rect(x, y, x+width, y+height)).(*ebiten.Image)
	subImageCache[cacheKey] = img
	return img
}

// Define a generic function to get a sub-image
func GetSubImageByTexture(texture *ebiten.Image, x, y, width, height int) *ebiten.Image {
	cacheKey := subCacheKey{name: "", x: x, y: y, width: width, height: height}
	if img, found := subImageCache[cacheKey]; found {
		return img
	}
	img := texture.SubImage(image.Rect(x, y, x+width, y+height)).(*ebiten.Image)
	subImageCache[cacheKey] = img
	return img
}
