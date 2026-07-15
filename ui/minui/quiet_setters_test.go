package minui

import (
	"testing"

	"github.com/mechanical-lich/mlge/event"
	"github.com/stretchr/testify/assert"
)

// countingListener records how many events of the registered type it received.
type countingListener struct{ count int }

func (l *countingListener) HandleEvent(event.EventData) error {
	l.count++
	return nil
}

// TestTabPanel_SetActiveTabQuiet verifies the loud setter emits a TabChangeEvent
// while the quiet setter changes the tab (and fires OnTabChange) without one.
func TestTabPanel_SetActiveTabQuiet(t *testing.T) {
	tp := NewTabPanel("tabs", 200, 100)
	tp.AddTab("a", "A", nil)
	tp.AddTab("b", "B", nil)
	tp.AddTab("c", "C", nil)

	var cbFires int
	tp.OnTabChange = func(string) { cbFires++ }

	listener := &countingListener{}
	q := event.GetQueuedInstance()
	q.RegisterListener(listener, EventTypeTabChange)
	defer q.UnregisterListener(listener, EventTypeTabChange)

	// Loud: emits.
	tp.SetActiveTab("b")
	q.HandleQueue()
	assert.Equal(t, "b", tp.ActiveTabID)
	assert.Equal(t, 1, listener.count, "SetActiveTab should emit a TabChangeEvent")

	// Quiet: changes tab + fires callback, but no event.
	tp.SetActiveTabQuiet("c")
	q.HandleQueue()
	assert.Equal(t, "c", tp.ActiveTabID)
	assert.Equal(t, 1, listener.count, "SetActiveTabQuiet must not emit a TabChangeEvent")
	assert.Equal(t, 2, cbFires, "OnTabChange should still fire for the quiet setter")
}

// TestRadioGroup_SelectByIDQuiet verifies the same split for radio groups.
func TestRadioGroup_SelectByIDQuiet(t *testing.T) {
	rg := NewRadioGroup("rg")
	rg.AddButton(NewRadioButton("r1", "One"))
	rg.AddButton(NewRadioButton("r2", "Two"))

	var cbFires int
	rg.OnSelectionChange = func(string, *RadioButton) { cbFires++ }

	listener := &countingListener{}
	q := event.GetQueuedInstance()
	q.RegisterListener(listener, EventTypeRadioGroupChange)
	defer q.UnregisterListener(listener, EventTypeRadioGroupChange)

	rg.SelectByID("r1") // loud
	q.HandleQueue()
	assert.Equal(t, "r1", rg.selectedID)
	assert.Equal(t, 1, listener.count, "SelectByID should emit a RadioGroupChangeEvent")

	rg.SelectByIDQuiet("r2") // quiet
	q.HandleQueue()
	assert.Equal(t, "r2", rg.selectedID)
	assert.Equal(t, 1, listener.count, "SelectByIDQuiet must not emit a RadioGroupChangeEvent")
	assert.Equal(t, 2, cbFires, "OnSelectionChange should still fire for the quiet setter")
}
