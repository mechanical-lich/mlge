// Package client provides the presentation layer of an mlge game in the
// client-server architecture introduced by the simulation and transport packages.
//
// The client runs inside Ebitengine's game loop (Update/Draw) and is responsible
// for rendering, input, and visual interpolation only. It communicates with
// the simulation.Server through a transport.ClientTransport.
//
// Key types:
//   - [RenderSystem]: client equivalent of ecs.SystemInterface. Runs at frame rate.
//   - [RenderSystemManager]: drives RenderSystems each frame.
//   - [ClientState]: client equivalent of state.StateInterface.
//     Receives the latest transport.Snapshot on Update; renders in Draw.
//   - [Client]: implements ebiten.Game. Wires transport, snapshot decode,
//     render systems, and the client state machine together.
package client
