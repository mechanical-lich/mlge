package transport

// CommandType identifies the kind of command being sent from client to server.
// Games should define their own CommandType constants (e.g., "move", "build", "attack").
// Built-in types are prefixed with "input." to avoid name collisions.
type CommandType string

const (
	// CommandInput carries raw input events translated from mlge's input package.
	// The Payload will be an [InputPayload].
	CommandInput CommandType = "input.raw"
)

// Command is a timestamped, typed message sent from the client to the server.
// It represents player intent: a key press, a mouse click, a game action.
//
// Tick is the client-estimated server tick this command targets.
// For Phase 1 (local transport, no prediction) this is informational.
// In a future TCP implementation it enables client-side prediction reconciliation.
type Command struct {
	Type    CommandType
	Tick    uint64
	Payload any
}

// InputPayload carries a translated input event as a command payload.
// It is set by the client when it converts an mlge input event to a Command.
// The EventType and EventData fields mirror mlge's event.EventType and event.EventData
// but are redeclared here to keep transport/ free of an mlge/event import.
type InputPayload struct {
	// EventType is the mlge event type string (e.g., "MouseLeftClick", "KeyPress").
	EventType string

	// Data is the original mlge event data struct.
	// For local transport this is the live struct pointer (zero copy).
	// A future network transport would serialize this to JSON/proto before sending.
	Data any
}
