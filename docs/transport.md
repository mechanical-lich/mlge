---
layout: default
title: Transport
nav_order: 16
---

# Transport

`github.com/mechanical-lich/mlge/transport`

The transport package defines the abstract boundary between the simulation (server) and presentation (client) layers of an mlge game. It provides interfaces and types for sending commands from client to server and delivering world state snapshots from server to client.

The design mirrors Quake's netcode abstraction: the same game code runs whether the transport is local channels or TCP. Swap the implementation at startup; server and client code are unaware of the difference.

| Implementation | Use case |
|----------------|----------|
| `LocalTransport` | Single-player or same-executable multiplayer. Zero serialization cost. |
| `TCPServerTransport` / `TCPClientTransport` | Networked multiplayer over TCP. JSON wire format with TCP_NODELAY. |

## ServerTransport

```go
type ServerTransport interface {
    ReceiveCommands() []*Command
    SendSnapshot(snapshot *Snapshot)
    Close()
}
```

Held by the simulation goroutine. Reads player commands and sends world state snapshots.

| Method | Description |
|--------|-------------|
| `ReceiveCommands` | Drains all pending commands. Returns empty slice (not nil) when idle. Non-blocking. |
| `SendSnapshot` | Sends a snapshot to clients. May drop if the client buffer is full. |
| `Close` | Shuts down the transport. Safe to call multiple times. |

## ClientTransport

```go
type ClientTransport interface {
    SendCommand(cmd *Command)
    ReceiveSnapshot() *Snapshot
    Close()
}
```

Held by the Ebitengine/render side. Sends input commands and polls for the latest world snapshot.

| Method | Description |
|--------|-------------|
| `SendCommand` | Queues a command for the server. Non-blocking; drops if buffer is full. |
| `ReceiveSnapshot` | Returns the most recent snapshot, or nil if none arrived. Non-blocking. |
| `Close` | Shuts down the transport. Safe to call multiple times. |

## Command

```go
type CommandType string

type Command struct {
    Type    CommandType
    Tick    uint64
    Payload any
}
```

A timestamped, typed message from client to server representing player intent (key press, mouse click, game action).

| Field | Description |
|-------|-------------|
| `Type` | Identifies the kind of command. Games define their own constants. |
| `Tick` | Client-estimated server tick this command targets (for future prediction). |
| `Payload` | The command data. Type depends on `CommandType`. |

Built-in command types:

| Constant | Value | Payload Type |
|----------|-------|--------------|
| `CommandInput` | `"input.raw"` | `InputPayload` |

## InputPayload

```go
type InputPayload struct {
    EventType string
    Data      any
}
```

Carries a translated mlge input event. `EventType` is the mlge event type string (e.g., `"MouseLeftClick"`). `Data` is the original event data struct.

## Snapshot

```go
type Snapshot struct {
    Tick      uint64
    Timestamp int64
    Entities  []*EntitySnapshot
}
```

Authoritative world state produced by the server each tick. Consumed by the client to update its local world for rendering.

| Field | Description |
|-------|-------------|
| `Tick` | Server tick number this snapshot was produced on. |
| `Timestamp` | When the snapshot was produced (UnixNano). |
| `Entities` | One snapshot per simulated entity. |

## EntitySnapshot

```go
type EntitySnapshot struct {
    ID         string
    Blueprint  string
    Components map[ecs.ComponentType]ComponentData
}
```

Captures the state of one entity at a given server tick. The `ID` is game-assigned (mlge entities have no built-in ID). The `Components` map holds whichever components the `SnapshotCodec` chose to include.

## SnapshotCodec

```go
type SnapshotCodec interface {
    Encode(tick uint64, entities []*ecs.Entity) *Snapshot
    Decode(snapshot *Snapshot, world any)
}
```

Implemented by the game to control snapshot serialization. Keeps the transport package free of game-specific component knowledge.

| Method | Description |
|--------|-------------|
| `Encode` | Called by the server to build a Snapshot from live entities. Include only what the client needs to render. |
| `Decode` | Called by the client to apply a Snapshot to its local world. Typically find-or-create entities by ID and overwrite components. |

## LocalTransport

```go
func NewLocalTransport() (ServerTransport, ClientTransport)
```

Returns a server/client pair backed by buffered Go channels. For single-player or same-executable games. Zero serialization overhead since values are passed as pointers.

**Buffer sizes:** 64 commands (client to server), 4 snapshots (server to client). When the snapshot buffer is full, the oldest snapshot is dropped so the client always gets the most recent state.

## TCPServerTransport

```go
func NewTCPServerTransport(addr string) (ServerTransport, error)
```

