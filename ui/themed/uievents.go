package ui

import (
	"github.com/mechanical-lich/mlge/event"
)

// UI Event Types
const (
	EventTypeUIClick       event.EventType = "ui_click"
	EventTypeUIValueChange event.EventType = "ui_value_change"
	EventTypeUIFocus       event.EventType = "ui_focus"
	EventTypeUIBlur        event.EventType = "ui_blur"
	EventTypeUIModalOpen   event.EventType = "ui_modal_open"
	EventTypeUIModalClose  event.EventType = "ui_modal_close"
	EventTypeUIHover       event.EventType = "ui_hover"
	EventTypeUISubmit      event.EventType = "ui_submit"
	EventTypeUITabChange   event.EventType = "ui_tab_change"
)

// ClickEventData represents a UI click event
type ClickEventData struct {
	SourceName string
	Data       map[string]interface{}
}

func (e ClickEventData) GetType() event.EventType {
	return EventTypeUIClick
}

// ValueChangeEventData represents a UI value change event
type ValueChangeEventData struct {
	SourceName string
	OldValue   interface{}
	NewValue   interface{}
}

func (e ValueChangeEventData) GetType() event.EventType {
	return EventTypeUIValueChange
}

// FocusEventData represents a UI focus event
type FocusEventData struct {
	SourceName string
}

func (e FocusEventData) GetType() event.EventType {
	return EventTypeUIFocus
}

// BlurEventData represents a UI blur event
type BlurEventData struct {
	SourceName string
}

func (e BlurEventData) GetType() event.EventType {
	return EventTypeUIBlur
}

// ModalOpenEventData represents a modal open event
type ModalOpenEventData struct {
	ModalName string
}

func (e ModalOpenEventData) GetType() event.EventType {
	return EventTypeUIModalOpen
}

// ModalCloseEventData represents a modal close event
type ModalCloseEventData struct {
	ModalName string
}

func (e ModalCloseEventData) GetType() event.EventType {
	return EventTypeUIModalClose
}

// SubmitEventData represents a form submit event
type SubmitEventData struct {
	SourceName string
	Data       map[string]interface{}
}

func (e SubmitEventData) GetType() event.EventType {
	return EventTypeUISubmit
}
