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

func AddMessage(x string) {
	MessageLog = append(MessageLog, x)
}
