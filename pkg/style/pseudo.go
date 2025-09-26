//go:build js && wasm

// Package style provides type-safe CSS styling for Vortex components.
// This file contains the functions to apply the pseudo-classes to a style.
// It is used to apply the pseudo-classes to a style.
//
// Basic Usage:
//
//   style := style.New(
//       style.OnHover(style.BackgroundColor("#f0f0f0")),
//   )
//
// For more information, see the style package documentation
//
// All properties available are (table with css equivalent)
// | Property | CSS Equivalent |
// | -------- | -------------- |
// | OnHover | :hover        |
// | OnActive | :active        |
// | OnFocus | :focus        |
// | OnFocusWithin | :focus-within        |
// | OnFocusVisible | :focus-visible        |
//
// For more information, see the style package documentation

package style

// --- Pseudo-classes

// OnHiver applies the given styles for the pseudo :hover
func OnHover(properties ...StyleOption) StyleOption {
	return func(s *Style) {
		// Create a temporary style object to collect the hover properties
		hoverStyle := New(properties...)
		s.Pseudos[":hover"] = hoverStyle.Base
	}
}

// OnActive applies the given styles for the pseudo :active
func OnActive(properties ...StyleOption) StyleOption {
	return func(s *Style) {
		activeStyle := New(properties...)
		s.Pseudos[":active"] = activeStyle.Base
	}
}

// OnFocus applies the given styles for the pseudo :focus
func OnFocus(properties ...StyleOption) StyleOption {
	return func(s *Style) {
		focusStyle := New(properties...)
		s.Pseudos[":focus"] = focusStyle.Base
	}
}

// OnFocusWithin applies the given styles for the pseudo :focus-within
func OnFocusWithin(properties ...StyleOption) StyleOption {
	return func(s *Style) {
		focusWithinStyle := New(properties...)
		s.Pseudos[":focus-within"] = focusWithinStyle.Base
	}
}

// OnFocusVisible applies the given styles for the pseudo :focus-visible
func OnFocusVisible(properties ...StyleOption) StyleOption {
	return func(s *Style) {
		focusVisibleStyle := New(properties...)
		s.Pseudos[":focus-visible"] = focusVisibleStyle.Base
	}
}
