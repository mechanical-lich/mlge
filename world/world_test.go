package world

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/mechanical-lich/mlge/path"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Test custom data types — simulating what a game would define
// ---------------------------------------------------------------------------

type testTileDefData struct {
	Door bool `json:"door"`
	Wall bool `json:"wall"`
}

type testTileData struct {
	Visited bool `json:"visited"`
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func testRegistry() *TileRegistry[testTileDefData] {
	reg := NewTileRegistry[testTileDefData]()
	reg.LoadFromDefinitions([]TileDefinition[testTileDefData]{
		{
			Name:     "grass",
			Solid:    false,
			Water:    false,
			Variants: []TileVariant{{Variant: 0, SpriteX: 0, SpriteY: 0}},
		},
		{
			Name:       "stone",
			Solid:      true,
			Water:      false,
			CustomData: testTileDefData{Wall: true},
		},
		{
			Name:  "water",
			Solid: false,
			Water: true,
		},
	})
	return reg
}

// ---------------------------------------------------------------------------
// TileRegistry tests
// ---------------------------------------------------------------------------

func TestNewTileRegistry(t *testing.T) {
	reg := NewTileRegistry[testTileDefData]()
	assert.NotNil(t, reg)
	assert.Empty(t, reg.Definitions)
	assert.NotNil(t, reg.NameToIndex)
	assert.Empty(t, reg.IndexToName)
}

func TestLoadFromDefinitions(t *testing.T) {
	reg := NewTileRegistry[testTileDefData]()
	reg.LoadFromDefinitions([]TileDefinition[testTileDefData]{
		{Name: "grass", Solid: false, Water: false, Variants: []TileVariant{{Variant: 0, SpriteX: 0, SpriteY: 0}}},
		{Name: "stone", Solid: true, Water: false, CustomData: testTileDefData{Wall: true}},
		{Name: "water", Solid: false, Water: true},
	})

	assert.Len(t, reg.Definitions, 3)
	assert.Equal(t, 0, reg.NameToIndex["grass"])
	assert.Equal(t, 1, reg.NameToIndex["stone"])
	assert.Equal(t, 2, reg.NameToIndex["water"])
	assert.Equal(t, "grass", reg.IndexToName[0])
}

func TestGetByName(t *testing.T) {
	reg := testRegistry()

	def, idx := reg.GetByName("stone")
	assert.Equal(t, 1, idx)
	assert.True(t, def.Solid)
	assert.True(t, def.CustomData.Wall)

	def, idx = reg.GetByName("nonexistent")
	assert.Nil(t, def)
	assert.Equal(t, -1, idx)
}

func TestGetByIndex(t *testing.T) {
	reg := testRegistry()

	def := reg.GetByIndex(0)
	assert.Equal(t, "grass", def.Name)

	def = reg.GetByIndex(999)
	assert.Nil(t, def)

	def = reg.GetByIndex(-1)
	assert.Nil(t, def)
}

func TestLoadFromFile(t *testing.T) {
	defs := []TileDefinition[testTileDefData]{
		{Name: "dirt", Solid: false, Variants: []TileVariant{{Variant: 0, SpriteX: 1, SpriteY: 2}}},
		{Name: "wall", Solid: true, CustomData: testTileDefData{Wall: true}},
	}
	data, err := json.Marshal(defs)
	require.NoError(t, err)

	tmpFile := filepath.Join(t.TempDir(), "tile_defs.json")
	require.NoError(t, os.WriteFile(tmpFile, data, 0644))

	reg := NewTileRegistry[testTileDefData]()
	err = reg.LoadFromFile(tmpFile)
	require.NoError(t, err)

	assert.Len(t, reg.Definitions, 2)
	assert.Equal(t, "dirt", reg.Definitions[0].Name)
	assert.True(t, reg.Definitions[1].CustomData.Wall)
	assert.Equal(t, 1, reg.Definitions[0].Variants[0].SpriteX)
}

// ---------------------------------------------------------------------------
// Level tests
// ---------------------------------------------------------------------------

func TestNewLevel(t *testing.T) {
	reg := testRegistry()
	level := NewLevel[testTileDefData, testTileData](10, 8, 3, reg)

	assert.Equal(t, 10, level.Width)
	assert.Equal(t, 8, level.Height)
	assert.Equal(t, 3, level.Depth)
	assert.Equal(t, 10*8*3, level.TileCount())
}

func TestNewLevel_MinDepth(t *testing.T) {
	reg := testRegistry()
	level := NewLevel[testTileDefData, testTileData](5, 5, 0, reg)
	assert.Equal(t, 1, level.Depth, "Depth should be clamped to 1")
}

func TestInBounds(t *testing.T) {
	reg := testRegistry()
	level := NewLevel[testTileDefData, testTileData](10, 10, 2, reg)

	assert.True(t, level.InBounds(0, 0, 0))
	assert.True(t, level.InBounds(9, 9, 1))
	assert.False(t, level.InBounds(-1, 0, 0))
	assert.False(t, level.InBounds(10, 0, 0))
	assert.False(t, level.InBounds(0, 10, 0))
	assert.False(t, level.InBounds(0, 0, 2))
	assert.False(t, level.InBounds(0, 0, -1))
}

func TestGetTileAt(t *testing.T) {
	reg := testRegistry()
	level := NewLevel[testTileDefData, testTileData](10, 10, 1, reg)

	tile := level.GetTileAt(3, 4, 0)
	require.NotNil(t, tile)
	assert.Equal(t, 3, tile.X)
	assert.Equal(t, 4, tile.Y)
	assert.Equal(t, 0, tile.Z)

	assert.Nil(t, level.GetTileAt(-1, 0, 0))
	assert.Nil(t, level.GetTileAt(10, 0, 0))
}

func TestSetTileAt(t *testing.T) {
	reg := testRegistry()
	level := NewLevel[testTileDefData, testTileData](10, 10, 1, reg)

	err := level.SetTileAt(5, 5, 0, "stone", 0)
	require.NoError(t, err)

	tile := level.GetTileAt(5, 5, 0)
	assert.Equal(t, 1, tile.Type)

	err = level.SetTileAt(5, 5, 0, "lava", 0)
	assert.Error(t, err)

	err = level.SetTileAt(100, 100, 0, "grass", 0)
	assert.Error(t, err)
}

func TestTileDefinitionAt(t *testing.T) {
	reg := testRegistry()
	level := NewLevel[testTileDefData, testTileData](10, 10, 1, reg)

	_ = level.SetTileAt(2, 2, 0, "water", 0)
	def := level.TileDefinitionAt(2, 2, 0)
	require.NotNil(t, def)
	assert.True(t, def.Water)
	assert.Equal(t, "water", def.Name)
}

func TestForEachTile(t *testing.T) {
	reg := testRegistry()
	level := NewLevel[testTileDefData, testTileData](3, 3, 1, reg)

	count := 0
	level.ForEachTile(func(tile *Tile[testTileData]) {
		count++
	})
	assert.Equal(t, 9, count)
}

func TestGetView(t *testing.T) {
	reg := testRegistry()
	level := NewLevel[testTileDefData, testTileData](20, 20, 1, reg)

	view := level.GetView(5, 5, 0, 3, 3, false)
	assert.Len(t, view, 3)
	assert.Len(t, view[0], 3)
	assert.Equal(t, 5, view[0][0].X)
	assert.Equal(t, 5, view[0][0].Y)

	view = level.GetView(10, 10, 0, 5, 5, true)
	assert.Equal(t, 8, view[0][0].X)
	assert.Equal(t, 8, view[0][0].Y)

	view = level.GetView(0, 0, 0, 3, 3, true)
	assert.Nil(t, view[0][0])
}

func TestTileCustomData(t *testing.T) {
	reg := testRegistry()
	level := NewLevel[testTileDefData, testTileData](5, 5, 1, reg)

	tile := level.GetTileAt(2, 2, 0)
	tile.CustomData.Visited = true

	same := level.GetTileAt(2, 2, 0)
	assert.True(t, same.CustomData.Visited)
}

func TestLevelCustomData(t *testing.T) {
	reg := testRegistry()
	level := NewLevel[testTileDefData, testTileData](5, 5, 1, reg)
	level.CustomData = map[string]int{"hour": 6, "day": 42}

	cd := level.CustomData.(map[string]int)
	assert.Equal(t, 6, cd["hour"])
}

// ---------------------------------------------------------------------------
// Pathfinding tests
// ---------------------------------------------------------------------------

func TestPathableTile_ImplementsPather(t *testing.T) {
	reg := testRegistry()
	level := NewLevel[testTileDefData, testTileData](5, 5, 1, reg)

	config := level.NewPathConfig(
		func(tile *Tile[testTileData]) bool { return true },
		func(from, to *Tile[testTileData]) float64 { return 1.0 },
		nil,
	)

	pt := level.GetPathableTile(0, 0, 0, config)
	require.NotNil(t, pt)

	var _ path.Pather = pt
}

func TestPathableTile_PathID(t *testing.T) {
	reg := testRegistry()
	level := NewLevel[testTileDefData, testTileData](10, 10, 2, reg)

	config := &PathConfig[testTileData]{}

	pt00 := level.GetPathableTile(0, 0, 0, config)
	assert.Equal(t, 0, pt00.PathID())

	pt32 := level.GetPathableTile(3, 2, 0, config)
	assert.Equal(t, 3+2*10, pt32.PathID())

	pt_z1 := level.GetPathableTile(1, 1, 1, config)
	assert.Equal(t, 1+1*10+1*10*10, pt_z1.PathID())
}

func TestPathableTile_Neighbors(t *testing.T) {
	reg := testRegistry()
	level := NewLevel[testTileDefData, testTileData](5, 5, 1, reg)

	config := level.NewPathConfig(
		func(tile *Tile[testTileData]) bool { return true },
		func(from, to *Tile[testTileData]) float64 { return 1.0 },
		nil,
	)

	pt := level.GetPathableTile(2, 2, 0, config)
	neighbors := pt.PathNeighborsAppend(nil)
	assert.Len(t, neighbors, 4)

	pt = level.GetPathableTile(0, 0, 0, config)
	neighbors = pt.PathNeighborsAppend(nil)
	assert.Len(t, neighbors, 2)

	pt = level.GetPathableTile(2, 0, 0, config)
	neighbors = pt.PathNeighborsAppend(nil)
	assert.Len(t, neighbors, 3)
}

func TestPathableTile_BlockedNeighbors(t *testing.T) {
	reg := testRegistry()
	level := NewLevel[testTileDefData, testTileData](5, 5, 1, reg)

	_ = level.SetTileAt(3, 2, 0, "stone", 0)
	_ = level.SetTileAt(1, 2, 0, "stone", 0)

	config := level.NewPathConfig(
		func(tile *Tile[testTileData]) bool {
			if tile.IsEmpty() {
				return true // empty tiles are passable (air)
			}
			def := reg.GetByIndex(tile.Type)
			return def != nil && !def.Solid
		},
		func(from, to *Tile[testTileData]) float64 { return 1.0 },
		nil,
	)

	pt := level.GetPathableTile(2, 2, 0, config)
	neighbors := pt.PathNeighborsAppend(nil)
	assert.Len(t, neighbors, 2)
}

func TestPathfinding_FullPath(t *testing.T) {
	reg := testRegistry()
	level := NewLevel[testTileDefData, testTileData](10, 10, 1, reg)

	config := level.NewPathConfig(
		func(tile *Tile[testTileData]) bool { return true },
		func(from, to *Tile[testTileData]) float64 { return 1.0 },
		nil,
	)

	from := level.GetPathableTile(0, 0, 0, config)
	to := level.GetPathableTile(4, 3, 0, config)

	result, distance, found := path.Path(from, to)
	assert.True(t, found)
	assert.Equal(t, 7.0, distance)
	assert.Len(t, result, 8)
}

// ---------------------------------------------------------------------------
// Utility tests
// ---------------------------------------------------------------------------

func TestSetTileTypeAndVariant(t *testing.T) {
	reg := testRegistry()
	level := NewLevel[testTileDefData, testTileData](5, 5, 1, reg)

	tile := level.GetTileAt(1, 1, 0)
	SetTileTypeAndVariant[testTileDefData](tile, reg, "water", 0)
	assert.Equal(t, 2, tile.Type)
}

func TestRandomVariant(t *testing.T) {
	reg := testRegistry()
	v := RandomVariant(reg, "grass")
	assert.Equal(t, 0, v)

	v = RandomVariant(reg, "nonexistent")
	assert.Equal(t, 0, v)
}

// ---------------------------------------------------------------------------
// Save/Load tests
// ---------------------------------------------------------------------------

func TestSaveAndLoad(t *testing.T) {
	reg := testRegistry()
	level := NewLevel[testTileDefData, testTileData](5, 5, 1, reg)

	_ = level.SetTileAt(2, 2, 0, "stone", 0)
	tile := level.GetTileAt(3, 3, 0)
	tile.CustomData.Visited = true

	tmpFile := filepath.Join(t.TempDir(), "test_save.json")
	gameData := map[string]string{"version": "1.0"}
	err := SaveToFile(level, tmpFile, gameData)
	require.NoError(t, err)

	loaded, rawCustom, err := LoadFromFile[testTileDefData, testTileData](tmpFile, reg)
	require.NoError(t, err)

	assert.Equal(t, 5, loaded.Width)
	assert.Equal(t, 5, loaded.Height)
	assert.Equal(t, 1, loaded.Depth)

	loadedTile := loaded.GetTileAt(2, 2, 0)
	require.NotNil(t, loadedTile)
	assert.Equal(t, 1, loadedTile.Type)

	loadedTile2 := loaded.GetTileAt(3, 3, 0)
	assert.True(t, loadedTile2.CustomData.Visited)

	assert.Equal(t, 2, loadedTile.X)
	assert.Equal(t, 2, loadedTile.Y)

	assert.NotNil(t, rawCustom)
	var gd map[string]string
	require.NoError(t, json.Unmarshal(rawCustom, &gd))
	assert.Equal(t, "1.0", gd["version"])

	config := loaded.NewPathConfig(
		func(tile *Tile[testTileData]) bool { return true },
		nil, nil,
	)
	pt := loaded.GetPathableTile(0, 0, 0, config)
	neighbors := pt.PathNeighborsAppend(nil)
	assert.Len(t, neighbors, 2)
}

func TestSaveLevel_Snapshot(t *testing.T) {
	reg := testRegistry()
	level := NewLevel[testTileDefData, testTileData](3, 3, 1, reg)
	_ = level.SetTileAt(1, 1, 0, "water", 0)

	sd := SaveLevel(level)
	assert.Equal(t, 3, sd.Width)
	assert.Equal(t, 3, sd.Height)
	assert.Len(t, sd.Tiles, 9)

	sd.Tiles[0].Type = 99
	assert.NotEqual(t, 99, level.GetTileAt(0, 0, 0).Type)
}

// ---------------------------------------------------------------------------
// 3D level tests
// ---------------------------------------------------------------------------

func TestLevel3D(t *testing.T) {
	reg := testRegistry()
	level := NewLevel[testTileDefData, testTileData](5, 5, 3, reg)

	_ = level.SetTileAt(2, 2, 0, "grass", 0)
	_ = level.SetTileAt(2, 2, 1, "stone", 0)
	_ = level.SetTileAt(2, 2, 2, "water", 0)

	assert.Equal(t, 0, level.GetTileAt(2, 2, 0).Type)
	assert.Equal(t, 1, level.GetTileAt(2, 2, 1).Type)
	assert.Equal(t, 2, level.GetTileAt(2, 2, 2).Type)
}

func TestPathfinding3D(t *testing.T) {
	reg := testRegistry()
	level := NewLevel[testTileDefData, testTileData](5, 5, 3, reg)

	config := level.NewPathConfig(
		func(tile *Tile[testTileData]) bool { return true },
		func(from, to *Tile[testTileData]) float64 { return 1.0 },
		nil,
	)
	config.Include3DNeighbors = true

	pt := level.GetPathableTile(2, 2, 1, config)
	neighbors := pt.PathNeighborsAppend(nil)
	assert.Len(t, neighbors, 6)

	pt = level.GetPathableTile(2, 2, 2, config)
	neighbors = pt.PathNeighborsAppend(nil)
	assert.Len(t, neighbors, 5)
}
