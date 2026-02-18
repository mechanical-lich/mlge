---
layout: default
title: Sensory System
nav_order: 11
---

# Sensory System

`github.com/mechanical-lich/mlge/sense`

A grid-based sensory simulation supporting sound, pheromone, and scent stimuli with decay and propagation.

## Stimulus Types

| Constant | Description |
|----------|-------------|
| `SoundStimuli` | Sound stimulus — propagates outward and decays |
| `PheremoneStimuli` | Pheromone stimulus — lingers and slowly decays |
| `ScentStimuli` | Scent stimulus — affected by draft/wind direction |

## Draft Directions

| Constant | Description |
|----------|-------------|
| `NorthernDraft` | Wind blowing north |
| `WesternDraft` | Wind blowing west |
| `EasternDraft` | Wind blowing east |
| `SouthernDraft` | Wind blowing south |

## Types

### Stimulus

```go
type Stimulus struct {
    Type      StimulusType
    Intensity float32
    Decay     float32
    ID        string
}
```

### StimuliTile

```go
type StimuliTile struct {
    Stimuli   []Stimulus
    Resonance float32
    Draft     float32
    DraftDir  DraftDirection
    Solid     bool
}
```

Tiles marked as `Solid` block stimulus propagation.

### SenseScape

A 2D grid of `StimuliTile` that manages stimulus propagation and decay.

## Usage

```go
import "github.com/mechanical-lich/mlge/sense"

// Create a 100x100 sensory grid
scape := sense.NewSenseScape(100, 100)

// Make a sound at position (50, 50)
scape.MakeSound(50, 50, "explosion", 100)

// Apply a custom stimulus
scape.ApplyStimulus(30, 30, sense.Stimulus{
    Type:      sense.PheremoneStimuli,
    Intensity: 50,
    Decay:     0.1,
    ID:        "trail",
})

// Each frame, update propagation and decay
scape.Update()

// Query stimuli at a position
stimuli, err := scape.GetStimuliAt(50, 50)
for _, s := range stimuli {
    if s.Type == sense.SoundStimuli && s.Intensity > 10 {
        // Entity detected a loud sound nearby
    }
}
```

## How It Works

Each call to `Update()`:
1. Stimuli propagate to neighboring tiles (reduced by decay)
2. Stimuli on each tile decay over time
3. Scent stimuli are influenced by draft direction and strength
4. Solid tiles block propagation
5. Tile resonance affects how stimuli persist

This system enables AI entities to "sense" their environment — hearing sounds, following pheromone trails, or detecting scents carried by wind.
