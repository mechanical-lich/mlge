// Package simulation provides the server-side (authoritative) game loop for mlge.
//
// In a Quake-style architecture the server owns the canonical world state and
// runs physics, AI, and all game logic at a fixed tick rate, independent of
// the client frame rate. The client only renders.
//
// Key types:
//   - [SimulationSystem]: server equivalent of ecs.SystemInterface. No rendering.
//   - [SimulationSystemManager]: runs SimulationSystems each server tick.
//   - [SimulationState]: server equivalent of state.StateInterface. No Draw().
//   - [Server]: runs the simulation loop in a goroutine at a configurable tick rate.
//
// This package has zero Ebitengine dependencies.
// If you see "github.com/hajimehoshi/ebiten" imported here, that is a bug.
package simulation
