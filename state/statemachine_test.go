package state

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stretchr/testify/assert"
)

type TestState struct {
	Value      string
	IsDone     bool
	UpdateFunc func() StateInterface
}

func (s *TestState) Update() StateInterface {
	if s.UpdateFunc != nil {
		return s.UpdateFunc()
	}
	return nil
}

func (s *TestState) Draw(screen *ebiten.Image) {
	s.Value = "Drawn"
}

func (s *TestState) Done() bool {
	return s.IsDone
}

func TestStateUpdates(t *testing.T) {
	sm := StateMachine{}

	s := &TestState{}
	s.UpdateFunc = func() StateInterface {
		s.Value = "Updated"
		return nil
	}
	sm.PushState(s)

	sm.Update()

	assert.Equal(t, "Updated", s.Value, "Failed to updated state")
}

func TestStateDraw(t *testing.T) {
	sm := StateMachine{}

	s := &TestState{}
	sm.PushState(s)

	sm.Draw(&ebiten.Image{})

	assert.Equal(t, "Drawn", s.Value, "Failed to draw state")
}
func TestStatePopsWhenDone(t *testing.T) {
	sm := StateMachine{}

	s := &TestState{}
	s.UpdateFunc = func() StateInterface {
		s.IsDone = true
		return nil
	}
	sm.PushState(s)

	sm.Update()

	assert.Empty(t, sm.states, "Failed to pop finished state")
}
func TestStatePushesNewState(t *testing.T) {
	sm := StateMachine{}

	s := &TestState{}
	s2 := &TestState{}
	s2.UpdateFunc = func() StateInterface {
		s2.Value = "Updated"
		return nil
	}

	s.UpdateFunc = func() StateInterface {
		return s2
	}

	sm.PushState(s)

	sm.Update()

	assert.Len(t, sm.states, 2, "Failed to push new state")
	assert.Equal(t, sm.states[sm.currentState], s2, "Current state is not the most recently pushed state")
	assert.Empty(t, s2.Value, "State data updated before update called")
	sm.Update()
	assert.Equal(t, "Updated", s2.Value, "Failed to update state 2 data")
}
