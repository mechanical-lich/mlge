package world

import (
	"encoding/json"
	"fmt"
	"os"
)

// TileVariant represents a single visual variant of a tile type,
// pointing to a specific sprite location.
type TileVariant struct {
	Variant int `json:"variant"`
	SpriteX int `json:"spriteX"`
	SpriteY int `json:"spriteY"`
}

// TileDefinition defines a type of tile with shared properties.
// TD is a generic type parameter for game-specific custom data
// (e.g., Door/Wall bools for a dungeon game, MixGroup for a platformer).
type TileDefinition[TD any] struct {
	Name       string        `json:"name"`
	Solid      bool          `json:"solid"`
	Water      bool          `json:"water"`
	Variants   []TileVariant `json:"variants"`
	CustomData TD            `json:"customData"`
}

// TileRegistry holds the loaded tile definitions and provides lookups
// by name or index. Since Go doesn't allow generic package-level variables,
// this struct holds the registry per-game instance.
type TileRegistry[TD any] struct {
	Definitions []TileDefinition[TD]
	NameToIndex map[string]int
	IndexToName []string
}

// NewTileRegistry creates an empty tile registry.
func NewTileRegistry[TD any]() *TileRegistry[TD] {
	return &TileRegistry[TD]{
		NameToIndex: make(map[string]int),
	}
}

// LoadFromFile reads a JSON file containing an array of tile definitions
// and populates the registry's Definitions, NameToIndex, and IndexToName fields.
func (r *TileRegistry[TD]) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("world: failed to read tile definitions file %s: %w", path, err)
	}

	var defs []TileDefinition[TD]
	if err := json.Unmarshal(data, &defs); err != nil {
		return fmt.Errorf("world: failed to parse tile definitions: %w", err)
	}

	r.Definitions = defs
	r.NameToIndex = make(map[string]int, len(defs))
	r.IndexToName = make([]string, len(defs))

	for i, def := range defs {
		r.NameToIndex[def.Name] = i
		r.IndexToName[i] = def.Name
	}

	return nil
}

// LoadFromDefinitions populates the registry from a pre-built slice of definitions.
func (r *TileRegistry[TD]) LoadFromDefinitions(defs []TileDefinition[TD]) {
	r.Definitions = defs
	r.NameToIndex = make(map[string]int, len(defs))
	r.IndexToName = make([]string, len(defs))

	for i, def := range defs {
		r.NameToIndex[def.Name] = i
		r.IndexToName[i] = def.Name
	}
}

// GetByName returns the tile definition and its index for the given name.
// Returns nil and -1 if not found.
func (r *TileRegistry[TD]) GetByName(name string) (*TileDefinition[TD], int) {
	idx, ok := r.NameToIndex[name]
	if !ok {
		return nil, -1
	}
	return &r.Definitions[idx], idx
}

// GetByIndex returns the tile definition at the given index.
// Returns nil if the index is out of range.
func (r *TileRegistry[TD]) GetByIndex(index int) *TileDefinition[TD] {
	if index < 0 || index >= len(r.Definitions) {
		return nil
	}
	return &r.Definitions[index]
}
