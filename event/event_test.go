package event

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	Test EventType = iota
)

type TestEventData struct {
}

func (t TestEventData) GetType() EventType {
	return Test
}

type TestListener struct {
}

func (t *TestListener) HandleEvent(data EventData) error {
	fmt.Println("Handling event: ", data)
	return nil
}

func TestRegisterListener(t *testing.T) {
	m := &QueuedEventManager{}

	testListener := &TestListener{}

	m.RegisterListener(testListener, Test)
	assert.Equal(t, 1, len(m.listeners[Test]))
}
func TestUnregisterListener(t *testing.T) {
	m := &QueuedEventManager{}

	testListener := &TestListener{}

	m.RegisterListener(testListener, Test)
	assert.Equal(t, 1, len(m.listeners[Test]))
	m.UnregisterListener(testListener, Test)
	assert.Equal(t, 0, len(m.listeners[Test]))
}
func TestUnregisterListenerFromAll(t *testing.T) {
	m := &QueuedEventManager{}

	testListener := &TestListener{}

	m.RegisterListener(testListener, Test)
	assert.Equal(t, 1, len(m.listeners[Test]))
	m.UnregisterListenerFromAll(testListener)
	assert.Equal(t, 0, len(m.listeners[Test]))
}
