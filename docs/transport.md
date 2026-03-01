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

Listens for incoming TCP connections on `addr` (e.g. `":7777"`). Accepts multiple clients. Each client gets a dedicated read goroutine. `Close` shuts down the listener and all connections.

`SendSnapshot` broadcasts to every connected client. Clients that fail a write are logged; other clients are unaffected.

`TCP_NODELAY` is set on every accepted connection to reduce latency.

## TCPClientTransport

```go
func NewTCPClientTransport(addr string) (ClientTransport, error)
```

Dials `addr` (e.g. `"localhost:7777"`) and starts a background read goroutine. `ReceiveSnapshot` returns the most recently received snapshot and clears it, so the caller always sees up-to-date state.

`TCP_NODELAY` is set on the connection.

## TCP Wire Format

All messages use a simple framing protocol:

```
[ 4-byte big-endian uint32 length ][ JSON bytes ]
```

The JSON body is an envelope:

```json
{"k": "cmd"|"snap", "p": <payload>}
```

Commands and snapshots are marshaled with `encoding/json`. This means `Command.Payload` and `EntitySnapshot.Components` values stored as bare `any`/`interface{}` will decode on the remote side as `map[string]interface{}`. Your game code must re-unmarshal those into concrete types, or store `json.RawMessage` directly in those fields.

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
server.Run(srvT)

// Client process
cliT, err := transport.NewTCPClientTransport("192.168.1.10:7777")
if err != nil {
    log.Fatal(err)
}
defer cliT.Close()
ebiten.RunGame(client.New(cliT))
```
