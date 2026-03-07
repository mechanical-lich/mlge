---
layout: default
title: Resource Manager
nav_order: 8
---

# Resource Manager

`github.com/mechanical-lich/mlge/resource`

Manages loading and caching of textures, fonts, and sounds from files or JSON manifests.

## Global Caches

All loaded assets are stored in global maps for easy access:

```go
var Textures map[string]*ebiten.Image
var Fonts    map[string]font.Face
var Sounds   map[string]audio.AudioResource
```

## Loading Assets from JSON

The primary way to load assets is from a JSON manifest:

```go
err := resource.LoadAssetsFromJSON("data/assets.json")
```

### JSON Format

The JSON file is a map where keys are asset names and values are file paths:

```json
{
    "player": "assets/player.png",
    "tileset": "assets/tiles.png",
    "main": "assets/fonts/main.ttf",
    "folders": ["assets/sprites", "assets/characters"]
}
```

- **Images** (`.png`, `.jpg`, `.jpeg`, `.gif`, `.webp`): Loaded as textures with the key as the name
- **Fonts** (`.ttf`): Loaded as fonts with the key as the name
- **Folders**: The special `folders` key accepts an array of directory paths to recursively load

When loading from folders, asset names are generated from the relative path with directory separators replaced by underscores and without the file extension. For example, `sprites/player/idle.png` becomes `player_idle`.

## Loading Individual Assets

```go
// Load a single texture
err := resource.LoadImageAsTexture("player", "assets/player.png")

// Load a raw image (not cached)
img, err := resource.LoadImage("assets/background.png")

// Load a font
err := resource.LoadFont("main", "assets/fonts/main.ttf")
```

## Sub-Images (Sprite Extraction)

Extract and cache sprite regions from loaded textures:

```go
// By texture name
sprite := resource.GetSubImage("tileset", 0, 0, 16, 16)

// By texture reference
sprite := resource.GetSubImageByTexture(texture, 32, 0, 16, 16)
```

Sub-images are cached internally, so repeated calls with the same parameters return the same image without re-extraction.

## Accessing Loaded Assets

```go
// Get a texture
playerImg := resource.Textures["player"]

// Get a font
mainFont := resource.Fonts["main"]

// Get a sound
bgm := resource.Sounds["bgm"]
```
