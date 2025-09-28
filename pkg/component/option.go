//go:build js && wasm

package component

import (
	"github.com/AureClai/vortex/pkg/vdom"
)

type OptionComponent struct {
	*vdom.FunctionalComponent
	value    string
	text     string
	selected bool
	disabled bool
}

func Option(value, text string) *OptionComponent {
	comp := &OptionComponent{
		value: value,
		text:  text,
	}

	comp.FunctionalComponent = vdom.NewFunctionalComponent(func() *vdom.VNode {
		props := map[string]interface{}{
			"value": comp.value,
		}
		if comp.selected {
			props["selected"] = true
		}
		if comp.disabled {
			props["disabled"] = true
		}

		return &vdom.VNode{
			Type: vdom.VNodeElement,
			Tag:  "option",
			Children: []*vdom.VNode{
				{
					Type: vdom.VNodeText,
					Text: comp.text,
				},
			},
			Props: props,
		}
	})

	return comp
}

func (o *OptionComponent) Selected(selected bool) *OptionComponent {
	o.selected = selected
	o.OnUpdate()
	return o
}
