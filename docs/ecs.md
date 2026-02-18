---
layout: default
title: ECS
nav_order: 3
---

# Entity Component System

`github.com/mechanical-lich/mlge/ecs`

A classic Entity Component System with blueprint-based entity creation, component maps, and a system manager.

## Core Concepts

- **Entity** — A container that holds a map of components
- **Component** — A data struct implementing the `Component` interface
- **System** — Logic that operates on entities with specific component requirements
- **Blueprint** — A JSON definition for creating entities with pre-configured components

## Types

### ComponentType

```go
type ComponentType string
```

String identifier for component types. Each component defines its own type constant.

### Component

```go
type Component interface {
    GetType() ComponentType
}
```

All components must implement this interface.

### Entity

```go
type Entity struct {
    Components map[ComponentType]Component
    Blueprint  string
}
```

**Methods:**

| Method | Signature | Description |
|--------|-----------|-------------|
| `AddComponent` | `(component Component)` | Adds a component to the entity |
| `HasComponent` | `(componentType ComponentType) bool` | Checks if entity has a component |
| `HasComponents` | `(componentTypes ...ComponentType) bool` | Checks for multiple components |
| `HasComponentsSlice` | `(componentTypes []ComponentType) bool` | Checks for components from a slice |
| `GetComponent` | `(componentType ComponentType) Component` | Retrieves a component by type |
| `RemoveComponent` | `(componentType ComponentType)` | Removes a component |

### SystemInterface

```go
type SystemInterface interface {
    UpdateSystem(params any) error
    UpdateEntity(params any, entity *Entity) error
    Requires() []ComponentType
}
```

Systems are called by the `SystemManager` for each entity that has the required components.

### SystemManager

```go
type SystemManager struct{}
```

**Methods:**

| Method | Signature | Description |
|--------|-----------|-------------|
| `AddSystem` | `(system SystemInterface)` | Registers a system |
| `UpdateSystems` | `(params any) error` | Runs all systems (calls `UpdateSystem`) |
| `UpdateSystemsForEntity` | `(params any, entity *Entity) error` | Runs systems for a single entity |
| `UpdateSystemsForEntities` | `(params any, entities []*Entity) error` | Runs systems for a slice of entities |

## Blueprints

Blueprints allow you to define entity templates in JSON and create entities from them at runtime.

### Registering Component Factories

Before loading blueprints, register factory functions that know how to create each component type:

```go
type ComponentAddFunction func([]string) (Component, error)

ecs.RegisterComponentAddFunction("position2d", func(args []string) (ecs.Component, error) {
    x, _ := strconv.ParseFloat(args[0], 64)
    y, _ := strconv.ParseFloat(args[1], 64)
    return &basecomponents.Position2dComponent{X: x, Y: y}, nil
})
```

### Blueprint JSON Format

```json
{
    "blueprints": {
        "player": {
            "components": {
                "position2d": ["100", "200"],
                "health": ["100"],
                "sprite": ["player_idle"]
            }
        },
        "enemy": {
            "components": {
                "position2d": ["0", "0"],
                "health": ["50"],
                "sprite": ["enemy_idle"],
                "ai": []
            }
        }
    }
}
```

### Loading and Creating

```go
// Load blueprints from file
err := ecs.LoadBlueprintsFromFile("data/blueprints.json")

// Or from an io.Reader
err := ecs.LoadFactoryFromStream(reader)

// Create an entity from a blueprint
entity, err := ecs.Create("player")
```

## Built-in Components

The `ecs/basecomponents` package provides common position components:

```go
import "github.com/mechanical-lich/mlge/ecs/basecomponents"
```

### Position2dComponent

```go
type Position2dComponent struct {
    X, Y float64
}
```

Type: `basecomponents.Position2d`

Methods: `GetX()`, `GetY()`, `SetPosition(x, y float64)`

### Position3dComponent

```go
type Position3dComponent struct {
    X, Y, Z float64
}
```

Type: `basecomponents.Position3d`

Methods: `GetX()`, `GetY()`, `GetZ()`, `SetPosition(x, y, z float64)`

## Inanimate Entities

Set `ecs.InanimateComponentType` to a component type to skip entities with that component during system updates. This is useful for static objects that don't need per-frame processing.

```go
ecs.InanimateComponentType = "inanimate"
```

## Example: Custom Component and System

```go
// Define a component
const HealthType ecs.ComponentType = "health"

type HealthComponent struct {
    Current int
    Max     int
}

func (h *HealthComponent) GetType() ecs.ComponentType {
    return HealthType
}

// Define a system
type RegenSystem struct{}

func (s *RegenSystem) Requires() []ecs.ComponentType {
    return []ecs.ComponentType{HealthType}
}

func (s *RegenSystem) UpdateSystem(params any) error {
    return nil
}

func (s *RegenSystem) UpdateEntity(params any, entity *ecs.Entity) error {
    health := entity.GetComponent(HealthType).(*HealthComponent)
    if health.Current < health.Max {
        health.Current++
    }
    return nil
}

// Register and use
sm := &ecs.SystemManager{}
sm.AddSystem(&RegenSystem{})
sm.UpdateSystemsForEntities(nil, entities)
```
