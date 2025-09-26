//go:build js && wasm

package styles

import "github.com/AureClai/vortex/pkg/style"

// Dev Panel Style
// On top of everything
// Full width and height
// Semi-transparent black background
var devPanelStyle = style.New(
	style.Position(style.PositionAbsolute),
	style.PositionSide(style.PositionTop, style.Percent(0)),
	style.PositionSide(style.PositionRight, style.Percent(0)),
	style.Width(style.Percent(100)),
	style.Height(style.Percent(100)),
	style.ZIndex(1000),                                                // Make sure it is on top of everything
	style.BackgroundColor(style.ColorValue(style.RGBA(0, 0, 0, 0.5))), // Semi-transparent black
)
