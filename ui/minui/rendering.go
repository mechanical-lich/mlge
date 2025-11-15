package minui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mechanical-lich/mlge/resource"
)

// DrawBackground draws the background of an element based on its style
func DrawBackground(screen *ebiten.Image, bounds Rect, style *Style) {
	if style == nil {
		return
	}

	// Get background color
	var bgColor color.Color = color.Transparent
	if style.BackgroundColor != nil {
		bgColor = *style.BackgroundColor
	}

	// If there's a background image, draw it
	if style.BackgroundImage != nil && *style.BackgroundImage != "" {
		img, ok := resource.Textures[*style.BackgroundImage]
		if ok && img != nil {
			// Draw tiled or stretched background
			op := &ebiten.DrawImageOptions{}

			// Scale to fit bounds
			imgBounds := img.Bounds()
			scaleX := float64(bounds.Width) / float64(imgBounds.Dx())
			scaleY := float64(bounds.Height) / float64(imgBounds.Dy())

			op.GeoM.Scale(scaleX, scaleY)
			op.GeoM.Translate(float64(bounds.X), float64(bounds.Y))

			// Apply opacity
			if style.Opacity != nil {
				op.ColorScale.ScaleAlpha(float32(*style.Opacity))
			}

			screen.DrawImage(img, op)
			return
		}
	}

	// Draw solid color background with border radius
	borderRadius := 0
	if style.BorderRadius != nil {
		borderRadius = *style.BorderRadius
	}

	opacity := float32(1.0)
	if style.Opacity != nil {
		opacity = float32(*style.Opacity)
	}

	// Apply opacity to background color
	r, g, b, a := bgColor.RGBA()
	bgColorWithOpacity := color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(float32(a>>8) * opacity),
	}

	if borderRadius > 0 {
		DrawRoundedRect(screen, bounds, borderRadius, bgColorWithOpacity)
	} else {
		DrawRect(screen, bounds, bgColorWithOpacity)
	}
}

// DrawBorder draws the border of an element based on its style
func DrawBorder(screen *ebiten.Image, bounds Rect, style *Style) {
	if style == nil {
		return
	}

	borderWidth := 0
	if style.BorderWidth != nil {
		borderWidth = *style.BorderWidth
	}

	if borderWidth <= 0 {
		return
	}

	borderColor := color.RGBA{80, 80, 90, 255}
	if style.BorderColor != nil {
		r, g, b, a := (*style.BorderColor).RGBA()
		borderColor = color.RGBA{
			R: uint8(r >> 8),
			G: uint8(g >> 8),
			B: uint8(b >> 8),
			A: uint8(a >> 8),
		}
	}

	borderRadius := 0
	if style.BorderRadius != nil {
		borderRadius = *style.BorderRadius
	}

	opacity := float32(1.0)
	if style.Opacity != nil {
		opacity = float32(*style.Opacity)
	}

	// Apply opacity
	r, g, b, a := borderColor.RGBA()
	borderColorWithOpacity := color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(float32(a>>8) * opacity),
	}

	if borderRadius > 0 {
		DrawRoundedRectStroke(screen, bounds, borderRadius, float32(borderWidth), borderColorWithOpacity)
	} else {
		DrawRectStroke(screen, bounds, float32(borderWidth), borderColorWithOpacity)
	}
}

// DrawRect draws a filled rectangle
func DrawRect(screen *ebiten.Image, bounds Rect, clr color.Color) {
	vector.DrawFilledRect(
		screen,
		float32(bounds.X),
		float32(bounds.Y),
		float32(bounds.Width),
		float32(bounds.Height),
		clr,
		false,
	)
}

// DrawRectStroke draws a rectangle outline
func DrawRectStroke(screen *ebiten.Image, bounds Rect, strokeWidth float32, clr color.Color) {
	vector.StrokeRect(
		screen,
		float32(bounds.X),
		float32(bounds.Y),
		float32(bounds.Width),
		float32(bounds.Height),
		strokeWidth,
		clr,
		false,
	)
}

