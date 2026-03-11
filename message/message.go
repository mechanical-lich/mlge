package message

import "github.com/mechanical-lich/mlge/event"

const MessageEventType event.EventType = "MessageEvent"

type MessageEvent struct {
	Sender  string
	Message string
}

func (e MessageEvent) GetType() event.EventType {
	return MessageEventType
}

func PostMessage(sender, message string) {
	event.GetQueuedInstance().QueueEvent(MessageEvent{
		Sender:  sender,
		Message: message,
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
