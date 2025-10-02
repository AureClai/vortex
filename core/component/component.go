//go:build js && wasm

// Package component provides a base for all components
// This file provides an interface for all components and an API to create them

package component

import (
	"syscall/js"

	"github.com/AureClai/vortex/core/style"
)

// PropsMapper interface for components that can map props to DOM attributes
type PropsMapper interface {
	MapPropsToDOM() map[string]interface{}
}

// Component interface
// It only has the render method
type Component interface {
	Render() *VNode
}

// Lifecycle methods
// this interface is used to define the lifecycle methods for a component
type LifecycleComponent interface {
	// Lifecycle methods
	OnMount()   // Called when the component is first mounted
	OnUpdate()  // Called when the component needs to update
	OnUnmount() // Called when the component is removed
}

// ComponentBase
type ComponentBase struct {
	vNode *VNode
	id    string
}

func NewComponent(tag string) *ComponentBase {
	return &ComponentBase{
		vNode: &VNode{
			Type:          VNodeElement,
			Tag:           tag,
			Attrs:         make(map[string]interface{}),
			EventHandlers: make(map[string]func(event js.Value)),
			Children:      make([]*VNode, 0),
		},
	}
}

// Implement Component interface
func (c *ComponentBase) Render() *VNode {
	return c.vNode
}

// GetVNode - direct access to underlying VNode (use with caution)
func (c *ComponentBase) GetVNode() *VNode {
	return c.vNode
}

func (c *ComponentBase) OnMount()   { /* Default: no-op */ }
func (c *ComponentBase) OnUpdate()  { /* Auto re-render on state changes */ }
func (c *ComponentBase) OnUnmount() { /* Default: no-op */ }

func (c *ComponentBase) Style(s *style.Style) *ComponentBase {
	c.vNode.AppliedStyle = s
	return c
}

func (c *ComponentBase) On(event string, handler func(event Event)) *ComponentBase {
	c.vNode.EventHandlers[event] = func(jsEvent js.Value) {
		// Convert js.Value to appropriate Event type based on event name
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

func (c *ComponentBase) GetText() string {
	return c.vNode.Text
}

// AddChild adds a single child component (renders immediately)
func (c *ComponentBase) AddChild(child Component) *ComponentBase {
	if child != nil {
		if vnode := child.Render(); vnode != nil {
			c.vNode.Children = append(c.vNode.Children, vnode)
		}
	}
	return c
}

// AddChildComponent - stores component reference (for VNodeComponent type)
// This allows the renderer to manage component lifecycle properly
func (c *ComponentBase) AddChildComponent(child Component) *ComponentBase {
	if child != nil {
		childVNode := &VNode{
			Type:      VNodeComponent,
			Component: child,
		}
		c.vNode.Children = append(c.vNode.Children, childVNode)
	}
	return c
}

// AddChildren adds multiple child components (renders immediately)
func (c *ComponentBase) AddChildren(children ...Component) *ComponentBase {
	for _, child := range children {
		if child != nil {
			if vnode := child.Render(); vnode != nil {
				c.vNode.Children = append(c.vNode.Children, vnode)
			}
		}
	}
	return c
}

// AddChildrenComponents - for component references
func (c *ComponentBase) AddChildrenComponents(children ...Component) *ComponentBase {
	for _, child := range children {
		if child != nil {
			c.AddChildComponent(child)
		}
	}
	return c
}

func (c *ComponentBase) AddChildrenMapping(mappingFunc func() []Component) *ComponentBase {
	if mappingFunc != nil {
		children := mappingFunc()
		c.AddChildren(children...)
	}
	return c
}

// ClearChildren - useful for dynamic updates
func (c *ComponentBase) ClearChildren() *ComponentBase {
	c.vNode.Children = make([]*VNode, 0)
	return c
}

func (c *ComponentBase) SetKey(key string) *ComponentBase {
	c.vNode.Key = key
	return c
}

func (c *ComponentBase) GetKey() string {
	return c.vNode.Key
}

func (c *ComponentBase) SetID(id string) *ComponentBase {
	c.id = id
	return c
}

func (c *ComponentBase) GetID() string {
	return c.id
}

func (c *ComponentBase) SetAttr(key string, value interface{}) *ComponentBase {
	c.vNode.Attrs[key] = value
	return c
}

func (c *ComponentBase) GetAttr(key string) (interface{}, bool) {
	val, ok := c.vNode.Attrs[key]
	return val, ok
}

func (c *ComponentBase) SetClass(class string) *ComponentBase {
	c.vNode.Attrs["class"] = class
	return c
}

func (c *ComponentBase) SetID_DOM(domID string) *ComponentBase {
	c.vNode.Attrs["id"] = domID
	return c
}

// ConditionalChild - adds child from mapping function if and else if condition is true
func (c *ComponentBase) ConditionalChild(condition bool, childIf Component) *ComponentBase {
	if condition && childIf != nil {
		c.AddChild(childIf)
	}
	return c
}

// ConditionalChildElse adds childIf if condition is true, otherwise adds childElse
func (c *ComponentBase) ConditionalChildElse(condition bool, childIf Component, childElse Component) *ComponentBase {
	if condition {
		if childIf != nil {
			c.AddChild(childIf)
		}
	} else {
		if childElse != nil {
			c.AddChild(childElse)
		}
	}
	return c
}

// ConditionalChildren - adds children only if condition is true
func (c *ComponentBase) ConditionalChildren(condition bool, children ...Component) *ComponentBase {
	if condition {
		c.AddChildren(children...)
	}
	return c
}

// ConditionalChildrenMapping adds children from mapping function based on condition
// If mappingFuncElse is nil and condition is false, no children are added
func (c *ComponentBase) ConditionalChildrenMapping(condition bool, mappingFuncIf func() []Component, mappingFuncElse func() []Component) *ComponentBase {
	if condition {
		if mappingFuncIf != nil {
			children := mappingFuncIf()
			c.AddChildren(children...)
		}
	} else if mappingFuncElse != nil {
		children := mappingFuncElse()
		c.AddChildren(children...)
	}
	return c
}
