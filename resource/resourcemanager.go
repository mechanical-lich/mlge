package resource

import (
	"errors"
	"fmt"
	"image"
	"log"
	"os"

	"github.com/go-fonts/liberation/liberationsansregular"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/audio"
	"github.com/mechanical-lich/mlge/config"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var Textures map[string]*ebiten.Image
var Fonts map[string]font.Face
var Sounds map[string]audio.AudioResource

// LoadImageAsTexture - Loads an image in the texture map with the given name and path.
func LoadImageAsTexture(name string, path string) error {
	if Textures == nil {
		log.Print("Initialize resource manager")
		Textures = make(map[string]*ebiten.Image)
	}
	img, err := LoadImage(path)
	if err != nil {
		return err
	}

	Textures[name] = img
	return nil
}

// LoadImage - Loads an image from the specified path and returns it as an ebiten.Image
func LoadImage(path string) (*ebiten.Image, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening tileset " + path)
		return nil, errors.New("error opening tileset " + path)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

// LoadFont - Loads a font into the font map with the given name and path.
func LoadFont(name string, path string) error {
	if Fonts == nil {
		log.Print("Initialize fonts")
		Fonts = make(map[string]font.Face)
	}
	raw := liberationsansregular.TTF
	tt, err := opentype.Parse(raw)
	if err != nil {
		return err
	}

	fontFace, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     config.DPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return err
	}

	Fonts[name] = fontFace
	return nil
}
