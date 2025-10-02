//go:build js && wasm

package components

import (
	"strconv"

	"github.com/AureClai/vortex/core/component"
)

// =============================================================================
// SELECT COMPONENT
// =============================================================================

type SelectState struct {
	Value    string
	Multiple bool
	Disabled bool
	Required bool
	Size     int
}

type SelectComponent struct {
	*component.StatefulComponentBase[SelectState]
	options []OptionComponent
}

func Select() *SelectComponent {
	return &SelectComponent{
		StatefulComponentBase: component.NewStatefulComponent[SelectState]("select", SelectState{}),
		options:               make([]OptionComponent, 0),
	}
}

func (s *SelectComponent) Render() *component.VNode {
	state := s.GetState()

	props := map[string]interface{}{}
	if state.Multiple {
		props["multiple"] = true
	}
	if state.Disabled {
		props["disabled"] = true
	}
	if state.Required {
		props["required"] = true
	}
	if state.Size > 0 {
		props["size"] = strconv.Itoa(state.Size)
	}

	// Convert options to VNodes
	children := make([]*component.VNode, len(s.options))
	for i, option := range s.options {
		children[i] = option.Render()
	}

	return &component.VNode{
		Type:     component.VNodeElement,
		Tag:      "select",
		Attrs:    props,
		Children: children,
	}
}
