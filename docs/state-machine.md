---
layout: default
title: State Machine
nav_order: 9
---

# State Machine

`github.com/mechanical-lich/mlge/state`

A stack-based state machine for managing game states such as menus, gameplay, pause screens, and cutscenes.

## StateInterface

```go
type StateInterface interface {
    Update() StateInterface
    Draw(screen *ebiten.Image)
    Done() bool
}
```

| Method | Description |
|--------|-------------|
| `Update` | Called each frame. Returns a new state to push, or `nil` to continue |
| `Draw` | Renders the state to the screen |
| `Done` | Returns `true` when the state should be popped |

## StateMachine

```go
type StateMachine struct{}
```

**Methods:**

| Method | Signature | Description |
|--------|-----------|-------------|
| `PushState` | `(state StateInterface)` | Push a new state onto the stack |
| `PopCurrentState` | `()` | Remove the top state |
| `Update` | `()` | Updates the top state; handles push/pop transitions |
| `Draw` | `(screen *ebiten.Image)` | Draws the top state |

## Usage

```go
sm := &state.StateMachine{}
sm.PushState(&MenuState{})

// In game loop
func (g *Game) Update() error {
    g.stateMachine.Update()
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    g.stateMachine.Draw(screen)
}
```

## Example: Game States

```go
type MenuState struct {
    done     bool
    nextState state.StateInterface
}

func (m *MenuState) Update() state.StateInterface {
    if startPressed {
        m.done = true
        return &PlayState{}  // Push PlayState
    }
    return nil
}

func (m *MenuState) Draw(screen *ebiten.Image) {
    // Draw menu UI
}

func (m *MenuState) Done() bool {
    return m.done
}
```

When `Update()` returns a non-nil state, the state machine pushes it onto the stack. When `Done()` returns `true`, the state is popped, revealing the previous state underneath.
