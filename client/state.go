package client

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/transport"
)

// ClientState is the presentation-layer equivalent of state.StateInterface.
//
// Differences from state.StateInterface:
//   - Update receives the latest [transport.Snapshot] (may be nil if no new
//     snapshot arrived this frame) so the state can apply server authority.
//   - Draw is identical to state.StateInterface.
//   - Done/push semantics are the same: return non-nil to push a new state;
//     return true from Done to pop.
//
// The state machine is driven by [Client.Update] and [Client.Draw].
type ClientState interface {
	// Update is called every Ebitengine frame.
	// snapshot is the latest decoded snapshot from the server, or nil.
	// Return a non-nil ClientState to push it on the state stack.
	Update(snapshot *transport.Snapshot) ClientState

	// Draw renders the current state to screen.
	Draw(screen *ebiten.Image)

	// Done signals the Client to pop this state.
	Done() bool
}

// clientStateMachine is the client-side state stack.
type clientStateMachine struct {
	states []ClientState
}

func (m *clientStateMachine) PushState(s ClientState) {
	m.states = append(m.states, s)
}

func (m *clientStateMachine) Current() ClientState {
	if len(m.states) == 0 {
		return nil
	}
	return m.states[len(m.states)-1]
}

func (m *clientStateMachine) Update(snapshot *transport.Snapshot) {
	if len(m.states) == 0 {
		return
	}
	top := m.states[len(m.states)-1]
	next := top.Update(snapshot)
	if top.Done() {
		m.states = m.states[:len(m.states)-1]
	}
	if next != nil {
		m.PushState(next)
	}
}

func (m *clientStateMachine) Draw(screen *ebiten.Image) {
	if len(m.states) == 0 {
		return
	}
	m.states[len(m.states)-1].Draw(screen)
}
