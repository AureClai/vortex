//go:build js && wasm

package component

import (
	"syscall/js"

	"github.com/AureClai/vortex/core/style"
)

// Functional Component - for simple stateless components
type FunctionalComponent struct {
	*ComponentBase
	renderFunc  func() *VNode
	cachedVNode *VNode
	dirty       bool
}

func NewFunctionalComponent(renderFunc func() *VNode) *FunctionalComponent {
	return &FunctionalComponent{
		ComponentBase: NewComponent("div"), // default tag is irrelevant; Render() supplies the actual VNode
		renderFunc:    renderFunc,
		dirty:         true, // render on first call
	}
}

func (f *FunctionalComponent) Render() *VNode {
	if f.dirty {
		f.cachedVNode = f.renderFunc()
		f.dirty = false
	}
	return f.cachedVNode
}

func (f *FunctionalComponent) OnMount()   {}
func (f *FunctionalComponent) OnUnmount() {}

// OnUpdate marks the component dirty and requests an invalidation so only this subtree re-renders.
func (f *FunctionalComponent) OnUpdate() {
	f.dirty = true
	if invalidate != nil {
		invalidate(f)
	}
}

func (f *FunctionalComponent) GetID() string   { return f.ComponentBase.GetID() }
func (f *FunctionalComponent) SetID(id string) { f.ComponentBase.SetID(id) }

// Fluent API methods - compatible with ComponentBase

func (f *FunctionalComponent) SetClass(class string) *FunctionalComponent {
	vnode := f.Render()
	if vnode.Attrs == nil {
		vnode.Attrs = make(map[string]interface{})
	}
	vnode.Attrs["class"] = class
	return f
}

func (f *FunctionalComponent) Style(s *style.Style) *FunctionalComponent {
	vnode := f.Render()
	vnode.AppliedStyle = s
	return f
}

func (f *FunctionalComponent) On(event string, handler func(event Event)) *FunctionalComponent {
	vnode := f.Render()
	if vnode.EventHandlers == nil {
		vnode.EventHandlers = make(map[string]func(event js.Value))
	}
	vnode.EventHandlers[event] = func(jsEvent js.Value) {
		// Convert js.Value to appropriate Event type based on event name
		var vdomEvent Event
		switch event {
		case "input", "change":
			vdomEvent = NewInputEvent(jsEvent, vnode)
		case "click", "mousedown", "mouseup", "mouseover", "mouseout":
			vdomEvent = NewMouseEvent(jsEvent, vnode)
		case "keydown", "keyup", "keypress":
			vdomEvent = NewKeyboardEvent(jsEvent, vnode)
		case "focus", "blur":
			vdomEvent = NewFocusEvent(jsEvent, vnode)
		default:
			vdomEvent = NewEvent(jsEvent, vnode)
		}
		handler(vdomEvent)
	}
	return f
}

func (f *FunctionalComponent) SetText(text string) *FunctionalComponent {
	vnode := f.Render()
	vnode.Text = text
	return f
}

func (f *FunctionalComponent) AddChild(child Component) *FunctionalComponent {
	if child != nil {
		vnode := f.Render()
		if c := child.Render(); c != nil {
			vnode.Children = append(vnode.Children, c)
		}
	}
	return f
}

func (f *FunctionalComponent) AddChildren(children ...Component) *FunctionalComponent {
	vnode := f.Render()
	for _, child := range children {
		if child == nil {
			continue
		}
		if c := child.Render(); c != nil {
			vnode.Children = append(vnode.Children, c)
		}
	}
	return f
}
