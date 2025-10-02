//go:build js && wasm

package components

import (
	"github.com/AureClai/vortex/core/component"
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
	*component.FunctionalComponent
	text       string
	buttonType ButtonType
}

func Button(text string) *ButtonComponent {
	comp := &ButtonComponent{
		text:       text,
		buttonType: ButtonTypeButton,
	}

	comp.FunctionalComponent = component.NewFunctionalComponent(func() *component.VNode {
		return &component.VNode{
			Type: component.VNodeElement,
			Tag:  "button",
			Children: []*component.VNode{
				{
					Type: component.VNodeText,
					Text: comp.text,
				},
			},
			Attrs: map[string]interface{}{
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

func (b *ButtonComponent) OnClick(handler func(event *component.Event)) *ButtonComponent {
	b.On("click", func(event component.Event) {
		handler(&event)
	})
	return b
}
