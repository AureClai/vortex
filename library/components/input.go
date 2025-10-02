//go:build js && wasm

// Input component - Stateful
// Usage :
// input := component.Input(component.InputTypeText).Placeholder("Enter your name").OnChange(func(event *component.InputEvent) {
// 	fmt.Println(event.Value())
// }).OnInput(func(event *component.InputEvent) {
// 	fmt.Println(event.Value())
// })

package components

import (
	"strconv"

	"github.com/AureClai/vortex/core/component"
)

type InputState struct {
	Value       string
	Placeholder string
	Disabled    bool
	ReadOnly    bool
	Required    bool
	Pattern     string
	MaxLength   int
	MinLength   int
	Name        string
	ID          string
}

type InputType string

const (
	InputTypeText          InputType = "text"
	InputTypePassword      InputType = "password"
	InputTypeEmail         InputType = "email"
	InputTypeNumber        InputType = "number"
	InputTypeTel           InputType = "tel"
	InputTypeUrl           InputType = "url"
	InputTypeSearch        InputType = "search"
	InputTypeDate          InputType = "date"
	InputTypeTime          InputType = "time"
	InputTypeDatetimeLocal InputType = "datetime-local"
	InputTypeMonth         InputType = "month"
	InputTypeWeek          InputType = "week"
	InputTypeColor         InputType = "color"
	InputTypeRange         InputType = "range"
	InputTypeFile          InputType = "file"
	InputTypeHidden        InputType = "hidden"
	InputTypeCheckbox      InputType = "checkbox"
	InputTypeRadio         InputType = "radio"
	InputTypeSubmit        InputType = "submit"
	InputTypeReset         InputType = "reset"
	InputTypeButton        InputType = "button"
)

type InputComponent struct {
	*component.StatefulComponentBase[InputState]
	inputType InputType
}

// Event handlers types
type InputChangeHandler func(*component.InputEvent)
type InputInputHandler func(*component.InputEvent)

// Constructor
func Input(inputType InputType) *InputComponent {
	initialState := InputState{
		Value:       "",
		Placeholder: "",
		Disabled:    false,
		ReadOnly:    false,
		Required:    false,
		Pattern:     "",
		MaxLength:   -1,
		MinLength:   -1,
		Name:        "",
		ID:          "",
	}

	component := &InputComponent{
		StatefulComponentBase: component.NewStatefulComponent[InputState]("input", initialState),
		inputType:             inputType,
	}

	return component
}

// Render the input component
func (i *InputComponent) Render() *component.VNode {
	state := i.GetState()

	props := map[string]interface{}{
		"type": string(i.inputType),
	}

	// Map state to HTML attributes
	if state.Value != "" {
		props["value"] = state.Value
	}
	if state.Placeholder != "" {
		props["placeholder"] = state.Placeholder
	}
	if state.Name != "" {
		props["name"] = state.Name
	}
	if state.ID != "" {
		props["id"] = state.ID
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
	if state.Pattern != "" {
		props["pattern"] = state.Pattern
	}
	if state.MaxLength > 0 {
		props["maxlength"] = strconv.Itoa(state.MaxLength)
	}
	if state.MinLength > 0 {
		props["minlength"] = strconv.Itoa(state.MinLength)
	}

	return &component.VNode{
		Type:  component.VNodeElement,
		Tag:   "input",
		Attrs: props,
	}
}

// Fluent API methods
func (i *InputComponent) Value(value string) *InputComponent {
	i.UpdateState(func(state InputState) InputState {
		state.Value = value
		return state
	})
	return i
}

func (i *InputComponent) Placeholder(placeholder string) *InputComponent {
	i.UpdateState(func(state InputState) InputState {
		state.Placeholder = placeholder
		return state
	})
	return i
}

func (i *InputComponent) Disabled(disabled bool) *InputComponent {
	i.UpdateState(func(state InputState) InputState {
		state.Disabled = disabled
		return state
	})
	return i
}

func (i *InputComponent) Required(required bool) *InputComponent {
	i.UpdateState(func(state InputState) InputState {
		state.Required = required
		return state
	})
	return i
}

func (i *InputComponent) Name(name string) *InputComponent {
	i.UpdateState(func(state InputState) InputState {
		state.Name = name
		return state
	})
	return i
}

func (i *InputComponent) ID(id string) *InputComponent {
	i.UpdateState(func(state InputState) InputState {
		state.ID = id
		return state
	})
	return i
}

// Event handlers
func (i *InputComponent) OnChange(handler InputChangeHandler) *InputComponent {
	i.On("change", func(event component.Event) {
		if inputEvent, ok := event.(*component.InputEvent); ok {
			// Update internal state
			i.UpdateState(func(state InputState) InputState {
				state.Value = inputEvent.Value()
				return state
			})
			// Call user handler
			handler(inputEvent)
		}
	})
	return i
}

func (i *InputComponent) OnInput(handler InputInputHandler) *InputComponent {
	i.On("input", func(event component.Event) {
		if inputEvent, ok := event.(*component.InputEvent); ok {
			// Update internal state in real-time
			i.UpdateState(func(state InputState) InputState {
				state.Value = inputEvent.Value()
				return state
			})
			// Call user handler
			handler(inputEvent)
		}
	})
	return i
}
