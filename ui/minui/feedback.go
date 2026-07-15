package minui

import "github.com/mechanical-lich/mlge/event"

// InteractionSound, if set, is called the instant an interactive widget
// registers a genuine user action — BEFORE the widget's OnClick/OnChange/OnSelect
// handler runs. It exists so audio (or other) feedback is immediate even when a
// handler does slow work (e.g. generating a world): the callback can start a
// sound before the handler blocks the frame.
//
// It is deliberately separate from the queued event bus. Bus events are for
// semantic game reactions, are dispatched a frame later, and also fire for
// programmatic setters; this hook fires synchronously and only from real input
// paths. kind matches the EventType the corresponding bus event uses (e.g.
// EventTypeButtonClick), so a handler can share one event→sound mapping.
//
// The hook must stay cheap and non-reentrant (it runs mid-Update): start a
// sound, don't mutate UI state.
var InteractionSound func(kind event.EventType, elementID string)

// playInteraction invokes the InteractionSound hook if one is registered.
func playInteraction(kind event.EventType, elementID string) {
	if InteractionSound != nil {
		InteractionSound(kind, elementID)
	}
}
