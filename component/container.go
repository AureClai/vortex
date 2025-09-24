//go:build js && wasm

package component

import (
	"github.com/AureClai/vortex/vdom"
)

type Container struct {
	vdom.ComponentBase // Integration of the ComponentBase
	Children           []vdom.Component
}

func NewContainer() *Container {
	// 1. Create the base with the tag "div"
	base := vdom.NewComponentBase("div")

	// 2. Return the container
	return &Container{
		ComponentBase: base,
		Children:      make([]vdom.Component, 0),
	}
}

func (c *Container) AddChild(child vdom.Component) *Container {
	c.Children = append(c.Children, child)
	c.ComponentBase.AddChildren(child.Render())
	return c
}
