//go:build js && wasm

package app

import (
	"syscall/js"

	"github.com/AureClai/vortex/core/component"
	"github.com/AureClai/vortex/core/renderer"
)

// App is the root application component that manages global state and rendering
type App[T any] struct {
	*component.StatefulComponentBase[T]
	renderer    *renderer.Renderer
	mounted     bool
	renderFunc  func(T) *component.VNode
	onMountedFn func()
}

// NewApp creates a new Vortex application
func New[T any](containerId string, initialState T) *App[T] {
	r := renderer.NewRenderer(containerId)

	app := &App[T]{
		StatefulComponentBase: component.NewStatefulComponent[T]("div", initialState),
		renderer:              r,
	}

	component.SetInvalidator(func(c component.Component) {
		app.renderer.Render(app.Render())
	})

	return app
}

func (a *App[T]) SetRender(render func(T) *component.VNode) {
	a.renderFunc = render
}

// SetOnMounted sets a callback to be called after the app is mounted
func (a *App[T]) SetOnMounted(callback func()) {
	a.onMountedFn = callback
}

// Override to use custom render function
func (a *App[T]) Render() *component.VNode {
	if a.renderFunc != nil {
		return a.renderFunc(a.StatefulComponentBase.GetState())
	}
	return a.StatefulComponentBase.Render()
}

// Mount renders the application to the DOM
func (a *App[T]) Mount() {
	a.mounted = true
	a.renderer.Render(a.Render())

	// Call post-mount hooks after DOM is ready
	a.onMounted()
}

// onMounted is called after the app is mounted and DOM elements are available
func (a *App[T]) onMounted() {
	if a.onMountedFn != nil {
		a.onMountedFn()
	}
}

// AddGlobalListener adds a global DOM event listener
func (a *App[T]) AddGlobalListener(event string, handler func()) {
	window := js.Global().Get("window")
	window.Call("addEventListener", event, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		handler()
		return nil
	}))
}
