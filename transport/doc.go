// Package transport defines the abstract boundary between the simulation (server)
// and presentation (client) layers of an mlge game.
//
// The key types are:
//   - [ServerTransport]: server side - reads commands from clients, sends snapshots.
//   - [ClientTransport]: client side - sends commands to the server, reads snapshots.
//   - [LocalTransport]: in-process implementation backed by Go channels for
//     single-player or same-executable multiplayer. Zero serialization overhead.
//   - [TCPServerTransport] / [TCPClientTransport]: network implementation using
//     length-prefixed JSON over TCP. TCP_NODELAY is set for lower latency.
//
// The design mirrors Quake's netcode abstraction: the same game code runs whether
// the transport is local channels or TCP. Swap the implementation at startup;
// server and client code are unaware of the difference.
//
// Serialization note for TCP: [Command.Payload] and [EntitySnapshot.Components]
// values are encoded with encoding/json. Concrete struct types round-trip
// correctly. Values stored as bare interface{} will decode as
// map[string]interface{} on the remote side; the game must handle that
// (e.g. re-unmarshal into a concrete type, or store json.RawMessage directly).
package transport
