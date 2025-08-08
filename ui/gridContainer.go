package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// GridContainer arranges child elements in a grid pattern.
type GridContainer struct {
	ElementBase
	Children    []ElementInterface // ordered children
	ChildLookup map[string]int     // name -> index in Children
	Columns     int
	CellWidth   int
	CellHeight  int
	SpacingX    int
	SpacingY    int

	MaxWidth  int
	MaxHeight int
	ScrollX   int
	ScrollY   int
	image     *ebiten.Image // Offscreen buffer for drawing

	vBarBg *ebiten.Image
	vThumb *ebiten.Image
	hBarBg *ebiten.Image
	hThumb *ebiten.Image
}

// NewGridContainer creates a new grid container.
func NewGridContainer(name string, x, y, columns, cellW, cellH, spacingX, spacingY int, maxWidth, maxHeight int) *GridContainer {
	g := &GridContainer{
		ElementBase: ElementBase{
			Name:   name,
			X:      x,
			Y:      y,
			Width:  0,
			Height: 0,
			op:     &ebiten.DrawImageOptions{},
		},
		Children:    []ElementInterface{},
		ChildLookup: map[string]int{},
		Columns:     columns,
		CellWidth:   cellW,
		CellHeight:  cellH,
		SpacingX:    spacingX,
		SpacingY:    spacingY,
		MaxWidth:    maxWidth,
		MaxHeight:   maxHeight,
		ScrollX:     0,
		ScrollY:     0,
	}

	clipW := g.MaxWidth
	clipH := g.MaxHeight
	if clipW == 0 {
		clipW = g.contentWidth()
	}
	if clipH == 0 {
		clipH = g.contentHeight()
	}

	g.image = ebiten.NewImage(clipW, clipH)

	return g
}

// AddChild adds a child element to the grid.
func (g *GridContainer) AddChild(child ElementInterface) {
	g.ChildLookup[child.GetName()] = len(g.Children)
	g.Children = append(g.Children, child)
}

// GetChild retrieves a child by name.
func (g *GridContainer) GetChild(name string) ElementInterface {
	if idx, ok := g.ChildLookup[name]; ok {
		return g.Children[idx]
	}
	return nil
}

// RemoveChild removes a child by name.
func (g *GridContainer) RemoveChild(name string) {
	if idx, ok := g.ChildLookup[name]; ok {
		g.Children = append(g.Children[:idx], g.Children[idx+1:]...)
		// Rebuild lookup map
		g.ChildLookup = map[string]int{}
		for i, c := range g.Children {
			g.ChildLookup[c.GetName()] = i
		}
	}
}

// Update positions children and updates them.
func (g *GridContainer) Update(parentX, parentY int) {
	// Handle scrolling (mouse wheel)
	if g.MaxHeight > 0 {
		_, yoff := ebiten.Wheel()
		if yoff != 0 {
			maxScrollY := g.contentHeight() - g.MaxHeight
			if maxScrollY < 0 {
				maxScrollY = 0
			}
			g.ScrollY -= int(yoff * 20)
			if g.ScrollY < 0 {
				g.ScrollY = 0
			}
			if g.ScrollY > maxScrollY {
				g.ScrollY = maxScrollY
			}
		}
	}
	if g.MaxWidth > 0 {
		xoff, _ := ebiten.Wheel()
		if xoff != 0 {
			maxScrollX := g.contentWidth() - g.MaxWidth
			if maxScrollX < 0 {
				maxScrollX = 0
			}
			g.ScrollX -= int(xoff * 20)
			if g.ScrollX < 0 {
				g.ScrollX = 0
			}
			if g.ScrollX > maxScrollX {
				g.ScrollX = maxScrollX
			}
		}
	}

	i := 0
	for _, child := range g.Children {
		col := i % g.Columns
		row := i / g.Columns
		x := g.X + col*(g.CellWidth+g.SpacingX) - g.ScrollX
		y := g.Y + row*(g.CellHeight+g.SpacingY) - g.ScrollY
		// Only update if within visible area
		if g.isVisible(x-g.X, y-g.Y) {
			child.SetPosition(x, y)
			child.Update(parentX, parentY)
		}
		i++
	}

}

