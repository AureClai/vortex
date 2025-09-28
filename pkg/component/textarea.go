//go:build js && wasm

package component

import (
	"strconv"

	"github.com/AureClai/vortex/pkg/vdom"
)

type TextareaState struct {
	Value       string
	Placeholder string
	Disabled    bool
	ReadOnly    bool
	Required    bool
	Rows        int
	Cols        int
	MaxLength   int
	MinLength   int
}

type TextareaComponent struct {
	*vdom.StatefulComponentBase[TextareaState]
}

func Textarea() *TextareaComponent {
	initialState := TextareaState{
		Rows: 4,
		Cols: 50,
	}

	return &TextareaComponent{
		StatefulComponentBase: vdom.NewStatefulComponent[TextareaState]("textarea", initialState),
	}
}

func (t *TextareaComponent) Render() *vdom.VNode {
	state := t.GetState()

	props := map[string]interface{}{}

	if state.Placeholder != "" {
		props["placeholder"] = state.Placeholder
	}
	if state.Disabled {
		props["disabled"] = true
	}
	if state.ReadOnly {
		props["readonly"] = true
	}
	if state.Required {
		props["required"] = true
	}
	if state.Rows > 0 {
		props["rows"] = strconv.Itoa(state.Rows)
	}
	if state.Cols > 0 {
		props["cols"] = strconv.Itoa(state.Cols)
	}

	return &vdom.VNode{
		Type: vdom.VNodeElement,
		Tag:  "textarea",
		Children: []*vdom.VNode{
			{
				Type: vdom.VNodeText,
				Text: state.Value, // Textarea content as child text node
			},
		},
		Props: props,
	}
}

// Fluent methods for Textarea
func (t *TextareaComponent) Value(value string) *TextareaComponent {
	t.UpdateState(func(state TextareaState) TextareaState {
		state.Value = value
		return state
	})
	return t
}

func (t *TextareaComponent) Placeholder(placeholder string) *TextareaComponent {
	t.UpdateState(func(state TextareaState) TextareaState {
		state.Placeholder = placeholder
		return state
	})
	return t
}

func (t *TextareaComponent) Rows(rows int) *TextareaComponent {
	t.UpdateState(func(state TextareaState) TextareaState {
		state.Rows = rows
		return state
	})
	return t
}
