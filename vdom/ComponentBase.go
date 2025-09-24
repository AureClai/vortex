//go:build js && wasm

package vdom

import (
	"syscall/js"

	"github.com/AureClai/vortex/style"
)

// ComponentBase contains the data and methods common to all components
type ComponentBase struct {
	vNode *VNode
}

// NewComponentBase creates a new component base
func NewComponentBase(tag string) ComponentBase {
	return ComponentBase{
		vNode: &VNode{
			Type:          VNodeElement,
			Tag:           tag,
			Props:         make(map[string]interface{}),
			EventHandlers: make(map[string]func(event js.Value)),
			Children:      make([]*VNode, 0),
		},
	}
}

// Render return simply the internal vNode
func (c *ComponentBase) Render() *VNode {
	return c.vNode
}

// --- Common methodes "Fluent" ---

func (c *ComponentBase) SetClass(class string) *ComponentBase {
	c.vNode.Props["class"] = class
	return c
}

func (c *ComponentBase) SetID(id string) *ComponentBase {
	c.vNode.Props["id"] = id
	return c
}

func (c *ComponentBase) SetKey(key string) *ComponentBase {
	c.vNode.Key = key
	return c
}

func (c *ComponentBase) Style(s *style.Style) *ComponentBase {
	c.vNode.AppliedStyle = s
	return c
}

func (c *ComponentBase) On(event string, handler func(event js.Value)) *ComponentBase {
	c.vNode.EventHandlers[event] = handler
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

func (c *ComponentBase) AddChildren(children ...*VNode) *ComponentBase {
	c.vNode.Children = append(c.vNode.Children, children...)
	return c
}
