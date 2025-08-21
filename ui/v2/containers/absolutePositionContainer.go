package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	elements "github.com/mechanical-lich/mlge/ui/v2/elements"
	theming "github.com/mechanical-lich/mlge/ui/v2/theming"
)

// AbsolutePositionContainer allows children to be placed at arbitrary positions and scrolls overflow.
type AbsolutePositionContainer struct {
	elements.ElementBase
	Children    []elements.ElementInterface
	ChildLookup map[string]int
	MaxWidth    int
	MaxHeight   int
	ScrollX     int
	ScrollY     int
	image       *ebiten.Image
	vBarBg      *ebiten.Image
	vThumb      *ebiten.Image
	hBarBg      *ebiten.Image
	hThumb      *ebiten.Image
}

func NewAbsolutePositionContainer(name string, x, y, maxWidth, maxHeight int) *AbsolutePositionContainer {
	c := &AbsolutePositionContainer{
		ElementBase: elements.ElementBase{
			Name:   name,
			X:      x,
			Y:      y,
			Width:  maxWidth,
			Height: maxHeight,
			Op:     &ebiten.DrawImageOptions{},
		},
		Children:    []elements.ElementInterface{},
		ChildLookup: map[string]int{},
		MaxWidth:    maxWidth,
		MaxHeight:   maxHeight,
		ScrollX:     0,
		ScrollY:     0,
	}
	clipW := c.MaxWidth
	clipH := c.MaxHeight
	if clipW == 0 {
		clipW = c.contentWidth()
	}
	if clipH == 0 {
		clipH = c.contentHeight()
	}
	c.image = ebiten.NewImage(clipW, clipH)
	return c
}

func (c *AbsolutePositionContainer) AddChild(child elements.ElementInterface) {
	child.SetParent(c)
	c.ChildLookup[child.GetName()] = len(c.Children)
	c.Children = append(c.Children, child)
}

func (c *AbsolutePositionContainer) GetChild(name string) elements.ElementInterface {
	if idx, ok := c.ChildLookup[name]; ok {
		return c.Children[idx]
	}
	return nil
}

func (c *AbsolutePositionContainer) RemoveChild(name string) {
	if idx, ok := c.ChildLookup[name]; ok {
		c.Children = append(c.Children[:idx], c.Children[idx+1:]...)
		c.ChildLookup = map[string]int{}
		for i, ch := range c.Children {
			c.ChildLookup[ch.GetName()] = i
		}
	}
}

func (c *AbsolutePositionContainer) Update() {
	// Handle scrolling (mouse wheel)
	if c.MaxHeight > 0 {
		_, yoff := ebiten.Wheel()
		if yoff != 0 {
			maxScrollY := c.contentHeight() - c.MaxHeight
			if maxScrollY < 0 {
				maxScrollY = 0
			}
			c.ScrollY -= int(yoff * 20)
			if c.ScrollY < 0 {
				c.ScrollY = 0
			}
			if c.ScrollY > maxScrollY {
				c.ScrollY = maxScrollY
			}
		}
	}
	if c.MaxWidth > 0 {
		xoff, _ := ebiten.Wheel()
		if xoff != 0 {
			maxScrollX := c.contentWidth() - c.MaxWidth
			if maxScrollX < 0 {
				maxScrollX = 0
			}
			c.ScrollX -= int(xoff * 20)
			if c.ScrollX < 0 {
				c.ScrollX = 0
			}
			if c.ScrollX > maxScrollX {
				c.ScrollX = maxScrollX
			}
		}
	}
	for _, child := range c.Children {
		x, y := child.GetPosition()
		if c.isVisible(x-c.ScrollX, y-c.ScrollY, child.GetWidth(), child.GetHeight()) {
			child.SetPosition(x-c.ScrollX, y-c.ScrollY)
			child.Update()
		}
	}
}

func (c *AbsolutePositionContainer) Draw(screen *ebiten.Image, theme *theming.Theme) {
	absX, absY := c.GetAbsolutePosition()
	clipW := c.MaxWidth
	clipH := c.MaxHeight
	if clipW == 0 {
		clipW = c.contentWidth()
	}
	if clipH == 0 {
		clipH = c.contentHeight()
	}
	c.image.Clear()
	for _, child := range c.Children {
		x, y := child.GetPosition()
		if c.isVisible(x, y, child.GetWidth(), child.GetHeight()) {
			child.Draw(c.image, theme)
		}
	}
	c.Op.GeoM.Reset()
	c.Op.GeoM.Translate(float64(c.X), float64(c.Y))
	screen.DrawImage(c.image, c.Op)
	if c.MaxHeight > 0 && c.contentHeight() > c.MaxHeight {
		c.drawVScrollbar(screen, absX+clipW-10, absY, clipH)
	}
	if c.MaxWidth > 0 && c.contentWidth() > c.MaxWidth {
		c.drawHScrollbar(screen, absX, absY+clipH-10, clipW)
	}
}

