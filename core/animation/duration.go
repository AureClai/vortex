//go:build js && wasm

// Package animation provides type-safe CSS animations for Vortex components.
// This file contains the functions to API for the duration of the animations with common durations
// Shortcuts for the time module of Go
//
// Basic Usage:
//
//   animation.Second(1)
//   animation.Millisecond(100)
//   animation.Microsecond(1000)
//   animation.Nanosecond(1000000)
//
// For more information, see the animation package documentation

package animation

import "time"

func Second(value float32) time.Duration {
	return time.Second * time.Duration(value)
}

func Millisecond(value float32) time.Duration {
	return time.Millisecond * time.Duration(value)
}

func Microsecond(value float32) time.Duration {
	return time.Microsecond * time.Duration(value)
}

func Nanosecond(value float32) time.Duration {
	return time.Nanosecond * time.Duration(value)
}
