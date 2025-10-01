//go:build js && wasm

// Container component is a div stateless component

// It is a wrapper for the FunctionalComponent
// It is used to create a container for the children components
package components

import (
	"github.com/AureClai/vortex/core/component"
)

type Container struct {
	*component.FunctionalComponent
}

func Div(children ...component.Component) *Container {
	childVNodes := make([]*component.VNode, 0, len(children))
	for _, child := range children {
		if child != nil {
			if vnode := child.Render(); vnode != nil {
				childVNodes = append(childVNodes, vnode)
			}
		}
	}

	return &Container{
		FunctionalComponent: component.NewFunctionalComponent(func() *component.VNode {
			return &component.VNode{
				Type:     component.VNodeElement,
				Tag:      "div",
				Children: childVNodes,
			}
		}),
	}
}

// Convenience constructors for different container types
func Section(children ...component.Component) *Container {
	return containerWithTag("section", children...)
}
func Article(children ...component.Component) *Container {
	return containerWithTag("article", children...)
}
func Header(children ...component.Component) *Container {
	return containerWithTag("header", children...)
}
func Footer(children ...component.Component) *Container {
	return containerWithTag("footer", children...)
}
func Main(children ...component.Component) *Container { return containerWithTag("main", children...) }
func Nav(children ...component.Component) *Container  { return containerWithTag("nav", children...) }

func containerWithTag(tag string, children ...component.Component) *Container {
	childVNodes := make([]*component.VNode, 0, len(children))
	for _, child := range children {
		if child != nil {
			if vnode := child.Render(); vnode != nil {
				childVNodes = append(childVNodes, vnode)
			}
		}
	}

	return &Container{
		FunctionalComponent: component.NewFunctionalComponent(func() *component.VNode {
			return &component.VNode{
				Type:     component.VNodeElement,
				Tag:      tag,
				Children: childVNodes,
			}
		}),
	}
}
