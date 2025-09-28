//go:build js && wasm

// Package component provides a base for all components
// This file provides an interface for all components and an API to create them

package vdom

import (
	"syscall/js"

	"github.com/AureClai/vortex/pkg/style"
)

type Component interface {
	Render() *VNode

	// Lifecycle methods
	OnMount()   // Called when the component is first mounted
	OnUpdate()  // Called when the component needs to update
	OnUnmount() // Called when the component is removed

	// Identiy
	GetID() string
	SetID(id string)
}

// ComponnentBase
type ComponentBase struct {
	id      string
	vNode   *VNode
	mounted bool
}

func NewComponent(tag string) *ComponentBase {
	return &ComponentBase{
		vNode: &VNode{
			Type:          VNodeElement,
			Tag:           tag,
			Props:         make(map[string]interface{}),
			EventHandlers: make(map[string]func(event js.Value)),
			Children:      make([]*VNode, 0),
		},
	}
}

// Implement Component interface
func (c *ComponentBase) Render() *VNode {
	return c.vNode
}

func (c *ComponentBase) OnMount()        { /* Default: no-op */ }
func (c *ComponentBase) OnUpdate()       { /* Default: no-op */ }
func (c *ComponentBase) OnUnmount()      { /* Default: no-op */ }
func (c *ComponentBase) GetID() string   { return c.id }
func (c *ComponentBase) SetID(id string) { c.id = id }

// Fluent API methods returns *ComponentBase for chaining
func (c *ComponentBase) SetClass(class string) *ComponentBase {
	c.vNode.Props["class"] = class
	return c
}

func (c *ComponentBase) Style(s *style.Style) *ComponentBase {
	c.vNode.AppliedStyle = s
	return c
}

func (c *ComponentBase) On(event string, handler func(event Event)) *ComponentBase {
	c.vNode.EventHandlers[event] = func(jsEvent js.Value) {
		// Convert js.Value to appropriate vdom.Event type based on event name
		var vdomEvent Event
		switch event {
		case "input", "change":
			vdomEvent = NewInputEvent(jsEvent, c.vNode)
		case "click", "mousedown", "mouseup", "mouseover", "mouseout":
			vdomEvent = NewMouseEvent(jsEvent, c.vNode)
		case "keydown", "keyup", "keypress":
			vdomEvent = NewKeyboardEvent(jsEvent, c.vNode)
		case "focus", "blur":
			vdomEvent = NewFocusEvent(jsEvent, c.vNode)
		default:
			vdomEvent = NewEvent(jsEvent, c.vNode) // Base event
		}
		handler(vdomEvent)
	}
	return c
}

func (c *ComponentBase) SetText(text string) *ComponentBase {
	c.vNode.Text = text
	return c
}

func (c *ComponentBase) AddChild(child Component) *ComponentBase {
	c.vNode.Children = append(c.vNode.Children, child.Render())
	return c
}

func (c *ComponentBase) AddChildren(children ...Component) *ComponentBase {
	for _, child := range children {
		c.vNode.Children = append(c.vNode.Children, child.Render())
	}
	return c
}

func (c *ComponentBase) AddChildrenMapping(mappingFunc func() []Component) *ComponentBase {
	children := mappingFunc()
	c.AddChildren(children...)
	return c
}
