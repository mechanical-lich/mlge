package text

import (
	"bytes"
	_ "embed"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/text/language"
)

//go:embed Roboto-Regular.ttf
var robotoRegularTTF []byte
var robotoRegularFaceSource *text.GoTextFaceSource

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(robotoRegularTTF))
	if err != nil {
		log.Fatal(err)
	}
	robotoRegularFaceSource = s
}

func Draw(dst *ebiten.Image, txt string, size float64, x int, y int, clr color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))

	r, g, b, a := clr.RGBA()
	op.ColorScale.SetR(float32(r>>8) / 255.0)
	op.ColorScale.SetG(float32(g>>8) / 255.0)
	op.ColorScale.SetB(float32(b>>8) / 255.0)
	op.ColorScale.SetA(float32(a>>8) / 255.0)

	f := &text.GoTextFace{
		Source:    robotoRegularFaceSource,
		Direction: text.DirectionLeftToRight,
		Size:      size,
		Language:  language.AmericanEnglish,
	}
	text.Draw(dst, txt, f, op)
}

func Measure(txt string, size float64) (float64, float64) {
	f := &text.GoTextFace{
		Source:    robotoRegularFaceSource,
		Direction: text.DirectionLeftToRight,
		Size:      size,
		Language:  language.AmericanEnglish,
	}
	return text.Measure(txt, f, 4)
}
