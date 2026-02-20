package world

import "fmt"

// Level represents a tile-based world grid. It stores tiles in a flat 1D array
// of size Width * Height * Depth. For 2D games, use Depth=1.
//
// TD is the type parameter for game-specific TileDefinition custom data.
// T is the type parameter for game-specific per-tile custom data.
type Level[TD any, T any] struct {
	tiles    []Tile[T]
	Width    int
	Height   int
	Depth    int
	Registry *TileRegistry[TD]

	// CustomData allows games to attach arbitrary level-wide state
	// (e.g., time of day, weather, turn counter).
	CustomData any
}

// NewLevel creates a new level with the given dimensions and tile registry.
// All tiles are initialized with their X, Y, Z coordinates and a back-pointer
// to the level for neighbor lookups.
func NewLevel[TD any, T any](width, height, depth int, registry *TileRegistry[TD]) *Level[TD, T] {
	if depth < 1 {
		depth = 1
	}

	level := &Level[TD, T]{
		tiles:    make([]Tile[T], width*height*depth),
		Width:    width,
		Height:   height,
		Depth:    depth,
		Registry: registry,
	}

	// Initialize tile coordinates and back-pointers; all tiles start empty.
	for z := 0; z < depth; z++ {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				idx := level.tileIndex(x, y, z)
				level.tiles[idx].Type = -1 // empty by default; IsEmpty() checks Type < 0
				level.tiles[idx].X = x
				level.tiles[idx].Y = y
				level.tiles[idx].Z = z
				level.tiles[idx].levelRef = level
			}
		}
	}

	return level
}

// tileIndex computes the flat array index for a 3D coordinate.
func (l *Level[TD, T]) tileIndex(x, y, z int) int {
	return x + y*l.Width + z*l.Width*l.Height
}

// InBounds returns true if the given coordinates are within the level bounds.
func (l *Level[TD, T]) InBounds(x, y, z int) bool {
	return x >= 0 && x < l.Width && y >= 0 && y < l.Height && z >= 0 && z < l.Depth
}

// GetTileAt returns a pointer to the tile at the given coordinates.
// Returns nil if the coordinates are out of bounds.
func (l *Level[TD, T]) GetTileAt(x, y, z int) *Tile[T] {
	if !l.InBounds(x, y, z) {
		return nil
	}
	return &l.tiles[l.tileIndex(x, y, z)]
}

// GetTileIndex returns a pointer to the tile at the given flat array index.
// Returns nil if the index is out of range.
func (l *Level[TD, T]) GetTileIndex(index int) *Tile[T] {
	if index < 0 || index >= len(l.tiles) {
		return nil
	}
	return &l.tiles[index]
}

// SetTileAt sets the tile type and variant at the given coordinates by name.
// Returns an error if the tile name is not found in the registry or coordinates
// are out of bounds.
func (l *Level[TD, T]) SetTileAt(x, y, z int, tileType string, variant int) error {
	if !l.InBounds(x, y, z) {
		return fmt.Errorf("world: coordinates (%d, %d, %d) out of bounds", x, y, z)
	}

	_, idx := l.Registry.GetByName(tileType)
	if idx < 0 {
		return fmt.Errorf("world: unknown tile type %q", tileType)
	}

	tile := &l.tiles[l.tileIndex(x, y, z)]
	tile.Type = idx
	tile.Variant = variant
	return nil
}

// ForEachTile iterates over every tile in the level, calling fn for each.
func (l *Level[TD, T]) ForEachTile(fn func(*Tile[T])) {
	for i := range l.tiles {
		fn(&l.tiles[i])
	}
}

// GetView extracts a 2D slice of tiles at a given Z level, centered on or
// starting from (aX, aY), with the specified viewport dimensions.
// If centered is true, (aX, aY) is the center of the view; otherwise it's
// the top-left corner.
func (l *Level[TD, T]) GetView(aX, aY, aZ, width, height int, centered bool) [][]*Tile[T] {
	startX := aX
	startY := aY
	if centered {
		startX = aX - width/2
		startY = aY - height/2
	}

	view := make([][]*Tile[T], height)
	for row := 0; row < height; row++ {
		view[row] = make([]*Tile[T], width)
		for col := 0; col < width; col++ {
			tx := startX + col
			ty := startY + row
			if l.InBounds(tx, ty, aZ) {
				view[row][col] = &l.tiles[l.tileIndex(tx, ty, aZ)]
			}
		}
	}

	return view
}

// TileDefinitionAt returns the tile definition for the tile at the given
// coordinates. Returns nil if out of bounds or the definition index is invalid.
func (l *Level[TD, T]) TileDefinitionAt(x, y, z int) *TileDefinition[TD] {
	tile := l.GetTileAt(x, y, z)
	if tile == nil {
		return nil
	}
	return l.Registry.GetByIndex(tile.Type)
}

// NewPathConfig creates a PathConfig with the given functions for use with
// pathfinding. If estimatedCost is nil, Manhattan distance is used as default.
func (l *Level[TD, T]) NewPathConfig(
	isPassable func(*Tile[T]) bool,
	cost func(from, to *Tile[T]) float64,
	estimatedCost func(from, to *Tile[T]) float64,
) *PathConfig[T] {
	return &PathConfig[T]{
		IsPassable:    isPassable,
		Cost:          cost,
		EstimatedCost: estimatedCost,
	}
}

// GetPathableTile wraps the tile at the given coordinates with a PathConfig
// so it can be used with path.Path(). Returns nil if coordinates are out of bounds.
func (l *Level[TD, T]) GetPathableTile(x, y, z int, config *PathConfig[T]) *PathableTile[T] {
	tile := l.GetTileAt(x, y, z)
	if tile == nil {
		return nil
	}
	return &PathableTile[T]{
		Tile:   tile,
		config: config,
	}
}

// TileCount returns the total number of tiles in the level.
func (l *Level[TD, T]) TileCount() int {
	return len(l.tiles)
}

// ClearTileAt resets the tile at the given coordinates to empty (Type = -1, Variant = 0).
// Returns an error if coordinates are out of bounds.
func (l *Level[TD, T]) ClearTileAt(x, y, z int) error {
	if !l.InBounds(x, y, z) {
		return fmt.Errorf("world: coordinates (%d, %d, %d) out of bounds", x, y, z)
	}
	tile := &l.tiles[l.tileIndex(x, y, z)]
	tile.Type = -1
	tile.Variant = 0
	return nil
}

// --- tileLevel interface implementation ---
// These unexported methods let Tile access the Level for neighbor lookups
// without needing the full generic type in the Tile struct.

func (l *Level[TD, T]) inBounds(x, y, z int) bool {
	return l.InBounds(x, y, z)
}

func (l *Level[TD, T]) tileAt(x, y, z int) tileAccessor {
	t := l.GetTileAt(x, y, z)
	if t == nil {
		return nil
	}
	return t
}

func (l *Level[TD, T]) levelWidth() int  { return l.Width }
func (l *Level[TD, T]) levelHeight() int { return l.Height }
func (l *Level[TD, T]) levelDepth() int  { return l.Depth }
