package world

import (
	"github.com/mechanical-lich/mlge/path"
	"github.com/mechanical-lich/mlge/utility"
)

// Tile represents a single tile in the world grid.
// T is a generic type parameter for game-specific per-tile custom data.
type Tile[T any] struct {
	Type       int `json:"type"`
	Variant    int `json:"variant"`
	X          int `json:"x"`
	Y          int `json:"y"`
	Z          int `json:"z"`
	CustomData T   `json:"customData"`

	// unexported back-pointer to the owning level, used for neighbor lookups.
	// Excluded from JSON serialization.
	levelRef tileLevel
}

// tileLevel is an interface that Tile uses internally to access its level
// for neighbor lookups and bounds checking, without needing the full generic
// Level[TD, T] type in the Tile struct (which would create a circular generic ref).
type tileLevel interface {
	inBounds(x, y, z int) bool
	tileAt(x, y, z int) tileAccessor
	levelWidth() int
	levelHeight() int
	levelDepth() int
}

// tileAccessor provides access to a tile without exposing the full generic type.
type tileAccessor interface {
	tileX() int
	tileY() int
	tileZ() int
}

func (t *Tile[T]) tileX() int { return t.X }
func (t *Tile[T]) tileY() int { return t.Y }
func (t *Tile[T]) tileZ() int { return t.Z }

// IsEmpty reports whether this tile has no type assigned.
// By convention, Type < 0 means empty (unset/air). NewLevel initializes
// all tiles this way; SetTileAt and ClearTileAt maintain the invariant.
func (t *Tile[T]) IsEmpty() bool { return t.Type < 0 }

// ---------------------------------------------------------------------------
// Pathfinding
// ---------------------------------------------------------------------------

// PathConfig holds game-specific functions that customize pathfinding behavior
// for tiles. Games provide these functions when they need A* pathfinding.
type PathConfig[T any] struct {
	// IsPassable returns true if the tile can be traversed.
	IsPassable func(tile *Tile[T]) bool

	// Cost returns the movement cost from one tile to an adjacent neighbor.
	Cost func(from, to *Tile[T]) float64

	// EstimatedCost returns the heuristic estimate from one tile to another
	// (typically Manhattan or Euclidean distance). Must be admissible.
	EstimatedCost func(from, to *Tile[T]) float64

	// Include3DNeighbors, when true, also considers tiles above/below (Z±1)
	// as neighbors for pathfinding.
	Include3DNeighbors bool
}

// PathableTile wraps a Tile to implement the path.Pather interface.
// It delegates passability, cost, and heuristic to the PathConfig functions.
type PathableTile[T any] struct {
	Tile   *Tile[T]
	config *PathConfig[T]
}

// Verify PathableTile implements path.Pather at compile time.
var _ path.Pather = (*PathableTile[struct{}])(nil)

// PathID returns a unique integer ID for this tile based on its grid position.
func (pt *PathableTile[T]) PathID() int {
	level := pt.Tile.levelRef
	w := level.levelWidth()
	h := level.levelHeight()
	return pt.Tile.Z*w*h + pt.Tile.Y*w + pt.Tile.X
}

// PathNeighborsAppend appends passable cardinal neighbors (and optionally
// vertical neighbors) to the provided slice and returns it.
func (pt *PathableTile[T]) PathNeighborsAppend(neighbors []path.Pather) []path.Pather {
	level := pt.Tile.levelRef
	x, y, z := pt.Tile.X, pt.Tile.Y, pt.Tile.Z

	// Cardinal directions (4-way)
	deltas := [][3]int{
		{x - 1, y, z},
		{x + 1, y, z},
		{x, y - 1, z},
		{x, y + 1, z},
	}

	// Optionally include vertical neighbors
	if pt.config.Include3DNeighbors {
		deltas = append(deltas, [3]int{x, y, z - 1}, [3]int{x, y, z + 1})
	}

	for _, d := range deltas {
		if !level.inBounds(d[0], d[1], d[2]) {
			continue
		}
		accessor := level.tileAt(d[0], d[1], d[2])
		if accessor == nil {
			continue
		}
		// We need to get the actual *Tile[T] from the level. Since we go through
		// the tileLevel interface, we cast back.
		neighborTile := accessor.(*Tile[T])
		if pt.config.IsPassable != nil && !pt.config.IsPassable(neighborTile) {
			continue
		}
		neighbors = append(neighbors, &PathableTile[T]{
			Tile:   neighborTile,
			config: pt.config,
		})
	}

	return neighbors
}

// PathNeighborCost returns the cost of moving from this tile to a neighbor.
func (pt *PathableTile[T]) PathNeighborCost(to path.Pather) float64 {
	if pt.config.Cost != nil {
		return pt.config.Cost(pt.Tile, to.(*PathableTile[T]).Tile)
	}
	return 1.0
}

// PathEstimatedCost returns the heuristic estimated cost to the target.
func (pt *PathableTile[T]) PathEstimatedCost(to path.Pather) float64 {
	if pt.config.EstimatedCost != nil {
		return pt.config.EstimatedCost(pt.Tile, to.(*PathableTile[T]).Tile)
	}
	// Default: Manhattan distance
	other := to.(*PathableTile[T]).Tile
	return float64(utility.Abs(pt.Tile.X-other.X) + utility.Abs(pt.Tile.Y-other.Y) + utility.Abs(pt.Tile.Z-other.Z))
}