Listens for incoming TCP connections on `addr` (e.g. `":7777"`). Returns an error if the listener cannot bind (port in use, permission denied, etc.).

### Multi-client broadcasting

`SendSnapshot` serializes the snapshot once and writes it to every connected peer. Each peer has its own send mutex so slow clients do not block others. A failed write to one peer is logged and that peer is removed; every other client continues unaffected.

### Peer lifecycle

- A background goroutine is started for each accepted connection to read incoming commands.
- When a peer disconnects (EOF or read error) its goroutine exits and the peer is removed from the broadcast list using a swap-and-nil pattern so the removed slot is eligible for garbage collection immediately.
- `Close` shuts down the TCP listener (stopping new accepts), closes every active peer connection, and waits for all peer goroutines to exit before returning.

`TCP_NODELAY` is set on every accepted connection to minimise command latency.

### Concurrency

`ReceiveCommands` drains a buffered channel (capacity 64). Commands from all connected clients are merged into this single channel. Calls from the simulation goroutine are safe without external locking.

## TCPClientTransport

```go
func NewTCPClientTransport(addr string) (ClientTransport, error)
```

Dials `addr` synchronously and starts a background read goroutine. Returns an error if the connection cannot be established. There is no built-in reconnection; if the server drops the connection a new `TCPClientTransport` must be created.

`ReceiveSnapshot` returns the most recently received snapshot under a mutex and clears it, so the caller always gets the latest state without buffering stale frames.

`TCP_NODELAY` is set on the connection.

## TCP Wire Format

All messages use a simple length-prefix framing protocol:

```
[ 4-byte big-endian uint32 length ][ JSON bytes ]
```

The JSON body is an envelope:

```json
{"k": "cmd"|"snap", "p": <payload>}
```

| `k` value | Direction | `p` payload |
|-----------|-----------|-------------|
| `"cmd"` | client → server | `transport.Command` |
| `"snap"` | server → client | `transport.Snapshot` |

**Maximum message size:** 4 MiB. A peer sending a length header larger than this causes the connection to be closed with an error. This prevents a misbehaving client from forcing unbounded allocation on the server.

### JSON decode and the any/interface{} problem

`encoding/json` decodes JSON objects into `map[string]interface{}` when the target field is typed `any`. This affects two fields that games commonly use:

| Field | Type | Round-trip behaviour |
|-------|------|---------------------|
| `Command.Payload` | `any` | Concrete type on sender; `map[string]interface{}` on receiver |
| `EntitySnapshot.Components` values | `any` | Same |

**Recommended pattern — re-unmarshal in the codec:**

```go
func (c *MyCodec) Decode(snap *transport.Snapshot, world any) {
    for _, es := range snap.Entities {
        raw := es.Components[TypePosition] // map[string]interface{} over TCP
        var pos PositionComponent
        if m, ok := raw.(map[string]interface{}); ok {
            // Extract fields manually or re-marshal to JSON and unmarshal again.
            pos.X, _ = m["X"].(float64)
            pos.Y, _ = m["Y"].(float64)
        } else {
            pos = raw.(PositionComponent) // typed (LocalTransport path)
        }
        entity.AddComponent(pos)
    }
}
```

The `examples/tcp/codec.go` demonstrates direct map-field extraction as the lower-allocation approach. See `examples/tcp/` for a complete runnable example (both local and TCP paths in one binary, switched with `--mode server|client`).

## Usage

**Local (single-player / same process):**

```go
srvT, cliT := transport.NewLocalTransport()
go server.Run(srvT)
ebiten.RunGame(client.New(cliT))
```

**TCP (networked multiplayer):**

```go
// Server process
srvT, err := transport.NewTCPServerTransport(":7777")
if err != nil {
    log.Fatal(err)
}
defer srvT.Close()
server.Run(srvT) // blocks until done

// Client process (separate binary / separate machine)
cliT, err := transport.NewTCPClientTransport("192.168.1.10:7777")
if err != nil {
    log.Fatal(err)
}
defer cliT.Close()
ebiten.RunGame(client.New(cliT))
```

**Swapping transports without changing game code:**

Because both `LocalTransport` and the TCP pair implement the same `ServerTransport` / `ClientTransport` interfaces, game code does not need to know which is in use. A typical pattern:

```go
var srvT transport.ServerTransport
var cliT transport.ClientTransport

if *network {
    srvT, _ = transport.NewTCPServerTransport(*addr)
    cliT, _ = transport.NewTCPClientTransport(*addr)
} else {
    srvT, cliT = transport.NewLocalTransport()
}
```
