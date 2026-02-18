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

```json
{
    "textures": [
        { "name": "player", "path": "assets/player.png" }
    ],
    "texture_folders": [
        { "path": "assets/sprites" }
    ],
    "fonts": [
        { "name": "main", "path": "assets/fonts/main.ttf" }
    ],
    "sounds": [
        { "name": "bgm", "path": "assets/audio/bgm.ogg", "type": "ogg" }
    ]
}
```

When loading from `texture_folders`, each image file is registered with its filename (without extension) as the key.

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
