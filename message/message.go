package message

import "github.com/mechanical-lich/mlge/event"

const MessageEventType event.EventType = "MessageEvent"

type MessageEvent struct {
	Sender  string
	Message string
	X, Y, Z int // Option location information
	Tag     string
}

func (e MessageEvent) GetType() event.EventType {
	return MessageEventType
}

// SuppressedTags contains message tags that should be silently dropped.
// Set a tag to true to suppress all messages with that matching tag.
var SuppressedTags = map[string]bool{}

// PostMessage posts an untagged message.
func PostMessage(sender, message string) {
	PostTaggedMessage("", sender, message)
}

// PostTaggedMessage posts a message with a category tag.
// Messages whose tag is in SuppressedTags are silently dropped.
func PostTaggedMessage(tag, sender, message string) {
	PostLocatedTaggedMessage(tag, sender, message, 0, 0, 0)
}

// PostLocatedTaggedMessage posts a tagged message with a world location.
// Listeners can use the location to filter messages by visibility.
func PostLocatedTaggedMessage(tag, sender, message string, x, y, z int) {
	if tag != "" && SuppressedTags[tag] {
		return
	}
	event.GetQueuedInstance().QueueEvent(MessageEvent{
		Sender:  sender,
		Message: message,
		Tag:     tag,
		X:       x,
		Y:       y,
		Z:       z,
	})
}

var MessageLog []string

const maxMessageLog = 1000

func AddMessage(x string) {
	MessageLog = append(MessageLog, x)
	if len(MessageLog) > maxMessageLog {
		// Trim oldest messages, copy to new slice to release old backing array
		trimmed := make([]string, maxMessageLog/2)
		copy(trimmed, MessageLog[len(MessageLog)-maxMessageLog/2:])
		MessageLog = trimmed
	}
}
