package minui

import (
	"github.com/mechanical-lich/mlge/event"
)

// UI Event Types
const (
	EventTypeButtonClick       event.EventType = "ui.button.click"
	EventTypeCheckboxChange    event.EventType = "ui.checkbox.change"
	EventTypeRadioButtonChange event.EventType = "ui.radiobutton.change"
	EventTypeRadioGroupChange  event.EventType = "ui.radiogroup.change"
	EventTypeTextInputChange   event.EventType = "ui.textinput.change"
	EventTypeTextInputSubmit   event.EventType = "ui.textinput.submit"
	EventTypeListBoxSelect     event.EventType = "ui.listbox.select"
	EventTypeSelectBoxChange   event.EventType = "ui.selectbox.change"
	EventTypeModalClose        event.EventType = "ui.modal.close"
	EventTypeElementFocus      event.EventType = "ui.element.focus"
	EventTypeElementBlur       event.EventType = "ui.element.blur"
	EventTypeElementHover      event.EventType = "ui.element.hover"
)

// ButtonClickEvent is fired when a button is clicked
type ButtonClickEvent struct {
	ButtonID string
	Button   *Button
}

func (e ButtonClickEvent) GetType() event.EventType {
	return EventTypeButtonClick
}

// CheckboxChangeEvent is fired when a checkbox state changes
type CheckboxChangeEvent struct {
	CheckboxID string
	Checkbox   *Checkbox
	Checked    bool
}

func (e CheckboxChangeEvent) GetType() event.EventType {
	return EventTypeCheckboxChange
}

// TextInputChangeEvent is fired when text input content changes
type TextInputChangeEvent struct {
	InputID string
	Input   *TextInput
	Text    string
	OldText string
}

func (e TextInputChangeEvent) GetType() event.EventType {
	return EventTypeTextInputChange
}

// TextInputSubmitEvent is fired when text input is submitted (Enter key)
type TextInputSubmitEvent struct {
	InputID string
	Input   *TextInput
	Text    string
}

func (e TextInputSubmitEvent) GetType() event.EventType {
	return EventTypeTextInputSubmit
}

// ListBoxSelectEvent is fired when a list item is selected
type ListBoxSelectEvent struct {
	ListBoxID     string
	ListBox       *ListBox
	SelectedIndex int
	SelectedItem  string
}

func (e ListBoxSelectEvent) GetType() event.EventType {
	return EventTypeListBoxSelect
}

// SelectBoxChangeEvent is fired when the SelectBox selection changes
type SelectBoxChangeEvent struct {
	SelectBoxID   string
	SelectBox     *SelectBox
	SelectedIndex int
	SelectedItem  string
}

func (e SelectBoxChangeEvent) GetType() event.EventType {
	return EventTypeSelectBoxChange
}

// ModalCloseEvent is fired when a modal is closed
type ModalCloseEvent struct {
	ModalID string
	Modal   *Modal
}

func (e ModalCloseEvent) GetType() event.EventType {
	return EventTypeModalClose
}

// ElementFocusEvent is fired when an element gains focus
type ElementFocusEvent struct {
	ElementID string
	Element   Element
}

func (e ElementFocusEvent) GetType() event.EventType {
	return EventTypeElementFocus
}

// ElementBlurEvent is fired when an element loses focus
type ElementBlurEvent struct {
	ElementID string
	Element   Element
}

func (e ElementBlurEvent) GetType() event.EventType {
	return EventTypeElementBlur
}

// ElementHoverEvent is fired when an element is hovered
type ElementHoverEvent struct {
	ElementID string
	Element   Element
	Hovered   bool // true when hover starts, false when hover ends
}

func (e ElementHoverEvent) GetType() event.EventType {
	return EventTypeElementHover
}

// RadioButtonChangeEvent is fired when a radio button state changes
type RadioButtonChangeEvent struct {
	RadioButtonID string
	RadioButton   *RadioButton
	Selected      bool
}

func (e RadioButtonChangeEvent) GetType() event.EventType {
	return EventTypeRadioButtonChange
}

// RadioGroupChangeEvent is fired when a radio group selection changes
type RadioGroupChangeEvent struct {
	RadioGroupID   string
	RadioGroup     *RadioGroup
	SelectedID     string
	SelectedButton *RadioButton
}

func (e RadioGroupChangeEvent) GetType() event.EventType {
	return EventTypeRadioGroupChange
}
