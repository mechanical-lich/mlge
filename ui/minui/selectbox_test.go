package minui

import (
	"testing"

	"github.com/mechanical-lich/mlge/event"
	"github.com/stretchr/testify/assert"
)

type testSelectListener struct {
	called bool
	index  int
	item   string
}

func (l *testSelectListener) HandleEvent(data event.EventData) error {
	switch v := data.(type) {
	case SelectBoxChangeEvent:
		l.called = true
		l.index = v.SelectedIndex
		l.item = v.SelectedItem
	}
	return nil
}

func TestSelectBox_SelectByIndex(t *testing.T) {
	sb := NewSelectBox("test-select", []string{"one", "two", "three"})
	listener := &testSelectListener{}
	q := event.GetQueuedInstance()
	q.RegisterListener(listener, EventTypeSelectBoxChange)

	sb.SelectByIndex(1)
	// Handle queued events
	q.HandleQueue()

	assert.True(t, listener.called)
	assert.Equal(t, 1, listener.index)
	assert.Equal(t, "two", listener.item)

	// Selecting out-of-range should clear selection and not fire
	listener.called = false
	sb.SelectByIndex(999)
	q.HandleQueue()
	assert.False(t, listener.called)
}

func TestSelectBox_SetItemsAndSelect(t *testing.T) {
	sb := NewSelectBox("test-select-2", []string{"a", "b"})
	sb.SetItems([]string{"x", "y", "z"})
	sb.SelectByIndex(2)
	idx, item := sb.GetSelectedItem()
	assert.Equal(t, 2, idx)
	assert.Equal(t, "z", item)
}
