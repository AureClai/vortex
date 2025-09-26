//go:build js && wasm

package main

import (
	"github.com/AureClai/vortex/dev_tool/components"
	"github.com/AureClai/vortex/pkg/renderer"
)

func main() {
	// Create dev tool renderer in its own container
	devRenderer := renderer.NewRenderer("vtx-dev-tool")

	// Create the main dev panel
	devPanel := components.NewDevPanel()
	devPanel.Style(styles.devPanelStyle)

}
