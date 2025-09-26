//go:build js && wasm

package component

import (
	"github.com/AureClai/vortex/pkg/vdom"
)

type Text struct {
	vdom.ComponentBase
}

func NewText(content string) *Text {
	base := vdom.NewComponentBase("span")
	base.AddChildren(&vdom.VNode{Type: vdom.VNodeText, Text: content})
	return &Text{
		ComponentBase: base,
	}
}

func NewHeading(content string, level int) *Text {
	tag := "h1"
	switch level {
	case 1:
		tag = "h1"
	case 2:
		tag = "h2"
	case 3:
		tag = "h3"
	case 4:
		tag = "h4"
	case 5:
		tag = "h5"
	case 6:
		tag = "h6"
	}
	base := vdom.NewComponentBase(tag)
	base.AddChildren(&vdom.VNode{Type: vdom.VNodeText, Text: content})

	return &Text{
		ComponentBase: base,
	}
}

func NewParagraph(content string) *Text {
	base := vdom.NewComponentBase("p")
	base.AddChildren(&vdom.VNode{Type: vdom.VNodeText, Text: content})
	return &Text{
		ComponentBase: base,
	}
}
