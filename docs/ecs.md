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
- **Blueprint** — A text-based definition for creating entities with pre-configured components

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

Blueprints allow you to define entity templates and create entities from them at runtime.

### Text-Based Factory (Legacy)

The original factory uses a line-based text format with string parameters.

#### Registering Component Factories

```go
type ComponentAddFunction func([]string) (Component, error)

ecs.RegisterComponentAddFunction("position2d", func(args []string) (ecs.Component, error) {
    x, _ := strconv.ParseFloat(args[0], 64)
    y, _ := strconv.ParseFloat(args[1], 64)
    return &basecomponents.Position2dComponent{X: x, Y: y}, nil
})
```

#### Blueprint Text Format

Each entity starts with its name on a new line, followed by component definitions in the format `componentType:param1,param2,...`. Entities are separated by blank lines.

```text
player
position2d:100,200
health:100
sprite:player_idle

enemy
position2d:0,0
health:50
sprite:enemy_idle
ai:
```

#### Loading and Creating

```go
// Load blueprints from file
err := ecs.LoadBlueprintsFromFile("data/blueprints.txt")

// Or from an io.Reader
err := ecs.LoadFactoryFromStream(reader)

// Create an entity from a blueprint
entity, err := ecs.Create("player")
```

### JSON Factory (Recommended)

The `JSONFactory` uses JSON blueprints and populates components via JSON unmarshalling, making it more flexible for complex component data.

#### Creating a Factory

```go
factory := ecs.NewJSONFactory()
```

#### Registering Components

Register component constructors that return zero-value component pointers:

```go
factory.RegisterComponent("HealthComponent", func() ecs.Component {
    return &components.HealthComponent{}
})
factory.RegisterComponent("AppearanceComponent", func() ecs.Component {
    return &components.AppearanceComponent{}
})
```

#### Blueprint JSON Format

Each JSON file contains a map of blueprint names to component definitions:

```json
{
    "player": {
        "HealthComponent": {"MaxHealth": 100, "Health": 100},
        "AppearanceComponent": {"SpriteName": "player_idle"}
    },
    "enemy": {
        "HealthComponent": {"MaxHealth": 50},
        "AppearanceComponent": {"SpriteName": "enemy_idle"},
        "HostileAIComponent": {}
    }
}
```

#### Loading Blueprints

```go
// Load from a single file
err := factory.LoadBlueprintsFromFile("data/units.json")

// Load all JSON files from a directory
err := factory.LoadBlueprintsFromDir("data/blueprints")
```

#### Creating Entities

```go
// Create an entity from a blueprint
entity, err := factory.Create("player")

// Create with a callback for custom initialization
entity, err := factory.CreateWithCallback("player", func(comp ecs.Component) error {
    // Auto-initialize health if not set
    if hc, ok := comp.(*HealthComponent); ok {
        if hc.Health == 0 {
            hc.Health = hc.MaxHealth
        }
    }
    return nil
})
```

#### Additional Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `BlueprintExists` | `(name string) bool` | Check if a blueprint is registered |
| `GetBlueprintNames` | `() []string` | Get all registered blueprint names |

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
