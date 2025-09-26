//go:build js && wasm

package vdom

import (
	"syscall/js"

	"github.com/AureClai/vortex/pkg/style"
)

type VNodeType int

const (
	VNodeElement VNodeType = iota
	VNodeText
)

// VNode represents a virtual node in the DOM
type VNode struct {
	Type          VNodeType
	Tag           string                          // HTML tag name
	Text          string                          // Text content
	Props         map[string]interface{}          // Attributes and properties
	Children      []*VNode                        // Child nodes
	EventHandlers map[string]func(event js.Value) // Event handlers
	Key           string                          // Key for list items
	Element       js.Value                        // Stck la référence à l'élément DOM
	AppliedStyle  *style.Style                    // the style to apply to this node
}

type Component interface {
	Render() *VNode
}
