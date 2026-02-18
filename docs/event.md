---
layout: default
title: Event System
nav_order: 4
---

# Event System

`github.com/mechanical-lich/mlge/event`

A publish/subscribe event system with both immediate and queued dispatch modes.

## Core Types

### EventType

```go
type EventType string
```

String identifier for event types.

### EventData

```go
type EventData interface {
    GetType() EventType
}
```

All event payloads must implement this interface.

### EventListener

```go
type EventListener interface {
    HandleEvent(EventData) error
}
```

Implement this interface to receive events.

## Event Managers

### EventManager (Immediate Dispatch)

Events are dispatched immediately to listeners when `SendEvent` is called. Access the singleton with:

```go
em := event.GetInstance()
```

**Methods:**

| Method | Signature | Description |
|--------|-----------|-------------|
| `RegisterListener` | `(eventType EventType, listener EventListener)` | Subscribe to an event type |
| `SendEvent` | `(event EventData) error` | Dispatch immediately to all listeners |
| `UnregisterListener` | `(eventType EventType, listener EventListener)` | Unsubscribe from a specific event type |
| `UnregisterListenerFromAll` | `(listener EventListener)` | Unsubscribe from all event types |

### QueuedEventManager

Extends `EventManager` with a queue. Events are collected during the frame and flushed at a controlled point. Access the singleton with:

```go
qem := event.GetQueuedInstance()
```

**Additional Methods:**

| Method | Signature | Description |
|--------|-----------|-------------|
| `QueueEvent` | `(event EventData)` | Add an event to the queue |
| `HandleQueue` | `()` | Flush and dispatch all queued events |

## Usage

### Defining an Event

```go
const DamageEventType event.EventType = "DamageEvent"

type DamageEvent struct {
    TargetID int
    Amount   int
}

func (e *DamageEvent) GetType() event.EventType {
    return DamageEventType
}
```

### Listening for Events

```go
type HealthSystem struct{}

func (h *HealthSystem) HandleEvent(data event.EventData) error {
    dmg := data.(*DamageEvent)
    // Apply damage to target
    return nil
}

// Register
em := event.GetInstance()
em.RegisterListener(DamageEventType, &HealthSystem{})
```

### Sending Events

```go
// Immediate dispatch
em.SendEvent(&DamageEvent{TargetID: 1, Amount: 10})

// Queued dispatch (preferred for game loops)
qem := event.GetQueuedInstance()
qem.QueueEvent(&DamageEvent{TargetID: 1, Amount: 10})

// Later in the update loop
qem.HandleQueue()
```

### Queued vs Immediate

Use **queued dispatch** when events are generated during system updates and should be processed after all systems have run. This prevents cascading side effects during a single frame update.

Use **immediate dispatch** for events that must be handled right away, such as critical input or synchronization events.
