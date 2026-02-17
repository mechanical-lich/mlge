package text

import (
	"bytes"
	_ "embed"
	"image/color"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/text/language"
)

//go:embed Roboto-Regular.ttf
var robotoRegularTTF []byte
var robotoRegularFaceSource *text.GoTextFaceSource

var op = &text.DrawOptions{}

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(robotoRegularTTF))
	if err != nil {
		log.Fatal(err)
	}
	robotoRegularFaceSource = s
}

func Draw(dst *ebiten.Image, txt string, size float64, x int, y int, clr color.Color) {
	op.GeoM.Reset()
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

// Wrap splits s into lines, each with at most maxChars characters.
// It tries to break at spaces, but will break long words if needed.
// Explicit '\n' in s will always start a new line.
func Wrap(s string, maxChars int, maxLines int) []string {
	if maxChars <= 0 {
		return []string{s}
	}
	rawLines := strings.Split(s, "\n")
	var lines []string

	for _, raw := range rawLines {
		words := strings.Fields(raw)
		var line string
		for _, word := range words {
			if len(line)+len(word)+1 > maxChars {
				if line != "" {
					lines = append(lines, line)
					if maxLines > 0 && len(lines) >= maxLines {
						return lines
					}
				}
				line = word
			} else {
				if line != "" {
					line += " "
				}
				line += word
			}
		}
		if line != "" {
			lines = append(lines, line)
			if maxLines > 0 && len(lines) >= maxLines {
				return lines
			}
		}
		// If the original line was empty (i.e., a blank line), preserve it
		if len(words) == 0 {
			lines = append(lines, "")
			if maxLines > 0 && len(lines) >= maxLines {
				return lines
			}
		}
	}
	if maxLines > 0 && len(lines) > maxLines {
		return lines[:maxLines]
	}
	return lines
}
