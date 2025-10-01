//go:build js && wasm

package component

// StatefulComponentBase implementation
type StatefulComponentBase[T any] struct {
	*ComponentBase // Integration of the ComponentBase
	state          T
}

// Global invalidation hook, set y the app/renderer
var invalidate func(Component)

func SetInvalidator(fn func(Component)) { invalidate = fn }

func NewStatefulComponent[T any](tag string, initialState T) *StatefulComponentBase[T] {
	base := NewComponent(tag)

	return &StatefulComponentBase[T]{
		ComponentBase: base,
		state:         initialState,
	}
}

// State management methods
func (c *StatefulComponentBase[T]) GetState() T {
	return c.state
}

func (c *StatefulComponentBase[T]) SetState(state T) {
	c.state = state
	if invalidate != nil {
		invalidate(c)
	}
}

func (c *StatefulComponentBase[T]) UpdateState(mutator func(T) T) {
	newState := mutator(c.state)
	c.SetState(newState)
}

func (c *StatefulComponentBase[T]) OnUpdate() {
	c.ComponentBase.Render()
}