func (c *AbsolutePositionContainer) isVisible(x, y, w, h int) bool {
	clipW := c.MaxWidth
	clipH := c.MaxHeight
	if clipW == 0 {
		clipW = c.contentWidth()
	}
	if clipH == 0 {
		clipH = c.contentHeight()
	}
	return x+w > 0 && x < clipW && y+h > 0 && y < clipH
}

func (c *AbsolutePositionContainer) contentWidth() int {
	maxX := 0
	for _, child := range c.Children {
		x, _ := child.GetPosition()
		w := child.GetWidth()
		if x+w > maxX {
			maxX = x + w
		}
	}
	return maxX
}

func (c *AbsolutePositionContainer) contentHeight() int {
	maxY := 0
	for _, child := range c.Children {
		_, y := child.GetPosition()
		h := child.GetHeight()
		if y+h > maxY {
			maxY = y + h
		}
	}
	return maxY
}

func (c *AbsolutePositionContainer) drawVScrollbar(screen *ebiten.Image, x, y, height int) {
	barW := 8
	barH := height
	totalH := c.contentHeight()
	if totalH <= c.MaxHeight {
		return
	}
	thumbH := c.MaxHeight * barH / totalH
	if thumbH < 16 {
		thumbH = 16
	}
	scrollRange := barH - thumbH
	thumbY := y
	if totalH > c.MaxHeight && scrollRange > 0 {
		thumbY = y + (scrollRange*c.ScrollY)/(totalH-c.MaxHeight)
	}
	if c.vBarBg == nil || c.vBarBg.Bounds().Dx() != barW || c.vBarBg.Bounds().Dy() != barH {
		c.vBarBg = ebiten.NewImage(barW, barH)
		c.vBarBg.Fill(color.RGBA{60, 60, 60, 180})
	}
	c.Op.GeoM.Reset()
	c.Op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(c.vBarBg, c.Op)
	if c.vThumb == nil || c.vThumb.Bounds().Dx() != barW || c.vThumb.Bounds().Dy() != thumbH {
		c.vThumb = ebiten.NewImage(barW, thumbH)
		c.vThumb.Fill(color.RGBA{160, 160, 160, 220})
	}
	c.Op.GeoM.Reset()
	c.Op.GeoM.Translate(float64(x), float64(thumbY))
	screen.DrawImage(c.vThumb, c.Op)
}

func (c *AbsolutePositionContainer) drawHScrollbar(screen *ebiten.Image, x, y, width int) {
	barH := 8
	barW := width
	totalW := c.contentWidth()
	if totalW <= c.MaxWidth {
		return
	}
	thumbW := c.MaxWidth * barW / totalW
	if thumbW < 16 {
		thumbW = 16
	}
	scrollRange := barW - thumbW
	thumbX := x
	if totalW > c.MaxWidth && scrollRange > 0 {
		thumbX = x + (scrollRange*c.ScrollX)/(totalW-c.MaxWidth)
	}
	if c.hBarBg == nil || c.hBarBg.Bounds().Dx() != barW || c.hBarBg.Bounds().Dy() != barH {
		c.hBarBg = ebiten.NewImage(barW, barH)
		c.hBarBg.Fill(color.RGBA{60, 60, 60, 180})
	}
	c.Op.GeoM.Reset()
	c.Op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(c.hBarBg, c.Op)
	if c.hThumb == nil || c.hThumb.Bounds().Dx() != thumbW || c.hThumb.Bounds().Dy() != barH {
		c.hThumb = ebiten.NewImage(thumbW, barH)
		c.hThumb.Fill(color.RGBA{160, 160, 160, 220})
	}
	c.Op.GeoM.Reset()
	c.Op.GeoM.Translate(float64(thumbX), float64(y))
	screen.DrawImage(c.hThumb, c.Op)
}

func (c *AbsolutePositionContainer) GetAbsolutePosition() (int, int) {
	return c.X, c.Y
}
