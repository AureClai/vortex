//go:build js && wasm

package components

import (
	"github.com/AureClai/vortex/core/component"
)

type OptionComponent struct {
	*component.FunctionalComponent
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

	comp.FunctionalComponent = component.NewFunctionalComponent(func() *component.VNode {
		props := map[string]interface{}{
			"value": comp.value,
		}
		if comp.selected {
			props["selected"] = true
		}
		if comp.disabled {
			props["disabled"] = true
		}

		return &component.VNode{
			Type: component.VNodeElement,
			Tag:  "option",
			Children: []*component.VNode{
				{
					Type: component.VNodeText,
					Text: comp.text,
				},
			},
			Attrs: props,
		}
	})

	return comp
}

func (o *OptionComponent) Selected(selected bool) *OptionComponent {
	o.selected = selected
	o.OnUpdate()
	return o
}
