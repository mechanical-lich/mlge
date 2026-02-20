---
layout: default
title: World
nav_order: 17
---

# World

`github.com/mechanical-lich/mlge/world`

A generic, tile-based world system using Go generics. Two type parameters let each game attach custom data to tile definitions and individual tiles without forking the package.

- **`TD`** — Custom data on tile *definitions* (e.g., `Door bool`, `Wall bool`, `MixGroup int`)
- **`T`** — Custom data on individual *tiles* (e.g., `Explored bool`, `Temperature float64`)

The package handles tile storage (3D grid with Depth=1 for 2D games), tile definition registries, A* pathfinding integration, and save/load — but intentionally excludes drawing and entity management, which are game-specific.

## Quick Start

```go
import "github.com/mechanical-lich/mlge/world"

// 1. Define your game-specific data types
type MyTileDef struct {
    Door bool `json:"door"`
    Wall bool `json:"wall"`
}

type MyTileExtra struct {
    Explored bool `json:"explored"`
}

// 2. Create a tile registry and load definitions
registry := world.NewTileRegistry[MyTileDef]()
err := registry.LoadFromFile("data/tile_definitions.json")

// 3. Create a level (10x10, single Z layer for 2D)
level := world.NewLevel[MyTileDef, MyTileExtra](10, 10, 1, registry)

// 4. Use it — type params are inferred on all method calls
_ = level.SetTileAt(5, 5, 0, "stone_wall", 0)
tile := level.GetTileAt(5, 5, 0)
tile.CustomData.Explored = true

def := level.TileDefinitionAt(5, 5, 0)
if def.CustomData.Door {
    // game-specific logic
}
```

## Tile Definitions

### TileVariant

```go
type TileVariant struct {
    Variant int `json:"variant"`
    SpriteX int `json:"spriteX"`
    SpriteY int `json:"spriteY"`
}
```

A single visual variant pointing to a sprite location.

### TileDefinition[TD]

```go
type TileDefinition[TD any] struct {
    Name       string        `json:"name"`
    Solid      bool          `json:"solid"`
    Water      bool          `json:"water"`
    Variants   []TileVariant `json:"variants"`
    CustomData TD            `json:"customData"`
}
```

Defines a type of tile. `Solid` and `Water` are common enough to be built-in; everything else goes in `CustomData`.

### TileRegistry[TD]

Holds the loaded definitions and provides name/index lookups.

```go
type TileRegistry[TD any] struct {
    Definitions []TileDefinition[TD]
    NameToIndex map[string]int
    IndexToName []string
}
```

**Methods:**

| Method | Signature | Description |
|--------|-----------|-------------|
| `NewTileRegistry` | `NewTileRegistry[TD any]() *TileRegistry[TD]` | Creates an empty registry |
| `LoadFromFile` | `(path string) error` | Reads a JSON array of tile definitions |
| `LoadFromDefinitions` | `(defs []TileDefinition[TD])` | Populates from a pre-built slice |
| `GetByName` | `(name string) (*TileDefinition[TD], int)` | Lookup by name; returns `nil, -1` if not found |
| `GetByIndex` | `(index int) *TileDefinition[TD]` | Lookup by index; returns `nil` if out of range |

### JSON Format

```json
[
  {
    "name": "grass",
    "solid": false,
    "water": false,
    "variants": [
      { "variant": 0, "spriteX": 0, "spriteY": 0 },
      { "variant": 1, "spriteX": 1, "spriteY": 0 }
    ],
    "customData": { "door": false, "wall": false }
  },
  {
    "name": "stone_wall",
    "solid": true,
    "water": false,
    "variants": [
      { "variant": 0, "spriteX": 3, "spriteY": 0 }
    ],
    "customData": { "door": false, "wall": true }
  }
]
```

## Tile

```go
type Tile[T any] struct {
    Type       int `json:"type"`       // index into TileRegistry.Definitions; -1 = empty
    Variant    int `json:"variant"`    // variant index
    X, Y, Z   int                     // grid coordinates
    CustomData T   `json:"customData"` // game-specific per-tile data
}
```

Individual tiles are stored by value in a flat array — no pointer indirection. Access them via `Level.GetTileAt()` which returns a pointer into the array.

### IsEmpty

```go
func (t *Tile[T]) IsEmpty() bool
```

Returns `true` when `Type < 0` (no tile type has been assigned). Newly created levels start all tiles empty. Use this instead of comparing `Type` to a sentinel constant:

```go
if tile.IsEmpty() {
    // skip unset tiles
}
```

## Level

### Level[TD, T]

```go
type Level[TD any, T any] struct {
    Width, Height, Depth int
    Registry             *TileRegistry[TD]
    CustomData           any  // arbitrary level-wide state (time of day, weather, etc.)
}
```

Tiles are stored internally in a flat 1D array of size `Width × Height × Depth`. For 2D games, use `Depth=1`. Coordinates are `(x, y, z)` where `z` is the vertical layer.

### Construction

```go
level := world.NewLevel[MyTileDef, MyTileExtra](width, height, depth, registry)
```

All tiles are initialized with correct `X, Y, Z` coordinates and internal back-pointers for neighbor lookups. Every tile starts empty — `tile.IsEmpty()` returns `true` until a tile type is assigned via `SetTileAt`.

### Tile Access

