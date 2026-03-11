package minui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mechanical-lich/mlge/resource"
)

// Reusable buffers for rounded-rect drawing to avoid per-frame allocations.
var (
	rrVertices []ebiten.Vertex
	rrIndices  []uint16
	rrOp       = &ebiten.DrawTrianglesOptions{AntiAlias: true, FillRule: ebiten.NonZero}
	rrWhiteImg *ebiten.Image // lazily initialised 1×1 white pixel

	// Stroke-specific reusable buffers
	rrStrokeVertices []ebiten.Vertex
	rrStrokeIndices  []uint16
	rrStrokeOp       = &ebiten.DrawTrianglesOptions{AntiAlias: true}
	rrStrokeOptions  = &vector.StrokeOptions{LineJoin: vector.LineJoinRound, LineCap: vector.LineCapRound}

	// Reusable DrawImageOptions for sprite/9-slice rendering
	spriteOp ebiten.DrawImageOptions
)

// getWhitePixel returns a cached 1×1 white ebiten.Image.
func getWhitePixel() *ebiten.Image {
	if rrWhiteImg == nil {
		rrWhiteImg = ebiten.NewImage(1, 1)
		rrWhiteImg.Fill(color.White)
	}
	return rrWhiteImg
}

// colorToRGBA converts a color.Color to color.RGBA
func colorToRGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}
}

// ---- Sprite-based rendering functions ----

// DrawSprite draws a sprite from a sprite sheet, scaled to fit the bounds
func DrawSprite(screen *ebiten.Image, spriteSheet string, coords *SpriteCoords, bounds Rect) {
	if coords == nil {
		return
	}

	img := resource.GetSubImage(spriteSheet, coords.SrcX, coords.SrcY, coords.Width, coords.Height)
	if img == nil {
		return
	}

	spriteOp.GeoM.Reset()
	spriteOp.ColorScale.Reset()
	scaleX := float64(bounds.Width) / float64(coords.Width)
	scaleY := float64(bounds.Height) / float64(coords.Height)
	spriteOp.GeoM.Scale(scaleX, scaleY)
	spriteOp.GeoM.Translate(float64(bounds.X), float64(bounds.Y))

	screen.DrawImage(img, &spriteOp)
}

// DrawSpriteWithOpacity draws a sprite with opacity
func DrawSpriteWithOpacity(screen *ebiten.Image, spriteSheet string, coords *SpriteCoords, bounds Rect, opacity float32) {
	if coords == nil {
		return
	}

	img := resource.GetSubImage(spriteSheet, coords.SrcX, coords.SrcY, coords.Width, coords.Height)
	if img == nil {
		return
	}

	spriteOp.GeoM.Reset()
	spriteOp.ColorScale.Reset()
	scaleX := float64(bounds.Width) / float64(coords.Width)
	scaleY := float64(bounds.Height) / float64(coords.Height)
	spriteOp.GeoM.Scale(scaleX, scaleY)
	spriteOp.GeoM.Translate(float64(bounds.X), float64(bounds.Y))
	spriteOp.ColorScale.ScaleAlpha(opacity)

	screen.DrawImage(img, &spriteOp)
}

