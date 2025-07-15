package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mechanical-lich/mlge/event"
)

// The InputManager checks for mouse and keyboard input and dispatches events accordingly.
type InputManager struct {
	eventManager   event.EventManagerInterface
	oldMouseX      int
	oldMouseY      int
	oldMouseWheelX float64
	oldMouseWheelY float64
}

func NewInputManager(eventManager event.EventManagerInterface) *InputManager {
	return &InputManager{
		eventManager: eventManager,
	}
}

// Check for mouse input and dispatch events.
func (im *InputManager) HandleInput() {
	cX, cY := ebiten.CursorPosition()

	// Mouse move
	if cX != im.oldMouseX || cY != im.oldMouseY {
		im.eventManager.SendEvent(MouseMoveEvent{X: cX, Y: cY, OldX: im.oldMouseX, OldY: im.oldMouseY})
		im.oldMouseX = cX
		im.oldMouseY = cY
	}

	// Mouse wheel
	wheelX, wheelY := ebiten.Wheel()
	if wheelX != im.oldMouseWheelX || wheelY != im.oldMouseWheelY {
		im.eventManager.SendEvent(MouseWheelEvent{X: wheelX, Y: wheelY, OldX: im.oldMouseWheelX, OldY: im.oldMouseWheelY})
		im.oldMouseWheelX = wheelX
		im.oldMouseWheelY = wheelY
	}

	// Mouse clicks
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		im.eventManager.SendEvent(MouseClickEvent{Button: ebiten.MouseButtonLeft, X: cX, Y: cY})
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		im.eventManager.SendEvent(MouseClickEvent{Button: ebiten.MouseButtonRight, X: cX, Y: cY})
	}

	// Mouse releases
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		im.eventManager.SendEvent(MouseReleasedEvent{Button: ebiten.MouseButtonLeft, X: cX, Y: cY})
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
		im.eventManager.SendEvent(MouseReleasedEvent{Button: ebiten.MouseButtonRight, X: cX, Y: cY})
	}

	// Key presses
	keys := inpututil.AppendPressedKeys(nil)
	if len(keys) > 0 {
		justPressed := inpututil.IsKeyJustPressed(keys[0])
		im.eventManager.SendEvent(KeyPressEvent{Keys: keys, JustPressed: justPressed})
	}

	// Key releases
	keys = inpututil.AppendJustReleasedKeys(nil)
	if len(keys) > 0 {
		im.eventManager.SendEvent(KeyReleaseEvent{Keys: keys})
	}
}
