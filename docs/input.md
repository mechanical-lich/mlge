---
layout: default
title: Input
nav_order: 5
---

# Input Manager

`github.com/mechanical-lich/mlge/input`

Translates raw Ebitengine mouse and keyboard input into events dispatched through the event system.

## Setup

```go
import (
    "github.com/mechanical-lich/mlge/event"
    "github.com/mechanical-lich/mlge/input"
)

inputManager := input.NewInputManager(event.GetInstance())
```

Call `HandleInput()` once per frame in your game's `Update()` method:

```go
func (g *Game) Update() error {
    g.inputManager.HandleInput()
    // ...
    return nil
}
```

## Event Types

The input manager automatically fires the following events:

### Mouse Events

| Event Type | Struct | Fields |
|------------|--------|--------|
| `MouseClickEventType` | `MouseClickEvent` | `Button`, `X`, `Y` |
| `MouseReleasedEventType` | `MouseReleasedEvent` | `Button`, `X`, `Y` |
| `MouseMoveEventType` | `MouseMoveEvent` | `X`, `Y`, `OldX`, `OldY` (int) |
| `MouseWheelEventType` | `MouseWheelEvent` | `X`, `Y`, `OldX`, `OldY` (float64) |

### Keyboard Events

| Event Type | Struct | Fields |
|------------|--------|--------|
| `KeyPressEventType` | `KeyPressEvent` | `Keys []ebiten.Key`, `JustPressed bool` |
| `KeyReleaseEventType` | `KeyReleaseEvent` | `Keys []ebiten.Key` |

## Handling Input Events

Register listeners through the event system:

```go
type PlayerController struct{}

func (p *PlayerController) HandleEvent(data event.EventData) error {
    switch e := data.(type) {
    case *input.KeyPressEvent:
        if e.KeyPressed(ebiten.KeyArrowUp) {
            // Move player up
        }
    case *input.MouseClickEvent:
        // Handle click at e.X, e.Y
    }
    return nil
}

// Register
em := event.GetInstance()
em.RegisterListener(input.KeyPressEventType, &PlayerController{})
em.RegisterListener(input.MouseClickEventType, &PlayerController{})
```

### KeyPressEvent Helpers

The `KeyPressEvent` provides a convenience method to check for specific keys:

```go
func (e *KeyPressEvent) KeyPressed(key ebiten.Key) bool
```

The `JustPressed` field distinguishes between a key that was just pressed this frame vs. being held down.
