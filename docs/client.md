---
layout: default
title: Client
nav_order: 18
---

# Client

`github.com/mechanical-lich/mlge/client`

The client package provides the presentation layer of an mlge game in the client-server architecture. It runs inside Ebitengine's game loop (Update/Draw) and is responsible for rendering, input, and visual interpolation only. It communicates with the `simulation.Server` through a `transport.ClientTransport`.

## RenderSystem

```go
type RenderSystem interface {
    UpdateRender(world any) error
    UpdateEntityRender(world any, entity *ecs.Entity) error
    Requires() []ecs.ComponentType
}
```

The client-side counterpart to `simulation.SimulationSystem`. Implement this for visual logic that runs every Ebitengine frame: sprite animation, position interpolation, particle effects, audio cues.

RenderSystems operate on the client's local (non-authoritative) world copy and must never mutate simulation state directly.

| Method | Description |
|--------|-------------|
| `UpdateRender` | Called once per frame for global render logic (animation timers, camera). |
| `UpdateEntityRender` | Called once per frame per entity matching `Requires()`. |
| `Requires` | Component types an entity must have. |

## RenderSystemManager

```go
type RenderSystemManager struct{}
```

Drives `RenderSystem`s each Ebitengine frame. Same pattern as `ecs.SystemManager` and `simulation.SimulationSystemManager`.

| Method | Signature | Description |
|--------|-----------|-------------|
| `AddSystem` | `(s RenderSystem)` | Append a system to execution order |
| `UpdateSystems` | `(world any) error` | Call `UpdateRender` on all systems |
| `UpdateSystemsForEntities` | `(world any, entities []*ecs.Entity) error` | Call `UpdateEntityRender` per entity per system |

## ClientState

```go
type ClientState interface {
    Update(snapshot *transport.Snapshot) ClientState
    Draw(screen *ebiten.Image)
    Done() bool
}
```

The presentation-layer equivalent of `state.StateInterface`.

Key differences from `state.StateInterface`:
- `Update` receives the latest `*transport.Snapshot` (may be nil if no new snapshot arrived this frame).
- Push/pop semantics are the same: return non-nil from `Update` to push, return `true` from `Done` to pop.

| Method | Description |
|--------|-------------|
| `Update` | Called every Ebitengine frame with the latest snapshot. Return non-nil to push a new state. |
| `Draw` | Renders the current state to screen. |
| `Done` | Returns `true` when this state should be popped. |

## InputMapper

```go
type InputMapper interface {
    MapInputEvent(e event.EventData) (*transport.Command, bool)
}
```

Translates mlge input events into transport `Command`s for forwarding to the server. Return `(nil, false)` to discard events handled locally (e.g., UI hotkeys).

## ClientConfig

```go
type ClientConfig struct {
    ScreenWidth, ScreenHeight int
    WindowTitle               string
}
```

Optional window configuration. If `ScreenWidth`/`ScreenHeight` are zero, the window size is not changed.

## Client

```go
type Client struct{}
```

Implements `ebiten.Game` and acts as the presentation layer. Create with `NewClient`.

### Constructor

```go
func NewClient(
    t transport.ClientTransport,
    codec transport.SnapshotCodec,
    initialState ClientState,
    world any,
    entitySource func() []*ecs.Entity,
    cfg ClientConfig,
) *Client
```

| Parameter | Description |
|-----------|-------------|
| `t` | Client side of a transport |
| `codec` | Same codec used by the server, for decoding snapshots into the local world |
| `initialState` | First `ClientState` pushed onto the state stack |
| `world` | Client's local (non-authoritative) copy of the game world |
| `entitySource` | Function returning the local entity slice for `RenderSystem`s |
| `cfg` | Window size and title |

### Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `SetInputMapper` | `(m InputMapper)` | Set an input mapper. Call before `Run`. |
| `AddRenderSystem` | `(s RenderSystem)` | Add a render system. Call before `Run`. |
| `Run` | `() error` | Start Ebitengine window loop. Blocks until close. |

### Frame Loop

Each Ebitengine frame the `Client` performs:

1. Poll OS input via `input.InputManager`
2. Drain queued events and forward them as `Command`s (if `InputMapper` is set)
3. Receive the latest `Snapshot` from the transport
4. Decode the snapshot into the local world via the `SnapshotCodec`
5. Run `RenderSystemManager` (animation, interpolation) at frame rate
6. Advance the `ClientState` machine with the snapshot

## Usage

```go
srvT, cliT := transport.NewLocalTransport()

// Server runs in background goroutine (see simulation package)

cliWorld := &World{}
c := client.NewClient(
    cliT,
    myCodec,
    &MainClientState{world: cliWorld},
    cliWorld,
    func() []*ecs.Entity { return cliWorld.Entities },
    client.ClientConfig{
        ScreenWidth:  800,
        ScreenHeight: 600,
        WindowTitle:  "My Game",
    },
)
c.AddRenderSystem(&AnimationRenderSystem{})
c.Run()
```