| Method | Signature | Description |
|--------|-----------|-------------|
| `GetTileAt` | `(x, y, z int) *Tile[T]` | Returns tile pointer, or `nil` if out of bounds |
| `SetTileAt` | `(x, y, z int, tileType string, variant int) error` | Sets tile type by name; errors on unknown type or OOB |
| `ClearTileAt` | `(x, y, z int) error` | Resets tile to empty (`Type=-1`, `Variant=0`); errors if OOB |
| `GetTileIndex` | `(index int) *Tile[T]` | Direct flat-array access by index |
| `TileDefinitionAt` | `(x, y, z int) *TileDefinition[TD]` | Returns the definition for the tile at the given position |
| `InBounds` | `(x, y, z int) bool` | Bounds check |
| `TileCount` | `() int` | Total number of tiles |

### Iteration

```go
level.ForEachTile(func(tile *world.Tile[MyTileExtra]) {
    tile.CustomData.Explored = false
})
```

### Viewport

```go
// Get a 2D slice of tiles at Z=0, starting at (cameraX, cameraY)
view := level.GetView(cameraX, cameraY, 0, viewWidth, viewHeight, false)

// Centered on the player
view := level.GetView(playerX, playerY, 0, viewWidth, viewHeight, true)
```

Returns `[][]*Tile[T]` — out-of-bounds positions are `nil`.

## Pathfinding

The world package integrates with `mlge/path` through a `PathableTile` wrapper. Games provide callback functions for passability, cost, and heuristics.

### PathConfig[T]

```go
type PathConfig[T any] struct {
    IsPassable         func(tile *Tile[T]) bool
    Cost               func(from, to *Tile[T]) float64
    EstimatedCost      func(from, to *Tile[T]) float64
    Include3DNeighbors bool
}
```

| Field | Description |
|-------|-------------|
| `IsPassable` | Returns true if the tile can be traversed |
| `Cost` | Movement cost between adjacent tiles (default: 1.0) |
| `EstimatedCost` | Heuristic estimate to target (default: Manhattan distance) |
| `Include3DNeighbors` | When true, also considers Z±1 as neighbors |

### Usage

{% raw %}
```go
config := level.NewPathConfig(
    func(tile *world.Tile[MyTileExtra]) bool {
        if tile.IsEmpty() {
            return false // empty tiles are not passable
        }
        def := level.Registry.GetByIndex(tile.Type)
        return def != nil && !def.Solid
    },
    func(from, to *world.Tile[MyTileExtra]) float64 {
        if level.Registry.GetByIndex(to.Type).Water {
            return 3.0 // water tiles cost more
        }
        return 1.0
    },
    nil, // default Manhattan distance heuristic
)

from := level.GetPathableTile(startX, startY, 0, config)
to := level.GetPathableTile(goalX, goalY, 0, config)

result, distance, found := path.Path(from, to)
if found {
    for _, node := range result {
        pt := node.(*world.PathableTile[MyTileExtra])
        fmt.Printf("Step: %d, %d\n", pt.Tile.X, pt.Tile.Y)
    }
}
```
{% endraw %}

### 3D Pathfinding

For games with multiple Z levels (dungeons, multi-story buildings):

```go
config.Include3DNeighbors = true
```

This adds Z±1 neighbors to the pathfinding graph alongside the 4 cardinal directions.

## Utilities

### SetTileTypeAndVariant

```go
world.SetTileTypeAndVariant(tile, registry, "grass", 0)
```

Sets the `Type` and `Variant` on a tile by looking up the name in the registry.

### RandomVariant

```go
variant := world.RandomVariant(registry, "grass")
```

Returns a random variant index for the named tile type.

### CreateTileCluster

```go
world.CreateTileCluster(level, x, y, z, size, "flowers", 0)
```

Places a random-walk cluster of tiles around the given position. Useful for terrain generation (flower patches, ore veins, etc.).

## Save / Load

The package provides JSON-based save/load that preserves all tile data including custom data. Games can attach additional save state (entities, settlements, etc.) via a `CustomData` field.

### SaveData[T]

```go
type SaveData[T any] struct {
    Width      int       `json:"width"`
    Height     int       `json:"height"`
    Depth      int       `json:"depth"`
    Tiles      []Tile[T] `json:"tiles"`
    CustomData any       `json:"customData,omitempty"`
}
```

### Saving

```go
// Save with game-specific extra data
gameState := map[string]any{"hour": 14, "day": 7}
err := world.SaveToFile(level, "save.json", gameState)

// Or get a snapshot without writing to disk
snapshot := world.SaveLevel(level)
```

`SaveLevel` returns a copy — modifying the snapshot won't affect the live level.

### Loading

```go
level, rawCustomData, err := world.LoadFromFile[MyTileDef, MyTileExtra]("save.json", registry)

// rawCustomData is json.RawMessage — unmarshal into your game's type
var gameState MyGameState
json.Unmarshal(rawCustomData, &gameState)
```

The loaded level has all tile coordinates and internal back-pointers restored, so pathfinding works immediately.

## Design Notes

- **No drawing** — Games implement rendering using `GetTileAt()` / `GetView()` with their own camera and sprite systems.
- **No entity storage** — Games manage their own entity lists and spatial indexes; the world package is tiles-only.
- **Flat 1D array** — Tiles are stored densely in `Width × Height × Depth` order. Every position has a tile (no sparse storage).
- **Generic registry instead of package vars** — Go doesn't support generic package-level variables, so `TileRegistry[TD]` is an explicit struct. This also allows multiple registries if needed.
- **`PathableTile` wrapper** — Rather than requiring `Tile` to implement `path.Pather` directly (which would bake in game-specific logic), a wrapper with function hooks keeps pathfinding configurable.
