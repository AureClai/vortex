//go:build js && wasm

// Package animation provides type-safe CSS animations for Vortex components.
// This file contains the functions to apply the common animations to a component.
//
// Basic Usage:
//
//   animation := animation.RotateContinuous() -
//
// For more information, see the animation package documentation

package animation

import (
	"time"
)

// RotateContinuous rotates a component continuously
// duration is the duration of the animation in seconds
func (g *GraphBuilder) RotateContinuous(duration time.Duration) {
	g.AddNode("rotate").WithClip(duration).
		Rotation(
			RotationKeyframe{Time: 0 * time.Millisecond, Value: 0},
			RotationKeyframe{Time: duration, Value: 360},
		).Done().StartWith("rotate")
}

// ScalePulse scales a component pulse
// amplitude is the amplitude of the animation
// duration is the duration of the animation in seconds
func (g *GraphBuilder) ScalePulse(amplitude float32, duration time.Duration) {
	g.AddNode("scale").WithClip(duration).
		Scale(
			ScaleKeyframe{Time: 0 * time.Millisecond, Value: 1.0},
			ScaleKeyframe{Time: duration, Value: amplitude},
		).Done().StartWith("scale")
}