// Draw9Slice draws a 9-slice scaled sprite to fill the bounds
func Draw9Slice(screen *ebiten.Image, spriteSheet string, coords *NineSliceCoords, bounds Rect) {
	if coords == nil {
		return
	}

	x := bounds.X
	y := bounds.Y
	w := bounds.Width
	h := bounds.Height
	srcX := coords.SrcX
	srcY := coords.SrcY
	tileSize := coords.TileSize
	tileScale := coords.TileScale
	scaledTile := tileSize * tileScale

	// Draw corners
	spriteOp.GeoM.Reset()
	spriteOp.ColorScale.Reset()
	spriteOp.GeoM.Scale(float64(tileScale), float64(tileScale))
	spriteOp.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(resource.GetSubImage(spriteSheet, srcX, srcY, tileSize, tileSize), &spriteOp)

	spriteOp.GeoM.Reset()
	spriteOp.GeoM.Scale(float64(tileScale), float64(tileScale))
	spriteOp.GeoM.Translate(float64(x+w-scaledTile), float64(y))
	screen.DrawImage(resource.GetSubImage(spriteSheet, srcX+2*tileSize, srcY, tileSize, tileSize), &spriteOp)

	spriteOp.GeoM.Reset()
	spriteOp.GeoM.Scale(float64(tileScale), float64(tileScale))
	spriteOp.GeoM.Translate(float64(x), float64(y+h-scaledTile))
	screen.DrawImage(resource.GetSubImage(spriteSheet, srcX, srcY+2*tileSize, tileSize, tileSize), &spriteOp)

	spriteOp.GeoM.Reset()
	spriteOp.GeoM.Scale(float64(tileScale), float64(tileScale))
	spriteOp.GeoM.Translate(float64(x+w-scaledTile), float64(y+h-scaledTile))
	screen.DrawImage(resource.GetSubImage(spriteSheet, srcX+2*tileSize, srcY+2*tileSize, tileSize, tileSize), &spriteOp)

	// Draw edges - Top and bottom
	for dx := scaledTile; dx < w-scaledTile; dx += scaledTile {
		spriteOp.GeoM.Reset()
		spriteOp.GeoM.Scale(float64(tileScale), float64(tileScale))
		spriteOp.GeoM.Translate(float64(x+dx), float64(y))
		screen.DrawImage(resource.GetSubImage(spriteSheet, srcX+tileSize, srcY, tileSize, tileSize), &spriteOp)

		spriteOp.GeoM.Reset()
		spriteOp.GeoM.Scale(float64(tileScale), float64(tileScale))
		spriteOp.GeoM.Translate(float64(x+dx), float64(y+h-scaledTile))
		screen.DrawImage(resource.GetSubImage(spriteSheet, srcX+tileSize, srcY+2*tileSize, tileSize, tileSize), &spriteOp)
	}

	// Left and right
	for dy := scaledTile; dy < h-scaledTile; dy += scaledTile {
		spriteOp.GeoM.Reset()
		spriteOp.GeoM.Scale(float64(tileScale), float64(tileScale))
		spriteOp.GeoM.Translate(float64(x), float64(y+dy))
		screen.DrawImage(resource.GetSubImage(spriteSheet, srcX, srcY+tileSize, tileSize, tileSize), &spriteOp)

		spriteOp.GeoM.Reset()
		spriteOp.GeoM.Scale(float64(tileScale), float64(tileScale))
		spriteOp.GeoM.Translate(float64(x+w-scaledTile), float64(y+dy))
		screen.DrawImage(resource.GetSubImage(spriteSheet, srcX+2*tileSize, srcY+tileSize, tileSize, tileSize), &spriteOp)
	}

	// Center
	for dx := scaledTile; dx < w-scaledTile; dx += scaledTile {
		for dy := scaledTile; dy < h-scaledTile; dy += scaledTile {
			spriteOp.GeoM.Reset()
			spriteOp.GeoM.Scale(float64(tileScale), float64(tileScale))
			spriteOp.GeoM.Translate(float64(x+dx), float64(y+dy))
			screen.DrawImage(resource.GetSubImage(spriteSheet, srcX+tileSize, srcY+tileSize, tileSize, tileSize), &spriteOp)
		}
	}
}

// Draw9SliceToImage draws a 9-slice to an offscreen image (for caching)
func Draw9SliceToImage(dst *ebiten.Image, spriteSheet string, coords *NineSliceCoords) {
	if coords == nil {
		return
	}

	bounds := dst.Bounds()
	Draw9Slice(dst, spriteSheet, coords, Rect{
		X:      0,
		Y:      0,
		Width:  bounds.Dx(),
		Height: bounds.Dy(),
	})
}

// ---- Vector-based rendering functions ----

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
			spriteOp.GeoM.Reset()
			spriteOp.ColorScale.Reset()

			// Scale to fit bounds
			imgBounds := img.Bounds()
			scaleX := float64(bounds.Width) / float64(imgBounds.Dx())
			scaleY := float64(bounds.Height) / float64(imgBounds.Dy())

			spriteOp.GeoM.Scale(scaleX, scaleY)
			spriteOp.GeoM.Translate(float64(bounds.X), float64(bounds.Y))

			// Apply opacity
			if style.Opacity != nil {
				spriteOp.ColorScale.ScaleAlpha(float32(*style.Opacity))
			}

			screen.DrawImage(img, &spriteOp)
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

