package event

import (
	"fmt"
)

type EventType string

// EventData - interface representing event data.  Each data struct needs to return its type as an EventType (int)
type EventData interface {
	GetType() EventType
}

// EventListener - interface representing an event listener.  Each event listener needs a HandleEvent function
type EventListener interface {
	HandleEvent(EventData) error
}

type EventManagerInterface interface {
	RegisterListener(EventListener, EventType)
	SendEvent(EventData)
	UnregisterListener(EventListener, EventType)
	UnregisterListenerFromAll(EventListener)
}

// EventManager - Entry point for registering event listeners and sending events
type EventManager struct {
	listeners map[EventType][]EventListener // Key is the event
}

// RegisterListener - registers a struct that can be represented by the EventListener
// interface to the specified event type
func (m *EventManager) RegisterListener(listener EventListener, eventType EventType) {
	if m.listeners == nil {
		m.listeners = make(map[EventType][]EventListener)
	}

	m.listeners[eventType] = append(m.listeners[eventType], listener)
}

// SendEvent - Sends the event data to the event listeners registered to its type.
func (m *EventManager) SendEvent(data EventData) {
	if m.listeners == nil {
		m.listeners = make(map[EventType][]EventListener)
		return
	}

	for _, v := range m.listeners[data.GetType()] {
		err := v.HandleEvent(data)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// UnregisterListener - unregisters a struct that can be represented by the EventListener
// interface from the specified event type's events
func (m *EventManager) UnregisterListener(listener EventListener, eventType EventType) {
	if m.listeners == nil {
		m.listeners = make(map[EventType][]EventListener)
		return
	}

	for k := 0; k < len(m.listeners[eventType]); k++ {
		if m.listeners[eventType][k] == listener {
			m.listeners[eventType] = append(m.listeners[eventType][:k], m.listeners[eventType][k+1:]...)
			k--
		}
	}
}

// UnregisterListenerFromAll - unregisters a struct that can be represented by the EventListener
// interface from all events
func (m *EventManager) UnregisterListenerFromAll(listener EventListener) {
	if m.listeners == nil {
		m.listeners = make(map[EventType][]EventListener)
		return
	}

	for eventType := range m.listeners {
		for k := 0; k < len(m.listeners[eventType]); k++ {
			if m.listeners[eventType][k] == listener {
				m.listeners[eventType] = append(m.listeners[eventType][:k], m.listeners[eventType][k+1:]...)
				k--
			}
		}
	}

}

type QueuedEventManager struct {
	EventManager
	eventQueue []EventData
}

// QueueEvent - Queues an event to be sent to listeners during the manger's handle call.
func (m *QueuedEventManager) QueueEvent(event EventData) {
	if m.eventQueue == nil {
		m.eventQueue = make([]EventData, 0)
	}
	m.eventQueue = append(m.eventQueue, event)
}

// HandleQueue - Sends all the current events in the queue, but doesn't send new events
func (m *QueuedEventManager) HandleQueue() {
	if m.eventQueue == nil {
		m.eventQueue = make([]EventData, 0)
		return
	}
	// We don't want to handle new events generated while handling these events or we'll get an endless loop.
	c := len(m.eventQueue)
	for len(m.eventQueue) > 0 && c > 0 {
		m.SendEvent(m.eventQueue[0])
		m.eventQueue = m.eventQueue[1:]
		c--
	}
}

// Singleton instance
var queuedEventSingleton *QueuedEventManager
var eventSingleton *EventManager

// GetInstance - Returns the singleton instance of EventManager
func GetInstance() *EventManager {
	if eventSingleton == nil {
		eventSingleton = &EventManager{}
	}
	return eventSingleton
}

// GetInstance - Returns the singleton instance of QueuedEventManager
func GetQueuedInstance() *QueuedEventManager {
	if queuedEventSingleton == nil {
		queuedEventSingleton = &QueuedEventManager{}
	}
	return queuedEventSingleton
}
