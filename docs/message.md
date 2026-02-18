---
layout: default
title: Message Log
nav_order: 15
---

# Message Log

`github.com/mechanical-lich/mlge/message`

A simple in-game message posting system that dispatches messages through the queued event system.

## Usage

### Posting Messages

```go
import "github.com/mechanical-lich/mlge/message"

message.PostMessage("System", "Game saved successfully")
message.PostMessage("Combat", "You dealt 15 damage to the goblin")
```

`PostMessage` creates a `MessageEvent` and queues it through the queued event manager.

### Message Log

A global log stores all messages:

```go
for _, msg := range message.MessageLog {
    fmt.Println(msg)
}
```

### Adding Messages Directly

To add a message to the log without dispatching an event:

```go
message.AddMessage("Direct log entry")
```

## MessageEvent

```go
type MessageEvent struct {
    Sender  string
    Message string
}
```

Event type: `"MessageEvent"`

Listen for message events through the event system to display them in your UI:

```go
em := event.GetInstance()
em.RegisterListener("MessageEvent", &MyMessageHandler{})
```
