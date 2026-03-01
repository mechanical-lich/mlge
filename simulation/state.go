package simulation

import "github.com/mechanical-lich/mlge/transport"

// SimulationState is the server-side equivalent of state.StateInterface.
//
// It has no Draw method - rendering is entirely the client's responsibility.
// The state machine mirrors the client's stack-based state machine but operates
// inside the simulation goroutine, never touching Ebitengine APIs.
//
// Typical implementations:
//   - A "loading" state that generates the world before handing off to gameplay.
//   - The main gameplay state that owns the Level and routes commands to it.
//   - A "round over" state that waits before resetting.
type SimulationState interface {
	// Tick advances the simulation by one server tick.
	// world is the authoritative game world (same type as what SimulationSystems receive).
	// Return a non-nil SimulationState to push a new state onto the stack.
	// When Done() is true the Server pops this state.
	Tick(world any) SimulationState

	// ProcessCommand handles a command received from a client this tick.
	// Called once per command, before Tick, in the order they arrived.
	ProcessCommand(cmd *transport.Command)

	// Done signals the Server to pop this state.
	Done() bool
}

// SimulationStateMachine is a minimal stack-based state machine for the
// server side. It mirrors state.StateMachine but has no Draw method.
type SimulationStateMachine struct {
	states []SimulationState
}

// PushState pushes a new state onto the stack.
func (m *SimulationStateMachine) PushState(s SimulationState) {
	m.states = append(m.states, s)
}

// Current returns the active state, or nil if the stack is empty.
func (m *SimulationStateMachine) Current() SimulationState {
	if len(m.states) == 0 {
		return nil
	}
	return m.states[len(m.states)-1]
}

// Tick advances the current state and handles push/pop transitions.
// Returns false when the stack is empty (simulation should stop).
func (m *SimulationStateMachine) Tick(world any) bool {
	if len(m.states) == 0 {
		return false
	}
	top := m.states[len(m.states)-1]
	next := top.Tick(world)
	if top.Done() {
		m.states = m.states[:len(m.states)-1]
	}
	if next != nil {
		m.PushState(next)
	}
	return len(m.states) > 0
}

// ProcessCommands routes all pending commands to the current state.
func (m *SimulationStateMachine) ProcessCommands(cmds []*transport.Command) {
	if len(m.states) == 0 {
		return
	}
	top := m.states[len(m.states)-1]
	for _, cmd := range cmds {
		top.ProcessCommand(cmd)
	}
}
