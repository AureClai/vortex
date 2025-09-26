//go:build js && wasm

// Package style provides type-safe CSS styling for Vortex components.
// This file contains the functions to apply the theme to a style.
// It is used to apply the theme to a style.
//
// The style.Theme function act as a wrapper for the whole App VDOM tree
// It is used to apply the theme to the whole App VDOM tree
// The theme is applied to the whole App VDOM tree through the style.
// For example, to apply a theme to the whole App VDOM tree, you can use the following code:
//
//   theme := style.New(
//       style.Theme(style.ThemeLight),
//   )
//
//
// Basic Usage:
//
//   style := style.New(
//       style.Theme(style.ThemeLight),
//   )
//
// For more information, see the style package documentation

package style
