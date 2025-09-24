//go:build js && wasm

package vdom

// StatefulComponent is a base struct for the components that need to manage their own state
// It uses the generic type T to store the state and it is compatible with any State Struct
type StatefulComponentBase[T any] struct {
	ComponentBase
	state    T
	reRender func() // the function to re-render the component
}

// NewStatefulComponent creates a new stateful component
func NewStatefulComponent[T any](tag string, initialState T, reRender func()) StatefulComponentBase[T] {
	return StatefulComponentBase[T]{
		ComponentBase: NewComponentBase(tag),
		state:         initialState,
		reRender:      reRender,
	}
}

// State return the state of the component
func (c *StatefulComponentBase[T]) State() T {
	return c.state
}

// SetState sets the state of the component
func (c *StatefulComponentBase[T]) SetState(state T) {
	c.state = state
	c.reRender()
}

// UpdateState is an easy way to update the state of the component
// It takes the "mutator" that receive the current state and return the new state
func (c *StatefulComponentBase[T]) UpdateState(mutator func(currentState T) T) {
	currentState := c.State()
	newState := mutator(currentState)
	c.SetState(newState)
}
