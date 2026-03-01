---
layout: default
title: Simulation
nav_order: 17
---

# Simulation

`github.com/mechanical-lich/mlge/simulation`

The simulation package provides the server-side (authoritative) game loop for mlge. In a Quake-style architecture the server owns the canonical world state and runs physics, AI, and all game logic at a fixed tick rate, independent of the client frame rate. The client only renders.

**This package has zero Ebitengine dependencies.** If you see `github.com/hajimehoshi/ebiten` imported by simulation code, that is a bug.

## SimulationSystem

```go
type SimulationSystem interface {
    UpdateSimulation(world any) error
    UpdateEntitySimulation(world any, entity *ecs.Entity) error
    Requires() []ecs.ComponentType
}
```

The server-side counterpart to `ecs.SystemInterface`. Implement this for authoritative logic: physics, AI, chemistry, pathfinding, collision resolution. Never put rendering code here.

| Method | Description |
|--------|-------------|
| `UpdateSimulation` | Runs once per server tick for global logic (clocks, spatial indices). |
| `UpdateEntitySimulation` | Runs once per entity per tick, for entities matching `Requires()`. |
| `Requires` | Component types an entity must have. Return nil to receive every entity. |

## SimulationSystemManager

```go
type SimulationSystemManager struct{}
```

Holds an ordered list of `SimulationSystem`s and drives them each server tick. Mirrors `ecs.SystemManager` but is typed to prevent accidentally adding render systems.

| Method | Signature | Description |
|--------|-----------|-------------|
| `AddSystem` | `(s SimulationSystem)` | Append a system to execution order |
| `UpdateSystems` | `(world any) error` | Call `UpdateSimulation` on all systems |
| `UpdateSystemsForEntities` | `(world any, entities []*ecs.Entity) error` | Call `UpdateEntitySimulation` per entity per system |

Entities with `ecs.InanimateComponentType` are automatically skipped, matching `ecs.SystemManager` behavior.

## SimulationState

```go
type SimulationState interface {
    Tick(world any) SimulationState
    ProcessCommand(cmd *transport.Command)
    Done() bool
}
```

The server-side equivalent of `state.StateInterface`. No `Draw()` method. Operates inside the simulation goroutine, never touching Ebitengine APIs.

| Method | Description |
|--------|-------------|
| `Tick` | Advances one server tick. Return non-nil to push a new state. |
| `ProcessCommand` | Handles one client command. Called before `Tick`, in arrival order. |
| `Done` | Returns `true` when this state should be popped. |

## SimulationStateMachine

```go
type SimulationStateMachine struct{}
```

Stack-based state machine for the server. Mirrors `state.StateMachine` but has no `Draw`.

| Method | Signature | Description |
|--------|-----------|-------------|
| `PushState` | `(s SimulationState)` | Push a new state onto the stack |
| `Current` | `() SimulationState` | Return the active state, or nil if empty |
| `Tick` | `(world any) bool` | Advance current state; returns false when stack is empty |
| `ProcessCommands` | `(cmds []*transport.Command)` | Route all pending commands to the current state |

## ServerConfig

```go
type ServerConfig struct {
    TickRate      int  // ticks per second (default: 20)
    SnapshotEvery int  // send snapshot every N ticks (default: 1)
}
```

## Server

```go
type Server struct{}
```

Runs the authoritative simulation loop in a dedicated goroutine. Owns the game world, drives `SimulationSystem`s at a fixed tick rate, processes client commands, and sends snapshots via the transport.

### Constructor

```go
func NewServer(
    config ServerConfig,
    world any,
    entitySource EntitySource,
    t transport.ServerTransport,
    codec transport.SnapshotCodec,
) *Server
```

| Parameter | Description |
|-----------|-------------|
| `config` | Tick rate and snapshot frequency |
| `world` | The authoritative game world (passed to systems as `world any`) |
| `entitySource` | Function returning the current entity slice (avoids stale pointers) |
| `t` | Server side of a transport |
| `codec` | Game-provided snapshot encoder |

### Methods

| Method | Signature | Description |
|--------|-----------|-------------|
| `AddSystem` | `(sys SimulationSystem)` | Register a system. Call before `Run` or `Step`. |
| `SetState` | `(state SimulationState)` | Set the initial state. Call before `Run` or `Step`. |
| `Run` | `()` | Start the loop. Blocks until `Stop()` or state machine empties. |
| `Step` | `() bool` | Advance exactly one tick. Returns false when the state machine is empty. |
| `Stop` | `()` | Signal the loop to exit cleanly. |
| `Tick` | `() uint64` | Current tick counter (read-only). |

### Driving Modes

The `Server` supports two modes of operation:

**Independent (decoupled)** -- call `Run()` in a goroutine. The server ticks at
its own fixed rate, independent of the client frame rate. Best for networked or
multi-threaded games.

```go
go srv.Run()
```

**Linked (lockstep)** -- call `Step()` from the caller's loop, typically inside
Ebitengine's `Update()`. The simulation ticks once per frame, locked to the
engine's TPS. Simpler architecture, no goroutines, no concurrency concerns.

```go
func (g *Game) Update() error {
    if !g.srv.Step() {
        return ebiten.Termination
    }
    // receive snapshot, render, etc.
    return nil
}
```

### Tick Loop

Each server tick the `Server` performs these steps in order:

1. Drain pending commands from the transport
2. Route commands to the current `SimulationState`
3. Run `SimulationSystemManager.UpdateSystems` (global pass)
4. Run `SimulationSystemManager.UpdateSystemsForEntities` (per-entity pass)
5. Advance the `SimulationStateMachine`
6. If this tick is a snapshot tick, encode and send a `Snapshot`

## Usage

### Independent mode (decoupled tick rate)

```go
srvT, cliT := transport.NewLocalTransport()

srv := simulation.NewServer(
    simulation.ServerConfig{TickRate: 20},
    myWorld,
    func() []*ecs.Entity { return myWorld.Entities },
    srvT,
    myCodec,
)
srv.AddSystem(&PhysicsSystem{})
srv.AddSystem(&AISystem{})
srv.SetState(&MainSimState{})
go srv.Run()

// Client runs in the main goroutine via Ebitengine
ebiten.RunGame(client.NewClient(cliT, myCodec, myState, myWorld, entityFunc, cfg))
```

See [examples/client_server](https://github.com/mechanical-lich/mlge/tree/main/examples/client_server) for a full runnable example.

### Linked mode (lockstep with Ebitengine)

```go
srvT, cliT := transport.NewLocalTransport()

srv := simulation.NewServer(
    simulation.ServerConfig{TickRate: 60, SnapshotEvery: 1},
    srvWorld,
    func() []*ecs.Entity { return srvWorld.Entities },
    srvT,
    myCodec,
)
srv.AddSystem(&PhysicsSystem{DT: 1.0 / 60.0})
srv.SetState(&MainSimState{})

// No goroutine -- Step() is called from Update().
type Game struct {
    srv   *simulation.Server
    cliT  transport.ClientTransport
    codec transport.SnapshotCodec
    world *World
    tick  uint64
}

func (g *Game) Update() error {
    if !g.srv.Step() {
        return ebiten.Termination
    }
    snap := g.cliT.ReceiveSnapshot()
    if snap != nil {
        g.codec.Decode(snap, g.world)
        g.tick = snap.Tick
    }
    return nil
}
```

See [examples/linked](https://github.com/mechanical-lich/mlge/tree/main/examples/linked) for a full runnable example.
