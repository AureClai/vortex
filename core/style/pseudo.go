//go:build js && wasm

// Package style provides type-safe CSS styling for Vortex components.
// This file contains the functions to apply the pseudo-classes to a style.
// It is used to apply the pseudo-classes to a style.
//
// Basic Usage:
//
//   style := style.New()
//       .OnHover(style.BackgroundColor("#f0f0f0")).
//       .OnActive(style.BackgroundColor("#f0f0f0"))
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
func (s *Style) OnHover(style *Style) *Style {
	// Create a temporary style object to collect the hover properties
	s.Pseudos[":hover"] = style.Base
	return s
}

// OnActive applies the given styles for the pseudo :active
func (s *Style) OnActive(style *Style) *Style {
	s.Pseudos[":active"] = style.Base
	return s
}

// OnFocus applies the given styles for the pseudo :focus
func (s *Style) OnFocus(style *Style) *Style {
	s.Pseudos[":focus"] = style.Base
	return s
}

// OnFocusWithin applies the given styles for the pseudo :focus-within
func (s *Style) OnFocusWithin(style *Style) *Style {
	s.Pseudos[":focus-within"] = style.Base
	return s
}

// OnFocusVisible applies the given styles for the pseudo :focus-visible
func (s *Style) OnFocusVisible(style *Style) *Style {
	s.Pseudos[":focus-visible"] = style.Base
	return s
}
