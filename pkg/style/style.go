//go:build js && wasm

// Package style provides type-safe CSS styling for Vortex components.
//
// This package includes comprehensive styling utilities organized into logical groups:
//
//   - Box Model: Display, Width, Height, Margin, Padding, Border (see box_model.go)
//   - Flexbox: Flex properties, direction, alignment, justify-content (see flex.go)
//   - Typography: Font, text color, alignment, decoration (see typography.go)
//   - Appearance: Background, opacity, visual effects (see appearance.go)
//   - Interactive: Hover, focus, active pseudo-classes (see pseudo.go)
//   - Layout: Position, overflow, cursor, z-index (see misc.go)
//   - Media Queries: Responsive design utilities (see media_query.go)
//
// Basic Usage:
//
//   style := style.New(
//       style.Display(style.DisplayFlex),
//       style.JustifyContent(style.JustifyContentCenter),
//       style.FontSize("16px"),
//       style.OnHover(style.BackgroundColor("#f0f0f0")),
//   )
//
// The style system uses functional options for composability and type safety.

// Core style creation and manipulation functions are defined in this file.
//
// For specific styling utilities, see:
//   - Display, margins, padding, borders: See box_model.go functions
//   - Flexbox layouts: See flex.go functions
//   - Text and font styling: See typography.go functions
//   - Hover and focus states: See pseudo.go functions
//   - Background and appearance: See appearance.go functions
//   - Positioning and cursor: See misc.go functions

package style

// Property is a simple alias for a CSS property
// Ex: "color" -> "blue"
type Property map[string]string

// Style is the object that contains the full definition of a the style of a component
type Style struct {
	Base         Property
	Pseudos      map[string]Property // ":hover", ":active", ":focus", etc.
	MediaQueries map[string]Property // "screen and (max-width: 768px)", "screen and (min-width: 769px)"

	// For caching
	className string
	css       string
}

// StyleOption is a function that modifies the style
type StyleOption func(*Style)

// New creates a new style object applying the given options
func New(options ...StyleOption) *Style {
	s := &Style{
		Base:         make(Property),
		Pseudos:      make(map[string]Property),
		MediaQueries: make(map[string]Property),
	}
	for _, option := range options {
		s.Update(option)
	}
	return s
}

func (s *Style) List() []StyleOption {
	options := []StyleOption{}
	for _, option := range s.Base {
		options = append(options, CustomStyle(option, s.Base[option]))
	}
	for _, pseudo := range s.Pseudos {
		for key, value := range pseudo {
			options = append(options, CustomStyle(key, value))
		}
	}
	for _, mediaQuery := range s.MediaQueries {
		for key, value := range mediaQuery {
			options = append(options, CustomStyle(key, value))
		}
	}
	return options
}

func Extend(baseStyle *Style, options ...StyleOption) *Style {
	// 1. Create a deep copy of the base style
	s := &Style{
		Base:         make(Property),
		Pseudos:      make(map[string]Property),
		MediaQueries: make(map[string]Property),
	}

	// 2. Deep copy the base properties
	for key, value := range baseStyle.Base {
		s.Base[key] = value
	}

	// 3. Deep copy the pseudo properties
	for pseudo, properties := range baseStyle.Pseudos {
		newPseudoProps := make(Property)
		for key, value := range properties {
			newPseudoProps[key] = value
		}
		s.Pseudos[pseudo] = newPseudoProps
	}

	// 4. Deep copy the media queries
	for mediaQuery, properties := range baseStyle.MediaQueries {
		newMediaProps := make(Property)
		for key, value := range properties {
			newMediaProps[key] = value
		}
		s.MediaQueries[mediaQuery] = newMediaProps
	}

	// 5. Apply the options
	s.Update(options...)
	return s
}

func (s *Style) Update(options ...StyleOption) {
	for _, option := range options {
		option(s)
	}
}

// Function to apply a style which is not in the function already defined
func CustomStyle(property string, value string) StyleOption {
	return func(s *Style) {
		s.Base[property] = value
	}
}
