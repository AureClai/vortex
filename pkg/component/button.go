//go:build js && wasm

package component

import (
	"github.com/AureClai/vortex/pkg/vdom"
)

// =============================================================================
// BUTTON COMPONENT
// =============================================================================

type ButtonType string

const (
	ButtonTypeButton ButtonType = "button"
	ButtonTypeSubmit ButtonType = "submit"
	ButtonTypeReset  ButtonType = "reset"
)

type ButtonComponent struct {
	*vdom.FunctionalComponent
	text       string
	buttonType ButtonType
}

func Button(text string) *ButtonComponent {
	comp := &ButtonComponent{
		text:       text,
		buttonType: ButtonTypeButton,
	}

	comp.FunctionalComponent = vdom.NewFunctionalComponent(func() *vdom.VNode {
		return &vdom.VNode{
			Type: vdom.VNodeElement,
			Tag:  "button",
			Children: []*vdom.VNode{
				{
					Type: vdom.VNodeText,
					Text: comp.text,
				},
			},
			Props: map[string]interface{}{
				"type": string(comp.buttonType),
			},
		}
	})

	return comp
}

func (b *ButtonComponent) Type(buttonType ButtonType) *ButtonComponent {
	b.buttonType = buttonType
	b.OnUpdate() // Trigger re-render
	return b
}
