//go:build js && wasm

package app

import (
	"syscall/js"

	"github.com/AureClai/vortex/pkg/renderer"
	"github.com/AureClai/vortex/pkg/vdom"
)

// App is the root application component that manages global state and rendering
type App[T any] struct {
	*vdom.StatefulComponentBase[T]
	renderer    *renderer.Renderer
	mounted     bool
	renderFunc  func(T) *vdom.VNode
	onMountedFn func()
}

// NewApp creates a new Vortex application
func New[T any](containerId string, initialState T) *App[T] {
	r := renderer.NewRenderer(containerId)

	app := &App[T]{
		StatefulComponentBase: vdom.NewStatefulComponent[T]("div", initialState),
		renderer:              r,
	}

	// Auto re-render on state changes
	app.SetOnStateChange(func(oldState, newState T) {
		if app.mounted {
			app.renderer.Render(app.Render())
		}
	})

	return app
}

func (a *App[T]) SetRender(render func(T) *vdom.VNode) {
	a.renderFunc = render
}

// SetOnMounted sets a callback to be called after the app is mounted
func (a *App[T]) SetOnMounted(callback func()) {
	a.onMountedFn = callback
}

// Override to use custom render function
func (a *App[T]) Render() *vdom.VNode {
	if a.renderFunc != nil {
		return a.renderFunc(a.GetState())
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
