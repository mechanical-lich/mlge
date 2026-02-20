package world

import (
	"encoding/json"
	"fmt"
	"os"
)

// SaveData holds the serializable state of a level.
// T is the per-tile custom data type.
type SaveData[T any] struct {
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	Depth      int       `json:"depth"`
	Tiles      []Tile[T] `json:"tiles"`
	CustomData any       `json:"customData,omitempty"`
}

// SaveLevel creates a SaveData snapshot of the level's current tile state.
func SaveLevel[TD any, T any](level *Level[TD, T]) SaveData[T] {
	// Copy tiles so the save data is independent of the live level
	tilesCopy := make([]Tile[T], len(level.tiles))
	copy(tilesCopy, level.tiles)

	// Clear unexported back-pointers so they don't leak into serialization
	for i := range tilesCopy {
		tilesCopy[i].levelRef = nil
	}

	return SaveData[T]{
		Width:  level.Width,
		Height: level.Height,
		Depth:  level.Depth,
		Tiles:  tilesCopy,
	}
}

// SaveToFile serializes the level's tile data and optional custom data to a JSON file.
func SaveToFile[TD any, T any](level *Level[TD, T], filename string, customData any) error {
	sd := SaveLevel[TD, T](level)
	sd.CustomData = customData

	data, err := json.MarshalIndent(sd, "", "  ")
	if err != nil {
		return fmt.Errorf("world: failed to marshal save data: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("world: failed to write save file %s: %w", filename, err)
	}

	return nil
}

// LoadFromFile reads a JSON save file and reconstructs a Level.
// It returns the level, the raw CustomData (as json.RawMessage for the game
// to unmarshal into its own type), and any error.
func LoadFromFile[TD any, T any](filename string, registry *TileRegistry[TD]) (*Level[TD, T], json.RawMessage, error) {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("world: failed to read save file %s: %w", filename, err)
	}

	// First pass: unmarshal into a structure with RawMessage for CustomData
	var raw struct {
		Width      int             `json:"width"`
		Height     int             `json:"height"`
		Depth      int             `json:"depth"`
		Tiles      []Tile[T]       `json:"tiles"`
		CustomData json.RawMessage `json:"customData,omitempty"`
	}

	if err := json.Unmarshal(fileData, &raw); err != nil {
		return nil, nil, fmt.Errorf("world: failed to parse save data: %w", err)
	}

	level := &Level[TD, T]{
		tiles:    raw.Tiles,
		Width:    raw.Width,
		Height:   raw.Height,
		Depth:    raw.Depth,
		Registry: registry,
	}

	// Restore back-pointers on tiles
	for i := range level.tiles {
		level.tiles[i].levelRef = level
	}

	return level, raw.CustomData, nil
}
