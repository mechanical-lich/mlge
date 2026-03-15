---
layout: default
title: Message Log
nav_order: 15
---

# Message Log

`github.com/mechanical-lich/mlge/message`

An in-game message posting system that dispatches messages through the queued event system. Messages can carry an optional category tag, and individual tags can be suppressed globally.

## Posting Messages

### PostMessage

```go
func PostMessage(sender, message string)
```

Posts an untagged message. Equivalent to calling `PostTaggedMessage` with an empty tag.

```go
import "github.com/mechanical-lich/mlge/message"

message.PostMessage("System", "Game saved successfully")
```

### PostTaggedMessage

```go
func PostTaggedMessage(tag, sender, message string)
```

Posts a message with a category tag. Messages whose tag is present in `SuppressedTags` are silently dropped without dispatching an event.

```go
message.PostTaggedMessage("combat", "Goblin", "hit you for 5 damage")
message.PostTaggedMessage("world", "System", "You entered the dungeon")
```

### SuppressedTags

```go
var SuppressedTags = map[string]bool{}
```

A global map of tags to suppress. Set a tag to `true` to silently drop all messages with that tag. Useful for disabling verbose categories (e.g. combat spam) during testing or in certain game states.

```go
message.SuppressedTags["combat"] = true  // suppress all combat messages
message.SuppressedTags["combat"] = false // re-enable combat messages
```

## MessageEvent

```go
type MessageEvent struct {
    Sender  string
    Message string
    Tag     string
}
```

Event type: `"MessageEvent"`

`Tag` is empty for messages posted via `PostMessage`. Listen for message events through the event system to display them in your UI:

```go
em := event.GetInstance()
em.RegisterListener("MessageEvent", &MyMessageHandler{})
```

Inside your handler you can filter by tag:

```go
func (h *MyMessageHandler) HandleEvent(e event.Event) {
    msg := e.(message.MessageEvent)
    if msg.Tag == "combat" {
        // display in combat log panel
    } else {
        // display in general log
    }
}
```

## Message Log

A global log stores up to 1000 raw string entries. The oldest half is trimmed when the cap is reached.

```go
for _, msg := range message.MessageLog {
    fmt.Println(msg)
}
```

### AddMessage

```go
func AddMessage(x string)
```

Appends a string directly to `MessageLog` without queuing an event. Use this when you want to record something in the log without notifying event listeners.
