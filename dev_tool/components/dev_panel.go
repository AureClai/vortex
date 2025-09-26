//go:build js && wasm

package components

import "github.com/AureClai/vortex/pkg/vdom"

type DevPanelState struct {
	IsOpen bool
}

type DevPanel struct {
	vdom.StatefulComponentBase[DevPanelState]
}

func NewDevPanel() DevPanel {
	return DevPanel{
		StatefulComponentBase: vdom.NewStatefulComponent("dev-panel", DevPanelState{IsOpen: false}, func() {}),
	}
}

func (d *DevPanel) Render() *vdom.VNode {
	return vdom.NewVNode("div", func(v *vdom.VNode) {
		v.SetText("Dev Panel")
	})
}
