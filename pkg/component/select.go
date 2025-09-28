//go:build js && wasm

package component

import (
	"strconv"

	"github.com/AureClai/vortex/pkg/vdom"
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
	*vdom.StatefulComponentBase[SelectState]
	options []OptionComponent
}

func Select() *SelectComponent {
	return &SelectComponent{
		StatefulComponentBase: vdom.NewStatefulComponent[SelectState]("select", SelectState{}),
		options:               make([]OptionComponent, 0),
	}
}

func (s *SelectComponent) Render() *vdom.VNode {
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
	children := make([]*vdom.VNode, len(s.options))
	for i, option := range s.options {
		children[i] = option.Render()
	}

	return &vdom.VNode{
		Type:     vdom.VNodeElement,
		Tag:      "select",
		Props:    props,
		Children: children,
	}
}