// DrawRoundedRect draws a filled rounded rectangle
func DrawRoundedRect(screen *ebiten.Image, bounds Rect, radius int, clr color.Color) {
	x := float32(bounds.X)
	y := float32(bounds.Y)
	w := float32(bounds.Width)
	h := float32(bounds.Height)
	r := float32(radius)

	// Create a path for rounded rectangle
	var path vector.Path

	// Top-left corner
	path.MoveTo(x+r, y)

	// Top edge and top-right corner
	path.LineTo(x+w-r, y)
	path.ArcTo(x+w, y, x+w, y+r, r)

	// Right edge and bottom-right corner
	path.LineTo(x+w, y+h-r)
	path.ArcTo(x+w, y+h, x+w-r, y+h, r)

	// Bottom edge and bottom-left corner
	path.LineTo(x+r, y+h)
	path.ArcTo(x, y+h, x, y+h-r, r)

	// Left edge and back to start
	path.LineTo(x, y+r)
	path.ArcTo(x, y, x+r, y, r)

	path.Close()

	vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)

	// Apply color to vertices
	for i := range vertices {
		vertices[i].ColorR = 1
		vertices[i].ColorG = 1
		vertices[i].ColorB = 1
		vertices[i].ColorA = 1
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = true
	op.FillRule = ebiten.NonZero

	// Create a 1x1 white image for coloring
	whiteImg := ebiten.NewImage(1, 1)
	whiteImg.Fill(clr)

	screen.DrawTriangles(vertices, indices, whiteImg, op)
}

// DrawRoundedRectStroke draws a rounded rectangle outline
func DrawRoundedRectStroke(screen *ebiten.Image, bounds Rect, radius int, strokeWidth float32, clr color.Color) {
	x := float32(bounds.X)
	y := float32(bounds.Y)
	w := float32(bounds.Width)
	h := float32(bounds.Height)
	r := float32(radius)

	// Create a path for rounded rectangle
	var path vector.Path

	// Top-left corner
	path.MoveTo(x+r, y)

	// Top edge and top-right corner
	path.LineTo(x+w-r, y)
	path.ArcTo(x+w, y, x+w, y+r, r)

	// Right edge and bottom-right corner
	path.LineTo(x+w, y+h-r)
	path.ArcTo(x+w, y+h, x+w-r, y+h, r)

	// Bottom edge and bottom-left corner
	path.LineTo(x+r, y+h)
	path.ArcTo(x, y+h, x, y+h-r, r)

	// Left edge and back to start
	path.LineTo(x, y+r)
	path.ArcTo(x, y, x+r, y, r)

	path.Close()

	vertices, indices := path.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width:    strokeWidth,
		LineJoin: vector.LineJoinRound,
		LineCap:  vector.LineCapRound,
	})

	// Apply color to vertices
	for i := range vertices {
		vertices[i].ColorR = 1
		vertices[i].ColorG = 1
		vertices[i].ColorB = 1
		vertices[i].ColorA = 1
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = true

	// Create a 1x1 image with the stroke color
	colorImg := ebiten.NewImage(1, 1)
	colorImg.Fill(clr)

	screen.DrawTriangles(vertices, indices, colorImg, op)
}

// GetContentBounds returns the bounds minus padding and border
func GetContentBounds(bounds Rect, style *Style) Rect {
	content := bounds

	if style == nil {
		return content
	}

	// Subtract border
	if style.BorderWidth != nil {
		bw := *style.BorderWidth
		content.X += bw
		content.Y += bw
		content.Width -= bw * 2
		content.Height -= bw * 2
	}

	// Subtract padding
	if style.Padding != nil {
		content.X += style.Padding.Left
		content.Y += style.Padding.Top
		content.Width -= style.Padding.Left + style.Padding.Right
		content.Height -= style.Padding.Top + style.Padding.Bottom
	}

	return content
}

// ApplyMargin applies margin to bounds
func ApplyMargin(bounds Rect, style *Style) Rect {
	if style == nil || style.Margin == nil {
		return bounds
	}

	return Rect{
		X:      bounds.X + style.Margin.Left,
		Y:      bounds.Y + style.Margin.Top,
		Width:  bounds.Width - style.Margin.Left - style.Margin.Right,
		Height: bounds.Height - style.Margin.Top - style.Margin.Bottom,
	}
}

// CreateSubImage creates a sub-image for clipping content
func CreateSubImage(screen *ebiten.Image, bounds Rect) *ebiten.Image {
	return screen.SubImage(image.Rect(
		bounds.X,
		bounds.Y,
		bounds.X+bounds.Width,
		bounds.Y+bounds.Height,
	)).(*ebiten.Image)
}
