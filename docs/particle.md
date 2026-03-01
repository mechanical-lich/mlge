---
layout: default
title: Particle
nav_order: 18
---

# Particle System

`github.com/mechanical-lich/mlge/particle`

The particle package provides a data-driven, ECS-integrated particle emitter. Each emitter is configured through an `EmitterComponent` attached to an entity. A single `ParticleSystem` manages all live particles, handles spawning, aging, movement, and gravity, and renders them each frame.

## Quick start

```go
ps := particle.NewParticleSystem(1.0 / 60.0)

fire := &ecs.Entity{Blueprint: "fire"}
fire.AddComponent(particle.FireEmitter(320, 240))

// Advance each frame:
ps.UpdateEntity(nil, fire)

// Render each frame:
ps.Draw(screen)
```

See `examples/particles/main.go` for a full runnable demo with fire, smoke, and click-triggered sparks.

## EmitterComponent

Attach one `EmitterComponent` per entity to configure that entity's emitter.

```go
type EmitterComponent struct {
    X, Y           float64     // world-space origin; update each tick to follow the entity
    StartColor     color.RGBA  // particle color at birth
    EndColor       color.RGBA  // particle color at death (alpha also interpolated)
    StartSize      float64     // particle radius in pixels at birth
    EndSize        float64     // particle radius in pixels at death
    EmitRate       float64     // particles per second (fractional rates supported)
    MaxParticles   int         // cap on live particles; 0 uses DefaultMaxParticles (256)
    DirectionAngle float64     // base emission angle in radians (0 = right)
    Spread         float64     // half-angle deviation in radians (math.Pi = omnidirectional)
    SpeedMin       float64     // minimum particle speed in px/s
    SpeedMax       float64     // maximum particle speed in px/s
    LifeMin        float64     // minimum particle lifetime in seconds
    LifeMax        float64     // maximum particle lifetime in seconds
    Gravity        float64     // downward acceleration in px/s² (negative = upward)
    Image          *ebiten.Image // optional sprite; nil draws filled circles
    Active         bool        // enables continuous emission
    BurstCount     int         // emit exactly N particles on the next tick, then reset to 0
}
```

### Continuous emission

Set `Active = true` and `EmitRate > 0`. The system accumulates a fractional counter each tick so low rates (e.g. `0.5` particles/second) fire accurately without drift.

### One-shot burst

Set `BurstCount` to a positive integer. The system fires exactly that many particles on the next update and clears `BurstCount` to zero automatically. Combine with `Active = false` for effects that fire once and go quiet.

### Moving emitters

Update `X` and `Y` on the `EmitterComponent` each tick to attach the emitter to a moving entity:

```go
pos := e.Components[yourPosType].(YourPosComponent)
cfg := e.Components[particle.ComponentType].(particle.EmitterComponent)
cfg.X, cfg.Y = pos.X, pos.Y
e.AddComponent(cfg)
```

## Preset constructors

Three ready-made configs are included:

| Constructor | Description |
|-------------|-------------|
| `FireEmitter(x, y)` | Upward flame, orange-to-red fade, slight spread |
| `SmokeEmitter(x, y)` | Slow upward drift, grey, expands over lifetime |
| `SparkBurst(x, y, count)` | One-shot omnidirectional explosion, yellow-to-orange |

## ParticleSystem

```go
func NewParticleSystem(dt float64) *ParticleSystem
```

`dt` is the fixed timestep in seconds. Typically `1.0/60.0` for a render-driven system or `1.0/tickRate` for a simulation-driven system.

### System interfaces

`ParticleSystem` satisfies all three mlge system interfaces through structural typing. Register it with whichever manager fits the game's architecture:

| Interface | Manager method | Typical use |
|-----------|---------------|-------------|
| `ecs.SystemInterface` | `systemManager.AddSystem(ps)` | Simple single-loop games |
| `simulation.SimulationSystem` | `simManager.AddSystem(ps)` | Server-authoritative particle state |
| `client.RenderSystem` | `client.AddRenderSystem(ps)` | Client-only visual effects (most common) |

### Draw pass

The system **does not** call `Draw` automatically — update and rendering are separated intentionally. Call `Draw` from the render state after drawing the world:

```go
func (s *MyState) Draw(screen *ebiten.Image) {
    drawWorld(screen)
    ps.Draw(screen) // particles on top
}
```

### Lifecycle helpers

```go
// Remove pools for entities that are no longer alive.
ps.Purge(world.Entities)

// Total live particle count across all emitters (useful for debug overlays).
n := ps.ActiveCount()
```

Call `Purge` after reconciling the entity list (e.g. after decoding a network snapshot) to prevent stale pools accumulating for destroyed entities.

## Memory notes

- Particle slices are compacted in place each tick with a swap-and-nil pattern so no allocations occur during normal operation once the pool reaches steady state.
- `MaxParticles` (default 256) bounds the per-emitter allocation. Keep it as small as the effect allows.
- Call `Purge` regularly if entities are frequently created and destroyed.
