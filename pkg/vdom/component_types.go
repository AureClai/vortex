//go:build js && wasm

// Component types

package vdom

import (
	"syscall/js"

	"github.com/AureClai/vortex/pkg/style"
)

// Functional Component - for simple stateless components
type FunctionalComponent struct {
	renderFunc func() *VNode
	id         string

	// Chache the rendered VNode for fluent modifications
	cachedVNode *VNode
	dirty       bool // Whether the cached VNode is dirty
}

func NewFunctionalComponent(renderFunc func() *VNode) *FunctionalComponent {
	return &FunctionalComponent{
		renderFunc: renderFunc,
		dirty:      true, // needs render on first call
	}
}

func (f *FunctionalComponent) Render() *VNode {
	if f.dirty {
		f.cachedVNode = f.renderFunc()
		f.dirty = false
	}
	return f.cachedVNode
}

func (f *FunctionalComponent) OnMount()        {}
func (f *FunctionalComponent) OnUpdate()       { f.dirty = true }
func (f *FunctionalComponent) OnUnmount()      {}
func (f *FunctionalComponent) GetID() string   { return f.id }
func (f *FunctionalComponent) SetID(id string) { f.id = id }

// Fluent API méthods - make it compatible with the ComponentBase
func (f *FunctionalComponent) SetClass(class string) *FunctionalComponent {
	vnode := f.Render()
	if vnode.Props == nil {
		vnode.Props = make(map[string]interface{})
	}
	vnode.Props["class"] = class
	return f
}

func (f *FunctionalComponent) Style(s *style.Style) *FunctionalComponent {
	vnode := f.Render()
	vnode.AppliedStyle = s
	return f
}

func (f *FunctionalComponent) On(event string, handler func(event js.Value)) *FunctionalComponent {
	vnode := f.Render()
	if vnode.EventHandlers == nil {
		vnode.EventHandlers = make(map[string]func(event js.Value))
	}
	vnode.EventHandlers[event] = handler
	return f
}

func (f *FunctionalComponent) SetText(text string) *FunctionalComponent {
	vnode := f.Render()
	vnode.Text = text
	return f
}

func (f *FunctionalComponent) AddChild(child Component) *FunctionalComponent {
	vnode := f.Render()
	vnode.Children = append(vnode.Children, child.Render())
	return f
}

func (f *FunctionalComponent) AddChildren(children ...Component) *FunctionalComponent {
	vnode := f.Render()
	for _, child := range children {
		vnode.Children = append(vnode.Children, child.Render())
	}
	return f
}
