package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mechanical-lich/mlge/event"
)

const (
	MouseClickEventType    event.EventType = "MouseClick"
	MouseReleasedEventType event.EventType = "MouseReleased"
	MouseMoveEventType     event.EventType = "MouseMove"
	MouseWheelEventType    event.EventType = "MouseWheel"
	KeyPressEventType      event.EventType = "KeyPress"
	KeyReleaseEventType    event.EventType = "KeyRelease"
)

// MouseClickEvent represents a mouse click event.
type MouseClickEvent struct {
	Button ebiten.MouseButton
	X, Y   int
}

func (e MouseClickEvent) GetType() event.EventType {
	return MouseClickEventType
}

// MouseReleasedEvent represents a mouse release event.
type MouseReleasedEvent struct {
	Button ebiten.MouseButton
	X, Y   int
}

func (e MouseReleasedEvent) GetType() event.EventType {
	return MouseReleasedEventType
}

// MouseMoveEvent represents a mouse move event.
type MouseWheelEvent struct {
	X, Y       float64
	OldX, OldY float64
}

func (e MouseWheelEvent) GetType() event.EventType {
	return MouseWheelEventType
}

// MouseMoveEvent represents a mouse move event.
type MouseMoveEvent struct {
	X, Y       int
	OldX, OldY int
}

func (e MouseMoveEvent) GetType() event.EventType {
	return MouseMoveEventType
}

// KeyPressEvent represents a key press event.
type KeyPressEvent struct {
	Keys        []ebiten.Key
	JustPressed bool
}

func (e KeyPressEvent) KeyPressed(key ebiten.Key) bool {
	for _, k := range e.Keys {
		if k == key {
			return true
		}
	}
	return false
}

func (e KeyPressEvent) GetType() event.EventType {
	return KeyPressEventType
}

// KeyReleaseEvent represents a key release event.
type KeyReleaseEvent struct {
	Keys []ebiten.Key
}

func (e KeyReleaseEvent) GetType() event.EventType {
	return KeyReleaseEventType
}
