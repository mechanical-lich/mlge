package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// GridContainer arranges child elements in a grid pattern.
type GridContainer struct {
	ElementBase
	Children   map[string]ElementInterface
	Columns    int
	CellWidth  int
	CellHeight int
	SpacingX   int
	SpacingY   int
}

// NewGridContainer creates a new grid container.
func NewGridContainer(name string, x, y, columns, cellW, cellH, spacingX, spacingY int) *GridContainer {
	return &GridContainer{
		ElementBase: ElementBase{
			Name:   name,
			X:      x,
			Y:      y,
			Width:  0, // can be set later if needed
			Height: 0,
			op:     &ebiten.DrawImageOptions{},
		},
		Children:   make(map[string]ElementInterface),
		Columns:    columns,
		CellWidth:  cellW,
		CellHeight: cellH,
		SpacingX:   spacingX,
		SpacingY:   spacingY,
	}
}

// AddChild adds a child element to the grid.
func (g *GridContainer) AddChild(child ElementInterface) {
	g.Children[child.GetName()] = child
}

// GetChild retrieves a child by name.
func (g *GridContainer) GetChild(name string) ElementInterface {
	return g.Children[name]
}

// RemoveChild removes a child by name.
func (g *GridContainer) RemoveChild(name string) {
	delete(g.Children, name)
}

// Update positions children and updates them.
func (g *GridContainer) Update(parentX, parentY int) {
	i := 0
	for _, child := range g.Children {
		col := i % g.Columns
		row := i / g.Columns
		x := g.X + col*(g.CellWidth+g.SpacingX)
		y := g.Y + row*(g.CellHeight+g.SpacingY)
		child.SetPosition(x, y)
		child.Update(parentX, parentY)
		i++
	}
}

// Draw draws all children.
func (g *GridContainer) Draw(screen *ebiten.Image, parentX, parentY int, theme *Theme) {
	for _, child := range g.Children {
		child.Draw(screen, parentX, parentY, theme)
	}
}
