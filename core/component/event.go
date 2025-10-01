//go:build js && wasm

// Event is a struct that contains the event data

package component

import (
	"strconv"
	"syscall/js"
)

type Event interface {
	RawEvent() js.Value
	Target() js.Value
	CurrentTarget() js.Value
	PreventDefault()
	StopPropagation()
	StopImmediatePropagation()

	// Event properties
	Type() string
	Bubbles() bool
	Cancelable() bool
	TimeStamp() float64
}

// Base event implementation
type BaseEvent struct {
	event js.Value
	vnode *VNode // The vnode that triggered the event
}

func NewEvent(jsEvent js.Value, vnode *VNode) *BaseEvent {
	return &BaseEvent{
		event: jsEvent,
		vnode: vnode,
	}
}

// Implement the Event interface
func (e *BaseEvent) RawEvent() js.Value        { return e.event }
func (e *BaseEvent) Target() js.Value          { return e.event.Get("target") }
func (e *BaseEvent) CurrentTarget() js.Value   { return e.event.Get("currentTarget") }
func (e *BaseEvent) PreventDefault()           { e.event.Call("preventDefault") }
func (e *BaseEvent) StopPropagation()          { e.event.Call("stopPropagation") }
func (e *BaseEvent) StopImmediatePropagation() { e.event.Call("stopImmediatePropagation") }

func (e *BaseEvent) Type() string       { return e.event.Get("type").String() }
func (e *BaseEvent) Bubbles() bool      { return e.event.Get("bubbles").Bool() }
func (e *BaseEvent) Cancelable() bool   { return e.event.Get("cancelable").Bool() }
func (e *BaseEvent) TimeStamp() float64 { return e.event.Get("timeStamp").Float() }

func (e *BaseEvent) VNode() *VNode { return e.vnode }

// ======
// SPECIALIZED EVENTS TYPES
// ======

// MouseEvent
type MouseEvent struct {
	*BaseEvent
}

func NewMouseEvent(jsEvent js.Value, vnode *VNode) *MouseEvent {
	return &MouseEvent{
		BaseEvent: NewEvent(jsEvent, vnode)}
}

func (m *MouseEvent) ClientX() int   { return m.event.Get("clientX").Int() }
func (m *MouseEvent) ClientY() int   { return m.event.Get("clientY").Int() }
func (m *MouseEvent) PageX() int     { return m.event.Get("pageX").Int() }
func (m *MouseEvent) PageY() int     { return m.event.Get("pageY").Int() }
func (m *MouseEvent) ScreenX() int   { return m.event.Get("screenX").Int() }
func (m *MouseEvent) ScreenY() int   { return m.event.Get("screenY").Int() }
func (m *MouseEvent) Button() int    { return m.event.Get("button").Int() }
func (m *MouseEvent) Buttons() int   { return m.event.Get("buttons").Int() }
func (m *MouseEvent) CtrlKey() bool  { return m.event.Get("ctrlKey").Bool() }
func (m *MouseEvent) ShiftKey() bool { return m.event.Get("shiftKey").Bool() }
func (m *MouseEvent) AltKey() bool   { return m.event.Get("altKey").Bool() }
func (m *MouseEvent) MetaKey() bool  { return m.event.Get("metaKey").Bool() }

// Keyboard Events
type KeyboardEvent struct {
	*BaseEvent
}

func NewKeyboardEvent(jsEvent js.Value, vnode *VNode) *KeyboardEvent {
	return &KeyboardEvent{BaseEvent: NewEvent(jsEvent, vnode)}
}

func (k *KeyboardEvent) Key() string    { return k.event.Get("key").String() }
func (k *KeyboardEvent) Code() string   { return k.event.Get("code").String() }
func (k *KeyboardEvent) KeyCode() int   { return k.event.Get("keyCode").Int() }
func (k *KeyboardEvent) CharCode() int  { return k.event.Get("charCode").Int() }
func (k *KeyboardEvent) CtrlKey() bool  { return k.event.Get("ctrlKey").Bool() }
func (k *KeyboardEvent) ShiftKey() bool { return k.event.Get("shiftKey").Bool() }
func (k *KeyboardEvent) AltKey() bool   { return k.event.Get("altKey").Bool() }
func (k *KeyboardEvent) MetaKey() bool  { return k.event.Get("metaKey").Bool() }
func (k *KeyboardEvent) Repeat() bool   { return k.event.Get("repeat").Bool() }

// Input Events (for form elements)
type InputEvent struct {
	*BaseEvent
}

func NewInputEvent(jsEvent js.Value, vnode *VNode) *InputEvent {
	return &InputEvent{BaseEvent: NewEvent(jsEvent, vnode)}
}

func (i *InputEvent) Value() string {
	target := i.Target()
	if target.Get("value").Type() != js.TypeUndefined {
		return target.Get("value").String()
	}
	return ""
}

func (i *InputEvent) ValueAsInt() (int, error) {
	return strconv.Atoi(i.Value())
}

func (i *InputEvent) ValueAsFloat() (float64, error) {
	return strconv.ParseFloat(i.Value(), 64)
}

func (i *InputEvent) Checked() bool {
	target := i.Target()
	if target.Get("checked").Type() != js.TypeUndefined {
		return target.Get("checked").Bool()
	}
	return false
}

// Focus Events
type FocusEvent struct {
	*BaseEvent
}

func NewFocusEvent(jsEvent js.Value, vnode *VNode) *FocusEvent {
	return &FocusEvent{BaseEvent: NewEvent(jsEvent, vnode)}
}

func (f *FocusEvent) RelatedTarget() js.Value { return f.event.Get("relatedTarget") }

// =============================================================================
// EVENT HANDLER TYPES
// =============================================================================

type ClickHandler func(*MouseEvent)
type KeyHandler func(*KeyboardEvent)
type InputHandler func(*InputEvent)
type FocusHandler func(*FocusEvent)
type ChangeHandler func(*InputEvent)

// Generic event handler (fallback)
type GenericHandler func(Event)
