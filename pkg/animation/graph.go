//go::build js && wasm

//

package animation

import "time"

// AnimationGraph manages the animation states and transitions
type Graph struct {
	// Graph structure
	Nodes map[string]*Node
	Edges map[string][]*Edge

	// runtime state
	CurrentNode *Node
	IsPlaying   bool
	GlobalTime  time.Duration

	// Callbacks
	OnNodeChanged func(from, to string)
	OnEdge        func(edge *Edge)

	// Parameters
	Parameters map[string]interface{}
}

// GetParameter returns the value of a parameter
func (sm *Graph) GetParameter(name string) (interface{}, bool) {
	param, exists := sm.Parameters[name]
	return param, exists
}

// AnimationNode represents a node in the animation graph
type Node struct {
	Name  string
	Clip  *Clip
	Loop  bool
	Speed float32

	// Node callbacks
	OnEnter  func()
	OnUpdate func(deltaTime time.Duration)
	OnExit   func()

	// Runtime data
	LocalTime time.Duration
	IsActive  bool
}

// AnimationEdge represents a transition between two nodes
type Edge struct {
	From, To  string
	Condition EdgeCondition
	Duration  time.Duration
	BlendType BlendType

	// Transition progress
	Progress float32
	IsActive bool
}

// EdgeCondition represents the condition that must be met for an edge to be taken
type EdgeCondition interface {
	Evaluate(sm *Graph) bool
}

// BlendType defines how two animations blend during transition
type BlendType int

const (
	BlendLinear BlendType = iota
	BlendEaseIn
	BlendEaseOut
	BlendEaseInOut
	BlendCustom
)