// DrawBackgroundWithTheme draws the background using theme colors as fallback
func DrawBackgroundWithTheme(screen *ebiten.Image, bounds Rect, style *Style, theme *Theme) {
	if style == nil && theme == nil {
		return
	}

	// Get background color - prefer explicit style, fall back to theme
	var bgColor color.Color = color.Transparent
	if style != nil && style.BackgroundColor != nil {
		bgColor = *style.BackgroundColor
	} else if theme != nil {
		bgColor = theme.Colors.Surface
	}

	// If there's a background image, draw it
	if style != nil && style.BackgroundImage != nil && *style.BackgroundImage != "" {
		img, ok := resource.Textures[*style.BackgroundImage]
		if ok && img != nil {
			spriteOp.GeoM.Reset()
			spriteOp.ColorScale.Reset()
			imgBounds := img.Bounds()
			scaleX := float64(bounds.Width) / float64(imgBounds.Dx())
			scaleY := float64(bounds.Height) / float64(imgBounds.Dy())
			spriteOp.GeoM.Scale(scaleX, scaleY)
			spriteOp.GeoM.Translate(float64(bounds.X), float64(bounds.Y))
			if style.Opacity != nil {
				spriteOp.ColorScale.ScaleAlpha(float32(*style.Opacity))
			}
			screen.DrawImage(img, &spriteOp)
			return
		}
	}

	// Draw solid color background with border radius
	borderRadius := 0
	if style != nil && style.BorderRadius != nil {
		borderRadius = *style.BorderRadius
	}

	opacity := float32(1.0)
	if style != nil && style.Opacity != nil {
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

// DrawBorderWithTheme draws the border using theme colors as fallback
func DrawBorderWithTheme(screen *ebiten.Image, bounds Rect, style *Style, theme *Theme) {
	if style == nil && theme == nil {
		return
	}

	borderWidth := 0
	if style != nil && style.BorderWidth != nil {
		borderWidth = *style.BorderWidth
	}

	if borderWidth <= 0 {
		return
	}

	// Get border color - prefer explicit style, fall back to theme
	var borderColor color.Color = color.RGBA{80, 80, 90, 255}
	if style != nil && style.BorderColor != nil {
		borderColor = *style.BorderColor
	} else if theme != nil {
		borderColor = theme.Colors.Border
	}

	borderRadius := 0
	if style != nil && style.BorderRadius != nil {
		borderRadius = *style.BorderRadius
	}

	opacity := float32(1.0)
	if style != nil && style.Opacity != nil {
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

	var path vector.Path
	path.MoveTo(x+r, y)
	path.LineTo(x+w-r, y)
	path.ArcTo(x+w, y, x+w, y+r, r)
	path.LineTo(x+w, y+h-r)
	path.ArcTo(x+w, y+h, x+w-r, y+h, r)
	path.LineTo(x+r, y+h)
	path.ArcTo(x, y+h, x, y+h-r, r)
	path.LineTo(x, y+r)
	path.ArcTo(x, y, x+r, y, r)
	path.Close()

	// Reuse backing arrays
	rrVertices = rrVertices[:0]
	rrIndices = rrIndices[:0]
	rrVertices, rrIndices = path.AppendVerticesAndIndicesForFilling(rrVertices, rrIndices)

	// Tint vertices with the fill color (white pixel × vertex color = desired color)
	cr, cg, cb, ca := clr.RGBA()
	rf := float32(cr) / 0xffff
	gf := float32(cg) / 0xffff
	bf := float32(cb) / 0xffff
	af := float32(ca) / 0xffff
	for i := range rrVertices {
		rrVertices[i].ColorR = rf
		rrVertices[i].ColorG = gf
		rrVertices[i].ColorB = bf
		rrVertices[i].ColorA = af
	}

	rrOp.AntiAlias = true
	rrOp.FillRule = ebiten.NonZero
	rrOp.ColorScaleMode = ebiten.ColorScaleModePremultipliedAlpha
	screen.DrawTriangles(rrVertices, rrIndices, getWhitePixel(), rrOp)
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

	// Reuse backing arrays
	rrStrokeVertices = rrStrokeVertices[:0]
	rrStrokeIndices = rrStrokeIndices[:0]
	rrStrokeOptions.Width = strokeWidth
	rrStrokeVertices, rrStrokeIndices = path.AppendVerticesAndIndicesForStroke(rrStrokeVertices, rrStrokeIndices, rrStrokeOptions)

	// Tint vertices with the stroke color (white pixel × vertex color = desired color)
	cr, cg, cb, ca := clr.RGBA()
	rf := float32(cr) / 0xffff
	gf := float32(cg) / 0xffff
	bf := float32(cb) / 0xffff
	af := float32(ca) / 0xffff
	for i := range rrStrokeVertices {
		rrStrokeVertices[i].ColorR = rf
		rrStrokeVertices[i].ColorG = gf
		rrStrokeVertices[i].ColorB = bf
		rrStrokeVertices[i].ColorA = af
	}

	rrStrokeOp.AntiAlias = true
	rrStrokeOp.ColorScaleMode = ebiten.ColorScaleModePremultipliedAlpha
	screen.DrawTriangles(rrStrokeVertices, rrStrokeIndices, getWhitePixel(), rrStrokeOp)
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
