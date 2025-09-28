//go:build js && wasm

// Container component is a div stateless component

// It is a wrapper for the FunctionalComponent
// It is used to create a container for the children components
package component

import (
	"github.com/AureClai/vortex/pkg/vdom"
)

type Container struct {
	*vdom.FunctionalComponent
}

func Div(children ...vdom.Component) *Container {
	childVNodes := make([]*vdom.VNode, 0, len(children))
	for _, child := range children {
		if child != nil {
			if vnode := child.Render(); vnode != nil {
				childVNodes = append(childVNodes, vnode)
			}
		}
	}

	return &Container{
		FunctionalComponent: vdom.NewFunctionalComponent(func() *vdom.VNode {
			return &vdom.VNode{
				Type:     vdom.VNodeElement,
				Tag:      "div",
				Children: childVNodes,
			}
		}),
	}
}

// Convenience constructors for different container types
func Section(children ...vdom.Component) *Container { return containerWithTag("section", children...) }
func Article(children ...vdom.Component) *Container { return containerWithTag("article", children...) }
func Header(children ...vdom.Component) *Container  { return containerWithTag("header", children...) }
func Footer(children ...vdom.Component) *Container  { return containerWithTag("footer", children...) }
func Main(children ...vdom.Component) *Container    { return containerWithTag("main", children...) }
func Nav(children ...vdom.Component) *Container     { return containerWithTag("nav", children...) }

func containerWithTag(tag string, children ...vdom.Component) *Container {
	childVNodes := make([]*vdom.VNode, 0, len(children))
	for _, child := range children {
		if child != nil {
			if vnode := child.Render(); vnode != nil {
				childVNodes = append(childVNodes, vnode)
			}
		}
	}

	return &Container{
		FunctionalComponent: vdom.NewFunctionalComponent(func() *vdom.VNode {
			return &vdom.VNode{
				Type:     vdom.VNodeElement,
				Tag:      tag,
				Children: childVNodes,
			}
		}),
	}
}
