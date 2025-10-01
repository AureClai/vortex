//go:build js && wasm

package component

import (
	"syscall/js"

	"github.com/AureClai/vortex/core/style"
)

type VNodeType int

const (
	VNodeElement VNodeType = iota
	VNodeText
	VNodeComponent
)

// Core VNode
type VNode struct {
	Type VNodeType
	Tag  string // HTML tag name
	Text string // Text content

	// DOM-related properties
	Attrs         map[string]interface{}          // was Props
	EventHandlers map[string]func(event js.Value) // Event handlers
	Element       js.Value                        // Stck la référence à l'élément DOM

	// VDOM-specific properties
	Key          string       // Key for list items
	Children     []*VNode     // Child nodes
	AppliedStyle *style.Style // the style to apply to this node

	// Optional logic-only metadata (ignored by renderer diff)
	Data map[string]interface{}

	// Component reference (for VNodeComponent type)
	Component Component
}
