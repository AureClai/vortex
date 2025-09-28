//go:build js && wasm

package vdom

// StatefulComponent interface extends Component with state management
type StatefulComponent[T any] interface {
	Component

	// State management
	GetState() T
	SetState(state T)
	UpdateState(updater func(state T) T)

	// State Lifecycle
	OnStateChange(callback func(state T))
}

// StatefulComponentBase implementation
type StatefulComponentBase[T any] struct {
	*ComponentBase
	state T

	// Callbacks
	onStateChange func(oldState T, newState T)
	onRender      func() *VNode
}

func NewStatefulComponent[T any](tag string, initialState T) *StatefulComponentBase[T] {
	base := NewComponent(tag)

	return &StatefulComponentBase[T]{
		ComponentBase: base,
		state:         initialState,
		onStateChange: func(oldState T, newState T) {}, // Default no-op
	}
}

// State management methods
func (c *StatefulComponentBase[T]) GetState() T {
	return c.state
}

func (c *StatefulComponentBase[T]) SetState(state T) {
	oldState := c.state
	c.state = state
	c.onStateChange(oldState, state)
	c.OnUpdate() // Trigger re-render
}

func (c *StatefulComponentBase[T]) UpdateState(mutator func(T) T) {
	newState := mutator(c.state)
	c.SetState(newState)
}

func (c *StatefulComponentBase[T]) OnStateChange(oldState T, newState T) {
	if c.onStateChange != nil {
		c.onStateChange(oldState, newState)
	}
}

// Allow custom change callbacks
func (c *StatefulComponentBase[T]) SetOnStateChange(callback func(T, T)) {
	c.onStateChange = callback
}
