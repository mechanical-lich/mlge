---
layout: default
title: Getting Started
nav_order: 2
---

# Getting Started

This guide covers installation, project setup, and a minimal example using MLGE.

## Prerequisites

- **Go 1.25+**
- Ebitengine system dependencies (see [Ebitengine install guide](https://ebitengine.org/en/documents/install.html))

## Installation

Add MLGE to your Go project:

```bash
go get github.com/mechanical-lich/mlge
```

## Minimal Example

A basic game loop that initializes the resource manager, input system, and state machine:

```go
package main

import (
    "log"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/mechanical-lich/mlge/event"
    "github.com/mechanical-lich/mlge/input"
    "github.com/mechanical-lich/mlge/resource"
    "github.com/mechanical-lich/mlge/state"
)

type Game struct {
    inputManager *input.InputManager
    stateMachine *state.StateMachine
}

func (g *Game) Update() error {
    g.inputManager.HandleInput()
    event.GetQueuedInstance().HandleQueue()
    g.stateMachine.Update()
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    g.stateMachine.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return 800, 600
}

func main() {
    // Load assets
    if err := resource.LoadAssetsFromJSON("data/assets.json"); err != nil {
        log.Fatal(err)
    }

    // Initialize input
    inputManager := input.NewInputManager(event.GetInstance())

    // Create state machine and push initial state
    sm := &state.StateMachine{}
    sm.PushState(&MyGameState{})

    game := &Game{
        inputManager: inputManager,
        stateMachine: sm,
    }

    ebiten.SetWindowSize(800, 600)
    ebiten.SetWindowTitle("My Game")
    if err := ebiten.RunGame(game); err != nil {
        log.Fatal(err)
    }
}
```

## Asset JSON Format

The `resource.LoadAssetsFromJSON` function expects a JSON file that declares textures, fonts, and sounds:

```json
{
    "textures": [
        { "name": "player", "path": "assets/player.png" },
        { "name": "tiles", "path": "assets/tileset.png" }
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

## Project Structure

A recommended project layout using MLGE:

```
my-game/
├── cmd/
│   └── game/
│       └── main.go
├── assets/
│   ├── sprites/
│   ├── fonts/
│   └── audio/
├── data/
│   ├── assets.json
│   └── blueprints.json
├── internal/
│   ├── components/
│   ├── systems/
│   └── states/
├── go.mod
└── Makefile
```

## Next Steps

- [ECS](ecs.md) — Learn how to create entities with components and systems
- [Event System](event.md) — Set up event-driven communication
- [Input](input.md) — Handle mouse and keyboard input
- [UI Framework](ui.md) — Build game interfaces
