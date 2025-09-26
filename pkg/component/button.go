//go:build js && wasm

package component

import (
	"github.com/AureClai/vortex/pkg/vdom"
)

type Button struct {
	vdom.ComponentBase // Integration of the ComponentBase
}

func NewButton(text string) *Button {
	// 1. Create the base with the tag "button"
	base := vdom.NewComponentBase("button")

	// 2. Set the text
	textNode := &vdom.VNode{Type: vdom.VNodeText, Text: text}
	base.AddChildren(textNode)

	// 3. Return the button
	return &Button{
		ComponentBase: base,
	}
}