// Draw draws all children, clipped to the scroll area if needed.
func (g *GridContainer) Draw(screen *ebiten.Image, parentX, parentY int, theme *Theme) {
	ox := g.X + parentX
	oy := g.Y + parentY

	clipW := g.MaxWidth
	clipH := g.MaxHeight
	if clipW == 0 {
		clipW = g.contentWidth()
	}
	if clipH == 0 {
		clipH = g.contentHeight()
	}

	g.image.Clear() // Clear the offscreen buffer before drawing

	// Draw to offscreen buffer for clipping

	i := 0
	for _, child := range g.Children {
		col := i % g.Columns
		row := i / g.Columns
		x := col*(g.CellWidth+g.SpacingX) - g.ScrollX
		y := row*(g.CellHeight+g.SpacingY) - g.ScrollY
		if g.isVisible(x, y) {
			child.SetPosition(x+g.X, y+g.Y)
			child.Draw(g.image, parentX, parentY, theme)
		}
		i++
	}

	g.op.GeoM.Reset()
	g.op.GeoM.Translate(float64(ox), float64(oy))
	screen.DrawImage(g.image, g.op)

	// Draw scrollbars if needed
	if g.MaxHeight > 0 && g.contentHeight() > g.MaxHeight {
		g.drawVScrollbar(screen, ox+clipW-10, oy, clipH)
	}
	if g.MaxWidth > 0 && g.contentWidth() > g.MaxWidth {
		g.drawHScrollbar(screen, ox, oy+clipH-10, clipW)
	}
}

// Helper: check if a cell is visible in the scroll area
func (g *GridContainer) isVisible(x, y int) bool {
	clipW := g.MaxWidth
	clipH := g.MaxHeight
	if clipW == 0 {
		clipW = g.contentWidth()
	}
	if clipH == 0 {
		clipH = g.contentHeight()
	}
	return x+g.CellWidth > 0 && x < clipW && y+g.CellHeight > 0 && y < clipH
}

// Helper: total content width
func (g *GridContainer) contentWidth() int {
	return g.Columns*g.CellWidth + (g.Columns-1)*g.SpacingX
}

// Helper: total content height
func (g *GridContainer) contentHeight() int {
	count := len(g.Children)
	rows := (count + g.Columns - 1) / g.Columns
	return rows*g.CellHeight + (rows-1)*g.SpacingY
}

// Draw vertical scrollbar
func (g *GridContainer) drawVScrollbar(screen *ebiten.Image, x, y, height int) {
	barW := 8
	barH := height
	totalH := g.contentHeight()
	if totalH <= g.MaxHeight {
		return
	}
	thumbH := g.MaxHeight * barH / totalH
	if thumbH < 16 {
		thumbH = 16
	}
	scrollRange := barH - thumbH
	thumbY := y
	if totalH > g.MaxHeight && scrollRange > 0 {
		thumbY = y + (scrollRange*g.ScrollY)/(totalH-g.MaxHeight)
	}

	// Cache or recreate vertical scrollbar background
	if g.vBarBg == nil || g.vBarBg.Bounds().Dx() != barW || g.vBarBg.Bounds().Dy() != barH {
		g.vBarBg = ebiten.NewImage(barW, barH)
		g.vBarBg.Fill(color.RGBA{60, 60, 60, 180})
	}
	g.op.GeoM.Reset()
	g.op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(g.vBarBg, g.op)

	// Cache or recreate vertical thumb
	if g.vThumb == nil || g.vThumb.Bounds().Dx() != barW || g.vThumb.Bounds().Dy() != thumbH {
		g.vThumb = ebiten.NewImage(barW, thumbH)
		g.vThumb.Fill(color.RGBA{160, 160, 160, 220})
	}
	g.op.GeoM.Reset()
	g.op.GeoM.Translate(float64(x), float64(thumbY))
	screen.DrawImage(g.vThumb, g.op)
}

// Draw horizontal scrollbar
func (g *GridContainer) drawHScrollbar(screen *ebiten.Image, x, y, width int) {
	barH := 8
	barW := width
	totalW := g.contentWidth()
	if totalW <= g.MaxWidth {
		return
	}
	thumbW := g.MaxWidth * barW / totalW
	if thumbW < 16 {
		thumbW = 16
	}
	scrollRange := barW - thumbW
	thumbX := x
	if totalW > g.MaxWidth && scrollRange > 0 {
		thumbX = x + (scrollRange*g.ScrollX)/(totalW-g.MaxWidth)
	}

	// Cache or recreate horizontal scrollbar background
	if g.hBarBg == nil || g.hBarBg.Bounds().Dx() != barW || g.hBarBg.Bounds().Dy() != barH {
		g.hBarBg = ebiten.NewImage(barW, barH)
		g.hBarBg.Fill(color.RGBA{60, 60, 60, 180})
	}
	g.op.GeoM.Reset()
	g.op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(g.hBarBg, g.op)

	// Cache or recreate horizontal thumb
	if g.hThumb == nil || g.hThumb.Bounds().Dx() != thumbW || g.hThumb.Bounds().Dy() != barH {
		g.hThumb = ebiten.NewImage(thumbW, barH)
		g.hThumb.Fill(color.RGBA{160, 160, 160, 220})
	}
	g.op.GeoM.Reset()
	g.op.GeoM.Translate(float64(thumbX), float64(y))
	screen.DrawImage(g.hThumb, g.op)
}
