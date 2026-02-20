package world

import (
	"github.com/mechanical-lich/mlge/utility"
)

// SetTileTypeAndVariant sets the Type and Variant on a tile by looking up
// the tile name in the provided registry.
func SetTileTypeAndVariant[TD any, T any](tile *Tile[T], registry *TileRegistry[TD], tileName string, variant int) {
	_, idx := registry.GetByName(tileName)
	if idx < 0 {
		return
	}
	tile.Type = idx
	tile.Variant = variant
}

// RandomVariant returns a random variant index for the named tile type.
// Returns 0 if the tile name is not found or has no variants.
func RandomVariant[TD any](registry *TileRegistry[TD], tileName string) int {
	def, _ := registry.GetByName(tileName)
	if def == nil || len(def.Variants) == 0 {
		return 0
	}
	return utility.GetRandom(0, len(def.Variants)-1)
}

// CreateTileCluster sets a cluster of tiles around (x, y, z) to the given type
// with a random walk pattern. Size determines how many tiles are placed.
func CreateTileCluster[TD any, T any](level *Level[TD, T], x, y, z, size int, tileName string, variant int) {
	_, idx := level.Registry.GetByName(tileName)
	if idx < 0 {
		return
	}

	cx, cy := x, y
	for i := 0; i < size; i++ {
		if level.InBounds(cx, cy, z) {
			tile := level.GetTileAt(cx, cy, z)
			tile.Type = idx
			tile.Variant = variant
		}

		// Random walk
		dir := utility.GetRandom(0, 3)
		switch dir {
		case 0:
			cx++
		case 1:
			cx--
		case 2:
			cy++
		case 3:
			cy--
		}
	}
}
